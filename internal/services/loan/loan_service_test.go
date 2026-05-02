package loan

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type noopRep struct{}

func (noopRep) RecordLoanCreated(ctx context.Context, studentID string) {}
func (noopRep) RecordReturn(ctx context.Context, loan *models.Loan)     {}

type fakeLoanRepo struct {
	createErr   error
	getByID     *models.Loan
	getJoin     *repositories.LoanJoinRow
	listRows    []repositories.LoanJoinRow
	listTotal   int64
	listErr     error
	returnErr   error
	checkedOut  int64
	countErr    error
	lastCreate  *models.Loan
	returnCalls int
}

func (f *fakeLoanRepo) Create(ctx context.Context, loan *models.Loan) error {
	_ = ctx
	f.lastCreate = loan
	return f.createErr
}

func (f *fakeLoanRepo) GetByID(ctx context.Context, id string) (*models.Loan, error) {
	_ = ctx
	_ = id
	return f.getByID, nil
}

func (f *fakeLoanRepo) GetJoinByID(ctx context.Context, id string) (*repositories.LoanJoinRow, error) {
	_ = ctx
	_ = id
	return f.getJoin, nil
}

func (f *fakeLoanRepo) List(ctx context.Context, filter repositories.LoanListFilter, limit, offset int) ([]repositories.LoanJoinRow, int64, error) {
	_ = ctx
	_ = filter
	_ = limit
	_ = offset
	if f.listErr != nil {
		return nil, 0, f.listErr
	}
	return f.listRows, f.listTotal, nil
}

func (f *fakeLoanRepo) Return(ctx context.Context, id string, returnedAt time.Time) error {
	_ = ctx
	_ = id
	_ = returnedAt
	f.returnCalls++
	return f.returnErr
}

func (f *fakeLoanRepo) CountCheckedOutByBookID(ctx context.Context, bookID string) (int64, error) {
	_ = ctx
	_ = bookID
	return f.checkedOut, f.countErr
}

func (f *fakeLoanRepo) MarkActiveLoansOverdue(ctx context.Context) (int64, error) {
	_ = ctx
	return 0, nil
}

type fakeBookGetter struct {
	book *models.Book
}

func (f *fakeBookGetter) GetByID(ctx context.Context, id string) (*models.Book, error) {
	_ = ctx
	_ = id
	return f.book, nil
}

type fakeStudentGetter struct {
	st *models.Student
}

func (f *fakeStudentGetter) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	_ = ctx
	_ = id
	return f.st, nil
}

func TestLoanService_ListInvalidStatus(t *testing.T) {
	s := NewService(&fakeLoanRepo{}, &fakeBookGetter{}, &fakeStudentGetter{}, noopRep{}, nil, nil, nil)
	_, _, err := s.List(context.Background(), repositories.LoanListFilter{Status: "nope"}, 20, 0)
	if err != ErrInvalidLoanStatus {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_CreateDueDatePast(t *testing.T) {
	s := NewService(&fakeLoanRepo{}, &fakeBookGetter{book: &models.Book{ID: "b", TotalCopies: 2}}, &fakeStudentGetter{st: &models.Student{ID: "s"}}, noopRep{}, nil, nil, nil)
	_, err := s.Create(context.Background(), models.CreateLoanRequest{
		StudentID: uuid.NewString(),
		BookID:    uuid.NewString(),
		DueDate:   "2000-01-01",
	}, "admin")
	if err != ErrDueDateInPast {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_CreateNoCopies(t *testing.T) {
	s := NewService(
		&fakeLoanRepo{checkedOut: 2},
		&fakeBookGetter{book: &models.Book{ID: "b", TotalCopies: 2}},
		&fakeStudentGetter{st: &models.Student{ID: "s"}},
		noopRep{},
		nil,
		nil,
		nil,
	)
	_, err := s.Create(context.Background(), models.CreateLoanRequest{
		StudentID: uuid.NewString(),
		BookID:    uuid.NewString(),
		DueDate:   time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02"),
	}, "admin")
	if err != ErrNoCopiesAvailable {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_ReturnAlreadyReturned(t *testing.T) {
	fr := &fakeLoanRepo{}
	st := NewService(fr, &fakeBookGetter{}, &fakeStudentGetter{}, noopRep{}, nil, nil, nil)
	fr.getByID = &models.Loan{ID: "L1", Status: models.LoanStatusReturned}
	_, err := st.Return(context.Background(), "L1", "admin")
	if err != ErrLoanAlreadyReturned {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_ReturnNotFound(t *testing.T) {
	fr := &fakeLoanRepo{getByID: nil}
	st := NewService(fr, &fakeBookGetter{}, &fakeStudentGetter{}, noopRep{}, nil, nil, nil)
	_, err := st.Return(context.Background(), uuid.NewString(), "admin")
	if err != ErrLoanNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_ReturnConflictOnRepo(t *testing.T) {
	fr := &fakeLoanRepo{
		getByID:   &models.Loan{ID: "L1", Status: models.LoanStatusActive},
		returnErr: gorm.ErrRecordNotFound,
	}
	st := NewService(fr, &fakeBookGetter{}, &fakeStudentGetter{}, noopRep{}, nil, nil, nil)
	_, err := st.Return(context.Background(), "L1", "admin")
	if err != ErrLoanAlreadyReturned {
		t.Fatalf("got %v", err)
	}
}
