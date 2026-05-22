package router

import (
	"time"

	"management-backend/internal/auth"
	"management-backend/internal/authz"
	"management-backend/internal/middleware"
	pkgmiddleware "management-backend/pkg/middleware"
	"management-backend/internal/module/system"
	"management-backend/internal/module/system/handler"
	"management-backend/internal/module/system/repository"
	"management-backend/internal/module/system/service"
	grpchandler "management-backend/internal/grpc/handler"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Setup 初始化路由和依赖注入（Repository → Service → Handler）
// 返回清理函数，调用方应在关闭时调用以释放资源
func Setup(engine *gin.Engine, db *gorm.DB, rdb *redis.Client, logger *zap.Logger) func() {
	// 初始化认证组件
	middleware.InitAuth(rdb)
	middleware.InitTenant()

	// 初始化 Casbin 权限管理（每小时淘汰 30 分钟未使用的 enforcer）
	enforcerMgr := authz.NewEnforcerManager(rdb)
	enforcerMgr.StartEviction(time.Hour, 30*time.Minute)
	middleware.InitPermission(enforcerMgr)

	loginLock := auth.NewLoginLock(rdb)
	blacklist := auth.NewBlacklist(rdb)

	// 依赖注入 — Repository
	userRepo := repository.NewUserRepo(db)
	tenantRepo := repository.NewTenantRepo(db)
	menuRepo := repository.NewMenuRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	deptRepo := repository.NewDeptRepo(db)
	postRepo := repository.NewPostRepo(db)
	dictTypeRepo := repository.NewDictTypeRepo(db)
	dictDataRepo := repository.NewDictDataRepo(db)
	paramRepo := repository.NewParamRepo(db)
	noticeRepo := repository.NewNoticeRepo(db)
	operLogRepo := repository.NewOperLogRepo(db)
	loginLogRepo := repository.NewLoginLogRepo(db)
	fileRepo := repository.NewFileRepo(db)

	// 依赖注入 — Service
	userSvc := service.NewUserService(userRepo)
	tenantSvc := service.NewTenantService(tenantRepo)
	menuSvc := service.NewMenuService(menuRepo)
	roleSvc := service.NewRoleService(roleRepo, menuRepo, enforcerMgr)
	deptSvc := service.NewDeptService(deptRepo)
	postSvc := service.NewPostService(postRepo)
	dictSvc := service.NewDictService(dictTypeRepo, dictDataRepo, rdb)
	paramSvc := service.NewParamService(paramRepo, rdb)
	noticeSvc := service.NewNoticeService(noticeRepo)
	logSvc := service.NewLogService(operLogRepo, loginLogRepo)
	fileSvc := service.NewFileService(fileRepo)

	// 初始化日志中间件
	middleware.InitLogService(logSvc)

	// 依赖注入 — Handler
	userHandler := handler.NewUserHandler(userSvc)
	tenantHandler := handler.NewTenantHandler(tenantSvc)
	menuHandler := handler.NewMenuHandler(menuSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)
	deptHandler := handler.NewDeptHandler(deptSvc)
	postHandler := handler.NewPostHandler(postSvc)
	dictHandler := handler.NewDictHandler(dictSvc)
	paramHandler := handler.NewParamHandler(paramSvc)
	noticeHandler := handler.NewNoticeHandler(noticeSvc)
	logHandler := handler.NewLogHandler(logSvc)
	profileHandler := handler.NewProfileHandler(userSvc, blacklist)
	authHandler := handler.NewAuthHandler(userSvc, loginLock, logger)
	fileHandler := handler.NewFileHandler(fileSvc)

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

		// 动态菜单 API
		authed.GET("/system/auth/get-permission-info", handler.NewPermissionInfoHandler(userSvc, roleSvc, menuSvc, roleRepo))

		// 操作日志中间件（写操作）
		authed.Use(middleware.OperationLog())

		system.RegisterRoutes(authed, userHandler, tenantHandler, roleHandler, menuHandler,
			deptHandler, postHandler, dictHandler, paramHandler, noticeHandler, logHandler, profileHandler, fileHandler)
	}

	return func() {
		enforcerMgr.Stop()
	}
}

// GetGRPCServices 创建 gRPC handler（复用相同 DI 链）
func GetGRPCServices(db *gorm.DB, rdb *redis.Client, logger *zap.Logger) (
	*grpchandler.UserGRPCHandler,
	*grpchandler.RoleGRPCHandler,
	*grpchandler.MenuGRPCHandler,
) {
	enforcerMgr := authz.NewEnforcerManager(rdb)

	userSvc := service.NewUserService(repository.NewUserRepo(db))
	roleSvc := service.NewRoleService(repository.NewRoleRepo(db), repository.NewMenuRepo(db), enforcerMgr)
	menuSvc := service.NewMenuService(repository.NewMenuRepo(db))

	return grpchandler.NewUserGRPCHandler(userSvc),
		grpchandler.NewRoleGRPCHandler(roleSvc),
		grpchandler.NewMenuGRPCHandler(menuSvc)
}
