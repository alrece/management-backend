package workflow

import (
	"management-backend/internal/middleware"
	"management-backend/internal/module/workflow/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册工作流模块路由
func RegisterRoutes(
	authed *gin.RouterGroup,
	public *gin.Engine,
	workflowHandler *handler.WorkflowHandler,
	executionHandler *handler.ExecutionHandler,
	webhookHandler *handler.WebhookHandler,
	categoryHandler *handler.CategoryHandler,
) {
	// 认证路由
	wf := authed.Group("/workflow/workflow")
	{
		wf.POST("", middleware.RequirePermission("workflow:workflow:create"), workflowHandler.Create)
		wf.PUT("", middleware.RequirePermission("workflow:workflow:update"), workflowHandler.Update)
		wf.DELETE("/:id", middleware.RequirePermission("workflow:workflow:delete"), workflowHandler.Delete)
		wf.GET("/:id", middleware.RequirePermission("workflow:workflow:query"), workflowHandler.Get)
		wf.GET("/page", middleware.RequirePermission("workflow:workflow:query"), workflowHandler.Page)
		wf.PUT("/:id/activate", middleware.RequirePermission("workflow:workflow:update"), workflowHandler.Activate)
		wf.PUT("/:id/deactivate", middleware.RequirePermission("workflow:workflow:update"), workflowHandler.Deactivate)
		wf.POST("/:id/execute", middleware.RequirePermission("workflow:workflow:execute"), workflowHandler.Execute)
		wf.GET("/nodes/schemas", middleware.RequirePermission("workflow:workflow:query"), workflowHandler.Schemas)
	}

	exec := authed.Group("/workflow/execution")
	{
		exec.GET("/page", middleware.RequirePermission("workflow:execution:query"), executionHandler.Page)
		exec.GET("/:id", middleware.RequirePermission("workflow:execution:query"), executionHandler.Get)
		exec.GET("/:id/logs", middleware.RequirePermission("workflow:execution:query"), executionHandler.Logs)
	}

	cat := authed.Group("/workflow/category")
	{
		cat.POST("", middleware.RequirePermission("workflow:category:create"), categoryHandler.Create)
		cat.PUT("", middleware.RequirePermission("workflow:category:update"), categoryHandler.Update)
		cat.DELETE("/:id", middleware.RequirePermission("workflow:category:delete"), categoryHandler.Delete)
		cat.GET("/tree", middleware.RequirePermission("workflow:category:query"), categoryHandler.Tree)
	}

	// Webhook 公开路由（无需认证，路径包含 tenantId 和 workflowId）
	public.POST("/api/wf/webhook/:tenantId/:workflowId/*path", webhookHandler.Handle)
}
