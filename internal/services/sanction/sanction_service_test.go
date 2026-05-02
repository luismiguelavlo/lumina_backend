package sanction

import (
	"context"
	"testing"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

type fakeStudentReader struct {
	st *models.Student
}

func (f *fakeStudentReader) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	_ = ctx
	_ = id
	return f.st, nil
}

func TestSanctionService_CreateStudentNotFound(t *testing.T) {
	mem := repositories.NewSanctionRepositoryInMemory()
	s := NewService(mem, &fakeStudentReader{st: nil}, nil)
	_, err := s.Create(context.Background(), models.CreateSanctionRequest{
		StudentID: uuid.NewString(),
		Reason:    "r",
	}, "admin-1")
	if err != ErrStudentNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestSanctionService_CreateSuccess(t *testing.T) {
	mem := repositories.NewSanctionRepositoryInMemory()
	sid := uuid.NewString()
	st := &models.Student{ID: sid, StudentIDCode: "S-1", FirstName: "A", LastName: "B", Email: strPtr("a@b.c")}
	s := NewService(mem, &fakeStudentReader{st: st}, nil)
	out, err := s.Create(context.Background(), models.CreateSanctionRequest{
		StudentID: sid,
		Reason:    "Late returns",
	}, "admin-uuid-1")
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "active" || out.Reason != "Late returns" || out.AppliedByID == nil || *out.AppliedByID != "admin-uuid-1" {
		t.Fatalf("%+v", out)
	}
}

func TestSanctionService_LiftNotFound(t *testing.T) {
	s := NewService(repositories.NewSanctionRepositoryInMemory(), &fakeStudentReader{st: &models.Student{ID: "x"}}, nil)
	_, err := s.Lift(context.Background(), uuid.NewString(), "admin")
	if err != ErrSanctionNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestSanctionService_LiftSuccess(t *testing.T) {
	mem := repositories.NewSanctionRepositoryInMemory()
	sid := uuid.NewString()
	st := &models.Student{ID: sid, StudentIDCode: "S-1", FirstName: "A", LastName: "B"}
	ab := "admin-1"
	_ = mem.Create(context.Background(), &models.Sanction{
		ID:        "san-1",
		StudentID: sid,
		Reason:    "x",
		Status:    models.SanctionStatusActive,
		AppliedBy: &ab,
	})
	svc := NewService(mem, &fakeStudentReader{st: st}, nil)
	out, err := svc.Lift(context.Background(), "san-1", "admin-2")
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "lifted" || out.LiftedByID == nil || *out.LiftedByID != "admin-2" {
		t.Fatalf("%+v", out)
	}
}

func strPtr(s string) *string { return &s }
