package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"library_back/internal/models"
	"library_back/internal/services/analytics"

	"github.com/gin-gonic/gin"
)

type analyticsServiceAPI interface {
	GetDashboard(ctx context.Context, limits models.DashboardLimits) (*models.DashboardResponse, error)
}

// AnalyticsHandler exposes GET /api/analytics/dashboard (staff JWT).
type AnalyticsHandler struct {
	svc analyticsServiceAPI
}

// NewAnalyticsHandler constructs AnalyticsHandler.
func NewAnalyticsHandler(svc *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// Dashboard GET /api/analytics/dashboard
func (h *AnalyticsHandler) Dashboard(c *gin.Context) {
	limits := models.DashboardLimits{
		MostBorrowedLimit:   parsePositiveIntQuery(c.Query("most_borrowed_limit"), 0),
		RecentActivityLimit: parsePositiveIntQuery(c.Query("recent_activity_limit"), 0),
		TopOverdueLimit:     parsePositiveIntQuery(c.Query("top_overdue_limit"), 0),
	}
	out, err := h.svc.GetDashboard(c.Request.Context(), limits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, out)
}

func parsePositiveIntQuery(raw string, zero int) int {
	s := strings.TrimSpace(raw)
	if s == "" {
		return zero
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return zero
	}
	return n
}

var _ analyticsServiceAPI = (*analytics.Service)(nil)
