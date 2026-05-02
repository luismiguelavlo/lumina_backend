package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"library_back/internal/repositories"
	"library_back/internal/services"

	"github.com/gin-gonic/gin"
)

func TestAuthBearerMiddleware_200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("m", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	access, _, _, err := tok.GeneratePair(ctx, "user-mw-1")
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/auth/me", AuthBearerMiddleware(tok), func(c *gin.Context) {
		v, _ := c.Get(ContextAuthUserID)
		c.String(http.StatusOK, v.(string))
	})
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "user-mw-1" {
		t.Fatalf("got %d body %q", w.Code, w.Body.String())
	}
}

func TestAuthBearerMiddleware_401NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("m", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	r.GET("/auth/me", AuthBearerMiddleware(tok), func(c *gin.Context) { c.Status(200) })
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}

func TestAuthBearerMiddleware_401Revoked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("m", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	access, _, _, err := tok.GeneratePair(ctx, "u-rev")
	if err != nil {
		t.Fatal(err)
	}
	pa, err := tok.ParseAccess(ctx, access)
	if err != nil {
		t.Fatal(err)
	}
	if err := tok.Revoke(ctx, pa.UserID, pa.JTI, pa.ExpiresAt); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/auth/me", AuthBearerMiddleware(tok), func(c *gin.Context) { c.Status(200) })
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}
