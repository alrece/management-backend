package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	wfmodel "management-backend/internal/module/workflow/model"
	"management-backend/internal/module/workflow/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类 API 处理器
type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// Create 创建分类 POST /workflow/category
func (h *CategoryHandler) Create(c *gin.Context) {
	var req wfmodel.CategoryCreateReq
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

// Update 更新分类 PUT /workflow/category
func (h *CategoryHandler) Update(c *gin.Context) {
	var req wfmodel.CategoryUpdateReq
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

// Delete 删除分类 DELETE /workflow/category/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
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

// Tree 分类树 GET /workflow/category/tree
func (h *CategoryHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, tree)
}
