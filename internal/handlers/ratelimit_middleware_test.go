package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"library_back/internal/pkg/ratelimit"

	"github.com/gin-gonic/gin"
)

func TestIPRateLimit_HealthSkipped(t *testing.T) {
	gin.SetMode(gin.TestMode)
	st := ratelimit.NewStore(ratelimit.Config{Enabled: true, GeneralRPM: 1, RankingRPM: 1})
	r := gin.New()
	r.Use(IPRateLimit(st))
	r.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("iter %d: %d", i, w.Code)
		}
	}
}

func TestIPRateLimit_429General(t *testing.T) {
	gin.SetMode(gin.TestMode)
	st := ratelimit.NewStore(ratelimit.Config{Enabled: true, GeneralRPM: 2, RankingRPM: 99})
	r := gin.New()
	r.Use(IPRateLimit(st))
	r.GET("/api/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/x", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("iter %d: %d", i, w.Code)
		}
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/x", nil))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("want 429 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestIPRateLimit_RankingBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	st := ratelimit.NewStore(ratelimit.Config{Enabled: true, GeneralRPM: 99, RankingRPM: 1})
	r := gin.New()
	r.Use(IPRateLimit(st))
	r.GET("/api/ranking/top3", func(c *gin.Context) { c.Status(http.StatusOK) })
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, httptest.NewRequest(http.MethodGet, "/api/ranking/top3", nil))
	if w1.Code != http.StatusOK {
		t.Fatalf("first %d", w1.Code)
	}
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/ranking/top3", nil))
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("want 429 got %d", w2.Code)
	}
}
