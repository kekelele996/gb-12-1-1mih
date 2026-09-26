package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerDocumentRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.DocumentHandler, limiter *middleware.RateLimiter) {
	docs := v1.Group("/documents", middleware.AuthRequired(cfg))
	docs.GET("/:id", h.Get)
	docs.PUT("/:id", limiter.Limit(), h.Save)
	docs.GET("/:id/versions", h.ListVersions)
	docs.POST("/:id/rollback", h.Rollback)
	docs.GET("/:id/annotations", h.ListAnnotations)
	docs.POST("/:id/annotations", middleware.RequireRole("counselor", "admin"), limiter.Limit(), h.AddAnnotation)
}
