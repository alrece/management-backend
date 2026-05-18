package system

import (
	"management-backend/internal/module/system/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册系统模块路由
func RegisterRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := rg.Group("/system/user")
	{
		users.POST("", userHandler.Create)
		users.PUT("", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.GET("/:id", userHandler.Get)
		users.GET("/page", userHandler.Page)
	}
}
