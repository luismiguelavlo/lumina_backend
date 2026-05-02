package sanction

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

type activeStudentReader interface {
	GetActiveByID(ctx context.Context, id string) (*models.Student, error)
}

// sanctionReputationHook is optional; when set, Create records a sanction_received reputation penalty once per sanction.
type sanctionReputationHook interface {
	RecordSanctionReceived(ctx context.Context, sanctionID, studentID string)
}

// Service handles sanctions for staff API.
type Service struct {
	sanctions repositories.SanctionRepository
	students  activeStudentReader
	rep       sanctionReputationHook
	activityLog repositories.ActivityLogRepository
}

// NewService constructs SanctionService. rep may be nil (tests).
func NewService(sanctions repositories.SanctionRepository, students activeStudentReader, rep sanctionReputationHook) *Service {
	return &Service{sanctions: sanctions, students: students, rep: rep}
}

// SetActivityLog configures optional activity log persistence for dashboard feed events.
func (s *Service) SetActivityLog(activityLog repositories.ActivityLogRepository) {
	s.activityLog = activityLog
}

// ListActive returns paginated active sanctions with student fields.
func (s *Service) ListActive(ctx context.Context, filter repositories.SanctionListFilter, limit, offset int) ([]models.SanctionListItem, int64, error) {
	limit, offset = normalizeListParams(limit, offset)
	if filter.StudentID != "" {
		if _, err := uuid.Parse(filter.StudentID); err != nil {
			return nil, 0, ErrInvalidStudentIDQuery
		}
	}
	rows, total, err := s.sanctions.ListActive(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.SanctionListItem, 0, len(rows))
	for i := range rows {
		out = append(out, rowToListItem(&rows[i]))
	}
	return out, total, nil
}

// Create registers an active sanction for an existing active student.
func (s *Service) Create(ctx context.Context, req models.CreateSanctionRequest, adminID string) (*models.SanctionResponse, error) {
	st, err := s.students.GetActiveByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}
	ab := adminID
	san := &models.Sanction{
		StudentID: req.StudentID,
		Reason:    req.Reason,
		Status:    models.SanctionStatusActive,
		AppliedBy: &ab,
	}
	if err := s.sanctions.Create(ctx, san); err != nil {
		return nil, err
	}
	refreshed, err := s.sanctions.GetByID(ctx, san.ID)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrSanctionNotFound
	}
	if s.rep != nil {
		s.rep.RecordSanctionReceived(ctx, refreshed.ID, refreshed.StudentID)
	}
	s.logSanctionCreated(ctx, refreshed, adminID)
	return toResponse(st, refreshed), nil
}

// Lift marks an active sanction as lifted.
func (s *Service) Lift(ctx context.Context, sanctionID, adminID string) (*models.SanctionResponse, error) {
	san, err := s.sanctions.GetByID(ctx, sanctionID)
	if err != nil {
		return nil, err
	}
	if san == nil || san.Status != models.SanctionStatusActive {
		return nil, ErrSanctionNotFound
	}
	now := time.Now().UTC()
	n, err := s.sanctions.Lift(ctx, sanctionID, now, adminID)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrSanctionNotFound
	}
	refreshed, err := s.sanctions.GetByID(ctx, sanctionID)
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, ErrSanctionNotFound
	}
	st, err := s.students.GetActiveByID(ctx, refreshed.StudentID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}
	s.logSanctionLifted(ctx, refreshed, adminID)
	return toResponse(st, refreshed), nil
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

func rowToListItem(row *repositories.SanctionListRow) models.SanctionListItem {
	return models.SanctionListItem{
		SanctionID:    row.SanctionID,
		StudentID:     row.StudentID,
		StudentIDCode: row.StudentIDCode,
		FirstName:     row.FirstName,
		LastName:      row.LastName,
		Email:         row.Email,
		Reason:        row.Reason,
		AppliedAt:     row.AppliedAt.UTC().Format(time.RFC3339Nano),
		AppliedByID:   row.AppliedByID,
	}
}

func toResponse(st *models.Student, san *models.Sanction) *models.SanctionResponse {
	var liftedAt, liftedBy *string
	if san.LiftedAt != nil {
		s := san.LiftedAt.UTC().Format(time.RFC3339Nano)
		liftedAt = &s
	}
	if san.LiftedBy != nil {
		lb := *san.LiftedBy
		liftedBy = &lb
	}
	return &models.SanctionResponse{
		SanctionID:    san.ID,
		StudentID:     st.ID,
		StudentIDCode: st.StudentIDCode,
		FirstName:     st.FirstName,
		LastName:      st.LastName,
		Email:         st.Email,
		Reason:        san.Reason,
		Status:        string(san.Status),
		AppliedAt:     san.AppliedAt.UTC().Format(time.RFC3339Nano),
		AppliedByID:   san.AppliedBy,
		LiftedAt:      liftedAt,
		LiftedByID:    liftedBy,
	}
}

func (s *Service) logSanctionCreated(ctx context.Context, sanction *models.Sanction, actorUserID string) {
	if s.activityLog == nil || sanction == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"sanction_id": sanction.ID,
		"student_id":  sanction.StudentID,
		"reason":      sanction.Reason,
		"status":      string(sanction.Status),
		"applied_at":  sanction.AppliedAt.UTC().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("sanction activity_log metadata: %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: "sanction_created",
		Title:     "Sanción aplicada",
		Metadata:  meta,
	}
	desc := "Se aplicó una sanción al estudiante"
	entry.Description = &desc
	sid := sanction.StudentID
	entry.StudentID = &sid
	aid := strings.TrimSpace(actorUserID)
	if aid != "" {
		entry.ActorID = &aid
	}
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("sanction activity_log create: %v", err)
	}
}

func (s *Service) logSanctionLifted(ctx context.Context, sanction *models.Sanction, actorUserID string) {
	if s.activityLog == nil || sanction == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"sanction_id": sanction.ID,
		"student_id":  sanction.StudentID,
		"reason":      sanction.Reason,
		"status":      string(sanction.Status),
		"lifted_at": func() string {
			if sanction.LiftedAt == nil {
				return ""
			}
			return sanction.LiftedAt.UTC().Format(time.RFC3339)
		}(),
	})
	if err != nil {
		log.Printf("sanction activity_log metadata (lifted): %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: "sanction_lifted",
		Title:     "Sanción levantada",
		Metadata:  meta,
	}
	desc := "Se levantó la sanción del estudiante"
	entry.Description = &desc
	sid := sanction.StudentID
	entry.StudentID = &sid
	aid := strings.TrimSpace(actorUserID)
	if aid != "" {
		entry.ActorID = &aid
	}
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("sanction activity_log create (lifted): %v", err)
	}
}
