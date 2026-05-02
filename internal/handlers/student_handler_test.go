package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"library_back/internal/services/student"
	"net/http"
	"net/http/httptest"
	"testing"

	"library_back/internal/models"

	"github.com/gin-gonic/gin"
)

type mockStudentSvc struct {
	createOut   *models.StudentResponse
	createErr   error
	getOut      *models.StudentResponse
	getErr      error
	listOut     []models.StudentResponse
	listTotal   int64
	listErr     error
	updateOut   *models.StudentResponse
	updateErr   error
	deactErr    error
	profileOut  *models.StudentProfileResponse
	profileErr  error
	awardOut    *models.AwardBadgeResponse
	awardErr    error
	lastAwardID string
	revokeErr   error
}

func (m *mockStudentSvc) Create(ctx context.Context, req models.CreateStudentRequest, adminID string) (*models.StudentResponse, error) {
	_ = ctx
	_ = req
	_ = adminID
	return m.createOut, m.createErr
}

func (m *mockStudentSvc) GetByID(ctx context.Context, id string) (*models.StudentResponse, error) {
	_ = ctx
	_ = id
	return m.getOut, m.getErr
}

func (m *mockStudentSvc) List(ctx context.Context, filter models.StudentFilter, limit, offset int) ([]models.StudentResponse, int64, error) {
	_ = ctx
	_ = filter
	_ = limit
	_ = offset
	return m.listOut, m.listTotal, m.listErr
}

func (m *mockStudentSvc) Update(ctx context.Context, id string, req models.UpdateStudentRequest) (*models.StudentResponse, error) {
	_ = ctx
	_ = id
	_ = req
	return m.updateOut, m.updateErr
}

func (m *mockStudentSvc) Deactivate(ctx context.Context, id string) error {
	_ = ctx
	_ = id
	return m.deactErr
}

func (m *mockStudentSvc) GetProfile(ctx context.Context, id string, loanLimit int) (*models.StudentProfileResponse, error) {
	_ = ctx
	_ = id
	_ = loanLimit
	return m.profileOut, m.profileErr
}

func (m *mockStudentSvc) AwardBadge(ctx context.Context, studentID, badgeID string) (*models.AwardBadgeResponse, error) {
	_ = ctx
	m.lastAwardID = badgeID
	return m.awardOut, m.awardErr
}

func (m *mockStudentSvc) RevokeBadge(ctx context.Context, studentID, badgeID string) error {
	_ = ctx
	_ = studentID
	_ = badgeID
	return m.revokeErr
}

func newTestStudentHandler(svc studentServiceAPI) *StudentHandler {
	return &StudentHandler{svc: svc}
}

func TestStudentHandler_Create201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{
		createOut: &models.StudentResponse{ID: "s1", StudentIDCode: "X", FirstName: "A", LastName: "B", IsActive: true},
	}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.POST("/students", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{"first_name": "A", "last_name": "B", "student_id_code": "CODE-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
}

func TestStudentHandler_Create201WithoutStudentIDCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{
		createOut: &models.StudentResponse{ID: "s1", StudentIDCode: "AUTO-20260418-abcdef12", FirstName: "A", LastName: "B", IsActive: true},
	}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.POST("/students", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{"first_name": "A", "last_name": "B"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
}

func TestStudentHandler_Create409Duplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{createErr: student.ErrDuplicateStudentEmail}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.POST("/students", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{"first_name": "A", "last_name": "B", "student_id_code": "CODE-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("want 409 got %d", w.Code)
	}
}

func TestStudentHandler_Create500NoAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{createOut: &models.StudentResponse{ID: "s1", StudentIDCode: "X", FirstName: "A", LastName: "B"}}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.POST("/students", h.Create)
	body := map[string]string{"first_name": "A", "last_name": "B", "student_id_code": "CODE-1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 got %d", w.Code)
	}
}

func TestStudentHandler_Get404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{getErr: student.ErrStudentNotFound}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.GET("/students/:id", h.Get)
	req := httptest.NewRequest(http.MethodGet, "/students/nope", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestStudentHandler_List200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{listOut: []models.StudentResponse{}, listTotal: 0}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.GET("/students", h.List)
	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", w.Code)
	}
}

func TestStudentHandler_AwardBadge201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{
		awardOut: &models.AwardBadgeResponse{ID: "b1", Slug: "x", Name: "X", EarnedAt: "2020-01-01T00:00:00Z"},
	}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.POST("/students/:id/badges", h.AwardBadge)
	badgeUUID := "11111111-1111-4111-8111-111111111111"
	body := map[string]string{"badge_id": badgeUUID}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/students/stu-1/badges", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("want 201 got %d %s", w.Code, w.Body.String())
	}
	if ms.lastAwardID != badgeUUID {
		t.Fatalf("badge id not passed")
	}
}

func TestStudentHandler_RevokeBadge204(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.DELETE("/students/:id/badges/:badgeId", h.RevokeBadge)
	badgeUUID := "22222222-2222-4222-8222-222222222222"
	req := httptest.NewRequest(http.MethodDelete, "/students/stu-1/badges/"+badgeUUID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204 got %d %s", w.Code, w.Body.String())
	}
}

func TestStudentHandler_RevokeBadge400BadUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestStudentHandler(&mockStudentSvc{})
	r := gin.New()
	r.DELETE("/students/:id/badges/:badgeId", h.RevokeBadge)
	req := httptest.NewRequest(http.MethodDelete, "/students/stu-1/badges/not-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestStudentHandler_RevokeBadge404NotEarned(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{revokeErr: student.ErrStudentDoesNotHaveBadge}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.DELETE("/students/:id/badges/:badgeId", h.RevokeBadge)
	req := httptest.NewRequest(http.MethodDelete, "/students/stu-1/badges/33333333-3333-4333-8333-333333333333", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestStudentHandler_DepartmentBadRequest400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockStudentSvc{createErr: student.ErrDepartmentNotFound}
	h := newTestStudentHandler(ms)
	r := gin.New()
	r.POST("/students", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{"first_name": "A", "last_name": "B", "student_id_code": "CODE-1", "department_id": "00000000-0000-0000-0000-000000000000"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d %s", w.Code, w.Body.String())
	}
}
