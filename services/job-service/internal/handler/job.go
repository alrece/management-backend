package handler

import (
	"net/http"
	"strconv"

	"management-backend/services/job-service/internal/model"
	"management-backend/services/job-service/internal/service"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	svc service.JobService
}

func NewJobHandler(svc service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

func (h *JobHandler) Create(c *gin.Context) {
	var req model.JobTaskCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败: " + err.Error()})
		return
	}
	userID := getUserID(c)
	tenantID := getTenantID(c)

	id, err := h.svc.Create(c.Request.Context(), &req, userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": id})
}

func (h *JobHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效ID"})
		return
	}
	var req model.JobTaskUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败"})
		return
	}
	if err := h.svc.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

func (h *JobHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效ID"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

func (h *JobHandler) Page(c *gin.Context) {
	var req model.JobTaskPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	tenantID := getTenantID(c)
	list, total, err := h.svc.Page(c.Request.Context(), &req, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": list, "total": total}})
}

func (h *JobHandler) Trigger(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效ID"})
		return
	}
	tenantID := getTenantID(c)
	logID, err := h.svc.Trigger(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": logID})
}

func (h *JobHandler) ExecLogPage(c *gin.Context) {
	var req model.ExecLogPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	tenantID := getTenantID(c)
	list, total, err := h.svc.ExecLogPage(c.Request.Context(), &req, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": list, "total": total}})
}

func getUserID(c *gin.Context) int64 {
	if v, ok := c.Get("userId"); ok {
		return v.(int64)
	}
	return 0
}

func getTenantID(c *gin.Context) int64 {
	if v, ok := c.Get("tenantId"); ok {
		return v.(int64)
	}
	return 0
}
