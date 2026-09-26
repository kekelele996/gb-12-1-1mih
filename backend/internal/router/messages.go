package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerMessageRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.MessageHandler, limiter *middleware.RateLimiter) {
	msgs := v1.Group("/messages", middleware.AuthRequired(cfg))
	msgs.GET("", h.List)
	msgs.POST("", limiter.Limit(), h.Send)
	msgs.PUT("/:id/read", h.MarkRead)
}
