package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// MenuHandler 菜单接口处理器
type MenuHandler struct {
	svc service.MenuService
}

// NewMenuHandler 构造函数
func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req smodel.MenuCreateReq
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

func (h *MenuHandler) Update(c *gin.Context) {
	var req smodel.MenuUpdateReq
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

func (h *MenuHandler) Delete(c *gin.Context) {
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

func (h *MenuHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, tree)
}
