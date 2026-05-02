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
	finesvc "library_back/internal/services/fine"

	"github.com/gin-gonic/gin"
)

type mockFineSvc struct {
	listItems []models.FineListItem
	listTotal int64
	listErr   error
	createOut *models.FineResponse
	createErr error
	getOut    *models.FineResponse
	getErr    error
	paidOut   *models.FineResponse
	paidErr   error
	waivedOut *models.FineResponse
	waivedErr error
}

func (m *mockFineSvc) List(ctx context.Context, filter repositories.FineListFilter, limit, offset int) ([]models.FineListItem, int64, error) {
	_ = filter
	_ = limit
	_ = offset
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.listItems, m.listTotal, nil
}

func (m *mockFineSvc) Create(ctx context.Context, req models.CreateFineRequest) (*models.FineResponse, error) {
	_ = ctx
	_ = req
	return m.createOut, m.createErr
}

func (m *mockFineSvc) GetByID(ctx context.Context, id string) (*models.FineResponse, error) {
	_ = ctx
	_ = id
	return m.getOut, m.getErr
}

func (m *mockFineSvc) MarkPaid(ctx context.Context, id string) (*models.FineResponse, error) {
	_ = ctx
	_ = id
	return m.paidOut, m.paidErr
}

func (m *mockFineSvc) MarkWaived(ctx context.Context, id string) (*models.FineResponse, error) {
	_ = ctx
	_ = id
	return m.waivedOut, m.waivedErr
}

var _ fineServiceAPI = (*mockFineSvc)(nil)

func TestFineHandler_List200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockFineSvc{
		listItems: []models.FineListItem{{ID: "f1", Amount: 10, Status: "pending"}},
		listTotal: 1,
	}
	h := &FineHandler{svc: ms}
	r := gin.New()
	r.GET("/fines", h.List)
	req := httptest.NewRequest(http.MethodGet, "/fines", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body models.ListFinesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total != 1 || len(body.Data) != 1 {
		t.Fatalf("%+v", body)
	}
}

func TestFineHandler_List400BadStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockFineSvc{listErr: finesvc.ErrInvalidStudentIDQuery}
	h := &FineHandler{svc: ms}
	r := gin.New()
	r.GET("/fines", h.List)
	req := httptest.NewRequest(http.MethodGet, "/fines?student_id=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestFineHandler_Create201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockFineSvc{createOut: &models.FineResponse{ID: "x", Status: "pending", Amount: 5}}
	h := &FineHandler{svc: ms}
	r := gin.New()
	r.POST("/fines", h.Create)
	body := `{"loan_id":"10000000-0000-4000-8000-000000000001","student_id":"20000000-0000-4000-8000-000000000002","amount":5}`
	req := httptest.NewRequest(http.MethodPost, "/fines", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
}

func TestFineHandler_Create400Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FineHandler{svc: &mockFineSvc{}}
	r := gin.New()
	r.POST("/fines", h.Create)
	body := `{"loan_id":"not-uuid","student_id":"20000000-0000-4000-8000-000000000002","amount":1}`
	req := httptest.NewRequest(http.MethodPost, "/fines", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestFineHandler_Create404Loan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockFineSvc{createErr: finesvc.ErrFineLoanNotFound}
	h := &FineHandler{svc: ms}
	r := gin.New()
	r.POST("/fines", h.Create)
	body := `{"loan_id":"10000000-0000-4000-8000-000000000001","student_id":"20000000-0000-4000-8000-000000000002","amount":1}`
	req := httptest.NewRequest(http.MethodPost, "/fines", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestFineHandler_Create409Duplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockFineSvc{createErr: finesvc.ErrFineDuplicatePending}
	h := &FineHandler{svc: ms}
	r := gin.New()
	r.POST("/fines", h.Create)
	body := `{"loan_id":"10000000-0000-4000-8000-000000000001","student_id":"20000000-0000-4000-8000-000000000002","amount":1}`
	req := httptest.NewRequest(http.MethodPost, "/fines", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("want 409 got %d", w.Code)
	}
}

func TestFineHandler_Get404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FineHandler{svc: &mockFineSvc{getErr: finesvc.ErrFineNotFound}}
	r := gin.New()
	r.GET("/fines/:id", h.Get)
	req := httptest.NewRequest(http.MethodGet, "/fines/nope", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestFineHandler_MarkPaid400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FineHandler{svc: &mockFineSvc{paidErr: finesvc.ErrFineNotPending}}
	r := gin.New()
	r.PATCH("/fines/:id/paid", h.MarkPaid)
	req := httptest.NewRequest(http.MethodPatch, "/fines/x/paid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestFineHandler_MarkWaived200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FineHandler{svc: &mockFineSvc{waivedOut: &models.FineResponse{ID: "f", Status: "waived"}}}
	r := gin.New()
	r.PATCH("/fines/:id/waived", h.MarkWaived)
	req := httptest.NewRequest(http.MethodPatch, "/fines/f/waived", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

var _ fineServiceAPI = (*finesvc.Service)(nil)
