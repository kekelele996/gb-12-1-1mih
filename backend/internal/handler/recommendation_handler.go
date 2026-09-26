package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// RecommendationHandler exposes school-selection plan endpoints.
type RecommendationHandler struct {
	svc    *service.RecommendationService
	logger *slog.Logger
}

// NewRecommendationHandler creates a RecommendationHandler.
func NewRecommendationHandler(svc *service.RecommendationService, logger *slog.Logger) *RecommendationHandler {
	return &RecommendationHandler{svc: svc, logger: logger}
}

// Create handles POST /recommendations (counselor).
func (h *RecommendationHandler) Create(c *gin.Context) {
	var req dto.RecommendationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, constants.MsgInvalidParam+": "+err.Error()))
		return
	}
	rec, err := h.svc.Create(middleware.GetUserID(c), req.StudentID, req.UniversityIDs, req.Reason)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(rec))
}

// ListByStudent handles GET /recommendations/student/:studentId.
func (h *RecommendationHandler) ListByStudent(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("studentId"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid student id"))
		return
	}
	items, err := h.svc.ListByStudent(uint(studentID))
	if err != nil {
		c.Error(err)
		return
	}
	result := make([]gin.H, 0, len(items))
	for _, rec := range items {
		unis, err := h.svc.ResolveUniversities(&rec)
		if err != nil {
			continue
		}
		result = append(result, gin.H{"id": rec.ID, "reason": rec.Reason, "created_at": rec.CreatedAt, "universities": unis})
	}
	c.JSON(http.StatusOK, dto.OK(result))
}
