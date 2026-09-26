package service

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// ---- in-memory fakes -------------------------------------------------------

type fakeAppRepo struct {
	apps map[uint]*model.ApplicationProject
}

func (r *fakeAppRepo) Create(a *model.ApplicationProject) error {
	if a.ID == 0 {
		a.ID = uint(len(r.apps) + 1)
	}
	r.apps[a.ID] = a
	return nil
}
func (r *fakeAppRepo) FindByID(id uint) (*model.ApplicationProject, error) {
	a, ok := r.apps[id]
	if !ok {
		return nil, repositoryErrNotFound
	}
	cp := *a
	return &cp, nil
}
func (r *fakeAppRepo) Update(a *model.ApplicationProject) error {
	r.apps[a.ID] = a
	return nil
}
func (r *fakeAppRepo) ListByStudent(studentID uint) ([]model.ApplicationProject, error) {
	return nil, nil
}
func (r *fakeAppRepo) ListByCounselor(counselorID uint) ([]model.ApplicationProject, error) {
	return nil, nil
}
func (r *fakeAppRepo) ListAll() ([]model.ApplicationProject, error) { return nil, nil }

type fakeUnivRepo struct{ exists map[uint]bool }

func (r *fakeUnivRepo) FindByID(id uint) (*model.University, error) {
	if r.exists[id] {
		return &model.University{ID: id}, nil
	}
	return nil, repositoryErrNotFound
}

type fakeMatRepo struct {
	items map[uint]*model.MaterialItem
}

func (r *fakeMatRepo) Create(m *model.MaterialItem) error {
	if m.ID == 0 {
		m.ID = uint(len(r.items) + 1)
	}
	r.items[m.ID] = m
	return nil
}
func (r *fakeMatRepo) FindByID(id uint) (*model.MaterialItem, error) {
	m, ok := r.items[id]
	if !ok {
		return nil, repositoryErrNotFound
	}
	cp := *m
	return &cp, nil
}
func (r *fakeMatRepo) Update(m *model.MaterialItem) error {
	r.items[m.ID] = m
	return nil
}
func (r *fakeMatRepo) ListByApplication(applicationID uint) ([]model.MaterialItem, error) {
	out := make([]model.MaterialItem, 0)
	for _, m := range r.items {
		if m.ApplicationID == applicationID {
			out = append(out, *m)
		}
	}
	return out, nil
}

func newServices() (*ApplicationService, *MaterialService, *model.ApplicationProject) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	appRepo := &fakeAppRepo{apps: map[uint]*model.ApplicationProject{}}
	univRepo := &fakeUnivRepo{exists: map[uint]bool{1: true}}
	matRepo := &fakeMatRepo{items: map[uint]*model.MaterialItem{}}
	appSvc := NewApplicationService(appRepo, univRepo, matRepo, logger)
	matSvc := NewMaterialService(matRepo, appRepo, logger)

	a := &model.ApplicationProject{StudentID: 10, CounselorID: 20, UniversityID: 1, Major: "CS", Status: constants.AppStatusPlanning}
	_ = appRepo.Create(a)
	return appSvc, matSvc, a
}

func codeOf(err error) int {
	var ae *util.AppError
	if errors.As(err, &ae) {
		return ae.Code
	}
	return -1
}
func isForbidden(err error) bool { return codeOf(err) == constants.CodeForbidden }
func isConflict(err error) bool  { return codeOf(err) == constants.CodeConflict }

// ---- scenarios -------------------------------------------------------------

