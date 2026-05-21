package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// DeptHandler 部门接口处理器
type DeptHandler struct {
	svc service.DeptService
}

// NewDeptHandler 构造函数
func NewDeptHandler(svc service.DeptService) *DeptHandler {
	return &DeptHandler{svc: svc}
}

func (h *DeptHandler) Create(c *gin.Context) {
	var req smodel.DeptCreateReq
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

func (h *DeptHandler) Update(c *gin.Context) {
	var req smodel.DeptUpdateReq
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

func (h *DeptHandler) Delete(c *gin.Context) {
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

func (h *DeptHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, tree)
}
