package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// RoleHandler 角色接口处理器
type RoleHandler struct {
	svc service.RoleService
}

// NewRoleHandler 构造函数
func NewRoleHandler(svc service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req smodel.RoleCreateReq
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

func (h *RoleHandler) Update(c *gin.Context) {
	var req smodel.RoleUpdateReq
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

func (h *RoleHandler) Delete(c *gin.Context) {
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

func (h *RoleHandler) Get(c *gin.Context) {
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

func (h *RoleHandler) Page(c *gin.Context) {
	var req smodel.RolePageReq
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

func (h *RoleHandler) AssignMenus(c *gin.Context) {
	var req smodel.RoleMenuAssignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	if err := h.svc.AssignMenus(c.Request.Context(), &req, middleware.GetUserID(c), middleware.GetTenantID(c)); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "分配成功")
}
