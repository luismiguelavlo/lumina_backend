package repositories

import (
	"context"
	"errors"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// BadgeRepository reads and persists badge definitions.
type BadgeRepository interface {
	GetByID(ctx context.Context, id string) (*models.Badge, error)
	GetBySlug(ctx context.Context, slug string) (*models.Badge, error)
	ListAll(ctx context.Context) ([]models.Badge, error)
	Create(ctx context.Context, b *models.Badge) error
	Save(ctx context.Context, b *models.Badge) error
	Delete(ctx context.Context, id string) error
	CountStudentBadgesByBadgeID(ctx context.Context, badgeID string) (int64, error)
}

type badgeRepositoryGorm struct {
	db *gorm.DB
}

// NewBadgeRepositoryGorm returns a GORM-backed BadgeRepository.
func NewBadgeRepositoryGorm(db *gorm.DB) BadgeRepository {
	return &badgeRepositoryGorm{db: db}
}

func (r *badgeRepositoryGorm) GetBySlug(ctx context.Context, slug string) (*models.Badge, error) {
	var b models.Badge
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *badgeRepositoryGorm) GetByID(ctx context.Context, id string) (*models.Badge, error) {
	var b models.Badge
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *badgeRepositoryGorm) ListAll(ctx context.Context) ([]models.Badge, error) {
	var rows []models.Badge
	err := r.db.WithContext(ctx).Model(&models.Badge{}).Order("slug ASC").Find(&rows).Error
	return rows, err
}

func (r *badgeRepositoryGorm) Create(ctx context.Context, b *models.Badge) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *badgeRepositoryGorm) Save(ctx context.Context, b *models.Badge) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *badgeRepositoryGorm) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Badge{}).Error
}

func (r *badgeRepositoryGorm) CountStudentBadgesByBadgeID(ctx context.Context, badgeID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.StudentBadge{}).Where("badge_id = ?", badgeID).Count(&n).Error
	return n, err
}
