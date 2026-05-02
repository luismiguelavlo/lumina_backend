package repositories

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestFineRepositoryInMemory_CreateList(t *testing.T) {
	r := NewFineRepositoryInMemory()
	sid := uuid.NewString()
	lid := uuid.NewString()
	f := &models.Fine{
		LoanID:    lid,
		StudentID: sid,
		Amount:    decimal.RequireFromString("12.50"),
		Status:    models.FineStatusPending,
	}
	if err := r.Create(context.Background(), f); err != nil {
		t.Fatal(err)
	}
	if f.ID == "" {
		t.Fatal("expected id")
	}
	items, total, err := r.List(context.Background(), FineListFilter{}, 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("list %+v total %d", items, total)
	}
	items2, total2, err := r.List(context.Background(), FineListFilter{StudentID: sid}, 20, 0)
	if err != nil || total2 != 1 {
		t.Fatalf("filter student %+v", items2)
	}
	_, total3, err := r.List(context.Background(), FineListFilter{StudentID: uuid.NewString()}, 20, 0)
	if err != nil || total3 != 0 {
		t.Fatalf("want0 got %d", total3)
	}
}

func TestFineRepositoryInMemory_UpdateStatusPaidWaived(t *testing.T) {
	r := NewFineRepositoryInMemory()
	f := &models.Fine{
		ID:        "fine-1",
		LoanID:    "loan-1",
		StudentID: "stu-1",
		Amount:    decimal.RequireFromString("5"),
		Status:    models.FineStatusPending,
	}
	if err := r.Create(context.Background(), f); err != nil {
		t.Fatal(err)
	}
	id := f.ID
	now := time.Now().UTC()
	n, err := r.UpdateStatus(context.Background(), id, models.FineStatusPaid, &now)
	if err != nil || n != 1 {
		t.Fatalf("paid n=%d err=%v", n, err)
	}
	got, _ := r.GetByID(context.Background(), id)
	if got.Status != models.FineStatusPaid || got.PaidAt == nil {
		t.Fatalf("state %+v", got)
	}
	n2, err := r.UpdateStatus(context.Background(), id, models.FineStatusWaived, nil)
	if err != nil || n2 != 0 {
		t.Fatalf("second update should affect 0 rows, n=%d", n2)
	}
}

func TestFineRepositoryInMemory_CountPendingByLoanID(t *testing.T) {
	r := NewFineRepositoryInMemory()
	lid := uuid.NewString()
	_ = r.Create(context.Background(), &models.Fine{
		LoanID:    lid,
		StudentID: uuid.NewString(),
		Amount:    decimal.NewFromInt(1),
		Status:    models.FineStatusPending,
	})
	n, err := r.CountPendingByLoanID(context.Background(), lid)
	if err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
}
