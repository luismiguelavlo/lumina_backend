package repositories

import (
	"context"
	"time"

	"library_back/internal/models"
)

// SanctionListRow is one active sanction with student fields (JOIN).
type SanctionListRow struct {
	SanctionID    string
	StudentID     string
	StudentIDCode string
	FirstName     string
	LastName      string
	Email         *string
	Reason        string
	Status        string
	AppliedAt     time.Time
	AppliedByID   *string
}

// SanctionListFilter represents query filters for active sanctions.
type SanctionListFilter struct {
	StudentID string
}

// SanctionRepository persists sanctions.
type SanctionRepository interface {
	Create(ctx context.Context, s *models.Sanction) error
	GetByID(ctx context.Context, id string) (*models.Sanction, error)
	ListActive(ctx context.Context, filter SanctionListFilter, limit, offset int) ([]SanctionListRow, int64, error)
	Lift(ctx context.Context, id string, liftedAt time.Time, liftedBy string) (int64, error)
	// CountActiveByStudentID counts sanctions with status active for the student.
	CountActiveByStudentID(ctx context.Context, studentID string) (int64, error)
}
