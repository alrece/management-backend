package handler

import (
	"fmt"
	"net/http"
	"strconv"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/errcode"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件管理处理器
type FileHandler struct {
	svc service.FileService
}

// NewFileHandler 构造函数
func NewFileHandler(svc service.FileService) *FileHandler {
	return &FileHandler{svc: svc}
}

// Upload 上传文件 POST /api/system/file
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := h.svc.Upload(c)
	if err != nil {
		response.Fail(c, errcode.FileUploadFail, err.Error())
		return
	}
	response.Ok(c, file)
}

// Delete 删除文件 DELETE /api/system/file/:id
func (h *FileHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, errcode.FileNotFound, err.Error())
		return
	}
	response.OkMsg(c, "删除成功")
}

// Get 获取文件详情 GET /api/system/file/:id
func (h *FileHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	resp, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, errcode.FileNotFound, err.Error())
		return
	}
	response.Ok(c, resp)
}

// Page 文件分页查询 GET /api/system/file/page
func (h *FileHandler) Page(c *gin.Context) {
	var req smodel.FilePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, 400, "参数校验失败: "+err.Error())
		return
	}
	list, total, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OkPage(c, list, total, req.Page, req.Limit())
}

// Download 下载文件 GET /api/system/file/:id/download
func (h *FileHandler) Download(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	reader, file, err := h.svc.Download(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, errcode.FileNotFound, err.Error())
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Header("Content-Type", file.ContentType)
	c.DataFromReader(http.StatusOK, file.FileSize, file.ContentType, reader, nil)
}
