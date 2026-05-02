package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"library_back/internal/services/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"library_back/internal/models"
	"library_back/internal/pkg/crypto"
	"library_back/internal/repositories"
	"library_back/internal/services"

	"github.com/gin-gonic/gin"
)

func TestAuth_Register201(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	authSvc := auth.NewAuthService(users, tok)
	h := NewAuthHandler(authSvc)
	r := gin.New()
	r.POST("/auth/register", h.Register)

	body := map[string]any{
		"first_name": "Ana",
		"last_name":  "López",
		"email":      "ana@example.com",
		"password":   "password1",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestAuth_Login401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	authSvc := auth.NewAuthService(users, tok)
	h := NewAuthHandler(authSvc)
	r := gin.New()
	r.POST("/auth/login", h.Login)

	body := map[string]string{"email": "nope@test.com", "password": "x"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", w.Code)
	}
}

func TestAuth_Logout204(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	hash, _ := crypto.HashPassword("pw123456")
	u := &models.User{
		ID: "u-log", FirstName: "L", LastName: "G", Email: "lg@test.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	authSvc := auth.NewAuthService(users, tok)
	h := NewAuthHandler(authSvc)
	r := gin.New()
	r.POST("/auth/logout", h.Logout)

	access, _, _, err := tok.GeneratePair(ctx, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+access)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204 got %d %s", w.Code, w.Body.String())
	}
}

func TestAuth_Register201PayloadAnd409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	authSvc := auth.NewAuthService(users, tok)
	h := NewAuthHandler(authSvc)
	r := gin.New()
	r.POST("/auth/register", h.Register)

	body := map[string]any{
		"first_name": "Ana",
		"last_name":  "López",
		"email":      "ana2@example.com",
		"password":   "password1",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("201: status %d %s", w.Code, w.Body.String())
	}
	var out models.RegisterResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.IsActive || out.ID == "" || out.Email != "ana2@example.com" {
		t.Fatalf("unexpected payload %+v", out)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusConflict {
		t.Fatalf("409: want conflict got %d %s", w2.Code, w2.Body.String())
	}
}

func TestAuth_Register400MalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	h := NewAuthHandler(auth.NewAuthService(users, tok))
	r := gin.New()
	r.POST("/auth/register", h.Register)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(`{"broken`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestAuth_Register400InvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	h := NewAuthHandler(auth.NewAuthService(users, tok))
	r := gin.New()
	r.POST("/auth/register", h.Register)
	body := map[string]any{
		"first_name": "A", "last_name": "B", "email": "r@test.com",
		"password": "password1", "role": "superuser",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d %s", w.Code, w.Body.String())
	}
}

func TestAuth_Login200And400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	hash, _ := crypto.HashPassword("secret123")
	u := &models.User{
		ID: "id-login", FirstName: "A", LastName: "B", Email: "log@example.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	h := NewAuthHandler(auth.NewAuthService(users, tok))
	r := gin.New()
	r.POST("/auth/login", h.Login)

	good := map[string]string{"email": "log@example.com", "password": "secret123"}
	bg, _ := json.Marshal(good)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bg))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login 200: got %d %s", w.Code, w.Body.String())
	}
	var ar models.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &ar); err != nil || ar.AccessToken == "" {
		t.Fatalf("tokens: %+v err %v", ar, err)
	}
	if ar.User.ID != "id-login" || ar.User.Email != "log@example.com" || ar.User.FirstName != "A" {
		t.Fatalf("user payload: %+v", ar.User)
	}

	bad := map[string]string{"email": "log@example.com"}
	bb, _ := json.Marshal(bad)
	req2 := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bb))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("login 400: got %d", w2.Code)
	}
}

func TestAuth_Refresh200And401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	hash, _ := crypto.HashPassword("pw123456")
	u := &models.User{
		ID: "id-ref", FirstName: "R", LastName: "F", Email: "ref@example.com",
		PasswordHash: hash, Role: models.UserRoleLibrarian, IsActive: true,
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	_, refresh, _, err := tok.GeneratePair(ctx, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	h := NewAuthHandler(auth.NewAuthService(users, tok))
	r := gin.New()
	r.POST("/auth/refresh", h.Refresh)

	body := map[string]string{"refresh_token": refresh}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("refresh 200: %d %s", w.Code, w.Body.String())
	}
	var lr models.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &lr); err != nil {
		t.Fatal(err)
	}
	if lr.AccessToken == "" || lr.User.ID != u.ID || lr.User.Email != u.Email {
		t.Fatalf("refresh payload: %+v", lr)
	}

	bad := map[string]string{"refresh_token": "not-a-jwt"}
	bb, _ := json.Marshal(bad)
	req2 := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(bb))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("refresh 401: got %d", w2.Code)
	}
}

func TestAuth_Me200And401AfterLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(strings.Repeat("k", 32), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	hash, _ := crypto.HashPassword("pw123456")
	u := &models.User{
		ID: "id-me", FirstName: "M", LastName: "E", Email: "me@example.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	authSvc := auth.NewAuthService(users, tok)
	h := NewAuthHandler(authSvc)
	access, _, _, err := tok.GeneratePair(ctx, u.ID)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/auth/me", AuthBearerMiddleware(tok), h.Me)
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("me 200: %d %s", w.Code, w.Body.String())
	}
	var me models.RegisterResponse
	if err := json.Unmarshal(w.Body.Bytes(), &me); err != nil || me.ID != u.ID || me.Email != u.Email {
		t.Fatalf("me body %+v err %v", me, err)
	}

	r2 := gin.New()
	r2.POST("/auth/logout", h.Logout)
	reqL := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewReader([]byte("{}")))
	reqL.Header.Set("Authorization", "Bearer "+access)
	wL := httptest.NewRecorder()
	r2.ServeHTTP(wL, reqL)
	if wL.Code != http.StatusNoContent {
		t.Fatalf("logout: %d", wL.Code)
	}

	r3 := gin.New()
	r3.GET("/auth/me", AuthBearerMiddleware(tok), h.Me)
	req3 := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req3.Header.Set("Authorization", "Bearer "+access)
	w3 := httptest.NewRecorder()
	r3.ServeHTTP(w3, req3)
	if w3.Code != http.StatusUnauthorized {
		t.Fatalf("me after revoke want 401 got %d", w3.Code)
	}
}
