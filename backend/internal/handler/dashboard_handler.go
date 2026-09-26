package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

// DashboardHandler aggregates stats for dashboards.
type DashboardHandler struct {
	appService  *service.ApplicationService
	univService *service.UniversityService
	matService  *service.MaterialService
	logger      *slog.Logger
}

// NewDashboardHandler creates a DashboardHandler.
func NewDashboardHandler(appService *service.ApplicationService, univService *service.UniversityService, matService *service.MaterialService, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{appService: appService, univService: univService, matService: matService, logger: logger}
}

// Stats handles GET /dashboard/stats.
func (h *DashboardHandler) Stats(c *gin.Context) {
	projects, err := h.appService.List(middleware.GetUserID(c), middleware.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}
	stats := service.ComputeAppStats(projects)
	// material average across projects
	totalProgress := 0
	count := 0
	for _, p := range projects {
		progress, err := h.matService.Progress(p.ID)
		if err == nil {
			totalProgress += progress
			count++
		}
	}
	if count > 0 {
		stats.MaterialAvg = totalProgress / count
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"stats": stats, "projects": projects}))
}

// AllTimeline handles GET /timeline/all via app list (simplified aggregation).
func (h *DashboardHandler) AllTimeline(c *gin.Context) {
	projects, err := h.appService.List(middleware.GetUserID(c), middleware.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}
	var nodes []model.TimelineNode
	_ = nodes
	c.JSON(http.StatusOK, dto.OK(gin.H{"projects": projects}))
}
