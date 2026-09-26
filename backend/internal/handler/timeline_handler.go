package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// TimelineHandler exposes timeline node endpoints.
type TimelineHandler struct {
	svc    *service.TimelineService
	logger *slog.Logger
}

// NewTimelineHandler creates a TimelineHandler.
func NewTimelineHandler(svc *service.TimelineService, logger *slog.Logger) *TimelineHandler {
	return &TimelineHandler{svc: svc, logger: logger}
}

// ListByApplication handles GET /applications/:appId/timeline.
func (h *TimelineHandler) ListByApplication(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid application id"))
		return
	}
	items, err := h.svc.ListByApplication(uint(appID))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// Create handles POST /applications/:appId/timeline.
func (h *TimelineHandler) Create(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid application id"))
		return
	}
	var req dto.TimelineCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	n := &model.TimelineNode{Title: req.Title, NodeType: req.NodeType, DueDate: req.DueDate}
	created, err := h.svc.Create(uint(appID), n)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// MarkDone handles PUT /timeline/:id/done.
func (h *TimelineHandler) MarkDone(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid timeline id"))
		return
	}
	n, err := h.svc.MarkDone(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(n))
}