func TestSubmitBlockedUntilRequiredMaterialsApproved(t *testing.T) {
	appSvc, matSvc, a := newServices()
	const studentID, counselorID = uint(10), uint(20)

	if _, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusPreparing); err != nil {
		t.Fatalf("student start preparing: %v", err)
	}

	m1 := &model.MaterialItem{ApplicationID: a.ID, Name: "成绩单", IsRequired: true, Status: constants.MaterialPending}
	m2 := &model.MaterialItem{ApplicationID: a.ID, Name: "语言成绩", IsRequired: true, Status: constants.MaterialPending}
	opt := &model.MaterialItem{ApplicationID: a.ID, Name: "作品集", IsRequired: false, Status: constants.MaterialPending}
	for _, m := range []*model.MaterialItem{m1, m2, opt} {
		if _, err := matSvc.Create(a.ID, counselorID, constants.RoleCounselor, m); err != nil {
			t.Fatalf("counselor create material: %v", err)
		}
	}

	_, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusSubmitted)
	if !isConflict(err) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if !strings.Contains(err.Error(), "成绩单") || !strings.Contains(err.Error(), "语言成绩") {
		t.Fatalf("error should list missing materials: %v", err)
	}

	if _, err := matSvc.UpdateStatus(studentID, m1.ID, constants.RoleStudent, constants.MaterialUploaded, "/files/a.pdf", ""); err != nil {
		t.Fatalf("student upload: %v", err)
	}
	if _, err := matSvc.UpdateStatus(studentID, m1.ID, constants.RoleStudent, constants.MaterialApproved, "", ""); !isForbidden(err) {
		t.Fatalf("student self-approve must be forbidden, got %v", err)
	}
	if _, err := matSvc.UpdateStatus(99, m1.ID, constants.RoleCounselor, constants.MaterialApproved, "", ""); !isForbidden(err) {
		t.Fatalf("other counselor approve must be forbidden, got %v", err)
	}
	if _, err := matSvc.UpdateStatus(counselorID, m1.ID, constants.RoleCounselor, constants.MaterialApproved, "", ""); err != nil {
		t.Fatalf("counselor approve: %v", err)
	}

	if _, err := matSvc.UpdateStatus(studentID, m2.ID, constants.RoleStudent, constants.MaterialUploaded, "/files/b.pdf", ""); err != nil {
		t.Fatalf("student upload m2: %v", err)
	}
	if _, err := matSvc.UpdateStatus(counselorID, m2.ID, constants.RoleCounselor, constants.MaterialRejected, "", "不清晰"); err != nil {
		t.Fatalf("counselor reject: %v", err)
	}
	if _, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusSubmitted); !isConflict(err) {
		t.Fatalf("submit with rejected required must be blocked, got %v", err)
	}

	if _, err := matSvc.UpdateStatus(studentID, m2.ID, constants.RoleStudent, constants.MaterialUploaded, "/files/b2.pdf", ""); err != nil {
		t.Fatalf("student re-upload after reject: %v", err)
	}
	if _, err := matSvc.UpdateStatus(1, m2.ID, constants.RoleAdmin, constants.MaterialApproved, "", ""); err != nil {
		t.Fatalf("admin approve: %v", err)
	}

	updated, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusSubmitted)
	if err != nil {
		t.Fatalf("submit after all required approved: %v", err)
	}
	if updated.Status != constants.AppStatusSubmitted {
		t.Fatalf("status = %s", updated.Status)
	}

	if p, _ := matSvc.Progress(a.ID); p != 66 {
		t.Fatalf("progress = %d, want 66", p)
	}
}

func TestPostSubmitLockAndCounselorReturn(t *testing.T) {
	appSvc, matSvc, a := newServices()
	const studentID, counselorID = uint(10), uint(20)

	m := &model.MaterialItem{ApplicationID: a.ID, Name: "成绩单", IsRequired: true}
	if _, err := matSvc.Create(a.ID, counselorID, constants.RoleCounselor, m); err != nil {
		t.Fatal(err)
	}
	if _, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusPreparing); err != nil {
		t.Fatal(err)
	}
	if _, err := matSvc.UpdateStatus(studentID, m.ID, constants.RoleStudent, constants.MaterialUploaded, "/f.pdf", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := matSvc.UpdateStatus(counselorID, m.ID, constants.RoleCounselor, constants.MaterialApproved, "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusSubmitted); err != nil {
		t.Fatal(err)
	}

	if _, err := matSvc.UpdateStatus(studentID, m.ID, constants.RoleStudent, constants.MaterialUploaded, "/f2.pdf", ""); !isConflict(err) {
		t.Fatalf("post-submit upload must be blocked, got %v", err)
	}
	if _, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusPreparing); !isForbidden(err) {
		t.Fatalf("student return must be forbidden, got %v", err)
	}
	if _, err := appSvc.UpdateStatus(a.ID, 99, constants.RoleCounselor, constants.AppStatusPreparing); !isForbidden(err) {
		t.Fatalf("other counselor return must be forbidden, got %v", err)
	}
	ret, err := appSvc.UpdateStatus(a.ID, counselorID, constants.RoleCounselor, constants.AppStatusPreparing)
	if err != nil {
		t.Fatalf("counselor return: %v", err)
	}
	if ret.Status != constants.AppStatusPreparing {
		t.Fatalf("status = %s", ret.Status)
	}
	// approved material stays student-locked even while preparing
	if _, err := matSvc.UpdateStatus(studentID, m.ID, constants.RoleStudent, constants.MaterialUploaded, "/f3.pdf", ""); !isConflict(err) {
		t.Fatalf("student replacing approved material must be blocked, got %v", err)
	}
	// staff can request replacement, clearing the review
	reset, err := matSvc.UpdateStatus(counselorID, m.ID, constants.RoleCounselor, constants.MaterialUploaded, "", "")
	if err != nil {
		t.Fatalf("staff reset approved to uploaded: %v", err)
	}
	if reset.ReviewedBy != 0 || reset.Status != constants.MaterialUploaded {
		t.Fatalf("review not cleared: %+v", reset)
	}
}

