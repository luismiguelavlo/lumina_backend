package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"library_back/internal/models"
	"library_back/internal/repositories"
	"library_back/internal/services/sanction"

	"github.com/gin-gonic/gin"
)

type mockSanctionSvc struct {
	listItems []models.SanctionListItem
	listTotal int64
	listErr   error
	createOut *models.SanctionResponse
	createErr error
	liftOut   *models.SanctionResponse
	liftErr   error
}

func (m *mockSanctionSvc) ListActive(ctx context.Context, filter repositories.SanctionListFilter, limit, offset int) ([]models.SanctionListItem, int64, error) {
	_ = filter
	_ = limit
	_ = offset
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.listItems, m.listTotal, nil
}

func (m *mockSanctionSvc) Create(ctx context.Context, req models.CreateSanctionRequest, adminID string) (*models.SanctionResponse, error) {
	_ = ctx
	_ = req
	_ = adminID
	return m.createOut, m.createErr
}

func (m *mockSanctionSvc) Lift(ctx context.Context, sanctionID, adminID string) (*models.SanctionResponse, error) {
	_ = ctx
	_ = sanctionID
	_ = adminID
	return m.liftOut, m.liftErr
}

func TestSanctionHandler_List200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SanctionHandler{svc: &mockSanctionSvc{
		listItems: []models.SanctionListItem{{SanctionID: "s1", Reason: "r", StudentIDCode: "X"}},
		listTotal: 1,
	}}
	rgin := gin.New()
	rgin.GET("/sanctions", h.List)
	w := httptest.NewRecorder()
	rgin.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/sanctions", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body models.ListSanctionsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total != 1 {
		t.Fatalf("%+v", body)
	}
}

func TestSanctionHandler_Create201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SanctionHandler{svc: &mockSanctionSvc{createOut: &models.SanctionResponse{SanctionID: "x", Status: "active"}}}
	body := `{"student_id":"10000000-0000-4000-8000-000000000001","reason":"Motivo"}`
	req := httptest.NewRequest(http.MethodPost, "/sanctions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rgin := gin.New()
	rgin.POST("/sanctions", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	rgin.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
}

func TestSanctionHandler_List400BadStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SanctionHandler{svc: &mockSanctionSvc{listErr: sanction.ErrInvalidStudentIDQuery}}
	rgin := gin.New()
	rgin.GET("/sanctions", h.List)
	w := httptest.NewRecorder()
	rgin.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/sanctions?student_id=bad", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestSanctionHandler_Create404Student(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SanctionHandler{svc: &mockSanctionSvc{createErr: sanction.ErrStudentNotFound}}
	rgin := gin.New()
	rgin.POST("/sanctions", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := `{"student_id":"10000000-0000-4000-8000-000000000001","reason":"x"}`
	req := httptest.NewRequest(http.MethodPost, "/sanctions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rgin.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestSanctionHandler_Lift200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SanctionHandler{svc: &mockSanctionSvc{liftOut: &models.SanctionResponse{SanctionID: "s", Status: "lifted"}}}
	rgin := gin.New()
	rgin.PATCH("/sanctions/:id", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Lift(c)
	})
	w := httptest.NewRecorder()
	rgin.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/sanctions/s", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestSanctionHandler_Lift404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &SanctionHandler{svc: &mockSanctionSvc{liftErr: sanction.ErrSanctionNotFound}}
	rgin := gin.New()
	rgin.PATCH("/sanctions/:id", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Lift(c)
	})
	w := httptest.NewRecorder()
	rgin.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/sanctions/missing", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

var _ sanctionServiceAPI = (*sanction.Service)(nil)
