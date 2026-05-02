package repositories

import (
	"context"
	"time"

	"library_back/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RevokedTokenRepository persists revoked JWT JTIs (blacklist).
type RevokedTokenRepository interface {
	Revoke(ctx context.Context, jti string, userID string, expiresAt time.Time) error
	RevokeMany(ctx context.Context, jtis []string, userID string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
	DeleteOlderThan(ctx context.Context, minAge time.Duration) (int64, error)
}

type revokedTokenRepositoryGorm struct {
	db *gorm.DB
}

// NewRevokedTokenRepositoryGorm returns a GORM-backed RevokedTokenRepository.
func NewRevokedTokenRepositoryGorm(db *gorm.DB) RevokedTokenRepository {
	return &revokedTokenRepositoryGorm{db: db}
}

func (r *revokedTokenRepositoryGorm) Revoke(ctx context.Context, jti string, userID string, expiresAt time.Time) error {
	row := models.RevokedToken{
		JTI:       jti,
		UserID:    userID,
		ExpiresAt: expiresAt,
		RevokedAt: time.Now().UTC(),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error
}

func (r *revokedTokenRepositoryGorm) RevokeMany(ctx context.Context, jtis []string, userID string, expiresAt time.Time) error {
	for _, jti := range jtis {
		if jti == "" {
			continue
		}
		if err := r.Revoke(ctx, jti, userID, expiresAt); err != nil {
			return err
		}
	}
	return nil
}

func (r *revokedTokenRepositoryGorm) IsRevoked(ctx context.Context, jti string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.RevokedToken{}).Where("jti = ?", jti).Count(&n).Error
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *revokedTokenRepositoryGorm) DeleteOlderThan(ctx context.Context, minAge time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-minAge)
	res := r.db.WithContext(ctx).Where("revoked_at < ?", cutoff).Delete(&models.RevokedToken{})
	return res.RowsAffected, res.Error
}
