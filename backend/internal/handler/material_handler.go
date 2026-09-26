package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// MaterialHandler exposes material checklist endpoints.
type MaterialHandler struct {
	svc    *service.MaterialService
	logger *slog.Logger
}

// NewMaterialHandler creates a MaterialHandler.
func NewMaterialHandler(svc *service.MaterialService, logger *slog.Logger) *MaterialHandler {
	return &MaterialHandler{svc: svc, logger: logger}
}

// ListByApplication handles GET /applications/:appId/materials.
func (h *MaterialHandler) ListByApplication(c *gin.Context) {
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
	progress, _ := h.svc.Progress(uint(appID))
	c.JSON(http.StatusOK, dto.OK(gin.H{"items": items, "progress": progress}))
}

// Create handles POST /applications/:appId/materials.
func (h *MaterialHandler) Create(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid application id"))
		return
	}
	var req dto.MaterialCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	m := &model.MaterialItem{Name: req.Name, Category: req.Category, IsRequired: req.IsRequired}
	created, err := h.svc.Create(uint(appID), m)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// UpdateStatus handles PUT /materials/:id/status.
func (h *MaterialHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid material id"))
		return
	}
	var req dto.MaterialStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	m, err := h.svc.UpdateStatus(middleware.GetUserID(c), uint(id), middleware.GetUserRole(c), req.Status, req.FileURL)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(m))
}
