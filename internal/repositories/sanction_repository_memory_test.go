package repositories

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
)

func TestSanctionRepositoryInMemory_CountActiveByStudentID(t *testing.T) {
	r := NewSanctionRepositoryInMemory()
	sid := uuid.NewString()
	ab := uuid.NewString()
	_ = r.Create(context.Background(), &models.Sanction{
		StudentID: sid,
		Reason:    "r1",
		Status:    models.SanctionStatusActive,
		AppliedBy: &ab,
	})
	_ = r.Create(context.Background(), &models.Sanction{
		StudentID: sid,
		Reason:    "r2",
		Status:    models.SanctionStatusActive,
		AppliedBy: &ab,
	})
	n, err := r.CountActiveByStudentID(context.Background(), sid)
	if err != nil || n != 2 {
		t.Fatalf("n=%d err=%v", n, err)
	}
}

func TestSanctionRepositoryInMemory_CreateListLift(t *testing.T) {
	r := NewSanctionRepositoryInMemory()
	sid := uuid.NewString()
	ab := uuid.NewString()
	s := &models.Sanction{
		StudentID: sid,
		Reason:    "Test reason",
		Status:    models.SanctionStatusActive,
		AppliedBy: &ab,
	}
	if err := r.Create(context.Background(), s); err != nil {
		t.Fatal(err)
	}
	if s.ID == "" {
		t.Fatal("expected id")
	}
	rows, total, err := r.ListActive(context.Background(), SanctionListFilter{}, 20, 0)
	if err != nil || total != 1 || len(rows) != 1 {
		t.Fatalf("list err=%v total=%d n=%d", err, total, len(rows))
	}
	n, err := r.Lift(context.Background(), s.ID, time.Now().UTC(), uuid.NewString())
	if err != nil || n != 1 {
		t.Fatalf("lift n=%d err=%v", n, err)
	}
	_, total2, _ := r.ListActive(context.Background(), SanctionListFilter{}, 20, 0)
	if total2 != 0 {
		t.Fatalf("want 0 active got %d", total2)
	}
}
