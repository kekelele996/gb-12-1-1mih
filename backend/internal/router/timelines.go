package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerTimelineRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.TimelineHandler, limiter *middleware.RateLimiter) {
	tl := v1.Group("/timeline", middleware.AuthRequired(cfg))
	tl.PUT("/:id/done", h.MarkDone)
}
