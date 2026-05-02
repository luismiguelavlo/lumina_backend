package loan

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

type mapBookGetter map[string]*models.Book

func (m mapBookGetter) GetByID(ctx context.Context, id string) (*models.Book, error) {
	_ = ctx
	return m[id], nil
}

type mapStudentGetter map[string]*models.Student

func (m mapStudentGetter) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	_ = ctx
	return m[id], nil
}

type spyReputation struct {
	loanCreatedFor []string
	returns        []*models.Loan
}

func (s *spyReputation) RecordLoanCreated(ctx context.Context, studentID string) {
	_ = ctx
	s.loanCreatedFor = append(s.loanCreatedFor, studentID)
}

func (s *spyReputation) RecordReturn(ctx context.Context, loan *models.Loan) {
	_ = ctx
	if loan != nil {
		cp := *loan
		s.returns = append(s.returns, &cp)
	}
}

type spyActivityLog struct {
	entries []*models.ActivityLog
	err     error
}

func (s *spyActivityLog) Create(ctx context.Context, e *models.ActivityLog) error {
	_ = ctx
	if s.err != nil {
		return s.err
	}
	cp := *e
	s.entries = append(s.entries, &cp)
	return nil
}

func TestLoanService_CreateSuccessAndIssuedBy(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	sid := uuid.NewString()
	bid := uuid.NewString()
	mem.RegisterBookForJoin(bid, "Title", "ISBN-9")
	mem.RegisterStudentForJoin(sid, "Ana", "López")

	books := mapBookGetter{
		bid: &models.Book{ID: bid, TotalCopies: 3, Title: "Title", ISBN: "ISBN-9"},
	}
	students := mapStudentGetter{
		sid: &models.Student{ID: sid, FirstName: "Ana", LastName: "López", IsActive: true},
	}
	spy := &spyReputation{}
	svc := NewService(mem, books, students, spy, nil, nil, nil)

	out, err := svc.Create(ctx, models.CreateLoanRequest{
		StudentID: sid,
		BookID:    bid,
		DueDate:   time.Now().UTC().AddDate(0, 0, 14).Format("2006-01-02"),
	}, "admin-user-7")
	if err != nil {
		t.Fatal(err)
	}
	if out.LoanID == "" || out.BookTitle != "Title" || out.Borrower != "Ana López" || out.Status != "active" {
		t.Fatalf("response %+v", out)
	}
	if len(spy.loanCreatedFor) != 1 || spy.loanCreatedFor[0] != sid {
		t.Fatalf("reputation: %+v", spy.loanCreatedFor)
	}
	// issued_by on persisted loan
	got, _ := mem.GetByID(ctx, out.LoanID)
	if got == nil || got.IssuedBy == nil || *got.IssuedBy != "admin-user-7" {
		t.Fatalf("issued_by: %+v", got)
	}
}

func TestLoanService_StudentNotFound(t *testing.T) {
	ctx := context.Background()
	sid := uuid.NewString()
	bid := uuid.NewString()
	svc := NewService(
		repositories.NewLoanRepositoryInMemory(),
		mapBookGetter{bid: &models.Book{ID: bid, TotalCopies: 1}},
		mapStudentGetter{}, // empty
		noopRep{},
		nil,
		nil,
		nil,
	)
	_, err := svc.Create(ctx, models.CreateLoanRequest{
		StudentID: sid,
		BookID:    bid,
		DueDate:   time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02"),
	}, "admin")
	if err != ErrStudentNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_BookNotFound(t *testing.T) {
	ctx := context.Background()
	sid := uuid.NewString()
	bid := uuid.NewString()
	svc := NewService(
		repositories.NewLoanRepositoryInMemory(),
		mapBookGetter{},
		mapStudentGetter{sid: &models.Student{ID: sid, IsActive: true}},
		noopRep{},
		nil,
		nil,
		nil,
	)
	_, err := svc.Create(ctx, models.CreateLoanRequest{
		StudentID: sid,
		BookID:    bid,
		DueDate:   time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02"),
	}, "admin")
	if err != ErrBookNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_ListInvalidStudentIDQuery(t *testing.T) {
	svc := NewService(&fakeLoanRepo{}, &fakeBookGetter{}, &fakeStudentGetter{}, noopRep{}, nil, nil, nil)
	_, _, err := svc.List(context.Background(), repositories.LoanListFilter{StudentID: "not-a-uuid"}, 20, 0)
	if err != ErrInvalidStudentIDQuery {
		t.Fatalf("got %v", err)
	}
}

func TestLoanService_ListTimeRemaining(t *testing.T) {
	due := time.Now().UTC().AddDate(0, 0, 3)
	due = time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, time.UTC)
	fr := &fakeLoanRepo{
		listRows: []repositories.LoanJoinRow{{
			LoanID: "L1", BookID: "b", BookTitle: "T", BookISBN: "x",
			BorrowerFirst: "A", BorrowerLast: "B",
			DueDate: due, Status: "active",
			BorrowedAt: time.Now().UTC(),
		}},
		listTotal: 1,
	}
	svc := NewService(fr, &fakeBookGetter{}, &fakeStudentGetter{}, noopRep{}, nil, nil, nil)
	items, total, err := svc.List(context.Background(), repositories.LoanListFilter{}, 20, 0)
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("list: %v %+v", err, items)
	}
	if items[0].TimeRemaining != 3 {
		t.Fatalf("time_remaining want3 got %d", items[0].TimeRemaining)
	}
}

