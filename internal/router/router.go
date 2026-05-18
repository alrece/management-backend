package router

import (
	"management-backend/internal/middleware"
	"management-backend/internal/module/system"
	"management-backend/internal/module/system/handler"
	"management-backend/internal/module/system/repository"
	"management-backend/internal/module/system/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Setup 初始化路由和依赖注入（Repository → Service → Handler）
func Setup(engine *gin.Engine, db *gorm.DB, rdb *redis.Client, logger *zap.Logger) {
	// 依赖注入
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	// 公开路由
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 认证路由组
	authed := engine.Group("/api")
	authed.Use(middleware.Auth())
	authed.Use(middleware.Tenant())
	{
		system.RegisterRoutes(authed, userHandler)
	}

	_ = rdb
	_ = logger
}
