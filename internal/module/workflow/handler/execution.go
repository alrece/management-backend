package handler

import (
	"strconv"

	wfmodel "management-backend/internal/module/workflow/model"
	"management-backend/internal/module/workflow/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExecutionHandler 执行实例 API 处理器
type ExecutionHandler struct {
	svc service.ExecutionService
}

func NewExecutionHandler(svc service.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{svc: svc}
}

// Page 执行实例分页 GET /workflow/execution/page
func (h *ExecutionHandler) Page(c *gin.Context) {
	var req wfmodel.ExecutionPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	list, total, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkPage(c, list, total, req.GetPage(), req.GetPageSize())
}

// Get 执行详情 GET /workflow/execution/:id
func (h *ExecutionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	data, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, data)
}

// Logs 节点执行日志 GET /workflow/execution/:id/logs
func (h *ExecutionHandler) Logs(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	logs, err := h.svc.GetLogs(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, logs)
}
