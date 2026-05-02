package auth

import (
	"context"
	"strings"
	"testing"

	"library_back/internal/models"
	"library_back/internal/pkg/crypto"
	"library_back/internal/repositories"
	"library_back/internal/services"
)

func testJWTSecret() string {
	return strings.Repeat("x", 32)
}

func TestAuthService_RegisterAndDuplicate(t *testing.T) {
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(testJWTSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(users, tok)
	ctx := context.Background()

	u, err := svc.Register(ctx, models.RegisterRequest{
		FirstName: "A",
		LastName:  "B",
		Email:     "a@b.com",
		Password:  "password1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if u == nil || u.IsActive {
		t.Fatalf("user: %+v", u)
	}
	_, err = svc.Register(ctx, models.RegisterRequest{
		FirstName: "C",
		LastName:  "D",
		Email:     "a@b.com",
		Password:  "password1",
	})
	if err != ErrEmailExists {
		t.Fatalf("want ErrEmailExists got %v", err)
	}
}

func TestAuthService_LoginActiveOnly(t *testing.T) {
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(testJWTSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(users, tok)
	ctx := context.Background()

	hash, err := crypto.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	inactive := &models.User{
		ID: "id1", FirstName: "X", LastName: "Y", Email: "x@y.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: false,
	}
	if err := users.Create(ctx, inactive); err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err = svc.Login(ctx, "x@y.com", "secret123")
	if err != ErrInvalidCredentials {
		t.Fatalf("inactive: want ErrInvalidCredentials got %v", err)
	}

	active := &models.User{
		ID: "id2", FirstName: "P", LastName: "Q", Email: "p@q.com",
		PasswordHash: hash, Role: models.UserRoleLibrarian, IsActive: true,
	}
	if err := users.Create(ctx, active); err != nil {
		t.Fatal(err)
	}
	lu, a, r, exp, err := svc.Login(ctx, "p@q.com", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	if lu == nil || lu.ID != active.ID || lu.Email != "p@q.com" {
		t.Fatalf("user: %+v", lu)
	}
	if a == "" || r == "" || exp != 36000 {
		t.Fatal("expected tokens")
	}
}

func TestAuthService_RefreshInactive(t *testing.T) {
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(testJWTSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(users, tok)
	ctx := context.Background()

	hash, _ := crypto.HashPassword("pw")
	u := &models.User{
		ID: "id3", FirstName: "I", LastName: "N", Email: "i@n.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: false,
	}
	_ = users.Create(ctx, u)
	_, refresh, _, err := tok.GeneratePair(ctx, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err = svc.Refresh(ctx, refresh)
	if err != ErrInvalidCredentials {
		t.Fatalf("want ErrInvalidCredentials got %v", err)
	}
}

func TestAuthService_LoginWrongPassword(t *testing.T) {
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(testJWTSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(users, tok)
	ctx := context.Background()
	hash, _ := crypto.HashPassword("correctPass9")
	u := &models.User{
		ID: "wp1", FirstName: "W", LastName: "P", Email: "wp@test.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err = svc.Login(ctx, "wp@test.com", "wrongPass9")
	if err != ErrInvalidCredentials {
		t.Fatalf("want ErrInvalidCredentials got %v", err)
	}
}

func TestAuthService_RefreshRevokedRefreshToken(t *testing.T) {
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(testJWTSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(users, tok)
	ctx := context.Background()
	hash, _ := crypto.HashPassword("pw")
	active := &models.User{
		ID: "rfv1", FirstName: "R", LastName: "V", Email: "rfv@test.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, active); err != nil {
		t.Fatal(err)
	}
	_, refresh, _, err := tok.GeneratePair(ctx, active.ID)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := tok.ParseRefresh(ctx, refresh)
	if err != nil {
		t.Fatal(err)
	}
	if err := tok.Revoke(ctx, pr.UserID, pr.JTI, pr.ExpiresAt); err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err = svc.Refresh(ctx, refresh)
	if err != ErrTokenRevoked {
		t.Fatalf("want ErrTokenRevoked got %v", err)
	}
}

func TestAuthService_GetProfile(t *testing.T) {
	users := repositories.NewUserRepositoryInMemory()
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	tok, err := services.NewTokenService(testJWTSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(users, tok)
	ctx := context.Background()

	_, err = svc.GetProfile(ctx, "missing-id")
	if err != ErrUserNotFound {
		t.Fatalf("missing: want ErrUserNotFound got %v", err)
	}

	hash, _ := crypto.HashPassword("pw")
	u := &models.User{
		ID: "gp1", FirstName: "G", LastName: "P", Email: "gp@test.com",
		PasswordHash: hash, Role: models.UserRoleAdmin, IsActive: true,
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	got, err := svc.GetProfile(ctx, u.ID)
	if err != nil || got == nil || got.Email != u.Email {
		t.Fatalf("profile: %+v err %v", got, err)
	}
}
