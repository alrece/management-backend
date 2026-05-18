package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// TenantHandler 租户接口处理器
type TenantHandler struct {
	svc service.TenantService
}

// NewTenantHandler 构造函数
func NewTenantHandler(svc service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

// Create 创建租户 POST /api/system/tenant
func (h *TenantHandler) Create(c *gin.Context) {
	var req smodel.TenantCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	id, err := h.svc.Create(c.Request.Context(), &req, middleware.GetUserID(c))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, id)
}

// Update 更新租户 PUT /api/system/tenant
func (h *TenantHandler) Update(c *gin.Context) {
	var req smodel.TenantUpdateReq
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

// Delete 删除租户 DELETE /api/system/tenant/:id
func (h *TenantHandler) Delete(c *gin.Context) {
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

// Get 获取租户详情 GET /api/system/tenant/:id
func (h *TenantHandler) Get(c *gin.Context) {
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

// Page 分页查询 GET /api/system/tenant/page
func (h *TenantHandler) Page(c *gin.Context) {
	var req smodel.TenantPageReq
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
