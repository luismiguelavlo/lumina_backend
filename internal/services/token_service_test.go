package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"library_back/internal/pkg/jwths256"
	"library_back/internal/pkg/tokenerr"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

func testSecret() string {
	return strings.Repeat("x", 32)
}

func TestTokenService_GeneratePairAndParse(t *testing.T) {
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	svc, err := NewTokenService(testSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	access, refresh, exp, err := svc.GeneratePair(ctx, "user-1")
	if err != nil {
		t.Fatal(err)
	}
	if exp != 36000 {
		t.Fatalf("expiresIn want 36000 got %d", exp)
	}
	pa, err := svc.ParseAccess(ctx, access)
	if err != nil {
		t.Fatal(err)
	}
	if pa.UserID != "user-1" || pa.JTI == "" {
		t.Fatalf("access claims: %+v", pa)
	}
	pr, err := svc.ParseRefresh(ctx, refresh)
	if err != nil {
		t.Fatal(err)
	}
	if pr.UserID != "user-1" {
		t.Fatalf("refresh sub: %s", pr.UserID)
	}
}

func TestTokenService_RevokedRejected(t *testing.T) {
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	svc, err := NewTokenService(testSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	access, _, _, err := svc.GeneratePair(ctx, "u2")
	if err != nil {
		t.Fatal(err)
	}
	pa, err := svc.ParseAccess(ctx, access)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Revoke(ctx, pa.UserID, pa.JTI, pa.ExpiresAt); err != nil {
		t.Fatal(err)
	}
	_, err = svc.ParseAccess(ctx, access)
	if err != tokenerr.ErrTokenRevoked {
		t.Fatalf("want ErrTokenRevoked got %v", err)
	}
}

func TestTokenService_ParseAccessExpiredAndGarbage(t *testing.T) {
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	svc, err := NewTokenService(testSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Unix()
	expired := jwths256.Claims{
		Typ: tokenTypeAccess,
		Sub: "sub-exp",
		Jti: uuid.NewString(),
		Exp: now - 120,
		Iat: now - 360,
		Nbf: now - 360,
	}
	tok, err := jwths256.Sign([]byte(testSecret()), expired)
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.ParseAccess(ctx, tok)
	if err != tokenerr.ErrTokenInvalid {
		t.Fatalf("expired: want ErrTokenInvalid got %v", err)
	}
	_, err = svc.ParseAccess(ctx, "not-a-jwt")
	if err != tokenerr.ErrTokenInvalid {
		t.Fatalf("garbage: want ErrTokenInvalid got %v", err)
	}
}

func TestTokenService_GeneratePairDistinctJTI(t *testing.T) {
	rev := repositories.NewRevokedTokenRepositoryInMemory()
	svc, err := NewTokenService(testSecret(), rev)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	access, refresh, _, err := svc.GeneratePair(ctx, "u-jti")
	if err != nil {
		t.Fatal(err)
	}
	pa, err := svc.ParseAccess(ctx, access)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := svc.ParseRefresh(ctx, refresh)
	if err != nil {
		t.Fatal(err)
	}
	if pa.JTI == "" || pr.JTI == "" || pa.JTI == pr.JTI {
		t.Fatalf("want distinct jti: %q %q", pa.JTI, pr.JTI)
	}
}
