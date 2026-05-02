package jobs

import (
	"context"
	"time"

	"library_back/internal/repositories"
)

// BlacklistCleanup removes old rows from revoked_tokens.
type BlacklistCleanup struct {
	Repo repositories.RevokedTokenRepository
}

// Run deletes entries with revoked_at older than minAge.
func (b *BlacklistCleanup) Run(ctx context.Context) (int64, error) {
	return b.Repo.DeleteOlderThan(ctx, 24*time.Hour)
}
