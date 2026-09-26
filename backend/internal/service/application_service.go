package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// ApplicationProjectStore is the persistence port for application projects.
type ApplicationProjectStore interface {
	Create(a *model.ApplicationProject) error
	FindByID(id uint) (*model.ApplicationProject, error)
	Update(a *model.ApplicationProject) error
	ListByStudent(studentID uint) ([]model.ApplicationProject, error)
	ListByCounselor(counselorID uint) ([]model.ApplicationProject, error)
	ListAll() ([]model.ApplicationProject, error)
}

// UniversityStore is the persistence port for universities.
type UniversityStore interface {
	FindByID(id uint) (*model.University, error)
}

// MaterialItemStore is the persistence port for material items.
type MaterialItemStore interface {
	Create(m *model.MaterialItem) error
	FindByID(id uint) (*model.MaterialItem, error)
	Update(m *model.MaterialItem) error
	ListByApplication(applicationID uint) ([]model.MaterialItem, error)
}

// ApplicationService implements the application project state machine.
type ApplicationService struct {
	repo     ApplicationProjectStore
	univRepo UniversityStore
	matRepo  MaterialItemStore
	logger   *slog.Logger
}

// NewApplicationService creates an ApplicationService.
func NewApplicationService(repo ApplicationProjectStore, univRepo UniversityStore, matRepo MaterialItemStore, logger *slog.Logger) *ApplicationService {
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
		if errors.Is(err, repositoryErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("ApplicationProject[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application get: %w", err)
	}
	if err := assertCanViewApplication("ApplicationProject", id, a, userID, role); err != nil {
		return nil, err
	}
	return a, nil
}

// MissingRequiredMaterials returns required items that have not passed review.
func (s *ApplicationService) MissingRequiredMaterials(id uint) ([]model.MaterialItem, error) {
	items, err := s.matRepo.ListByApplication(id)
	if err != nil {
		return nil, fmt.Errorf("application required materials: %w", err)
	}
	var missing []model.MaterialItem
	for _, m := range items {
		if m.IsRequired && m.Status != constants.MaterialApproved {
			missing = append(missing, m)
		}
	}
	return missing, nil
}

