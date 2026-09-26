package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerRecommendationRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.RecommendationHandler, limiter *middleware.RateLimiter) {
	recs := v1.Group("/recommendations")
	recs.GET("/student/:studentId", middleware.AuthRequired(cfg), h.ListByStudent)
	recs.POST("", middleware.AuthRequired(cfg), middleware.RequireRole("counselor"), limiter.Limit(), h.Create)
}
