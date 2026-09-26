package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// MaterialService handles material checklist items.
type MaterialService struct {
	repo   *repository.MaterialItemRepository
	logger *slog.Logger
}

// NewMaterialService creates a MaterialService.
func NewMaterialService(repo *repository.MaterialItemRepository, logger *slog.Logger) *MaterialService {
	return &MaterialService{repo: repo, logger: logger}
}

// Create adds a checklist item.
func (s *MaterialService) Create(applicationID uint, m *model.MaterialItem) (*model.MaterialItem, error) {
	m.ApplicationID = applicationID
	if m.Status == "" {
		m.Status = constants.MaterialPending
	}
	if err := s.repo.Create(m); err != nil {
		return nil, fmt.Errorf("material create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMaterialCreateSuccess, applicationID, m.Name), "id", m.ID)
	return m, nil
}

// ListByApplication returns items of an application.
func (s *MaterialService) ListByApplication(applicationID uint) ([]model.MaterialItem, error) {
	return s.repo.ListByApplication(applicationID)
}

// UpdateStatus updates an item status (upload or approve).
func (s *MaterialService) UpdateStatus(userID, id uint, role, status, fileURL string) (*model.MaterialItem, error) {
	if !constants.IsValidMaterialStatus(status) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("MaterialItem[id=%d] status=%s invalid", id, status))
	}
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("MaterialItem[id=%d] not found", id))
		}
		return nil, fmt.Errorf("material status find: %w", err)
	}
	m.Status = status
	if fileURL != "" {
		m.FileURL = fileURL
	}
	if status == constants.MaterialUploaded {
		now := time.Now()
		m.UploadedAt = &now
	}
	if err := s.repo.Update(m); err != nil {
		return nil, fmt.Errorf("material status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMaterialStatusChanged, id, status), "id", id)
	return m, nil
}

// Progress computes material completion percentage for an application.
func (s *MaterialService) Progress(applicationID uint) (int, error) {
	items, err := s.repo.ListByApplication(applicationID)
	if err != nil {
		return 0, fmt.Errorf("material progress: %w", err)
	}
	if len(items) == 0 {
		return 0, nil
	}
	done := 0
	for _, m := range items {
		if m.Status == constants.MaterialUploaded || m.Status == constants.MaterialApproved {
			done++
		}
	}
	return done * 100 / len(items), nil
}
