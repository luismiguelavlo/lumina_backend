package loan

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

type stubSanctionCount struct {
	n   int64
	err error
}

func (s *stubSanctionCount) CountActiveByStudentID(ctx context.Context, studentID string) (int64, error) {
	_ = ctx
	_ = studentID
	return s.n, s.err
}

func TestLoanService_CreateRejectedWhenActiveSanction(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	sid := uuid.NewString()
	bid := uuid.NewString()
	mem.RegisterBookForJoin(bid, "T", "ISBN")
	mem.RegisterStudentForJoin(sid, "A", "B")
	svc := NewService(
		mem,
		mapBookGetter{bid: &models.Book{ID: bid, TotalCopies: 2}},
		mapStudentGetter{sid: &models.Student{ID: sid, IsActive: true}},
		noopRep{},
		&stubSanctionCount{n: 1},
		nil,
		nil,
	)
	_, err := svc.Create(ctx, models.CreateLoanRequest{
		StudentID: sid,
		BookID:    bid,
		DueDate:   time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02"),
	}, "admin")
	if err != ErrStudentHasActiveSanction {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_CreateAllowedWhenNoActiveSanction(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	sid := uuid.NewString()
	bid := uuid.NewString()
	mem.RegisterBookForJoin(bid, "T", "ISBN")
	mem.RegisterStudentForJoin(sid, "A", "B")
	svc := NewService(
		mem,
		mapBookGetter{bid: &models.Book{ID: bid, TotalCopies: 2}},
		mapStudentGetter{sid: &models.Student{ID: sid, IsActive: true}},
		noopRep{},
		&stubSanctionCount{n: 0},
		nil,
		nil,
	)
	out, err := svc.Create(ctx, models.CreateLoanRequest{
		StudentID: sid,
		BookID:    bid,
		DueDate:   time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02"),
	}, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "active" {
		t.Fatalf("%+v", out)
	}
}