func TestLoanService_ReturnSuccess(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	sid := uuid.NewString()
	bid := uuid.NewString()
	mem.RegisterBookForJoin(bid, "Book", "ISBN")
	mem.RegisterStudentForJoin(sid, "X", "Y")
	lid := uuid.NewString()
	_ = mem.Create(ctx, &models.Loan{
		ID: lid, BookID: bid, StudentID: sid,
		DueDate:    time.Now().UTC().AddDate(0, 0, 1),
		Status:     models.LoanStatusActive,
		BorrowedAt: time.Now().UTC(),
	})
	spy := &spyReputation{}
	svc := NewService(
		mem,
		mapBookGetter{bid: &models.Book{ID: bid, TotalCopies: 2}},
		mapStudentGetter{sid: &models.Student{ID: sid, IsActive: true}},
		spy,
		nil,
		nil,
		nil,
	)
	out, err := svc.Return(ctx, lid, "admin-1")
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "returned" || out.BookTitle != "Book" {
		t.Fatalf("out %+v", out)
	}
	if len(spy.returns) != 1 || spy.returns[0].ID != lid || spy.returns[0].Status != models.LoanStatusReturned {
		t.Fatalf("spy %+v", spy.returns)
	}
}

func TestLoanService_ReturnWritesActivityLog(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	sid := uuid.NewString()
	bid := uuid.NewString()
	mem.RegisterBookForJoin(bid, "Book", "ISBN")
	mem.RegisterStudentForJoin(sid, "X", "Y")
	lid := uuid.NewString()
	_ = mem.Create(ctx, &models.Loan{
		ID: lid, BookID: bid, StudentID: sid,
		DueDate:    time.Now().UTC().AddDate(0, 0, 1),
		Status:     models.LoanStatusActive,
		BorrowedAt: time.Now().UTC(),
	})
	logSpy := &spyActivityLog{}
	svc := NewService(
		mem,
		mapBookGetter{bid: &models.Book{ID: bid, TotalCopies: 2}},
		mapStudentGetter{sid: &models.Student{ID: sid, IsActive: true}},
		noopRep{},
		nil,
		logSpy,
		nil,
	)
	actor := "00000000-0000-0000-0000-000000000099"
	_, err := svc.Return(ctx, lid, actor)
	if err != nil {
		t.Fatal(err)
	}
	if len(logSpy.entries) != 1 {
		t.Fatalf("entries %d", len(logSpy.entries))
	}
	e := logSpy.entries[0]
	if e.EventType != activityEventBookReturned || e.Title != "Libro devuelto" {
		t.Fatalf("entry %+v", e)
	}
	if e.StudentID == nil || *e.StudentID != sid || e.ActorID == nil || *e.ActorID != actor {
		t.Fatalf("actor/student %+v", e)
	}
	var meta map[string]any
	if err := json.Unmarshal(e.Metadata, &meta); err != nil {
		t.Fatal(err)
	}
	if meta["loan_id"] != lid || meta["book_id"] != bid {
		t.Fatalf("metadata %+v", meta)
	}
	if meta["returned_on_time"] != true {
		t.Fatalf("want on_time metadata %+v", meta)
	}
}

func TestLoanService_ReturnActivityLogErrorIgnored(t *testing.T) {
	ctx := context.Background()
	mem := repositories.NewLoanRepositoryInMemory()
	sid := uuid.NewString()
	bid := uuid.NewString()
	mem.RegisterBookForJoin(bid, "Book", "ISBN")
	mem.RegisterStudentForJoin(sid, "X", "Y")
	lid := uuid.NewString()
	_ = mem.Create(ctx, &models.Loan{
		ID: lid, BookID: bid, StudentID: sid,
		DueDate:    time.Now().UTC().AddDate(0, 0, 1),
		Status:     models.LoanStatusActive,
		BorrowedAt: time.Now().UTC(),
	})
	logSpy := &spyActivityLog{err: errors.New("db down")}
	svc := NewService(
		mem,
		mapBookGetter{bid: &models.Book{ID: bid, TotalCopies: 2}},
		mapStudentGetter{sid: &models.Student{ID: sid, IsActive: true}},
		noopRep{},
		nil,
		logSpy,
		nil,
	)
	_, err := svc.Return(ctx, lid, "admin")
	if err != nil {
		t.Fatalf("return should succeed despite activity log: %v", err)
	}
}

func TestLoanService_InvalidDueDateFormat(t *testing.T) {
	svc := NewService(&fakeLoanRepo{}, &fakeBookGetter{book: &models.Book{TotalCopies: 1}}, &fakeStudentGetter{st: &models.Student{ID: "s"}}, noopRep{}, nil, nil, nil)
	_, err := svc.Create(context.Background(), models.CreateLoanRequest{
		StudentID: uuid.NewString(),
		BookID:    uuid.NewString(),
		DueDate:   "02-01-2026",
	}, "admin")
	if err != ErrInvalidDueDateFormat {
		t.Fatalf("got %v", err)
	}
}
