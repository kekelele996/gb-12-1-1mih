package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerApplicationRoutes(v1 *gin.RouterGroup, cfg *config.Config, ah *handler.ApplicationHandler, dh *handler.DocumentHandler, mh *handler.MaterialHandler, th *handler.TimelineHandler, limiter *middleware.RateLimiter) {
	apps := v1.Group("/applications", middleware.AuthRequired(cfg))
	apps.GET("", ah.List)
	apps.GET("/:id", ah.Get)
	apps.POST("", middleware.RequireRole("student"), limiter.Limit(), ah.Create)
	apps.PUT("/:id/status", ah.UpdateStatus)
	apps.GET("/:id/documents", dh.ListByApplication)
	apps.POST("/:id/documents", limiter.Limit(), dh.Create)
	apps.GET("/:id/materials", mh.ListByApplication)
	apps.POST("/:id/materials", mh.Create)
	apps.GET("/:id/timeline", th.ListByApplication)
	apps.POST("/:id/timeline", limiter.Limit(), th.Create)
	v1.PUT("/materials/:id/status", middleware.AuthRequired(cfg), mh.UpdateStatus)
}
