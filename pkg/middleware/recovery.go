package middleware

import (
	"net/http"

	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 使用 Zap 记录 panic，替代默认 recovery
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		logger.Error("panic recovered",
			zap.String("request_id", GetRequestID(c.Request.Context())),
			zap.Any("error", err),
			zap.Stack("stack"),
		)
		response.FailWithStatus(c, http.StatusInternalServerError, 500, "服务器内部错误")
		c.Abort()
	})
}
