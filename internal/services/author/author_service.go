package author

import (
	"context"
	"strings"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
)

type authorWriter interface {
	Create(ctx context.Context, a *models.Author) error
}

// Service creates catalog authors (staff API).
type Service struct {
	repo authorWriter
}

// NewService constructs AuthorService.
func NewService(repo authorWriter) *Service {
	return &Service{repo: repo}
}

// Create persists a new author and returns the public payload.
func (s *Service) Create(ctx context.Context, req models.CreateAuthorRequest) (*models.AuthorCreatedResponse, error) {
	name := strings.TrimSpace(req.Name)
	a := &models.Author{
		ID:   uuid.NewString(),
		Name: name,
		Bio:  req.Bio,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return &models.AuthorCreatedResponse{
		ID:        a.ID,
		Name:      a.Name,
		Bio:       a.Bio,
		CreatedAt: a.CreatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}
