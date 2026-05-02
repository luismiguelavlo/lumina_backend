package repositories

import (
	"context"
	"time"

	"library_back/internal/models"
)

// FineListFilter drives GET /api/fines.
type FineListFilter struct {
	StudentID string // optional UUID
	Status    string // optional: pending | paid | waived
}

// FineWithStudent is a fine row plus borrower display name from JOIN.
type FineWithStudent struct {
	Fine        models.Fine
	StudentName string
}

// FineRepository persists fines.
type FineRepository interface {
	Create(ctx context.Context, f *models.Fine) error
	GetByID(ctx context.Context, id string) (*models.Fine, error)
	List(ctx context.Context, filter FineListFilter, limit, offset int) ([]FineWithStudent, int64, error)
	UpdateStatus(ctx context.Context, id string, status models.FineStatus, paidAt *time.Time) (int64, error)
	CountPendingByLoanID(ctx context.Context, loanID string) (int64, error)
}
