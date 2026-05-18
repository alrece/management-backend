package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] %v\n%s", r, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					gin.H{"code": 500, "msg": "服务器内部错误"})
			}
		}()
		c.Next()
	}
}
