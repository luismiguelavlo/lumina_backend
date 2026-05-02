package repositories

import (
	"context"
	"sort"
	"sync"
	"time"

	"library_back/internal/models"

	"gorm.io/gorm"
)

var _ LoanRepository = (*LoanRepositoryInMemory)(nil)

// LoanRepositoryInMemory is a test/dev implementation of LoanRepository (no PostgreSQL).
type LoanRepositoryInMemory struct {
	mu       sync.Mutex
	loans    []models.Loan
	books    map[string]loanBookMeta
	students map[string]loanStudentMeta
}

type loanBookMeta struct {
	Title string
	ISBN  string
}

type loanStudentMeta struct {
	First, Last string
}

// NewLoanRepositoryInMemory constructs an empty in-memory loan store.
func NewLoanRepositoryInMemory() *LoanRepositoryInMemory {
	return &LoanRepositoryInMemory{
		books:    make(map[string]loanBookMeta),
		students: make(map[string]loanStudentMeta),
	}
}

// RegisterBookForJoin supplies title/isbn for GetJoinByID and List joins.
func (r *LoanRepositoryInMemory) RegisterBookForJoin(bookID, title, isbn string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[bookID] = loanBookMeta{Title: title, ISBN: isbn}
}

// RegisterStudentForJoin supplies borrower names for GetJoinByID and List joins.
func (r *LoanRepositoryInMemory) RegisterStudentForJoin(studentID, firstName, lastName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.students[studentID] = loanStudentMeta{First: firstName, Last: lastName}
}

func (r *LoanRepositoryInMemory) Create(ctx context.Context, loan *models.Loan) error {
	_ = ctx
	if loan == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *loan
	r.loans = append(r.loans, cp)
	return nil
}

func (r *LoanRepositoryInMemory) GetByID(ctx context.Context, id string) (*models.Loan, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.loans {
		if r.loans[i].ID == id {
			l := r.loans[i]
			return &l, nil
		}
	}
	return nil, nil
}

func (r *LoanRepositoryInMemory) GetJoinByID(ctx context.Context, id string) (*LoanJoinRow, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.loans {
		if r.loans[i].ID != id {
			continue
		}
		return r.joinRowLocked(&r.loans[i]), nil
	}
	return nil, nil
}

func (r *LoanRepositoryInMemory) joinRowLocked(l *models.Loan) *LoanJoinRow {
	bm, okB := r.books[l.BookID]
	sm, okS := r.students[l.StudentID]
	if !okB || !okS {
		return nil
	}
	var ret *time.Time
	if l.ReturnedAt != nil {
		t := *l.ReturnedAt
		ret = &t
	}
	return &LoanJoinRow{
		LoanID:        l.ID,
		BookID:        l.BookID,
		BookTitle:     bm.Title,
		BookISBN:      bm.ISBN,
		StudentID:     l.StudentID,
		BorrowerFirst: sm.First,
		BorrowerLast:  sm.Last,
		DueDate:       l.DueDate,
		Status:        string(l.Status),
		BorrowedAt:    l.BorrowedAt,
		ReturnedAt:    ret,
	}
}

func (r *LoanRepositoryInMemory) List(ctx context.Context, filter LoanListFilter, limit, offset int) ([]LoanJoinRow, int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	status := filter.Status
	if status == "" {
		status = string(models.LoanStatusActive)
	}

	var idx []int
	for i := range r.loans {
		l := &r.loans[i]
		if string(l.Status) != status {
			continue
		}
		if filter.StudentID != "" && l.StudentID != filter.StudentID {
			continue
		}
		idx = append(idx, i)
	}

	sort.Slice(idx, func(a, b int) bool {
		ia, ib := idx[a], idx[b]
		da := r.loans[ia].DueDate
		db := r.loans[ib].DueDate
		if da.Equal(db) {
			return r.loans[ia].ID < r.loans[ib].ID
		}
		return da.Before(db)
	})

	total := int64(len(idx))
	if offset > len(idx) {
		offset = len(idx)
	}
	end := offset + limit
	if end > len(idx) {
		end = len(idx)
	}
	slice := idx[offset:end]

	out := make([]LoanJoinRow, 0, len(slice))
	for _, i := range slice {
		row := r.joinRowLocked(&r.loans[i])
		if row == nil {
			continue
		}
		out = append(out, *row)
	}
	return out, total, nil
}

func (r *LoanRepositoryInMemory) Return(ctx context.Context, id string, returnedAt time.Time) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.loans {
		if r.loans[i].ID != id {
			continue
		}
		st := r.loans[i].Status
		if st != models.LoanStatusActive && st != models.LoanStatusOverdue {
			return gorm.ErrRecordNotFound
		}
		r.loans[i].Status = models.LoanStatusReturned
		t := returnedAt
		r.loans[i].ReturnedAt = &t
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (r *LoanRepositoryInMemory) CountCheckedOutByBookID(ctx context.Context, bookID string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var n int64
	for i := range r.loans {
		if r.loans[i].BookID != bookID {
			continue
		}
		st := r.loans[i].Status
		if st == models.LoanStatusActive || st == models.LoanStatusOverdue {
			n++
		}
	}
	return n, nil
}

func (r *LoanRepositoryInMemory) MarkActiveLoansOverdue(ctx context.Context) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	var n int64
	for i := range r.loans {
		if r.loans[i].Status != models.LoanStatusActive {
			continue
		}
		d := time.Date(r.loans[i].DueDate.Year(), r.loans[i].DueDate.Month(), r.loans[i].DueDate.Day(), 0, 0, 0, 0, time.UTC)
		if d.Before(today) {
			r.loans[i].Status = models.LoanStatusOverdue
			n++
		}
	}
	return n, nil
}
