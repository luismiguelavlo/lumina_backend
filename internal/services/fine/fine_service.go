package fine

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type loanByIDReader interface {
	GetByID(ctx context.Context, id string) (*models.Loan, error)
}

// finePaidReputationHook is optional; when set, MarkPaid records a fine_paid reputation event.
type finePaidReputationHook interface {
	RecordFinePaid(ctx context.Context, fineID, studentID string)
}

// Service handles fines for staff API.
type Service struct {
	fines repositories.FineRepository
	loans loanByIDReader
	rep   finePaidReputationHook
	activityLog repositories.ActivityLogRepository
}

// NewService constructs FineService. rep may be nil (tests); when set, paid fines grant reputation points once per fine.
func NewService(fines repositories.FineRepository, loans loanByIDReader, rep finePaidReputationHook) *Service {
	return &Service{fines: fines, loans: loans, rep: rep}
}

// SetActivityLog configures optional activity log persistence for dashboard feed events.
func (s *Service) SetActivityLog(activityLog repositories.ActivityLogRepository) {
	s.activityLog = activityLog
}

// List returns paginated fines.
func (s *Service) List(ctx context.Context, filter repositories.FineListFilter, limit, offset int) ([]models.FineListItem, int64, error) {
	if strings.TrimSpace(filter.StudentID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(filter.StudentID)); err != nil {
			return nil, 0, ErrInvalidStudentIDQuery
		}
	}
	if st := strings.TrimSpace(strings.ToLower(filter.Status)); st != "" {
		if !isValidFineStatus(st) {
			return nil, 0, ErrInvalidFineStatus
		}
		filter.Status = st
	}
	limit, offset = normalizeListParams(limit, offset)
	rows, total, err := s.fines.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.FineListItem, 0, len(rows))
	for i := range rows {
		out = append(out, fineWithStudentToListItem(&rows[i]))
	}
	return out, total, nil
}

// Create registers a new pending fine.
func (s *Service) Create(ctx context.Context, req models.CreateFineRequest) (*models.FineResponse, error) {
	n, err := s.fines.CountPendingByLoanID(ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	if n > 0 {
		return nil, ErrFineDuplicatePending
	}
	loan, err := s.loans.GetByID(ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, ErrFineLoanNotFound
	}
	if loan.StudentID != req.StudentID {
		return nil, ErrLoanStudentMismatch
	}
	amt := decimal.NewFromFloat(req.Amount).Round(2)
	f := &models.Fine{
		LoanID:    req.LoanID,
		StudentID: req.StudentID,
		Amount:    amt,
		Status:    models.FineStatusPending,
		Reason:    req.Reason,
	}
	if err := s.fines.Create(ctx, f); err != nil {
		return nil, err
	}
	refreshed, err := s.fines.GetByID(ctx, f.ID)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrFineNotFound
	}
	s.logFineCreated(ctx, refreshed)
	return fineToResponse(refreshed), nil
}

// GetByID returns one fine.
func (s *Service) GetByID(ctx context.Context, id string) (*models.FineResponse, error) {
	f, err := s.fines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, ErrFineNotFound
	}
	return fineToResponse(f), nil
}

// MarkPaid sets status paid and paid_at = now.
func (s *Service) MarkPaid(ctx context.Context, id string) (*models.FineResponse, error) {
	f, err := s.fines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, ErrFineNotFound
	}
	if f.Status != models.FineStatusPending {
		return nil, ErrFineNotPending
	}
	now := time.Now().UTC()
	n, err := s.fines.UpdateStatus(ctx, id, models.FineStatusPaid, &now)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrFineNotPending
	}
	refreshed, err := s.fines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrFineNotFound
	}
	if s.rep != nil {
		s.rep.RecordFinePaid(ctx, id, refreshed.StudentID)
	}
	s.logFinePaid(ctx, refreshed)
	return fineToResponse(refreshed), nil
}

// MarkWaived sets status waived (paid_at stays null).
func (s *Service) MarkWaived(ctx context.Context, id string) (*models.FineResponse, error) {
	f, err := s.fines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, ErrFineNotFound
	}
	if f.Status != models.FineStatusPending {
		return nil, ErrFineNotPending
	}
	n, err := s.fines.UpdateStatus(ctx, id, models.FineStatusWaived, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrFineNotPending
	}
	refreshed, err := s.fines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrFineNotFound
	}
	s.logFineWaived(ctx, refreshed)
	return fineToResponse(refreshed), nil
}

