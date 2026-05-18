package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Hub 管理所有 WebSocket 连接
type Hub struct {
	mu         sync.RWMutex
	clients    map[int64]*Client // userID -> Client
	register   chan *Client
	unregister chan *Client
	rdb        *redis.Client
	logger     *zap.Logger
}

// NewHub 创建 Hub
func NewHub(rdb *redis.Client, logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		rdb:        rdb,
		logger:     logger,
	}
}

// Run 启动 Hub 事件循环
func (h *Hub) Run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.removeAll()
			return
		case client := <-h.register:
			h.mu.Lock()
			if old, ok := h.clients[client.userID]; ok {
				old.Close()
			}
			h.clients[client.userID] = client
			h.mu.Unlock()
			h.addOnlineUser(ctx, client.userID, client.tenantID)
		case client := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[client.userID]; ok && existing == client {
				delete(h.clients, client.userID)
			}
			h.mu.Unlock()
			h.removeOnlineUser(ctx, client.userID, client.tenantID)
		case <-ticker.C:
			h.heartbeat()
		}
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	select {
	case h.unregister <- client:
	default:
	}
}

// KickUser 强制踢出用户
func (h *Hub) KickUser(userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client, ok := h.clients[userID]; ok {
		client.Close()
		delete(h.clients, userID)
	}
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// OnlineCount 在线人数
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Broadcast 广播消息
func (h *Hub) Broadcast(msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		c.Send(msg)
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, msg []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if c, ok := h.clients[userID]; ok {
		c.Send(msg)
		return true
	}
	return false
}

// GetOnlineUserIDs 获取所有在线用户 ID
func (h *Hub) GetOnlineUserIDs() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]int64, 0, len(h.clients))
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

func (h *Hub) heartbeat() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for userID, client := range h.clients {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := client.conn.Ping(ctx)
		cancel()
		if err != nil {
			client.Close()
			delete(h.clients, userID)
		}
	}
}

func (h *Hub) removeAll() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for userID, client := range h.clients {
		client.Close()
		delete(h.clients, userID)
	}
}

func (h *Hub) addOnlineUser(ctx context.Context, userID, _ int64) {
	if h.rdb != nil {
		h.rdb.SAdd(ctx, "ws:online", userID)
	}
}

func (h *Hub) removeOnlineUser(ctx context.Context, userID, _ int64) {
	if h.rdb != nil {
		h.rdb.SRem(ctx, "ws:online", userID)
	}
}
