package handler

import (
	"strconv"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// LogHandler 日志处理器
type LogHandler struct {
	svc service.LogService
}

func NewLogHandler(svc service.LogService) *LogHandler {
	return &LogHandler{svc: svc}
}

func (h *LogHandler) PageOperLog(c *gin.Context) {
	var req smodel.OperLogPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	list, total, err := h.svc.PageOperLog(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkPage(c, list, total, req.GetPage(), req.GetPageSize())
}

func (h *LogHandler) DeleteOperLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	if err := h.svc.DeleteOperLog(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "删除成功")
}

func (h *LogHandler) CleanOperLog(c *gin.Context) {
	if err := h.svc.CleanOperLog(c.Request.Context()); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "清空成功")
}

func (h *LogHandler) PageLoginLog(c *gin.Context) {
	var req smodel.LoginLogPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	list, total, err := h.svc.PageLoginLog(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkPage(c, list, total, req.GetPage(), req.GetPageSize())
}

func (h *LogHandler) DeleteLoginLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	if err := h.svc.DeleteLoginLog(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "删除成功")
}

func (h *LogHandler) CleanLoginLog(c *gin.Context) {
	if err := h.svc.CleanLoginLog(c.Request.Context()); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkMsg(c, "清空成功")
}
