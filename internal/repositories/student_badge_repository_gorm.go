package repositories

import (
	"context"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// StudentBadgeRepository manages student_badges rows.
type StudentBadgeRepository interface {
	Exists(ctx context.Context, studentID, badgeID string) (bool, error)
	Create(ctx context.Context, sb *models.StudentBadge) error
	// DeleteByStudentAndBadge removes the earned badge row; returns rows affected (0 if none).
	DeleteByStudentAndBadge(ctx context.Context, studentID, badgeID string) (int64, error)
}

type studentBadgeRepositoryGorm struct {
	db *gorm.DB
}

// NewStudentBadgeRepositoryGorm returns a GORM-backed StudentBadgeRepository.
func NewStudentBadgeRepositoryGorm(db *gorm.DB) StudentBadgeRepository {
	return &studentBadgeRepositoryGorm{db: db}
}

func (r *studentBadgeRepositoryGorm) Exists(ctx context.Context, studentID, badgeID string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.StudentBadge{}).
		Where("student_id = ? AND badge_id = ?", studentID, badgeID).
		Count(&n).Error
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *studentBadgeRepositoryGorm) Create(ctx context.Context, sb *models.StudentBadge) error {
	return r.db.WithContext(ctx).Create(sb).Error
}

func (r *studentBadgeRepositoryGorm) DeleteByStudentAndBadge(ctx context.Context, studentID, badgeID string) (int64, error) {
	res := r.db.WithContext(ctx).Where("student_id = ? AND badge_id = ?", studentID, badgeID).Delete(&models.StudentBadge{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
