package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerUniversityRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.UniversityHandler, limiter *middleware.RateLimiter) {
	unis := v1.Group("/universities")
	unis.GET("", h.List)
	unis.GET("/:id", h.Get)
	admin := unis.Group("", middleware.AuthRequired(cfg), middleware.RequireRole("admin"))
	admin.POST("", limiter.Limit(), h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}
