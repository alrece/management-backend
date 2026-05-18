package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

// Handler WebSocket 处理器
type Handler struct {
	hub    *Hub
	logger *zap.Logger
}

// NewHandler 创建 WebSocket Handler
func NewHandler(hub *Hub, logger *zap.Logger) *Handler {
	return &Handler{hub: hub, logger: logger}
}

// HandleConnection 处理 WebSocket 连接升级 GET /api/ws
func (h *Handler) HandleConnection(c *gin.Context) {
	userIDVal, exists := c.Get("userId")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	tenantIDVal, _ := c.Get("tenantId")
	userID := userIDVal.(int64)
	tenantID := int64(0)
	if tid, ok := tenantIDVal.(int64); ok {
		tenantID = tid
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		h.logger.Error("WebSocket 升级失败", zap.Error(err))
		return
	}

	client := NewClient(h.hub, conn, userID, tenantID)
	h.hub.Register(client)

	go client.WritePump(c.Request.Context())
	go client.ReadPump(c.Request.Context())
}

// GetOnlineUsers 获取在线用户列表 GET /api/system/online/list
func (h *Handler) GetOnlineUsers(c *gin.Context) {
	ids := h.hub.GetOnlineUserIDs()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": ids,
		"msg":  "success",
	})
}

// KickUser 强制踢出用户 DELETE /api/system/online/:id
func (h *Handler) KickUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效用户ID"})
		return
	}
	h.hub.KickUser(id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已踢出"})
}
