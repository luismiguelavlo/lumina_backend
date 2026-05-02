package repositories

import (
	"context"

	"library_back/internal/models"

	"gorm.io/gorm"
)

type activityLogRepositoryGorm struct {
	db *gorm.DB
}

// NewActivityLogRepositoryGorm builds an ActivityLogRepository backed by GORM.
func NewActivityLogRepositoryGorm(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepositoryGorm{db: db}
}

func (r *activityLogRepositoryGorm) Create(ctx context.Context, entry *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(entry).Error
}
