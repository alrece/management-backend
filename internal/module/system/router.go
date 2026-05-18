package system

import (
	"management-backend/internal/middleware"
	"management-backend/internal/module/system/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册系统模块路由
func RegisterRoutes(rg *gin.RouterGroup,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	deptHandler *handler.DeptHandler,
	postHandler *handler.PostHandler,
	dictHandler *handler.DictHandler,
	paramHandler *handler.ParamHandler,
	noticeHandler *handler.NoticeHandler,
	logHandler *handler.LogHandler,
	profileHandler *handler.ProfileHandler,
	fileHandler *handler.FileHandler,
) {
	users := rg.Group("/system/user")
	{
		users.POST("", userHandler.Create)
		users.PUT("", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.GET("/:id", userHandler.Get)
		users.GET("/page", userHandler.Page)
		users.GET("/profile", profileHandler.GetProfile)
		users.PUT("/profile", profileHandler.UpdateProfile)
		users.PUT("/password", profileHandler.UpdatePassword)
	}

	tenants := rg.Group("/system/tenant")
	{
		tenants.POST("", tenantHandler.Create)
		tenants.PUT("", tenantHandler.Update)
		tenants.DELETE("/:id", tenantHandler.Delete)
		tenants.GET("/:id", tenantHandler.Get)
		tenants.GET("/page", tenantHandler.Page)
	}

	roles := rg.Group("/system/role")
	{
		roles.POST("", middleware.RequirePermission("system:role:create"), roleHandler.Create)
		roles.PUT("", middleware.RequirePermission("system:role:update"), roleHandler.Update)
		roles.DELETE("/:id", middleware.RequirePermission("system:role:delete"), roleHandler.Delete)
		roles.GET("/:id", middleware.RequirePermission("system:role:query"), roleHandler.Get)
		roles.GET("/page", middleware.RequirePermission("system:role:query"), roleHandler.Page)
		roles.PUT("/assign-menus", middleware.RequirePermission("system:role:update"), roleHandler.AssignMenus)
	}

	menus := rg.Group("/system/menu")
	{
		menus.POST("", middleware.RequirePermission("system:menu:create"), menuHandler.Create)
		menus.PUT("", middleware.RequirePermission("system:menu:update"), menuHandler.Update)
		menus.DELETE("/:id", middleware.RequirePermission("system:menu:delete"), menuHandler.Delete)
		menus.GET("/tree", middleware.RequirePermission("system:menu:query"), menuHandler.Tree)
	}

	depts := rg.Group("/system/dept")
	{
		depts.POST("", middleware.RequirePermission("system:dept:create"), deptHandler.Create)
		depts.PUT("", middleware.RequirePermission("system:dept:update"), deptHandler.Update)
		depts.DELETE("/:id", middleware.RequirePermission("system:dept:delete"), deptHandler.Delete)
		depts.GET("/tree", middleware.RequirePermission("system:dept:query"), deptHandler.Tree)
	}

	posts := rg.Group("/system/post")
	{
		posts.POST("", middleware.RequirePermission("system:post:create"), postHandler.Create)
		posts.PUT("", middleware.RequirePermission("system:post:update"), postHandler.Update)
		posts.DELETE("/:id", middleware.RequirePermission("system:post:delete"), postHandler.Delete)
		posts.GET("/page", middleware.RequirePermission("system:post:query"), postHandler.Page)
	}

	dictTypes := rg.Group("/system/dict-type")
	{
		dictTypes.POST("", dictHandler.CreateType)
		dictTypes.PUT("", dictHandler.UpdateType)
		dictTypes.DELETE("/:id", dictHandler.DeleteType)
		dictTypes.GET("/:id", dictHandler.GetType)
		dictTypes.GET("/page", dictHandler.PageType)
	}

	dictData := rg.Group("/system/dict-data")
	{
		dictData.POST("", dictHandler.CreateData)
		dictData.PUT("", dictHandler.UpdateData)
		dictData.DELETE("/:id", dictHandler.DeleteData)
		dictData.GET("/page", dictHandler.PageData)
		dictData.GET("/type/:type", dictHandler.ListDataByType)
	}

	params := rg.Group("/system/param")
	{
		params.POST("", paramHandler.Create)
		params.PUT("", paramHandler.Update)
		params.DELETE("/:id", paramHandler.Delete)
		params.GET("/:id", paramHandler.Get)
		params.GET("/key", paramHandler.GetByKey)
		params.GET("/page", paramHandler.Page)
	}

	notices := rg.Group("/system/notice")
	{
		notices.POST("", noticeHandler.Create)
		notices.PUT("", noticeHandler.Update)
		notices.DELETE("/:id", noticeHandler.Delete)
		notices.GET("/:id", noticeHandler.Get)
		notices.GET("/page", noticeHandler.Page)
	}

	operLogs := rg.Group("/system/oper-log")
	{
		operLogs.GET("/page", logHandler.PageOperLog)
		operLogs.DELETE("/:id", logHandler.DeleteOperLog)
		operLogs.DELETE("/clean", logHandler.CleanOperLog)
	}

	loginLogs := rg.Group("/system/login-log")
	{
		loginLogs.GET("/page", logHandler.PageLoginLog)
		loginLogs.DELETE("/:id", logHandler.DeleteLoginLog)
		loginLogs.DELETE("/clean", logHandler.CleanLoginLog)
	}

	files := rg.Group("/system/file")
	{
		files.POST("", fileHandler.Upload)
		files.DELETE("/:id", fileHandler.Delete)
		files.GET("/:id", fileHandler.Get)
		files.GET("/page", fileHandler.Page)
		files.GET("/:id/download", fileHandler.Download)
	}
}
