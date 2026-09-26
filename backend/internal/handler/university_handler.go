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

// UniversityHandler exposes university endpoints.
type UniversityHandler struct {
	svc    *service.UniversityService
	logger *slog.Logger
}

// NewUniversityHandler creates a UniversityHandler.
func NewUniversityHandler(svc *service.UniversityService, logger *slog.Logger) *UniversityHandler {
	return &UniversityHandler{svc: svc, logger: logger}
}

// List handles GET /universities.
func (h *UniversityHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	country := c.Query("country")
	keyword := c.Query("keyword")
	rankMin, _ := strconv.Atoi(c.DefaultQuery("rank_min", "0"))
	rankMax, _ := strconv.Atoi(c.DefaultQuery("rank_max", "0"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := h.svc.List(country, rankMin, rankMax, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: items, Total: total, Page: page, Size: pageSize}))
}

// Get handles GET /universities/:id.
func (h *UniversityHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid university id"))
		return
	}
	u, err := h.svc.Get(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(u))
}

// Create handles POST /universities (admin).
func (h *UniversityHandler) Create(c *gin.Context) {
	var req dto.UniversityCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	u := &model.University{
		Name: req.Name, Country: req.Country, City: req.City, Ranking: req.Ranking,
		TopMajors: req.TopMajors, ApplicationDeadline: req.ApplicationDeadline,
		TuitionRange: req.TuitionRange, Requirements: req.Requirements,
	}
	created, err := h.svc.Create(u)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(created))
}

// Update handles PUT /universities/:id (admin).
func (h *UniversityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid university id"))
		return
	}
	var req dto.UniversityCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam))
		return
	}
	u := &model.University{
		Name: req.Name, Country: req.Country, City: req.City, Ranking: req.Ranking,
		TopMajors: req.TopMajors, ApplicationDeadline: req.ApplicationDeadline,
		TuitionRange: req.TuitionRange, Requirements: req.Requirements,
	}
	updated, err := h.svc.Update(uint(id), u)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(updated))
}

// Delete handles DELETE /universities/:id (admin).
func (h *UniversityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid university id"))
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"deleted": true}))
}
