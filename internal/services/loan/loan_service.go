package loan

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	activityEventLoanReturned = "loan_returned"
	activityEventBookReturned = activityEventLoanReturned // backward-compatible alias for tests/references
	activityEventLoanCreated  = "loan_created"
)

type reputationHook interface {
	RecordLoanCreated(ctx context.Context, studentID string)
	RecordReturn(ctx context.Context, loan *models.Loan)
}

type catalogBookReader interface {
	GetByID(ctx context.Context, id string) (*models.Book, error)
}

type activeStudentReader interface {
	GetActiveByID(ctx context.Context, id string) (*models.Student, error)
}

// activeSanctionChecker counts active sanctions for a student (optional; nil skips the check).
type activeSanctionChecker interface {
	CountActiveByStudentID(ctx context.Context, studentID string) (int64, error)
}

// lateReturnFineHook optionally creates a pending fine when a return is after due_date (calendar).
type lateReturnFineHook interface {
	EnsureLateReturnFine(ctx context.Context, loan *models.Loan) error
}

// Service handles loan workflows for staff API.
type Service struct {
	loans       repositories.LoanRepository
	books       catalogBookReader
	students    activeStudentReader
	rep         reputationHook
	sanctions   activeSanctionChecker
	activityLog repositories.ActivityLogRepository
	lateFines   lateReturnFineHook
}

// NewService constructs LoanService. sanctions may be nil (tests); when set, Create rejects students with active sanctions.
// activityLog may be nil (tests); when set, Return writes a row to activity_log for the analytics feed.
// lateFines may be nil; when set, Return may create an automatic pending fine for late calendar returns (see fine.Service).
func NewService(
	loans repositories.LoanRepository,
	books catalogBookReader,
	students activeStudentReader,
	rep reputationHook,
	sanctions activeSanctionChecker,
	activityLog repositories.ActivityLogRepository,
	lateFines lateReturnFineHook,
) *Service {
	return &Service{
		loans: loans, books: books, students: students, rep: rep, sanctions: sanctions,
		activityLog: activityLog, lateFines: lateFines,
	}
}

// List returns paginated loans for the control panel.
func (s *Service) List(ctx context.Context, filter repositories.LoanListFilter, limit, offset int) ([]models.LoanListItem, int64, error) {
	st := strings.TrimSpace(strings.ToLower(filter.Status))
	if st == "" {
		st = string(models.LoanStatusActive)
	}
	if !isValidLoanStatus(st) {
		return nil, 0, ErrInvalidLoanStatus
	}
	filter.Status = st
	if strings.TrimSpace(filter.StudentID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(filter.StudentID)); err != nil {
			return nil, 0, ErrInvalidStudentIDQuery
		}
	}
	limit, offset = normalizeListParams(limit, offset)
	rows, total, err := s.loans.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	now := time.Now().UTC()
	out := make([]models.LoanListItem, 0, len(rows))
	for i := range rows {
		out = append(out, joinRowToListItem(&rows[i], now))
	}
	return out, total, nil
}

// Create registers a new active loan.
func (s *Service) Create(ctx context.Context, req models.CreateLoanRequest, adminID string) (*models.LoanDetailResponse, error) {
	due, err := parseDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}
	if beforeCalendarDay(due, time.Now().UTC()) {
		return nil, ErrDueDateInPast
	}

	st, err := s.students.GetActiveByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}

	if s.sanctions != nil {
		n, err := s.sanctions.CountActiveByStudentID(ctx, req.StudentID)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			return nil, ErrStudentHasActiveSanction
		}
	}

	b, err := s.books.GetByID(ctx, req.BookID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookNotFound
	}

	outCount, err := s.loans.CountCheckedOutByBookID(ctx, req.BookID)
	if err != nil {
		return nil, err
	}
	if b.TotalCopies <= 0 || int(outCount) >= b.TotalCopies {
		return nil, ErrNoCopiesAvailable
	}

	issued := adminID
	loan := &models.Loan{
		ID:         uuid.NewString(),
		BookID:     req.BookID,
		StudentID:  req.StudentID,
		IssuedBy:   &issued,
		BorrowedAt: time.Now().UTC(),
		DueDate:    due,
		Status:     models.LoanStatusActive,
	}
	if err := s.loans.Create(ctx, loan); err != nil {
		return nil, err
	}
	s.logLoanCreated(ctx, loan, adminID)

	s.rep.RecordLoanCreated(ctx, req.StudentID)

	row, err := s.loans.GetJoinByID(ctx, loan.ID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, errors.New("loan created but not found for detail")
	}
	return joinRowToDetail(row, time.Now().UTC()), nil
}

