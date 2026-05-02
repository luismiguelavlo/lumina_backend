package repositories

import (
	"context"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// GenreCatalogRepository lists genres and checks existence (catalog module).
type GenreCatalogRepository interface {
	List(ctx context.Context) ([]models.Genre, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

type genreCatalogRepositoryGorm struct {
	db *gorm.DB
}

// NewGenreCatalogRepositoryGorm returns a GORM-backed GenreCatalogRepository.
func NewGenreCatalogRepositoryGorm(db *gorm.DB) GenreCatalogRepository {
	return &genreCatalogRepositoryGorm{db: db}
}

func (r *genreCatalogRepositoryGorm) List(ctx context.Context) ([]models.Genre, error) {
	var rows []models.Genre
	err := r.db.WithContext(ctx).Model(&models.Genre{}).Order("name ASC").Find(&rows).Error
	return rows, err
}

func (r *genreCatalogRepositoryGorm) ExistsByID(ctx context.Context, id string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Genre{}).Where("id = ?", id).Count(&n).Error
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
