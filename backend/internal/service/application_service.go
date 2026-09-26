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

// ApplicationService implements the application project state machine.
type ApplicationService struct {
	repo    *repository.ApplicationProjectRepository
	univRepo *repository.UniversityRepository
	logger  *slog.Logger
}

// NewApplicationService creates an ApplicationService.
func NewApplicationService(repo *repository.ApplicationProjectRepository, univRepo *repository.UniversityRepository, logger *slog.Logger) *ApplicationService {
	return &ApplicationService{repo: repo, univRepo: univRepo, logger: logger}
}

// Create creates an application project for a student.
func (s *ApplicationService) Create(studentID uint, a *model.ApplicationProject) (*model.ApplicationProject, error) {
	if _, err := s.univRepo.FindByID(a.UniversityID); err != nil {
		return nil, util.NewAppError(404, constants.CodeNotFound,
			fmt.Sprintf("University[id=%d] not found", a.UniversityID))
	}
	a.StudentID = studentID
	if a.Status == "" {
		a.Status = constants.AppStatusPlanning
	}
	if err := s.repo.Create(a); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppCreateFailed, a.Major), "error", err)
		return nil, fmt.Errorf("application create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppCreateSuccess, a.ID, a.Major), "student_id", studentID)
	return a, nil
}

// Get returns a project, verifying access.
func (s *ApplicationService) Get(id, userID uint, role string) (*model.ApplicationProject, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("ApplicationProject[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application get: %w", err)
	}
	if role == constants.RoleStudent && a.StudentID != userID {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("ApplicationProject[id=%d] get failed: user_id=%d not student owner", id, userID))
	}
	return a, nil
}

// UpdateStatus transitions a project along the state machine.
func (s *ApplicationService) UpdateStatus(id, userID uint, role, next string) (*model.ApplicationProject, error) {
	if !constants.IsValidApplicationStatus(next) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("ApplicationProject[id=%d] status=%s invalid", id, next))
	}
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("application status find: %w", err)
	}
	if role == constants.RoleStudent && a.StudentID != userID {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("ApplicationProject[id=%d] status change failed: not owner", id))
	}
	allowed := false
	for _, s2 := range constants.NextApplicationStatuses(a.Status) {
		if s2 == next {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("ApplicationProject[id=%d] status change failed: %s -> %s not allowed", id, a.Status, next))
	}
	a.Status = next
	if err := s.repo.Update(a); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppStatusChangeFailed, id), "error", err)
		return nil, fmt.Errorf("application status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppStatusChanged, id, next), "id", id)
	return a, nil
}

// List returns projects visible to the caller.
func (s *ApplicationService) List(userID uint, role string) ([]model.ApplicationProject, error) {
	switch role {
	case constants.RoleStudent:
		return s.repo.ListByStudent(userID)
	case constants.RoleCounselor:
		return s.repo.ListByCounselor(userID)
	default:
		return s.repo.ListAll()
	}
}
