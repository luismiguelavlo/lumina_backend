package badgedef

import (
	"context"
	"errors"
	"strings"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service manages badge definitions (staff CRUD).
type Service struct {
	badges repositories.BadgeRepository
}

// NewService constructs Service.
func NewService(badges repositories.BadgeRepository) *Service {
	return &Service{badges: badges}
}

// List returns all badge definitions ordered by slug.
func (s *Service) List(ctx context.Context) ([]models.BadgeDefinitionResponse, error) {
	rows, err := s.badges.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]models.BadgeDefinitionResponse, 0, len(rows))
	for i := range rows {
		out = append(out, badgeToDefResp(&rows[i]))
	}
	return out, nil
}

// GetByID returns one definition.
func (s *Service) GetByID(ctx context.Context, id string) (*models.BadgeDefinitionResponse, error) {
	b, err := s.badges.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound
	}
	r := badgeToDefResp(b)
	return &r, nil
}

// Create inserts a new badge definition.
func (s *Service) Create(ctx context.Context, req models.CreateBadgeDefinitionRequest) (*models.BadgeDefinitionResponse, error) {
	slug := strings.TrimSpace(req.Slug)
	existing, err := s.badges.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSlugTaken
	}
	b := &models.Badge{
		ID:          uuid.NewString(),
		Slug:        slug,
		Name:        strings.TrimSpace(req.Name),
		IconURL:     req.IconURL,
		Description: req.Description,
		Criteria:    req.Criteria,
	}
	if err := s.badges.Create(ctx, b); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	refreshed, err := s.badges.GetByID(ctx, b.ID)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrNotFound
	}
	r := badgeToDefResp(refreshed)
	return &r, nil
}

// Patch updates fields on an existing badge.
func (s *Service) Patch(ctx context.Context, id string, req models.PatchBadgeDefinitionRequest) (*models.BadgeDefinitionResponse, error) {
	b, err := s.badges.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound
	}
	if req.Slug != nil {
		ns := strings.TrimSpace(*req.Slug)
		if ns != "" && ns != b.Slug {
			other, err := s.badges.GetBySlug(ctx, ns)
			if err != nil {
				return nil, err
			}
			if other != nil && other.ID != b.ID {
				return nil, ErrSlugTaken
			}
			b.Slug = ns
		}
	}
	if req.Name != nil {
		b.Name = strings.TrimSpace(*req.Name)
	}
	if req.IconURL != nil {
		b.IconURL = req.IconURL
	}
	if req.Description != nil {
		b.Description = req.Description
	}
	if req.Criteria != nil {
		b.Criteria = req.Criteria
	}
	if err := s.badges.Save(ctx, b); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	refreshed, err := s.badges.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrNotFound
	}
	r := badgeToDefResp(refreshed)
	return &r, nil
}

// Delete removes a badge definition if no student has earned it.
func (s *Service) Delete(ctx context.Context, id string) error {
	b, err := s.badges.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if b == nil {
		return ErrNotFound
	}
	n, err := s.badges.CountStudentBadgesByBadgeID(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrInUse
	}
	return s.badges.Delete(ctx, id)
}

func badgeToDefResp(b *models.Badge) models.BadgeDefinitionResponse {
	return models.BadgeDefinitionResponse{
		ID:          b.ID,
		Slug:        b.Slug,
		Name:        b.Name,
		IconURL:     b.IconURL,
		Description: b.Description,
		Criteria:    b.Criteria,
		CreatedAt:   b.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}