// UpdateStatus transitions a project along the state machine.
//
// Role rules:
//   - student owner: planning -> preparing, preparing -> submitted (the
//     submit is blocked while any required material is not approved)
//   - responsible counselor / admin: submitted -> waiting (and onward to
//     admitted/rejected/waitlisted), and submitted -> preparing (return
//     for revision, which unlocks material/document editing for the student)
func (s *ApplicationService) UpdateStatus(id, userID uint, role, next string) (*model.ApplicationProject, error) {
	if !constants.IsValidApplicationStatus(next) {
		return nil, util.NewAppError(422, constants.CodeValidationError,
			fmt.Sprintf("ApplicationProject[id=%d] status=%s invalid", id, next))
	}
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repositoryErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("ApplicationProject[id=%d] not found", id))
		}
		return nil, fmt.Errorf("application status find: %w", err)
	}
	if err := assertCanViewApplication("ApplicationProject", id, a, userID, role); err != nil {
		return nil, err
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

	// Who is allowed to perform this transition.
	staffOnly := next != constants.AppStatusPreparing && next != constants.AppStatusSubmitted
	returnToPreparing := next == constants.AppStatusPreparing && a.Status == constants.AppStatusSubmitted
	isSubmit := next == constants.AppStatusSubmitted

	switch role {
	case constants.RoleStudent:
		if a.StudentID != userID {
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgAppSubmitOwnerOnly)
		}
		if staffOnly || returnToPreparing {
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgAppReturnStaffOnly)
		}
	case constants.RoleCounselor, constants.RoleAdmin:
		if isSubmit {
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgAppSubmitOwnerOnly)
		}
		if !isManagingStaff(a, userID, role) {
			return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgNotResponsibleCounselor)
		}
	default:
		return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgForbidden)
	}

	// Submit gate: every required material must be student-uploaded and approved.
	if isSubmit {
		missing, err := s.MissingRequiredMaterials(id)
		if err != nil {
			return nil, err
		}
		if len(missing) > 0 {
			names := make([]string, 0, len(missing))
			for _, m := range missing {
				names = append(names, m.Name)
			}
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("%s：%s", constants.MsgSubmitMissingMaterials, strings.Join(names, "、"))).
				WithDetails(map[string]interface{}{
					"reason":            "required_materials_not_approved",
					"missing_materials": names,
				})
		}
	}

	a.Status = next
	if err := s.repo.Update(a); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogAppStatusChangeFailed, id), "error", err)
		return nil, fmt.Errorf("application status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogAppStatusChanged, id, next), "id", id, "by", userID, "role", role)
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

// MaterialBlock describes one required material blocking submission.
type MaterialBlock struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// BlockedAction explains why an action is unavailable for the current viewer.
type BlockedAction struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

// ApplicationAbilities is the role-aware capability/block view used by the
// application detail page to decide which actions to show and why others are
// blocked.
type ApplicationAbilities struct {
	IsStudentOwner     bool            `json:"is_student_owner"`
	IsManagingStaff    bool            `json:"is_managing_staff"`
	Editable           bool            `json:"editable"`
	CanSubmit          bool            `json:"can_submit"`
	CanReturn          bool            `json:"can_return"`
	CanAdvance         bool            `json:"can_advance"`
	CanReviewMaterials bool            `json:"can_review_materials"`
	CanManageMaterials bool            `json:"can_manage_materials"`
	CanEditDocuments   bool            `json:"can_edit_documents"`
	MissingRequired    []MaterialBlock `json:"missing_required_materials"`
	Blocked            []BlockedAction `json:"blocked"`
}

// Abilities computes the role-aware capability view for a project.
func (s *ApplicationService) Abilities(a *model.ApplicationProject, userID uint, role string) (*ApplicationAbilities, error) {
	isOwner := role == constants.RoleStudent && a.StudentID == userID
	isStaff := isManagingStaff(a, userID, role)
	editable := constants.IsEditableApplicationStatus(a.Status)

	missing, err := s.MissingRequiredMaterials(a.ID)
	if err != nil {
		return nil, err
	}
	blocks := make([]MaterialBlock, 0, len(missing))
	for _, m := range missing {
		blocks = append(blocks, MaterialBlock{ID: m.ID, Name: m.Name, Status: m.Status})
	}

	ab := &ApplicationAbilities{
		IsStudentOwner:     isOwner,
		IsManagingStaff:    isStaff,
		Editable:           editable && isOwner,
		CanSubmit:          false,
		CanReturn:          false,
		CanAdvance:         false,
		CanReviewMaterials: isStaff,
		CanManageMaterials: isStaff,
		CanEditDocuments:   editable && isOwner,
		MissingRequired:    blocks,
		Blocked:            []BlockedAction{},
	}

	// Submit: only the student owner while preparing, gated on materials.
	if isOwner && a.Status == constants.AppStatusPreparing {
		if len(blocks) == 0 {
			ab.CanSubmit = true
		} else {
			ab.Blocked = append(ab.Blocked, BlockedAction{
				Action: "submit",
				Reason: fmt.Sprintf("%s：%s", constants.MsgSubmitMissingMaterials, materialNames(blocks)),
			})
		}
	} else if isOwner {
		ab.Blocked = append(ab.Blocked, BlockedAction{Action: "submit", Reason: "当前状态不可提交申请"})
	}

	// Return for revision: responsible staff, only after submission.
	if isStaff && a.Status == constants.AppStatusSubmitted {
		ab.CanReturn = true
	} else if role == constants.RoleCounselor && !isStaff {
		ab.Blocked = append(ab.Blocked, BlockedAction{Action: "return", Reason: constants.MsgNotResponsibleCounselor})
	}

	// Forward staff transitions from the current status.
	if isStaff {
		for _, nxt := range constants.NextApplicationStatuses(a.Status) {
			if nxt != constants.AppStatusPreparing && nxt != constants.AppStatusSubmitted {
				ab.CanAdvance = true
				break
			}
		}
	}

	// Editing lock explanation for the student after submission.
	if isOwner && !editable {
		ab.Blocked = append(ab.Blocked,
			BlockedAction{Action: "edit_documents", Reason: constants.MsgAppNotEditable},
			BlockedAction{Action: "upload_materials", Reason: constants.MsgAppNotEditable},
		)
	}

	return ab, nil
}

func materialNames(blocks []MaterialBlock) string {
	names := make([]string, 0, len(blocks))
	for _, b := range blocks {
		names = append(names, b.Name)
	}
	return strings.Join(names, "、")
}
