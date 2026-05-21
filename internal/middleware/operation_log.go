package middleware

import (
	"bytes"
	"io"
	"time"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/pkg/snowflake"

	"github.com/gin-gonic/gin"
)

var logSvc service.LogService

// InitLogService 初始化日志服务
func InitLogService(svc service.LogService) {
	logSvc = svc
}

// OperationLog 操作日志中间件（仅记录写操作）
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			c.Next()
			return
		}

		start := time.Now()

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		if logSvc != nil {
			go func() {
				userID := GetUserID(c)
				username := GetString(c, "username")
				tenantID := GetTenantID(c)
				requestID := GetString(c, "request_id")

				businessType := 0
				switch method {
				case "POST":
					businessType = 1
				case "PUT":
					businessType = 2
				case "DELETE":
					businessType = 3
				}

				status := 0
				if c.Writer.Status() >= 400 {
					status = 1
				}

				log := &smodel.SysOperLog{
					Title:        c.FullPath(),
					BusinessType: businessType,
					Method:       method,
					RequestURL:   c.Request.URL.String(),
					OperIP:       c.ClientIP(),
					OperUserID:   userID,
					OperName:     username,
					TenantID:     tenantID,
					RequestID:    requestID,
					Status:       status,
					OperTime:     start,
				}
				log.ID = snowflake.NextID()

				_ = logSvc.CreateOperLog(c.Request.Context(), log)
			}()
		}
	}
}

// GetString 从 context 获取字符串值
func GetString(c *gin.Context, key string) string {
	val, _ := c.Get(key)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
