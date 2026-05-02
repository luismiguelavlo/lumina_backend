package repositories

import (
	"context"
	"sync"

	"library_back/internal/models"
)

// NewUserRepositoryInMemory returns an in-memory UserRepository (tests).
func NewUserRepositoryInMemory() UserRepository {
	return &userRepositoryInMemory{
		byEmail: make(map[string]*models.User),
		byID:    make(map[string]*models.User),
	}
}

type userRepositoryInMemory struct {
	mu      sync.RWMutex
	byEmail map[string]*models.User
	byID    map[string]*models.User
}

func (r *userRepositoryInMemory) Create(ctx context.Context, u *models.User) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byEmail[u.Email]; ok {
		return ErrDuplicateEmail
	}
	r.byEmail[u.Email] = u
	r.byID[u.ID] = u
	return nil
}

func (r *userRepositoryInMemory) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *userRepositoryInMemory) FindByID(ctx context.Context, id string) (*models.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}
