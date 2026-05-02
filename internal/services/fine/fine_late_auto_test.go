package fine

import (
	"context"
	"os"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

func TestEnsureLateReturnFine_Disabled(t *testing.T) {
	_ = os.Setenv("AUTO_LATE_RETURN_FINE_ENABLED", "0")
	defer os.Unsetenv("AUTO_LATE_RETURN_FINE_ENABLED")

	mem := repositories.NewFineRepositoryInMemory()
	loanID := uuid.NewString()
	sid := uuid.NewString()
	due := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	ret := time.Date(2026, 1, 12, 12, 0, 0, 0, time.UTC)
	loan := &models.Loan{
		ID: loanID, StudentID: sid, Status: models.LoanStatusReturned,
		DueDate: due, ReturnedAt: &ret,
	}
	s := NewService(mem, &fakeLoanReader{loan: loan}, nil)
	if err := s.EnsureLateReturnFine(context.Background(), loan); err != nil {
		t.Fatal(err)
	}
	n, _ := mem.CountPendingByLoanID(context.Background(), loanID)
	if n != 0 {
		t.Fatalf("expected no fine, got %d", n)
	}
}

func TestEnsureLateReturnFine_CreatesWhenLate(t *testing.T) {
	_ = os.Unsetenv("AUTO_LATE_RETURN_FINE_ENABLED")
	_ = os.Unsetenv("AUTO_LATE_RETURN_FINE_AMOUNT")

	mem := repositories.NewFineRepositoryInMemory()
	loanID := uuid.NewString()
	sid := uuid.NewString()
	due := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	ret := time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC)
	loan := &models.Loan{
		ID: loanID, StudentID: sid, Status: models.LoanStatusReturned,
		DueDate: due, ReturnedAt: &ret,
	}
	s := NewService(mem, &fakeLoanReader{loan: loan}, nil)
	if err := s.EnsureLateReturnFine(context.Background(), loan); err != nil {
		t.Fatal(err)
	}
	n, _ := mem.CountPendingByLoanID(context.Background(), loanID)
	if n != 1 {
		t.Fatalf("want 1 pending fine got %d", n)
	}
}

func TestEnsureLateReturnFine_SkipsOnTime(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	loanID := uuid.NewString()
	sid := uuid.NewString()
	due := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	ret := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	loan := &models.Loan{
		ID: loanID, StudentID: sid, Status: models.LoanStatusReturned,
		DueDate: due, ReturnedAt: &ret,
	}
	s := NewService(mem, &fakeLoanReader{loan: loan}, nil)
	if err := s.EnsureLateReturnFine(context.Background(), loan); err != nil {
		t.Fatal(err)
	}
	n, _ := mem.CountPendingByLoanID(context.Background(), loanID)
	if n != 0 {
		t.Fatalf("on-time return should not create fine, got %d", n)
	}
}
