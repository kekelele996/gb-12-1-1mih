package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// ApplicationService implements the application project state machine.
type ApplicationService struct {
	repo     *repository.ApplicationProjectRepository
	univRepo *repository.UniversityRepository
	matRepo  *repository.MaterialItemRepository
	logger   *slog.Logger
}

// NewApplicationService creates an ApplicationService.
func NewApplicationService(repo *repository.ApplicationProjectRepository, univRepo *repository.UniversityRepository, matRepo *repository.MaterialItemRepository, logger *slog.Logger) *ApplicationService {
	return &ApplicationService{repo: repo, univRepo: univRepo, matRepo: matRepo, logger: logger}
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
// Forward transitions follow NextApplicationStatuses; submitted->preparing is the
// counselor return transition. Submitting requires every required material approved.
func (s *ApplicationService) UpdateStatus(id, userID uint, role, next string) (*model.ApplicationProject, error) {
	if !constants.IsValidApplicationStatus(next) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("ApplicationProject[id=%d] status=%s invalid", id, next))
	}
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("ApplicationProject[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application status find: %w", err)
	}
	isOwner := role == constants.RoleStudent && a.StudentID == userID
	isResponsible := role == constants.RoleCounselor && a.CounselorID == userID
	isAdmin := role == constants.RoleAdmin
	if role == constants.RoleStudent && !isOwner {
		return nil, util.NewAppError(403, constants.CodeForbidden,
			fmt.Sprintf("ApplicationProject[id=%d] status change failed: not owner", id))
	}
	if constants.IsReturnTransition(a.Status, next) {
		// 退回准备中：仅负责顾问或管理员
		if !isResponsible && !isAdmin {
			return nil, util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("ApplicationProject[id=%d] return to preparing failed: require responsible counselor or admin", id))
		}
	} else {
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
		switch next {
		case constants.AppStatusPreparing, constants.AppStatusSubmitted:
			// 学生本人、负责顾问或管理员
			if !isOwner && !isResponsible && !isAdmin {
				return nil, util.NewAppError(403, constants.CodeForbidden,
					fmt.Sprintf("ApplicationProject[id=%d] status change failed: require owner, responsible counselor or admin", id))
			}
		default:
			// waiting/admitted/rejected/waitlisted：学生无权操作
			if role == constants.RoleStudent {
				return nil, util.NewAppError(403, constants.CodeForbidden,
					fmt.Sprintf("ApplicationProject[id=%d] status change failed: student cannot set status=%s", id, next))
			}
		}
		if next == constants.AppStatusSubmitted {
			if err := s.checkRequiredMaterials(id); err != nil {
				return nil, err
			}
		}
	}
	a.Status = next
	if err := s.repo.Update(a); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppStatusChangeFailed, id), "error", err)
		return nil, fmt.Errorf("application status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppStatusChanged, id, next), "id", id)
	return a, nil
}

// checkRequiredMaterials blocks submission while required materials are not approved.
func (s *ApplicationService) checkRequiredMaterials(applicationID uint) error {
	items, err := s.matRepo.ListByApplication(applicationID)
	if err != nil {
		return fmt.Errorf("application submit material check: %w", err)
	}
	if missing := MissingRequiredMaterials(items); len(missing) > 0 {
		return util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("ApplicationProject[id=%d] submit blocked: required materials not approved: %s",
				applicationID, strings.Join(missing, ", ")))
	}
	return nil
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
