package jobs

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"
)

func TestLoanOverdueSync_Run(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	past := time.Now().UTC().AddDate(0, 0, -2)
	_ = mem.Create(ctx, &models.Loan{
		ID: "loan-1", BookID: "b", StudentID: "s",
		DueDate:    past,
		Status:     models.LoanStatusActive,
		BorrowedAt: time.Now().UTC(),
	})
	j := &LoanOverdueSync{Loans: mem}
	n, err := j.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("rows %d", n)
	}
	got, _ := mem.GetByID(ctx, "loan-1")
	if got.Status != models.LoanStatusOverdue {
		t.Fatalf("status %s", got.Status)
	}
}
