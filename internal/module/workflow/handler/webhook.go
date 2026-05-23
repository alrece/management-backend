package handler

import (
	"encoding/json"
	"strconv"

	"management-backend/internal/module/workflow/service"
	pkgmiddleware "management-backend/pkg/middleware"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// WebhookHandler Webhook 触发入口处理器
type WebhookHandler struct {
	workflowSvc service.WorkflowService
}

func NewWebhookHandler(workflowSvc service.WorkflowService) *WebhookHandler {
	return &WebhookHandler{workflowSvc: workflowSvc}
}

// Handle Webhook 触发入口 POST /api/wf/webhook/:tenantId/:workflowId/*path
func (h *WebhookHandler) Handle(c *gin.Context) {
	workflowID, err := strconv.ParseInt(c.Param("workflowId"), 10, 64)
	if err != nil {
		response.FailWithStatus(c, 400, 400, "无效的工作流ID")
		return
	}

	// 读取请求体作为工作流输入
	var input json.RawMessage
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&input); err != nil {
			input = nil
		}
	}

	// 注入租户上下文（从 URL 参数获取）
	tenantID, _ := strconv.ParseInt(c.Param("tenantId"), 10, 64)
	ctx := pkgmiddleware.WithTenantID(c.Request.Context(), tenantID)
	c.Request = c.Request.WithContext(ctx)

	execution, err := h.workflowSvc.Execute(ctx, workflowID, input)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, execution)
}
