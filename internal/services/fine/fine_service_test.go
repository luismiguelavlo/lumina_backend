package fine

import (
	"context"
	"testing"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type fakeLoanReader struct {
	loan *models.Loan
	err  error
}

func (f *fakeLoanReader) GetByID(ctx context.Context, id string) (*models.Loan, error) {
	_ = ctx
	_ = id
	if f.err != nil {
		return nil, f.err
	}
	return f.loan, nil
}

func TestFineService_CreateLoanNotFound(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	s := NewService(mem, &fakeLoanReader{loan: nil}, nil)
	_, err := s.Create(context.Background(), models.CreateFineRequest{
		LoanID:    uuid.NewString(),
		StudentID: uuid.NewString(),
		Amount:    10,
	})
	if err != ErrFineLoanNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_CreateStudentMismatch(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	loanID := uuid.NewString()
	studentOnLoan := uuid.NewString()
	s := NewService(mem, &fakeLoanReader{loan: &models.Loan{ID: loanID, StudentID: studentOnLoan}}, nil)
	_, err := s.Create(context.Background(), models.CreateFineRequest{
		LoanID:    loanID,
		StudentID: uuid.NewString(),
		Amount:    5,
	})
	if err != ErrLoanStudentMismatch {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_CreateDuplicatePending(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	loanID := uuid.NewString()
	sid := uuid.NewString()
	_ = mem.Create(context.Background(), &models.Fine{
		LoanID:    loanID,
		StudentID: sid,
		Amount:    decimal.NewFromInt(1),
		Status:    models.FineStatusPending,
	})
	s := NewService(mem, &fakeLoanReader{loan: &models.Loan{ID: loanID, StudentID: sid}}, nil)
	_, err := s.Create(context.Background(), models.CreateFineRequest{
		LoanID:    loanID,
		StudentID: sid,
		Amount:    2,
	})
	if err != ErrFineDuplicatePending {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_CreateSuccess(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	loanID := uuid.NewString()
	sid := uuid.NewString()
	s := NewService(mem, &fakeLoanReader{loan: &models.Loan{ID: loanID, StudentID: sid}}, nil)
	out, err := s.Create(context.Background(), models.CreateFineRequest{
		LoanID:    loanID,
		StudentID: sid,
		Amount:    15.5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "pending" || out.Amount != 15.5 {
		t.Fatalf("%+v", out)
	}
}

func TestFineService_ListInvalidStatus(t *testing.T) {
	s := NewService(repositories.NewFineRepositoryInMemory(), &fakeLoanReader{}, nil)
	_, _, err := s.List(context.Background(), repositories.FineListFilter{Status: "nope"}, 20, 0)
	if err != ErrInvalidFineStatus {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_ListInvalidStudentID(t *testing.T) {
	s := NewService(repositories.NewFineRepositoryInMemory(), &fakeLoanReader{}, nil)
	_, _, err := s.List(context.Background(), repositories.FineListFilter{StudentID: "bad"}, 20, 0)
	if err != ErrInvalidStudentIDQuery {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_MarkPaidNotPending(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	id := uuid.NewString()
	_ = mem.Create(context.Background(), &models.Fine{
		ID:        id,
		LoanID:    uuid.NewString(),
		StudentID: uuid.NewString(),
		Amount:    decimal.NewFromInt(1),
		Status:    models.FineStatusPaid,
	})
	s := NewService(mem, &fakeLoanReader{}, nil)
	_, err := s.MarkPaid(context.Background(), id)
	if err != ErrFineNotPending {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_MarkPaidSuccess(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	id := uuid.NewString()
	_ = mem.Create(context.Background(), &models.Fine{
		ID:        id,
		LoanID:    uuid.NewString(),
		StudentID: uuid.NewString(),
		Amount:    decimal.NewFromInt(3),
		Status:    models.FineStatusPending,
	})
	s := NewService(mem, &fakeLoanReader{}, nil)
	out, err := s.MarkPaid(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "paid" || out.PaidAt == nil {
		t.Fatalf("%+v", out)
	}
}

func TestFineService_MarkWaivedSuccess(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	id := uuid.NewString()
	_ = mem.Create(context.Background(), &models.Fine{
		ID:        id,
		LoanID:    uuid.NewString(),
		StudentID: uuid.NewString(),
		Amount:    decimal.NewFromInt(2),
		Status:    models.FineStatusPending,
	})
	s := NewService(mem, &fakeLoanReader{}, nil)
	out, err := s.MarkWaived(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "waived" || out.PaidAt != nil {
		t.Fatalf("%+v", out)
	}
}

func TestFineService_GetByIDNotFound(t *testing.T) {
	s := NewService(repositories.NewFineRepositoryInMemory(), &fakeLoanReader{}, nil)
	_, err := s.GetByID(context.Background(), uuid.NewString())
	if err != ErrFineNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestFineService_GetByIDSuccess(t *testing.T) {
	mem := repositories.NewFineRepositoryInMemory()
	id := uuid.NewString()
	_ = mem.Create(context.Background(), &models.Fine{
		ID:        id,
		LoanID:    uuid.NewString(),
		StudentID: uuid.NewString(),
		Amount:    decimal.NewFromInt(7),
		Status:    models.FineStatusPending,
	})
	s := NewService(mem, &fakeLoanReader{}, nil)
	out, err := s.GetByID(context.Background(), id)
	if err != nil || out.ID != id {
		t.Fatalf("%+v err=%v", out, err)
	}
}
