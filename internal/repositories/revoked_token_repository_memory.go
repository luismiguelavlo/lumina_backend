package repositories

import (
	"context"
	"sync"
	"time"

	"library_back/internal/models"
)

// NewRevokedTokenRepositoryInMemory returns an in-memory RevokedTokenRepository (tests).
func NewRevokedTokenRepositoryInMemory() RevokedTokenRepository {
	return &revokedTokenRepositoryMemory{jtis: make(map[string]models.RevokedToken)}
}

type revokedTokenRepositoryMemory struct {
	mu   sync.RWMutex
	jtis map[string]models.RevokedToken
}

func (r *revokedTokenRepositoryMemory) Revoke(ctx context.Context, jti string, userID string, expiresAt time.Time) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.jtis[jti]; ok {
		return nil
	}
	r.jtis[jti] = models.RevokedToken{
		JTI:       jti,
		UserID:    userID,
		ExpiresAt: expiresAt,
		RevokedAt: time.Now().UTC(),
	}
	return nil
}

func (r *revokedTokenRepositoryMemory) RevokeMany(ctx context.Context, jtis []string, userID string, expiresAt time.Time) error {
	for _, jti := range jtis {
		if err := r.Revoke(ctx, jti, userID, expiresAt); err != nil {
			return err
		}
	}
	return nil
}

func (r *revokedTokenRepositoryMemory) IsRevoked(ctx context.Context, jti string) (bool, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.jtis[jti]
	return ok, nil
}

func (r *revokedTokenRepositoryMemory) DeleteOlderThan(ctx context.Context, minAge time.Duration) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := time.Now().UTC().Add(-minAge)
	var n int64
	for jti, row := range r.jtis {
		if row.RevokedAt.Before(cutoff) {
			delete(r.jtis, jti)
			n++
		}
	}
	return n, nil
}
