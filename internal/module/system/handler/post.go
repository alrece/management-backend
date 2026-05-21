package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// PostHandler 岗位接口处理器
type PostHandler struct {
	svc service.PostService
}

// NewPostHandler 构造函数
func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req smodel.PostCreateReq
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

func (h *PostHandler) Update(c *gin.Context) {
	var req smodel.PostUpdateReq
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

func (h *PostHandler) Delete(c *gin.Context) {
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

func (h *PostHandler) Page(c *gin.Context) {
	var req smodel.PostPageReq
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