// Return marks a loan as returned. actorUserID is the staff user id from JWT (for activity_log.actor_id).
func (s *Service) Return(ctx context.Context, loanID string, actorUserID string) (*models.LoanDetailResponse, error) {
	loan, err := s.loans.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, ErrLoanNotFound
	}
	if loan.Status == models.LoanStatusReturned {
		return nil, ErrLoanAlreadyReturned
	}

	retAt := time.Now().UTC()
	if err := s.loans.Return(ctx, loanID, retAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLoanAlreadyReturned
		}
		return nil, err
	}

	loan, err = s.loans.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, ErrLoanNotFound
	}

	s.rep.RecordReturn(ctx, loan)
	if s.lateFines != nil {
		if err := s.lateFines.EnsureLateReturnFine(ctx, loan); err != nil {
			log.Printf("loan late return fine: %v", err)
		}
	}
	s.logBookReturned(ctx, loan, actorUserID)

	row, err := s.loans.GetJoinByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		log.Printf("loan return: missing join row for %s", loanID)
		return &models.LoanDetailResponse{
			LoanID:        loan.ID,
			BookID:        loan.BookID,
			Status:        string(loan.Status),
			BorrowedAt:    loan.BorrowedAt.UTC().Format(time.RFC3339),
			DueDate:       loan.DueDate.UTC().Format("2006-01-02"),
			TimeRemaining: daysRemaining(loan.DueDate, time.Now().UTC()),
		}, nil
	}
	return joinRowToDetail(row, time.Now().UTC()), nil
}

func isValidLoanStatus(s string) bool {
	switch models.LoanStatus(s) {
	case models.LoanStatusActive, models.LoanStatusReturned, models.LoanStatusOverdue:
		return true
	default:
		return false
	}
}

func parseDueDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return time.Time{}, ErrInvalidDueDateFormat
	}
	return t, nil
}

func (s *Service) logBookReturned(ctx context.Context, loan *models.Loan, actorUserID string) {
	if s.activityLog == nil || loan == nil || loan.ReturnedAt == nil {
		return
	}
	onTime := calendarReturnedOnTime(loan.DueDate, *loan.ReturnedAt)
	meta, err := json.Marshal(map[string]any{
		"loan_id":             loan.ID,
		"book_id":             loan.BookID,
		"returned_on_time":    onTime,
		"returned_at_rfc3339": loan.ReturnedAt.UTC().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("loan activity_log metadata: %v", err)
		return
	}
	desc := "Devolución a tiempo"
	if !onTime {
		desc = "Devolución fuera de plazo"
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: activityEventLoanReturned,
		Title:     "Libro devuelto",
		Metadata:  meta,
	}
	entry.Description = &desc
	sid := loan.StudentID
	entry.StudentID = &sid
	aid := strings.TrimSpace(actorUserID)
	if aid != "" {
		entry.ActorID = &aid
	}
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("loan activity_log create: %v", err)
	}
}

func (s *Service) logLoanCreated(ctx context.Context, loan *models.Loan, actorUserID string) {
	if s.activityLog == nil || loan == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"loan_id":       loan.ID,
		"book_id":       loan.BookID,
		"student_id":    loan.StudentID,
		"due_date":      loan.DueDate.UTC().Format("2006-01-02"),
		"borrowed_at":   loan.BorrowedAt.UTC().Format(time.RFC3339),
		"loan_status":   string(loan.Status),
	})
	if err != nil {
		log.Printf("loan activity_log metadata: %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: activityEventLoanCreated,
		Title:     "Préstamo creado",
		Metadata:  meta,
	}
	desc := "Se registró un nuevo préstamo"
	entry.Description = &desc
	sid := loan.StudentID
	entry.StudentID = &sid
	aid := strings.TrimSpace(actorUserID)
	if aid != "" {
		entry.ActorID = &aid
	}
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("loan activity_log create: %v", err)
	}
}

func calendarReturnedOnTime(due, returned time.Time) bool {
	dueDay := time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, time.UTC)
	retDay := time.Date(returned.Year(), returned.Month(), returned.Day(), 0, 0, 0, 0, time.UTC)
	return !retDay.After(dueDay)
}

func beforeCalendarDay(due time.Time, now time.Time) bool {
	d := time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, time.UTC)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d.Before(n)
}

func daysRemaining(due time.Time, now time.Time) int {
	dueDay := time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, time.UTC)
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return int(dueDay.Sub(nowDay).Hours() / 24)
}

func joinRowToListItem(row *repositories.LoanJoinRow, now time.Time) models.LoanListItem {
	borrower := strings.TrimSpace(row.BorrowerFirst + " " + row.BorrowerLast)
	return models.LoanListItem{
		LoanID:        row.LoanID,
		BookID:        row.BookID,
		BookTitle:     row.BookTitle,
		Borrower:      borrower,
		DueDate:       row.DueDate.UTC().Format("2006-01-02"),
		TimeRemaining: daysRemaining(row.DueDate, now),
		ISBN:          row.BookISBN,
		Status:        row.Status,
	}
}

func joinRowToDetail(row *repositories.LoanJoinRow, now time.Time) *models.LoanDetailResponse {
	borrower := strings.TrimSpace(row.BorrowerFirst + " " + row.BorrowerLast)
	return &models.LoanDetailResponse{
		LoanID:        row.LoanID,
		BookID:        row.BookID,
		BookTitle:     row.BookTitle,
		Borrower:      borrower,
		BorrowedAt:    row.BorrowedAt.UTC().Format(time.RFC3339),
		DueDate:       row.DueDate.UTC().Format("2006-01-02"),
		TimeRemaining: daysRemaining(row.DueDate, now),
		ISBN:          row.BookISBN,
		Status:        row.Status,
	}
}

func normalizeListParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