func isValidFineStatus(s string) bool {
	switch models.FineStatus(s) {
	case models.FineStatusPending, models.FineStatusPaid, models.FineStatusWaived:
		return true
	default:
		return false
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

// EnsureLateReturnFine creates one pending fine for a returned loan when the return is after due_date (calendar)
// and AUTO_LATE_RETURN_FINE_ENABLED is not disabled. Idempotent with manual fines: skips if any pending fine exists for the loan.
func (s *Service) EnsureLateReturnFine(ctx context.Context, loan *models.Loan) error {
	if loan == nil || loan.ReturnedAt == nil || loan.Status != models.LoanStatusReturned {
		return nil
	}
	if !autoLateReturnFineEnabled() {
		return nil
	}
	if !calendarReturnIsLate(loan.DueDate, *loan.ReturnedAt) {
		return nil
	}
	n, err := s.fines.CountPendingByLoanID(ctx, loan.ID)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	amt := autoLateReturnFineAmount()
	reason := "Devolución fuera de plazo (automática)"
	_, err = s.Create(ctx, models.CreateFineRequest{
		LoanID:    loan.ID,
		StudentID: loan.StudentID,
		Amount:    amt,
		Reason:    &reason,
	})
	return err
}

func autoLateReturnFineEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("AUTO_LATE_RETURN_FINE_ENABLED")))
	switch v {
	case "0", "false", "no", "off":
		return false
	default:
		return true
	}
}

func autoLateReturnFineAmount() float64 {
	raw := strings.TrimSpace(os.Getenv("AUTO_LATE_RETURN_FINE_AMOUNT"))
	if raw == "" {
		return 5
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil || f < 0 {
		return 5
	}
	return f
}

func calendarReturnIsLate(due, returned time.Time) bool {
	dueDay := time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, time.UTC)
	retDay := time.Date(returned.Year(), returned.Month(), returned.Day(), 0, 0, 0, 0, time.UTC)
	return retDay.After(dueDay)
}

func fineToResponse(f *models.Fine) *models.FineResponse {
	var paid *string
	if f.PaidAt != nil {
		s := f.PaidAt.UTC().Format(time.RFC3339Nano)
		paid = &s
	}
	return &models.FineResponse{
		ID:        f.ID,
		LoanID:    f.LoanID,
		StudentID: f.StudentID,
		Amount:    f.Amount.Round(2).InexactFloat64(),
		Status:    string(f.Status),
		Reason:    f.Reason,
		CreatedAt: f.CreatedAt.UTC().Format(time.RFC3339Nano),
		PaidAt:    paid,
	}
}

func fineWithStudentToListItem(row *repositories.FineWithStudent) models.FineListItem {
	f := &row.Fine
	var studentName *string
	if n := strings.TrimSpace(row.StudentName); n != "" {
		studentName = &n
	}
	var paid *string
	if f.PaidAt != nil {
		s := f.PaidAt.UTC().Format(time.RFC3339Nano)
		paid = &s
	}
	return models.FineListItem{
		ID:          f.ID,
		LoanID:      f.LoanID,
		StudentID:   f.StudentID,
		Amount:      f.Amount.Round(2).InexactFloat64(),
		Status:      string(f.Status),
		Reason:      f.Reason,
		CreatedAt:   f.CreatedAt.UTC().Format(time.RFC3339Nano),
		PaidAt:      paid,
		StudentName: studentName,
	}
}

func (s *Service) logFineCreated(ctx context.Context, fine *models.Fine) {
	if s.activityLog == nil || fine == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"fine_id":      fine.ID,
		"loan_id":      fine.LoanID,
		"student_id":   fine.StudentID,
		"amount":       fine.Amount.Round(2).InexactFloat64(),
		"status":       string(fine.Status),
		"created_at":   fine.CreatedAt.UTC().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("fine activity_log metadata: %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: "fine_created",
		Title:     "Multa creada",
		Metadata:  meta,
	}
	desc := "Se registró una multa pendiente"
	entry.Description = &desc
	sid := fine.StudentID
	entry.StudentID = &sid
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("fine activity_log create: %v", err)
	}
}

func (s *Service) logFinePaid(ctx context.Context, fine *models.Fine) {
	if s.activityLog == nil || fine == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"fine_id":    fine.ID,
		"loan_id":    fine.LoanID,
		"student_id": fine.StudentID,
		"amount":     fine.Amount.Round(2).InexactFloat64(),
		"status":     string(fine.Status),
		"paid_at": func() string {
			if fine.PaidAt == nil {
				return ""
			}
			return fine.PaidAt.UTC().Format(time.RFC3339)
		}(),
	})
	if err != nil {
		log.Printf("fine activity_log metadata (paid): %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: "fine_paid",
		Title:     "Multa pagada",
		Metadata:  meta,
	}
	desc := "La multa fue marcada como pagada"
	entry.Description = &desc
	sid := fine.StudentID
	entry.StudentID = &sid
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("fine activity_log create (paid): %v", err)
	}
}

func (s *Service) logFineWaived(ctx context.Context, fine *models.Fine) {
	if s.activityLog == nil || fine == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"fine_id":    fine.ID,
		"loan_id":    fine.LoanID,
		"student_id": fine.StudentID,
		"amount":     fine.Amount.Round(2).InexactFloat64(),
		"status":     string(fine.Status),
		"waived_at":  time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("fine activity_log metadata (waived): %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: "fine_waived",
		Title:     "Multa condonada",
		Metadata:  meta,
	}
	desc := "La multa fue condonada"
	entry.Description = &desc
	sid := fine.StudentID
	entry.StudentID = &sid
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("fine activity_log create (waived): %v", err)
	}
}
