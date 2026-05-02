package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"library_back/internal/models"

	"github.com/gin-gonic/gin"
)

type mockAnalyticsSvc struct {
	out *models.DashboardResponse
	err error
	got models.DashboardLimits
}

func (m *mockAnalyticsSvc) GetDashboard(ctx context.Context, limits models.DashboardLimits) (*models.DashboardResponse, error) {
	_ = ctx
	m.got = limits
	if m.err != nil {
		return nil, m.err
	}
	return m.out, nil
}

func TestAnalyticsHandler_Dashboard200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockAnalyticsSvc{
		out: &models.DashboardResponse{
			TotalBooks: 1, ActiveStudents: 2, OverdueFines: 3.5,
			MostBorrowedBooks: []models.MostBorrowedItem{},
			RecentActivity:    []models.RecentActivityItem{{ID: "x", EventType: "e", Title: "t", CreatedAt: time.Now().UTC()}},
		},
	}
	h := &AnalyticsHandler{svc: ms}
	r := gin.New()
	r.GET("/analytics/dashboard", h.Dashboard)
	req := httptest.NewRequest(http.MethodGet, "/analytics/dashboard?most_borrowed_limit=7&recent_activity_limit=15", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
	if ms.got.MostBorrowedLimit != 7 || ms.got.RecentActivityLimit != 15 {
		t.Fatalf("limits %+v", ms.got)
	}
	var body models.DashboardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.TotalBooks != 1 || len(body.RecentActivity) != 1 {
		t.Fatalf("%+v", body)
	}
}

func TestAnalyticsHandler_Dashboard500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockAnalyticsSvc{err: errors.New("db")}
	h := &AnalyticsHandler{svc: ms}
	r := gin.New()
	r.GET("/analytics/dashboard", h.Dashboard)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/analytics/dashboard", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 got %d", w.Code)
	}
}
