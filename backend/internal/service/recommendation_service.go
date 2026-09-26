package service

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
)

// RecommendationService handles counselor school-selection plans.
type RecommendationService struct {
	repo     *repository.RecommendationRepository
	univRepo *repository.UniversityRepository
	logger   *slog.Logger
}

// NewRecommendationService creates a RecommendationService.
func NewRecommendationService(repo *repository.RecommendationRepository, univRepo *repository.UniversityRepository, logger *slog.Logger) *RecommendationService {
	return &RecommendationService{repo: repo, univRepo: univRepo, logger: logger}
}

// Create makes a recommendation for a student.
func (s *RecommendationService) Create(counselorID, studentID uint, universityIDs []uint, reason string) (*model.Recommendation, error) {
	ids := make([]string, 0, len(universityIDs))
	for _, id := range universityIDs {
		ids = append(ids, strconv.FormatUint(uint64(id), 10))
	}
	rec := &model.Recommendation{
		StudentID: studentID, CounselorID: counselorID,
		UniversityIDs: "[" + strings.Join(ids, ",") + "]", Reason: reason,
	}
	if err := s.repo.Create(rec); err != nil {
		return nil, fmt.Errorf("recommendation create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogRecommendationCreate, studentID), "id", rec.ID)
	return rec, nil
}

// ListByStudent returns recommendations with resolved universities.
func (s *RecommendationService) ListByStudent(studentID uint) ([]model.Recommendation, error) {
	return s.repo.ListByStudent(studentID)
}

// ResolveUniversities parses university ids and loads the schools.
func (s *RecommendationService) ResolveUniversities(rec *model.Recommendation) ([]model.University, error) {
	var ids []uint
	raw := strings.Trim(rec.UniversityIDs, "[]")
	if raw == "" {
		return nil, nil
	}
	for _, part := range strings.Split(raw, ",") {
		if n, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64); err == nil {
			ids = append(ids, uint(n))
		}
	}
	return s.univRepo.ListByIDs(ids)
}
