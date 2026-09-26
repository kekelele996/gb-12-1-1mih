package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// MaterialService handles material checklist items.
type MaterialService struct {
	repo    MaterialItemStore
	appRepo ApplicationProjectStore
	logger  *slog.Logger
}

// NewMaterialService creates a MaterialService.
func NewMaterialService(repo MaterialItemStore, appRepo ApplicationProjectStore, logger *slog.Logger) *MaterialService {
	return &MaterialService{repo: repo, appRepo: appRepo, logger: logger}
}

// Create adds a checklist item. Only the responsible counselor or admin may
// maintain the checklist.
func (s *MaterialService) Create(applicationID, userID uint, role string, m *model.MaterialItem) (*model.MaterialItem, error) {
	a, err := s.loadApplication(applicationID)
	if err != nil {
		return nil, err
	}
	if role != constants.RoleCounselor && role != constants.RoleAdmin || !isManagingStaff(a, userID, role) {
		return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgMaterialManageStaffOnly)
	}
	m.ApplicationID = applicationID
	if m.Status == "" {
		m.Status = constants.MaterialPending
	}
	if err := s.repo.Create(m); err != nil {
		return nil, fmt.Errorf("material create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMaterialCreateSuccess, applicationID, m.Name), "id", m.ID, "by", userID)
	return m, nil
}

// ListByApplication returns items of an application, verifying access.
func (s *MaterialService) ListByApplication(applicationID, userID uint, role string) ([]model.MaterialItem, error) {
	a, err := s.loadApplication(applicationID)
	if err != nil {
		return nil, err
	}
	if err := assertCanViewApplication("ApplicationProject", applicationID, a, userID, role); err != nil {
		return nil, err
	}
	return s.repo.ListByApplication(applicationID)
}

// UpdateStatus runs an item through the checklist state machine with strict
// role separation:
//   - student owner uploads/re-uploads (-> uploaded), only while the project
//     is still editable (planning/preparing) and the item is not approved;
//   - responsible counselor/admin reviews uploads (-> approved/rejected).
func (s *MaterialService) UpdateStatus(userID, id uint, role, status, fileURL, remark string) (*model.MaterialItem, error) {
	if !constants.IsValidMaterialStatus(status) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("MaterialItem[id=%d] status=%s invalid", id, status))
	}
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositoryErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("MaterialItem[id=%d] not found", id))
		}
		return nil, fmt.Errorf("material status find: %w", err)
	}
	a, err := s.loadApplication(m.ApplicationID)
	if err != nil {
		return nil, err
	}
	if err := assertCanViewApplication("MaterialItem", id, a, userID, role); err != nil {
		return nil, err
	}

	transitionOK := false
	for _, nxt := range constants.NextMaterialStatuses(m.Status) {
		if nxt == status {
			transitionOK = true
			break
		}
	}
	if !transitionOK {
		return nil, util.NewAppError(409, constants.CodeConflict,
			fmt.Sprintf("MaterialItem[id=%d] status change failed: %s -> %s not allowed", id, m.Status, status))
	}

	editable := constants.IsEditableApplicationStatus(a.Status)
	switch role {
	case constants.RoleStudent:
		if a.StudentID != userID {
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgForbidden)
		}
		if status != constants.MaterialUploaded {
			// A student attempting to approve/reject is an over-privilege.
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgMaterialReviewStaffOnly)
		}
		if !editable {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgAppNotEditable)
		}
		if m.Status == constants.MaterialApproved {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgMaterialApprovedLocked)
		}
		if fileURL == "" && m.FileURL == "" {
			return nil, util.NewAppError(422, constants.CodeValidationError, constants.MsgMaterialRequiresUpload)
		}
	case constants.RoleCounselor, constants.RoleAdmin:
		if !isManagingStaff(a, userID, role) {
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgNotResponsibleCounselor)
		}
		switch status {
		case constants.MaterialApproved:
			// Only a student-uploaded item can be approved.
			if m.FileURL == "" {
				return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgMaterialRequiresUpload)
			}
			now := time.Now()
			m.Status = status
			m.ReviewedBy = userID
			m.ReviewedAt = &now
			m.ReviewRemark = remark
		case constants.MaterialRejected:
			now := time.Now()
			m.Status = status
			m.ReviewedBy = userID
			m.ReviewedAt = &now
			m.ReviewRemark = remark
		case constants.MaterialUploaded:
			// Staff requesting a replacement clears the previous review.
			m.Status = status
			m.ReviewedBy = 0
			m.ReviewedAt = nil
			m.ReviewRemark = ""
			if fileURL != "" {
				m.FileURL = fileURL
			}
		}
	default:
		return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgForbidden)
	}

	if role == constants.RoleStudent {
		m.Status = status
		if fileURL != "" {
			m.FileURL = fileURL
		}
		now := time.Now()
		m.UploadedAt = &now
	}

	if err := s.repo.Update(m); err != nil {
		return nil, fmt.Errorf("material status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogMaterialStatusChanged, id, status), "id", id, "by", userID, "role", role)
	return m, nil
}

// Progress computes material completion percentage for an application.
// A required item counts as complete only after it passes review (approved).
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
		if m.Status == constants.MaterialApproved {
			done++
		}
	}
	return done * 100 / len(items), nil
}

func (s *MaterialService) loadApplication(applicationID uint) (*model.ApplicationProject, error) {
	a, err := s.appRepo.FindByID(applicationID)
	if err != nil {
		if errors.Is(err, repositoryErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("ApplicationProject[id=%d] not found", applicationID))
		}
		return nil, fmt.Errorf("material load application: %w", err)
	}
	return a, nil
}
