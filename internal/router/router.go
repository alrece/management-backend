package router

import (
	"management-backend/internal/auth"
	"management-backend/internal/middleware"
	pkgmiddleware "management-backend/pkg/middleware"
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
	// 初始化认证组件
	middleware.InitAuth(rdb)
	middleware.InitTenant()
	loginLock := auth.NewLoginLock(rdb)

	// 依赖注入
	userRepo := repository.NewUserRepo(db)
	tenantRepo := repository.NewTenantRepo(db)
	userSvc := service.NewUserService(userRepo)
	tenantSvc := service.NewTenantService(tenantRepo)
	userHandler := handler.NewUserHandler(userSvc)
	tenantHandler := handler.NewTenantHandler(tenantSvc)
	authHandler := handler.NewAuthHandler(userSvc, loginLock, logger)

	// 全局中间件
	engine.Use(pkgmiddleware.SecurityHeaders())

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 公开路由（无需认证）
	public := engine.Group("/api/auth")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
	}

	// 认证路由组
	authed := engine.Group("/api")
	authed.Use(middleware.Auth())
	authed.Use(middleware.Tenant())
	{
		authed.POST("/auth/logout", authHandler.Logout)
		system.RegisterRoutes(authed, userHandler, tenantHandler)
	}
}
