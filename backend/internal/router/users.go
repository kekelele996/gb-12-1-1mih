package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
)

func registerUserRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.UserHandler, limiter *middleware.RateLimiter) {
	users := v1.Group("/users")
	users.POST("/register", limiter.Limit(), h.Register)
	users.POST("/login", limiter.Limit(), h.Login)
	me := users.Group("/me", middleware.AuthRequired(cfg))
	me.GET("", h.GetProfile)
	me.PUT("", h.UpdateProfile)
	users.GET("/students", middleware.AuthRequired(cfg), middleware.RequireRole("counselor", "admin"), h.ListStudents)
}
