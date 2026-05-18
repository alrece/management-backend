package handler

import (
	"net/http"

	"management-backend/internal/module/system/service"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	svc *service.CaptchaService
}

// NewCaptchaHandler 构造函数
func NewCaptchaHandler(svc *service.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{svc: svc}
}

// GetCaptcha 获取图形验证码 GET /api/auth/captcha
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	result, err := h.svc.Generate(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, "验证码生成失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result, "msg": "success"})
}
