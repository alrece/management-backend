package handler

import (
	"strconv"

	"management-backend/internal/middleware"
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// ParamHandler 参数处理器
type ParamHandler struct {
	svc service.ParamService
}

func NewParamHandler(svc service.ParamService) *ParamHandler {
	return &ParamHandler{svc: svc}
}

func (h *ParamHandler) Create(c *gin.Context) {
	var req smodel.ParamCreateReq
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

func (h *ParamHandler) Update(c *gin.Context) {
	var req smodel.ParamUpdateReq
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

func (h *ParamHandler) Delete(c *gin.Context) {
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

func (h *ParamHandler) Get(c *gin.Context) {
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

func (h *ParamHandler) GetByKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.Fail(c, 400, "参数key不能为空")
		return
	}
	val, err := h.svc.GetByKey(c.Request.Context(), key)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Ok(c, val)
}

func (h *ParamHandler) Page(c *gin.Context) {
	var req smodel.ParamPageReq
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
