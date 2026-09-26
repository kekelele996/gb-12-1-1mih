package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// UniversityService handles the university catalog.
type UniversityService struct {
	repo   *repository.UniversityRepository
	logger *slog.Logger
}

// NewUniversityService creates a UniversityService.
func NewUniversityService(repo *repository.UniversityRepository, logger *slog.Logger) *UniversityService {
	return &UniversityService{repo: repo, logger: logger}
}

// Create adds a university (admin).
func (s *UniversityService) Create(u *model.University) (*model.University, error) {
	if u.TopMajors == "" {
		u.TopMajors = "[]"
	}
	if u.Requirements == "" {
		u.Requirements = "{}"
	}
	if err := s.repo.Create(u); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("University[name=%s] create failed: name exists", u.Name))
		}
		s.logger.Error(fmt.Sprintf(constants.LogUniversityCreateFailed, u.Name), "error", err)
		return nil, fmt.Errorf("university create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUniversityCreateSuccess, u.Name), "id", u.ID)
	return u, nil
}

// Get returns a university by id.
func (s *UniversityService) Get(id uint) (*model.University, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("University[id=%d] not found", id))
		}
		return nil, fmt.Errorf("university get: %w", err)
	}
	return u, nil
}

// Update edits a university (admin).
func (s *UniversityService) Update(id uint, u *model.University) (*model.University, error) {
	exist, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("university update find: %w", err)
	}
	if u.Name != "" {
		exist.Name = u.Name
	}
	if u.Country != "" {
		exist.Country = u.Country
	}
	if u.City != "" {
		exist.City = u.City
	}
	if u.Ranking > 0 {
		exist.Ranking = u.Ranking
	}
	if u.TopMajors != "" {
		exist.TopMajors = u.TopMajors
	}
	if u.ApplicationDeadline != "" {
		exist.ApplicationDeadline = u.ApplicationDeadline
	}
	if u.TuitionRange != "" {
		exist.TuitionRange = u.TuitionRange
	}
	if u.Requirements != "" {
		exist.Requirements = u.Requirements
	}
	if err := s.repo.Update(exist); err != nil {
		return nil, fmt.Errorf("university update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUniversityUpdateSuccess, id), "id", id)
	return exist, nil
}

// Delete removes a university (admin).
func (s *UniversityService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("university delete: %w", err)
	}
	return nil
}

// List filters universities.
func (s *UniversityService) List(country string, rankMin, rankMax int, keyword string, page, pageSize int) ([]model.University, int64, error) {
	items, total, err := s.repo.List(country, rankMin, rankMax, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("university list: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUniversityListSuccess, country), "total", total)
	return items, total, nil
}
