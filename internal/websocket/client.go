package websocket

import (
	"context"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

// Client WebSocket 客户端连接
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	userID   int64
	tenantID int64
	send     chan []byte
	mu       sync.Mutex
}

// NewClient 创建 WebSocket 客户端
func NewClient(hub *Hub, conn *websocket.Conn, userID, tenantID int64) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		userID:   userID,
		tenantID: tenantID,
		send:     make(chan []byte, 256),
	}
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	c.conn.SetReadLimit(4096)
	for {
		_, _, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
	}
}

// WritePump 向客户端发送消息
func (c *Client) WritePump(ctx context.Context) {
	defer c.conn.Close(websocket.StatusNormalClosure, "")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

// Send 发送消息到客户端
func (c *Client) Send(msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case c.send <- msg:
	default:
		close(c.send)
	}
}

// Close 关闭连接
func (c *Client) Close() {
	c.conn.Close(websocket.StatusNormalClosure, "")
}

// UserID 返回用户 ID
func (c *Client) UserID() int64 {
	return c.userID
}
