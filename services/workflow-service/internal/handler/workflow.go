package handler

import (
	"net/http"
	"strconv"

	"management-backend/services/workflow-service/internal/model"
	"management-backend/services/workflow-service/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	wfSvc service.WorkflowService
	catSvc service.CategoryService
}

func NewWorkflowHandler(wfSvc service.WorkflowService, catSvc service.CategoryService) *WorkflowHandler {
	return &WorkflowHandler{wfSvc: wfSvc, catSvc: catSvc}
}

// ===== 分类 =====

func (h *WorkflowHandler) CreateCategory(c *gin.Context) {
	var req model.CategoryCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	userID, _ := c.Get("userId")
	tenantID, _ := c.Get("tenantId")
	id, err := h.catSvc.Create(c.Request.Context(), &req, toInt64(userID), toInt64(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": id})
}

func (h *WorkflowHandler) UpdateCategory(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	var req model.CategoryUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.catSvc.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *WorkflowHandler) DeleteCategory(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	if err := h.catSvc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *WorkflowHandler) CategoryTree(c *gin.Context) {
	tenantID, _ := c.Get("tenantId")
	tree, err := h.catSvc.Tree(c.Request.Context(), toInt64(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tree})
}

// ===== 工作流 =====

func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var req model.WorkflowCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	userID, _ := c.Get("userId")
	tenantID, _ := c.Get("tenantId")
	id, err := h.wfSvc.Create(c.Request.Context(), &req, toInt64(userID), toInt64(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": id})
}

func (h *WorkflowHandler) UpdateWorkflow(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	var req model.WorkflowUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.wfSvc.Update(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	if err := h.wfSvc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	resp, err := h.wfSvc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": resp})
}

func (h *WorkflowHandler) PageWorkflow(c *gin.Context) {
	var req model.WorkflowPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	tenantID, _ := c.Get("tenantId")
	list, total, err := h.wfSvc.Page(c.Request.Context(), &req, toInt64(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": list, "total": total}})
}

func (h *WorkflowHandler) ActivateWorkflow(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	if err := h.wfSvc.Activate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *WorkflowHandler) DeactivateWorkflow(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	if err := h.wfSvc.Deactivate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *WorkflowHandler) ExecuteWorkflow(c *gin.Context) {
	id := toInt64FromParam(c, "id")
	var req model.ExecuteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	userID, _ := c.Get("userId")
	tenantID, _ := c.Get("tenantId")
	instID, err := h.wfSvc.Execute(c.Request.Context(), id, &req, toInt64(userID), toInt64(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": instID})
}

// ===== 辅助 =====

func toInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	default:
		return 0
	}
}

func toInt64FromParam(c *gin.Context, key string) int64 {
	s := c.Param(key)
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
