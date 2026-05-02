package repositories

import (
	"context"
	"errors"
	"strings"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// UserRepository persists staff users (PostgreSQL via GORM).
type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

type userRepositoryGorm struct {
	db *gorm.DB
}

// NewUserRepositoryGorm returns a GORM-backed UserRepository.
func NewUserRepositoryGorm(db *gorm.DB) UserRepository {
	return &userRepositoryGorm{db: db}
}

func (r *userRepositoryGorm) Create(ctx context.Context, u *models.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return ErrDuplicateEmail
	}
	return err
}

func (r *userRepositoryGorm) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepositoryGorm) FindByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