func TestStudentCannotManageChecklist(t *testing.T) {
	_, matSvc, a := newServices()
	m := &model.MaterialItem{ApplicationID: a.ID, Name: "x", IsRequired: true}
	if _, err := matSvc.Create(a.ID, 10, constants.RoleStudent, m); !isForbidden(err) {
		t.Fatalf("student material create must be forbidden, got %v", err)
	}
}

func TestAbilitiesReflectRoleAndBlocks(t *testing.T) {
	appSvc, matSvc, a := newServices()
	const studentID, counselorID = uint(10), uint(20)
	m := &model.MaterialItem{ApplicationID: a.ID, Name: "成绩单", IsRequired: true}
	if _, err := matSvc.Create(a.ID, counselorID, constants.RoleCounselor, m); err != nil {
		t.Fatal(err)
	}
	if _, err := appSvc.UpdateStatus(a.ID, studentID, constants.RoleStudent, constants.AppStatusPreparing); err != nil {
		t.Fatal(err)
	}

	stAb, err := appSvc.Abilities(a, studentID, constants.RoleStudent)
	if err != nil {
		t.Fatal(err)
	}
	if stAb.CanSubmit {
		t.Fatal("student should not be able to submit with missing materials")
	}
	if len(stAb.MissingRequired) != 1 || stAb.MissingRequired[0].Name != "成绩单" {
		t.Fatalf("missing list wrong: %+v", stAb.MissingRequired)
	}
	if !stAb.CanEditDocuments || !stAb.Editable {
		t.Fatal("student should be able to edit while preparing")
	}

	a.Status = constants.AppStatusSubmitted
	stAb2, err := appSvc.Abilities(a, studentID, constants.RoleStudent)
	if err != nil {
		t.Fatal(err)
	}
	if stAb2.CanEditDocuments || stAb2.Editable {
		t.Fatal("student editing must be locked after submission")
	}
	if stAb2.CanReturn || !stAb2.IsStudentOwner {
		t.Fatal("abilities wrong after submission")
	}

	cAb, err := appSvc.Abilities(a, counselorID, constants.RoleCounselor)
	if err != nil {
		t.Fatal(err)
	}
	if !cAb.CanReturn || !cAb.CanReviewMaterials || cAb.IsStudentOwner {
		t.Fatalf("counselor abilities wrong: %+v", cAb)
	}
}

func TestUploadRequiresFileAndPendingCannotBeApproved(t *testing.T) {
	_, matSvc, a := newServices()
	const studentID, counselorID = uint(10), uint(20)
	m := &model.MaterialItem{ApplicationID: a.ID, Name: "成绩单", IsRequired: true}
	if _, err := matSvc.Create(a.ID, counselorID, constants.RoleCounselor, m); err != nil {
		t.Fatal(err)
	}
	if _, err := matSvc.UpdateStatus(studentID, m.ID, constants.RoleStudent, constants.MaterialUploaded, "", ""); err == nil {
		t.Fatal("upload without file must fail")
	}
	if _, err := matSvc.UpdateStatus(counselorID, m.ID, constants.RoleCounselor, constants.MaterialApproved, "", ""); !isConflict(err) {
		t.Fatalf("approve pending must be conflict, got %v", err)
	}
}

func TestStudentCannotAccessForeignApplication(t *testing.T) {
	appSvc, _, a := newServices()
	// another student tries to read / transition
	if _, err := appSvc.Get(a.ID, 99, constants.RoleStudent); !isForbidden(err) {
		t.Fatalf("foreign student read must be forbidden, got %v", err)
	}
}
