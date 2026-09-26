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
	repo    *repository.MaterialItemRepository
	appRepo *repository.ApplicationProjectRepository
	logger  *slog.Logger
}

// NewMaterialService creates a MaterialService.
func NewMaterialService(repo *repository.MaterialItemRepository, appRepo *repository.ApplicationProjectRepository, logger *slog.Logger) *MaterialService {
	return &MaterialService{repo: repo, appRepo: appRepo, logger: logger}
}

// MissingRequiredMaterials returns names of required items not yet approved.
func MissingRequiredMaterials(items []model.MaterialItem) []string {
	var missing []string
	for _, m := range items {
		if m.IsRequired && m.Status != constants.MaterialApproved {
			missing = append(missing, m.Name)
		}
	}
	return missing
}

// Create adds a checklist item.
func (s *MaterialService) Create(userID uint, role string, applicationID uint, m *model.MaterialItem) (*model.MaterialItem, error) {
	a, err := s.appRepo.FindByID(applicationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("ApplicationProject[id=%d] not found", applicationID))
		}
		return nil, fmt.Errorf("material application find: %w", err)
	}
	if role == constants.RoleStudent {
		if a.StudentID != userID {
			return nil, util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("MaterialItem create failed: user_id=%d not student owner of ApplicationProject[id=%d]", userID, applicationID))
		}
		if !constants.IsApplicationEditable(a.Status) {
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("MaterialItem create failed: ApplicationProject[id=%d] already submitted, materials locked", applicationID))
		}
	}
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

// UpdateStatus updates an item status. Upload is student-only while the
// application is editable; approving requires the responsible counselor or admin.
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
	a, err := s.appRepo.FindByID(m.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("material application find: %w", err)
	}
	if status == constants.MaterialApproved {
		// 进入已审核：仅负责顾问或管理员
		isResponsible := role == constants.RoleCounselor && a.CounselorID == userID
		if !isResponsible && role != constants.RoleAdmin {
			return nil, util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("MaterialItem[id=%d] approve failed: require responsible counselor or admin", id))
		}
	} else {
		// 上传/重新上传：仅学生本人，且申请未提交锁定
		if role != constants.RoleStudent || a.StudentID != userID {
			return nil, util.NewAppError(403, constants.CodeForbidden,
				fmt.Sprintf("MaterialItem[id=%d] upload failed: only the owning student can upload", id))
		}
		if !constants.IsApplicationEditable(a.Status) {
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("MaterialItem[id=%d] update failed: application already submitted, materials locked", id))
		}
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
