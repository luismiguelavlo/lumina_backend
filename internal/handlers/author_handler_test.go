package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"library_back/internal/models"

	"github.com/gin-gonic/gin"
)

type mockAuthorSvc struct {
	out *models.AuthorCreatedResponse
	err error
}

func (m *mockAuthorSvc) Create(ctx context.Context, req models.CreateAuthorRequest) (*models.AuthorCreatedResponse, error) {
	_ = ctx
	_ = req
	if m.err != nil {
		return nil, m.err
	}
	return m.out, nil
}

func TestAuthorHandler_Create201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{svc: &mockAuthorSvc{
		out: &models.AuthorCreatedResponse{ID: "a1", Name: "Jane Doe", CreatedAt: "2026-01-01T00:00:00Z"},
	}}
	r := gin.New()
	r.POST("/authors", h.Create)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewBufferString(`{"name":"Jane Doe"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
}

func TestAuthorHandler_Create400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{svc: &mockAuthorSvc{}}
	r := gin.New()
	r.POST("/authors", h.Create)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestAuthorHandler_Create500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{svc: &mockAuthorSvc{err: errors.New("db")}}
	r := gin.New()
	r.POST("/authors", h.Create)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewBufferString(`{"name":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 got %d", w.Code)
	}
}
