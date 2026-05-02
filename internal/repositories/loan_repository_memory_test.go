package repositories

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestLoanRepositoryInMemory_CreateGetByID(t *testing.T) {
	ctx := context.Background()
	r := NewLoanRepositoryInMemory()
	id := uuid.NewString()
	l := &models.Loan{
		ID: id, BookID: "b1", StudentID: "s1",
		BorrowedAt: time.Now().UTC(),
		DueDate:    time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		Status:     models.LoanStatusActive,
	}
	if err := r.Create(ctx, l); err != nil {
		t.Fatal(err)
	}
	got, err := r.GetByID(ctx, id)
	if err != nil || got == nil || got.ID != id || got.Status != models.LoanStatusActive {
		t.Fatalf("GetByID: %+v %v", got, err)
	}
	miss, err := r.GetByID(ctx, uuid.NewString())
	if err != nil || miss != nil {
		t.Fatalf("miss: %+v", miss)
	}
}

func TestLoanRepositoryInMemory_ListFilterAndOrder(t *testing.T) {
	ctx := context.Background()
	r := NewLoanRepositoryInMemory()
	r.RegisterBookForJoin("b1", "Book", "ISBN")
	r.RegisterStudentForJoin("s1", "A", "B")

	d1 := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	_ = r.Create(ctx, &models.Loan{ID: "L1", BookID: "b1", StudentID: "s1", DueDate: d1, Status: models.LoanStatusActive, BorrowedAt: time.Now().UTC()})
	_ = r.Create(ctx, &models.Loan{ID: "L2", BookID: "b1", StudentID: "s1", DueDate: d2, Status: models.LoanStatusActive, BorrowedAt: time.Now().UTC()})

	rows, total, err := r.List(ctx, LoanListFilter{Status: "active"}, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Fatalf("total %d", total)
	}
	if len(rows) != 2 || rows[0].LoanID != "L2" || rows[1].LoanID != "L1" {
		t.Fatalf("order want L2,L1 got %+v", rows)
	}

	rows, total, err = r.List(ctx, LoanListFilter{Status: "active", StudentID: "s1"}, 1, 1)
	if err != nil || total != 2 || len(rows) != 1 || rows[0].LoanID != "L1" {
		t.Fatalf("paginate: %+v total=%d", rows, total)
	}
}

func TestLoanRepositoryInMemory_Return(t *testing.T) {
	ctx := context.Background()
	r := NewLoanRepositoryInMemory()
	id := uuid.NewString()
	_ = r.Create(ctx, &models.Loan{
		ID: id, BookID: "b1", StudentID: "s1",
		DueDate:    time.Now().UTC().AddDate(0, 0, 1),
		Status:     models.LoanStatusActive,
		BorrowedAt: time.Now().UTC(),
	})
	retAt := time.Now().UTC()
	if err := r.Return(ctx, id, retAt); err != nil {
		t.Fatal(err)
	}
	got, _ := r.GetByID(ctx, id)
	if got == nil || got.Status != models.LoanStatusReturned || got.ReturnedAt == nil {
		t.Fatalf("after return: %+v", got)
	}
	if err := r.Return(ctx, id, retAt); err != gorm.ErrRecordNotFound {
		t.Fatalf("second return: %v", err)
	}
	if err := r.Return(ctx, uuid.NewString(), retAt); err != gorm.ErrRecordNotFound {
		t.Fatalf("missing id: %v", err)
	}
}

func TestLoanRepositoryInMemory_CountCheckedOutAndMarkOverdue(t *testing.T) {
	ctx := context.Background()
	r := NewLoanRepositoryInMemory()
	past := time.Now().UTC().AddDate(0, 0, -5)
	fut := time.Now().UTC().AddDate(0, 0, 5)
	_ = r.Create(ctx, &models.Loan{ID: "a", BookID: "b1", StudentID: "s1", DueDate: past, Status: models.LoanStatusActive, BorrowedAt: time.Now().UTC()})
	_ = r.Create(ctx, &models.Loan{ID: "b", BookID: "b1", StudentID: "s1", DueDate: fut, Status: models.LoanStatusActive, BorrowedAt: time.Now().UTC()})

	n, err := r.CountCheckedOutByBookID(ctx, "b1")
	if err != nil || n != 2 {
		t.Fatalf("count %d %v", n, err)
	}

	marked, err := r.MarkActiveLoansOverdue(ctx)
	if err != nil || marked != 1 {
		t.Fatalf("marked %d %v", marked, err)
	}
	got, _ := r.GetByID(ctx, "a")
	if got.Status != models.LoanStatusOverdue {
		t.Fatalf("status %s", got.Status)
	}
	n, err = r.CountCheckedOutByBookID(ctx, "b1")
	if err != nil || n != 2 {
		t.Fatalf("count after overdue %d", n)
	}
}
