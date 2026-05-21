package router

import (
	"management-backend/services/workflow-service/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册工作流服务路由
func RegisterRoutes(rg *gin.RouterGroup, wfHandler *handler.WorkflowHandler, webhookHandler *handler.WebhookHandler) {
	// 分类
	categories := rg.Group("/workflow/categories")
	{
		categories.POST("", wfHandler.CreateCategory)
		categories.PUT("/:id", wfHandler.UpdateCategory)
		categories.DELETE("/:id", wfHandler.DeleteCategory)
		categories.GET("/tree", wfHandler.CategoryTree)
	}

	// 工作流
	workflows := rg.Group("/workflow/workflows")
	{
		workflows.POST("", wfHandler.CreateWorkflow)
		workflows.PUT("/:id", wfHandler.UpdateWorkflow)
		workflows.DELETE("/:id", wfHandler.DeleteWorkflow)
		workflows.GET("/:id", wfHandler.GetWorkflow)
		workflows.GET("/page", wfHandler.PageWorkflow)
		workflows.PUT("/:id/activate", wfHandler.ActivateWorkflow)
		workflows.PUT("/:id/deactivate", wfHandler.DeactivateWorkflow)
		workflows.POST("/:id/execute", wfHandler.ExecuteWorkflow)
	}

	// Webhook 回调（无需认证）
	rg.POST("/workflow/callback", webhookHandler.Callback)
}
