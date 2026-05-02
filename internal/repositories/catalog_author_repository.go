package repositories

import (
	"context"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// AuthorCatalogRepository lists authors, creates rows, and checks existence (catalog module).
type AuthorCatalogRepository interface {
	List(ctx context.Context) ([]models.Author, error)
	Create(ctx context.Context, a *models.Author) error
	ExistsByID(ctx context.Context, id string) (bool, error)
}

type authorCatalogRepositoryGorm struct {
	db *gorm.DB
}

// NewAuthorCatalogRepositoryGorm returns a GORM-backed AuthorCatalogRepository.
func NewAuthorCatalogRepositoryGorm(db *gorm.DB) AuthorCatalogRepository {
	return &authorCatalogRepositoryGorm{db: db}
}

func (r *authorCatalogRepositoryGorm) List(ctx context.Context) ([]models.Author, error) {
	var rows []models.Author
	err := r.db.WithContext(ctx).Model(&models.Author{}).Order("name ASC").Find(&rows).Error
	return rows, err
}

func (r *authorCatalogRepositoryGorm) Create(ctx context.Context, a *models.Author) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *authorCatalogRepositoryGorm) ExistsByID(ctx context.Context, id string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Author{}).Where("id = ?", id).Count(&n).Error
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
