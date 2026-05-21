package router

import (
	"management-backend/services/job-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *handler.JobHandler) {
	tasks := r.Group("/job/tasks")
	{
		tasks.POST("", h.Create)
		tasks.GET("/page", h.Page)
		tasks.PUT("/:id", h.Update)
		tasks.DELETE("/:id", h.Delete)
		tasks.POST("/:id/trigger", h.Trigger)
	}
	r.GET("/job/execution-logs/page", h.ExecLogPage)
}
