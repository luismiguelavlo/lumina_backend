package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRevokedTokenRepositoryInMemory_RevokeAndIsRevoked(t *testing.T) {
	ctx := context.Background()
	r := NewRevokedTokenRepositoryInMemory()
	jti := uuid.NewString()
	ok, err := r.IsRevoked(ctx, jti)
	if err != nil || ok {
		t.Fatalf("before revoke: ok=%v err=%v", ok, err)
	}
	exp := time.Now().UTC().Add(time.Hour)
	if err := r.Revoke(ctx, jti, "user-1", exp); err != nil {
		t.Fatal(err)
	}
	ok, err = r.IsRevoked(ctx, jti)
	if err != nil || !ok {
		t.Fatalf("after revoke: ok=%v err=%v", ok, err)
	}
}

func TestRevokedTokenRepositoryInMemory_RevokeMany(t *testing.T) {
	ctx := context.Background()
	r := NewRevokedTokenRepositoryInMemory()
	j1, j2 := uuid.NewString(), uuid.NewString()
	exp := time.Now().UTC().Add(time.Hour)
	if err := r.RevokeMany(ctx, []string{j1, j2}, "u", exp); err != nil {
		t.Fatal(err)
	}
	for _, j := range []string{j1, j2} {
		ok, err := r.IsRevoked(ctx, j)
		if err != nil || !ok {
			t.Fatalf("jti %s ok=%v err=%v", j, ok, err)
		}
	}
}

func TestRevokedTokenRepositoryInMemory_DeleteOlderThan(t *testing.T) {
	ctx := context.Background()
	r := NewRevokedTokenRepositoryInMemory()
	j := uuid.NewString()
	exp := time.Now().UTC().Add(time.Hour)
	if err := r.Revoke(ctx, j, "u", exp); err != nil {
		t.Fatal(err)
	}
	n, err := r.DeleteOlderThan(ctx, 24*time.Hour)
	if err != nil || n != 0 {
		t.Fatalf("recent row should remain: n=%d err=%v", n, err)
	}
	n, err = r.DeleteOlderThan(ctx, -time.Millisecond)
	if err != nil || n != 1 {
		t.Fatalf("want delete 1 with negative minAge trick, got n=%d err=%v", n, err)
	}
	ok, _ := r.IsRevoked(ctx, j)
	if ok {
		t.Fatal("jti should be gone")
	}
}
