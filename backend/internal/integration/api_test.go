package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
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
		return nil, errNotFound
	}
	cp := *a
	return &cp, nil
}
func (r *fakeAppRepo) Update(a *model.ApplicationProject) error                 { r.apps[a.ID] = a; return nil }
func (r *fakeAppRepo) ListByStudent(uint) ([]model.ApplicationProject, error)   { return nil, nil }
func (r *fakeAppRepo) ListByCounselor(uint) ([]model.ApplicationProject, error) { return nil, nil }
func (r *fakeAppRepo) ListAll() ([]model.ApplicationProject, error)             { return nil, nil }

type fakeUnivRepo struct{}

func (r *fakeUnivRepo) FindByID(id uint) (*model.University, error) {
	return &model.University{ID: id}, nil
}

type fakeMatRepo struct{ items map[uint]*model.MaterialItem }

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
		return nil, errNotFound
	}
	cp := *m
	return &cp, nil
}
func (r *fakeMatRepo) Update(m *model.MaterialItem) error { r.items[m.ID] = m; return nil }
func (r *fakeMatRepo) ListByApplication(appID uint) ([]model.MaterialItem, error) {
	out := make([]model.MaterialItem, 0)
	for _, m := range r.items {
		if m.ApplicationID == appID {
			out = append(out, *m)
		}
	}
	return out, nil
}

var errNotFound = repository.ErrNotFound

// ---- test engine -----------------------------------------------------------

func setupRouter(t *testing.T) (*gin.Engine, *fakeAppRepo, *fakeMatRepo) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	appRepo := &fakeAppRepo{apps: map[uint]*model.ApplicationProject{}}
	matRepo := &fakeMatRepo{items: map[uint]*model.MaterialItem{}}
	appSvc := service.NewApplicationService(appRepo, &fakeUnivRepo{}, matRepo, logger)
	matSvc := service.NewMaterialService(matRepo, appRepo, logger)
	appH := handler.NewApplicationHandler(appSvc, logger)
	matH := handler.NewMaterialHandler(matSvc, logger)

	cfg := &config.Config{JWTSecret: "integration-secret"}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(logger))
	v1 := r.Group("/api/v1")
	auth := v1.Group("", middleware.AuthRequired(cfg))
	auth.GET("/applications/:id/abilities", appH.Abilities)
	auth.GET("/applications/:id", appH.Get)
	auth.PUT("/applications/:id/status", appH.UpdateStatus)
	auth.GET("/applications/:id/materials", matH.ListByApplication)
	auth.POST("/applications/:id/materials", matH.Create)
	v1.PUT("/materials/:id/status", middleware.AuthRequired(cfg), matH.UpdateStatus)
	return r, appRepo, matRepo
}

