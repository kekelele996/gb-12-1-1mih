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

// ApplicationHandler exposes application project endpoints.
type ApplicationHandler struct {
	svc    *service.ApplicationService
	logger *slog.Logger
}

// NewApplicationHandler creates an ApplicationHandler.
func NewApplicationHandler(svc *service.ApplicationService, logger *slog.Logger) *ApplicationHandler {
	return &ApplicationHandler{svc: svc, logger: logger}
}

// List handles GET /applications.
func (h *ApplicationHandler) List(c *gin.Context) {
	items, err := h.svc.List(middleware.GetUserID(c), middleware.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// Get handles GET /applications/:id.
func (h *ApplicationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid application id"))
		return
	}
	a, err := h.svc.Get(uint(id), middleware.GetUserID(c), middleware.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(a))
}

// Create handles POST /applications (student).
func (h *ApplicationHandler) Create(c *gin.Context) {
	var req dto.ApplicationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	a := &model.ApplicationProject{UniversityID: req.UniversityID, Major: req.Major, Round: req.Round}
	created, err := h.svc.Create(middleware.GetUserID(c), a)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// UpdateStatus handles PUT /applications/:id/status.
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid application id"))
		return
	}
	var req dto.ApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	a, err := h.svc.UpdateStatus(uint(id), middleware.GetUserID(c), middleware.GetUserRole(c), req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(a))
}
