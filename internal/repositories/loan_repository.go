package repositories

import (
	"context"
	"time"

	"library_back/internal/models"
)

// LoanListFilter drives GET /api/loans.
type LoanListFilter struct {
	Status    string // active | returned | overdue (empty → active)
	StudentID string // optional UUID
}

// LoanJoinRow is a loan row with book and student fields for API mapping.
type LoanJoinRow struct {
	LoanID        string
	BookID        string
	BookTitle     string
	BookISBN      string
	StudentID     string
	BorrowerFirst string
	BorrowerLast  string
	DueDate       time.Time
	Status        string
	BorrowedAt    time.Time
	ReturnedAt    *time.Time
}

// LoanRepository persists loans.
type LoanRepository interface {
	Create(ctx context.Context, loan *models.Loan) error
	GetByID(ctx context.Context, id string) (*models.Loan, error)
	GetJoinByID(ctx context.Context, id string) (*LoanJoinRow, error)
	List(ctx context.Context, filter LoanListFilter, limit, offset int) ([]LoanJoinRow, int64, error)
	Return(ctx context.Context, id string, returnedAt time.Time) error
	CountCheckedOutByBookID(ctx context.Context, bookID string) (int64, error)
	// MarkActiveLoansOverdue sets status = overdue for active loans past due_date (calendar).
	MarkActiveLoansOverdue(ctx context.Context) (int64, error)
}