func token(t *testing.T, userID uint, role string) string {
	t.Helper()
	tok, err := util.GenerateToken(userID, role, role, "integration-secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

func do(t *testing.T, r *gin.Engine, method, path, tok string, body interface{}) (int, map[string]interface{}) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var out map[string]interface{}
	if w.Body.Len() > 0 {
		_ = json.Unmarshal(w.Body.Bytes(), &out)
	}
	return w.Code, out
}

// ---- scenarios -------------------------------------------------------------

func TestHTTPSubmitGateAndRoleFlow(t *testing.T) {
	r, appRepo, matRepo := setupRouter(t)
	const appID, student, counselor = uint(1), uint(10), uint(20)
	appRepo.apps[appID] = &model.ApplicationProject{ID: appID, StudentID: student, CounselorID: counselor, UniversityID: 1, Status: "preparing"}
	req1 := &model.MaterialItem{ID: 1, ApplicationID: appID, Name: "成绩单", IsRequired: true, Status: "pending"}
	req2 := &model.MaterialItem{ID: 2, ApplicationID: appID, Name: "语言成绩", IsRequired: true, Status: "pending"}
	matRepo.items[1] = req1
	matRepo.items[2] = req2

	stTok := token(t, student, "student")
	coTok := token(t, counselor, "counselor")
	otherCoTok := token(t, 99, "counselor")

	// abilities show blocked submit with reasons
	code, body := do(t, r, "GET", "/api/v1/applications/1/abilities", stTok, nil)
	if code != 200 {
		t.Fatalf("abilities status = %d body=%v", code, body)
	}
	data := body["data"].(map[string]interface{})
	if data["can_submit"].(bool) {
		t.Fatal("can_submit must be false with pending required materials")
	}
	if data["editable"].(bool) != true || data["can_edit_documents"].(bool) != true {
		t.Fatal("student editing should be allowed while preparing")
	}

	// submit blocked: 409 with details listing missing materials
	code, body = do(t, r, "PUT", "/api/v1/applications/1/status", stTok, map[string]string{"status": "submitted"})
	if code != 409 {
		t.Fatalf("submit status = %d body=%v", code, body)
	}
	if body["code"].(float64) != 40900 {
		t.Fatalf("code = %v", body["code"])
	}
	details, _ := body["details"].(map[string]interface{})
	missing, _ := details["missing_materials"].([]interface{})
	if len(missing) != 2 {
		t.Fatalf("details.missing_materials = %v", missing)
	}

	// student upload without file_url rejected (422)
	code, _ = do(t, r, "PUT", "/api/v1/materials/1/status", stTok, map[string]string{"status": "uploaded"})
	if code != 422 {
		t.Fatalf("upload without file must be 422, got %d", code)
	}
	// student upload ok
	code, _ = do(t, r, "PUT", "/api/v1/materials/1/status", stTok, map[string]string{"status": "uploaded", "file_url": "/api/v1/files/a.pdf"})
	if code != 200 {
		t.Fatalf("student upload = %d", code)
	}
	// student self-approval rejected (now in uploaded state)
	code, _ = do(t, r, "PUT", "/api/v1/materials/1/status", stTok, map[string]string{"status": "approved"})
	if code != 403 {
		t.Fatalf("student approve must be 403, got %d", code)
	}
	// another counselor cannot approve
	code, _ = do(t, r, "PUT", "/api/v1/materials/1/status", otherCoTok, map[string]string{"status": "approved"})
	if code != 403 {
		t.Fatalf("other counselor approve must be 403, got %d", code)
	}
	// responsible counselor approves both (student uploads 2 first)
	code, _ = do(t, r, "PUT", "/api/v1/materials/2/status", stTok, map[string]string{"status": "uploaded", "file_url": "/api/v1/files/b.pdf"})
	if code != 200 {
		t.Fatalf("upload 2 = %d", code)
	}
	code, _ = do(t, r, "PUT", "/api/v1/materials/1/status", coTok, map[string]string{"status": "approved"})
	if code != 200 {
		t.Fatalf("counselor approve 1 = %d", code)
	}
	code, _ = do(t, r, "PUT", "/api/v1/materials/2/status", coTok, map[string]string{"status": "approved", "review_remark": "ok"})
	if code != 200 {
		t.Fatalf("counselor approve 2 = %d", code)
	}

	// submit now succeeds and locks editing
	code, body = do(t, r, "PUT", "/api/v1/applications/1/status", stTok, map[string]string{"status": "submitted"})
	if code != 200 {
		t.Fatalf("submit = %d body=%v", code, body)
	}
	code, body = do(t, r, "GET", "/api/v1/applications/1/abilities", stTok, nil)
	data = body["data"].(map[string]interface{})
	if data["editable"].(bool) || data["can_edit_documents"].(bool) {
		t.Fatal("editing must be locked after submission")
	}
	// post-submit student upload rejected (409)
	code, _ = do(t, r, "PUT", "/api/v1/materials/1/status", stTok, map[string]string{"status": "uploaded", "file_url": "/x.pdf"})
	if code != 409 {
		t.Fatalf("post-submit upload must be 409, got %d", code)
	}

	// responsible counselor returns -> preparing, editing restored
	code, _ = do(t, r, "PUT", "/api/v1/applications/1/status", coTok, map[string]string{"status": "preparing"})
	if code != 200 {
		t.Fatalf("counselor return = %d", code)
	}
	code, body = do(t, r, "GET", "/api/v1/applications/1/abilities", stTok, nil)
	data = body["data"].(map[string]interface{})
	if !data["editable"].(bool) {
		t.Fatal("editing must be restored after return")
	}
}

func TestHTTPForeignStudentForbidden(t *testing.T) {
	r, appRepo, _ := setupRouter(t)
	appRepo.apps[1] = &model.ApplicationProject{ID: 1, StudentID: 10, CounselorID: 20, UniversityID: 1, Status: "preparing"}
	tok := token(t, 77, "student")
	code, body := do(t, r, "GET", "/api/v1/applications/1/abilities", tok, nil)
	if code != 403 {
		t.Fatalf("foreign student = %d body=%v", code, body)
	}
}

func TestHTTPStudentCannotCreateMaterial(t *testing.T) {
	r, appRepo, _ := setupRouter(t)
	appRepo.apps[1] = &model.ApplicationProject{ID: 1, StudentID: 10, CounselorID: 20, UniversityID: 1, Status: "preparing"}
	tok := token(t, 10, "student")
	code, _ := do(t, r, "POST", "/api/v1/applications/1/materials", tok, map[string]interface{}{"name": "新必交", "is_required": true})
	if code != 403 {
		t.Fatalf("student create material must be 403, got %d", code)
	}
}

func TestHTTPUnauthorized(t *testing.T) {
	r, _, _ := setupRouter(t)
	code, _ := do(t, r, "GET", "/api/v1/applications/1/abilities", "", nil)
	if code != http.StatusUnauthorized {
		t.Fatalf("no token = %d", code)
	}
}
