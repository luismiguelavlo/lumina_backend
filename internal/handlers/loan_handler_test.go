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
	"library_back/internal/services/loan"

	"github.com/gin-gonic/gin"
)

type mockLoanSvc struct {
	listItems []models.LoanListItem
	listTotal int64
	listErr   error
	createOut *models.LoanDetailResponse
	createErr error
	returnOut *models.LoanDetailResponse
	returnErr error
}

func (m *mockLoanSvc) List(ctx context.Context, filter repositories.LoanListFilter, limit, offset int) ([]models.LoanListItem, int64, error) {
	_ = filter
	_ = limit
	_ = offset
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.listItems, m.listTotal, nil
}

func (m *mockLoanSvc) Create(ctx context.Context, req models.CreateLoanRequest, adminID string) (*models.LoanDetailResponse, error) {
	_ = ctx
	_ = req
	_ = adminID
	return m.createOut, m.createErr
}

func (m *mockLoanSvc) Return(ctx context.Context, loanID string, actorUserID string) (*models.LoanDetailResponse, error) {
	_ = ctx
	_ = loanID
	_ = actorUserID
	return m.returnOut, m.returnErr
}

var _ loanServiceAPI = (*mockLoanSvc)(nil)

func TestLoanHandler_List200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{
		listItems: []models.LoanListItem{{LoanID: "l1", BookTitle: "T", Borrower: "A B", DueDate: "2026-05-01", Status: "active"}},
		listTotal: 1,
	}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.GET("/loans", h.List)
	req := httptest.NewRequest(http.MethodGet, "/loans", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body models.ListLoansResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total != 1 || len(body.Data) != 1 {
		t.Fatalf("payload %+v", body)
	}
}

func TestLoanHandler_List400BadStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{listErr: loan.ErrInvalidStudentIDQuery}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.GET("/loans", h.List)
	req := httptest.NewRequest(http.MethodGet, "/loans?student_id=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want400 got %d %s", w.Code, w.Body.String())
	}
}

func TestLoanHandler_List400InvalidStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{listErr: loan.ErrInvalidLoanStatus}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.GET("/loans", h.List)
	req := httptest.NewRequest(http.MethodGet, "/loans?status=invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestLoanHandler_Create201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{
		createOut: &models.LoanDetailResponse{
			LoanID: "L1", BookID: "b1", BookTitle: "T", Borrower: "A B",
			DueDate: "2026-06-01", Status: "active", ISBN: "x",
		},
	}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2026-12-31",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
}

func TestLoanHandler_Create404Book(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{createErr: loan.ErrBookNotFound}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2026-12-31",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestLoanHandler_Create404Student(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{createErr: loan.ErrStudentNotFound}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2026-12-31",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestLoanHandler_Create400DueDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{createErr: loan.ErrDueDateInPast}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2000-01-01",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestLoanHandler_Create500NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{createOut: &models.LoanDetailResponse{LoanID: "L1"}}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", h.Create)
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2026-12-31",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 got %d", w.Code)
	}
}

func TestLoanHandler_Return200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{
		returnOut: &models.LoanDetailResponse{LoanID: "L1", Status: "returned"},
	}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.PATCH("/loans/:id/return", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Return(c)
	})
	req := httptest.NewRequest(http.MethodPatch, "/loans/L1/return", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestLoanHandler_Return404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{returnErr: loan.ErrLoanNotFound}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.PATCH("/loans/:id/return", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Return(c)
	})
	req := httptest.NewRequest(http.MethodPatch, "/loans/nope/return", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}

func TestLoanHandler_Return409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{returnErr: loan.ErrLoanAlreadyReturned}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.PATCH("/loans/:id/return", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Return(c)
	})
	req := httptest.NewRequest(http.MethodPatch, "/loans/L1/return", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("want 409 got %d", w.Code)
	}
}

func TestLoanHandler_Return500NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{returnOut: &models.LoanDetailResponse{LoanID: "L1"}}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.PATCH("/loans/:id/return", h.Return)
	req := httptest.NewRequest(http.MethodPatch, "/loans/L1/return", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 got %d", w.Code)
	}
}

func TestLoanHandler_Create400ValidationMissingField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{"book_id": "00000000-0000-0000-0000-000000000002", "due_date": "2026-12-31"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d %s", w.Code, w.Body.String())
	}
}

func TestLoanHandler_Create409NoCopies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{createErr: loan.ErrNoCopiesAvailable}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2026-12-31",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("want 409 got %d", w.Code)
	}
}

func TestLoanHandler_Create403ActiveSanction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockLoanSvc{createErr: loan.ErrStudentHasActiveSanction}
	h := &LoanHandler{svc: ms}
	r := gin.New()
	r.POST("/loans", func(c *gin.Context) {
		c.Set(ContextAuthUserID, "admin-1")
		h.Create(c)
	})
	body := map[string]string{
		"student_id": "00000000-0000-0000-0000-000000000001",
		"book_id":    "00000000-0000-0000-0000-000000000002",
		"due_date":   "2026-12-31",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("want 403 got %d %s", w.Code, w.Body.String())
	}
}
