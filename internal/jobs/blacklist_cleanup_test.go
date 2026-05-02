package jobs

import (
	"context"
	"testing"
	"time"

	"library_back/internal/repositories"
)

func TestBlacklistCleanup_Run(t *testing.T) {
	repo := repositories.NewRevokedTokenRepositoryInMemory()
	_, err := (&BlacklistCleanup{Repo: repo}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

type trackRevokedRepo struct {
	repositories.RevokedTokenRepository
	lastDur time.Duration
}

func (t *trackRevokedRepo) DeleteOlderThan(ctx context.Context, minAge time.Duration) (int64, error) {
	t.lastDur = minAge
	return t.RevokedTokenRepository.DeleteOlderThan(ctx, minAge)
}

func TestBlacklistCleanup_UsesTwentyFourHourWindow(t *testing.T) {
	inner := repositories.NewRevokedTokenRepositoryInMemory()
	tr := &trackRevokedRepo{RevokedTokenRepository: inner}
	_, err := (&BlacklistCleanup{Repo: tr}).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tr.lastDur != 24*time.Hour {
		t.Fatalf("want DeleteOlderThan(24h), got %v", tr.lastDur)
	}
}
