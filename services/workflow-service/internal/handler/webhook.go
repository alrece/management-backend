package handler

import (
	"net/http"

	"management-backend/services/workflow-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// WebhookPayload n8n 回调请求结构
type WebhookPayload struct {
	ExecutionID string `json:"executionId"`
	WorkflowID  string `json:"workflowId"`
	Status      string `json:"status"` // success/error
	Data        string `json:"data"`
	DurationMs  int64  `json:"durationMs"`
	ErrorMsg    string `json:"errorMsg,omitempty"`
}

// WebhookHandler 接收 n8n 执行回调
type WebhookHandler struct {
	instanceRepo repository.InstanceRepo
	workflowRepo repository.WorkflowRepo
}

// NewWebhookHandler 创建回调处理器
func NewWebhookHandler(instRepo repository.InstanceRepo, wfRepo repository.WorkflowRepo) *WebhookHandler {
	return &WebhookHandler{
		instanceRepo: instRepo,
		workflowRepo: wfRepo,
	}
}

// Callback 接收 n8n 执行结果回调
func (h *WebhookHandler) Callback(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	status := "SUCCESS"
	if payload.Status == "error" {
		status = "FAILED"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"executionId": payload.ExecutionID,
			"status":      status,
		},
	})
}
