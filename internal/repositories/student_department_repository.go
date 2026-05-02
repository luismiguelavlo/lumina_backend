package repositories

import (
	"context"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// StudentDepartmentRepository lists departments and validates FK for students.
type StudentDepartmentRepository interface {
	List(ctx context.Context) ([]models.Department, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

type studentDepartmentRepositoryGorm struct {
	db *gorm.DB
}

// NewStudentDepartmentRepositoryGorm returns a GORM-backed StudentDepartmentRepository.
func NewStudentDepartmentRepositoryGorm(db *gorm.DB) StudentDepartmentRepository {
	return &studentDepartmentRepositoryGorm{db: db}
}

func (r *studentDepartmentRepositoryGorm) List(ctx context.Context) ([]models.Department, error) {
	var rows []models.Department
	err := r.db.WithContext(ctx).Model(&models.Department{}).Order("name ASC").Find(&rows).Error
	return rows, err
}

func (r *studentDepartmentRepositoryGorm) ExistsByID(ctx context.Context, id string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Department{}).Where("id = ?", id).Count(&n).Error
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
