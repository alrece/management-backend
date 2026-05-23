package handler

import (
	"encoding/json"
	"strconv"

	"management-backend/internal/middleware"
	wfmodel "management-backend/internal/module/workflow/model"
	"management-backend/internal/module/workflow/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// WorkflowHandler 工作流 API 处理器
type WorkflowHandler struct {
	svc service.WorkflowService
}

func NewWorkflowHandler(svc service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{svc: svc}
}

// Create 创建工作流 POST /workflow/workflow
func (h *WorkflowHandler) Create(c *gin.Context) {
	var req wfmodel.WorkflowCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	id, err := h.svc.Create(c.Request.Context(), &req,
		middleware.GetUserID(c), middleware.GetTenantID(c))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, id)
}

// Update 更新工作流 PUT /workflow/workflow
func (h *WorkflowHandler) Update(c *gin.Context) {
	var req wfmodel.WorkflowUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	if err := h.svc.Update(c.Request.Context(), &req, middleware.GetUserID(c)); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "更新成功")
}

// Delete 删除工作流 DELETE /workflow/workflow/:id
func (h *WorkflowHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "删除成功")
}

// Get 获取工作流详情 GET /workflow/workflow/:id
func (h *WorkflowHandler) Get(c *gin.Context) {
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

// Page 分页查询 GET /workflow/workflow/page
func (h *WorkflowHandler) Page(c *gin.Context) {
	var req wfmodel.WorkflowPageReq
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

// Activate 激活工作流 PUT /workflow/workflow/:id/activate
func (h *WorkflowHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	if err := h.svc.Activate(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "激活成功")
}

// Deactivate 停用工作流 PUT /workflow/workflow/:id/deactivate
func (h *WorkflowHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	if err := h.svc.Deactivate(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "停用成功")
}

// Execute 手动执行 POST /workflow/workflow/:id/execute
func (h *WorkflowHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	var input json.RawMessage
	c.ShouldBindJSON(&input)

	execution, err := h.svc.Execute(c.Request.Context(), id, input)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, execution)
}

// Schemas 获取所有节点 Schema GET /workflow/workflow/nodes/schemas
func (h *WorkflowHandler) Schemas(c *gin.Context) {
	schemas := h.svc.GetAllSchemas(c.Request.Context())
	response.Ok(c, schemas)
}
