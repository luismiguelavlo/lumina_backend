package reputation

import (
	"context"
	"fmt"
	"log"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

const (
	pointsReturnedOnTime   = 10
	pointsReturnedLate     = -5
	pointsFinePaid         = 5
	pointsSanctionReceived = -10
	pointsBadgeEarned      = 15
	pointsStreakMilestone  = 3
)

// Service records reputation events and updates student_stats (invoked from LoanService).
type Service struct {
	store repositories.ReputationStore
}

// NewService constructs ReputationService.
func NewService(store repositories.ReputationStore) *Service {
	return &Service{store: store}
}

// RecordLoanCreated increments active_loans for the student. Errors are logged only.
func (s *Service) RecordLoanCreated(ctx context.Context, studentID string) {
	if err := s.store.IncrementStudentActiveLoans(ctx, studentID, 1); err != nil {
		log.Printf("reputation RecordLoanCreated: %v", err)
	}
}

// RecordReturn creates a return event and updates stats if not already done for this loan.
func (s *Service) RecordReturn(ctx context.Context, loan *models.Loan) {
	if loan == nil || loan.ReturnedAt == nil || loan.Status != models.LoanStatusReturned {
		return
	}
	exists, err := s.store.ExistsReturnEventForLoan(ctx, loan.ID)
	if err != nil {
		log.Printf("reputation RecordReturn exists check: %v", err)
		return
	}
	if exists {
		return
	}

	onTime := returnedOnTime(loan.DueDate, *loan.ReturnedAt)
	points := pointsReturnedOnTime
	evType := models.ReputationEventBookReturnedOnTime
	if !onTime {
		points = pointsReturnedLate
		evType = models.ReputationEventBookReturnedLate
	}

	ref := loan.ID
	desc := "Book return"
	ev := &models.ReputationEvent{
		ID:          uuid.NewString(),
		StudentID:   loan.StudentID,
		EventType:   evType,
		Points:      points,
		Description: &desc,
		ReferenceID: &ref,
	}
	if err := s.store.CreateReputationEvent(ctx, ev); err != nil {
		log.Printf("reputation RecordReturn create event: %v", err)
		return
	}
	newStreak, err := s.store.ApplyStudentReturn(ctx, loan.StudentID, points, onTime)
	if err != nil {
		log.Printf("reputation RecordReturn stats: %v", err)
		return
	}
	if onTime {
		s.tryStreakMilestones(ctx, loan.StudentID, newStreak)
	}
}

// RecordFinePaid logs a fine_paid reputation delta (idempotent per fine_id).
func (s *Service) RecordFinePaid(ctx context.Context, fineID, studentID string) {
	if s == nil || s.store == nil || fineID == "" || studentID == "" {
		return
	}
	desc := "Multa pagada"
	if err := s.store.CreateReputationEventWithPointsIfNew(ctx, studentID, pointsFinePaid, models.ReputationEventFinePaid, fineID, &desc); err != nil {
		log.Printf("reputation RecordFinePaid: %v", err)
	}
}

// RecordSanctionReceived logs a sanction_received penalty (idempotent per sanction_id).
func (s *Service) RecordSanctionReceived(ctx context.Context, sanctionID, studentID string) {
	if s == nil || s.store == nil || sanctionID == "" || studentID == "" {
		return
	}
	desc := "Sanción registrada"
	if err := s.store.CreateReputationEventWithPointsIfNew(ctx, studentID, pointsSanctionReceived, models.ReputationEventSanctionReceived, sanctionID, &desc); err != nil {
		log.Printf("reputation RecordSanctionReceived: %v", err)
	}
}

// RecordBadgeEarned logs badge_earned (idempotent per student_badge row id).
func (s *Service) RecordBadgeEarned(ctx context.Context, studentBadgeID, studentID string) {
	if s == nil || s.store == nil || studentBadgeID == "" || studentID == "" {
		return
	}
	desc := "Insignia obtenida"
	if err := s.store.CreateReputationEventWithPointsIfNew(ctx, studentID, pointsBadgeEarned, models.ReputationEventBadgeEarned, studentBadgeID, &desc); err != nil {
		log.Printf("reputation RecordBadgeEarned: %v", err)
	}
}

func streakMilestoneReferenceID(studentID string, days int) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("streak:%s:%d", studentID, days))).String()
}

func (s *Service) tryStreakMilestones(ctx context.Context, studentID string, newStreak int) {
	for _, m := range []int{7, 14, 30} {
		if newStreak != m {
			continue
		}
		ref := streakMilestoneReferenceID(studentID, m)
		desc := fmt.Sprintf("Racha de devoluciones a tiempo: %d días", m)
		if err := s.store.CreateReputationEventWithPointsIfNew(ctx, studentID, pointsStreakMilestone, models.ReputationEventStreakMilestone, ref, &desc); err != nil {
			log.Printf("reputation streak milestone: %v", err)
		}
	}
}

func returnedOnTime(due time.Time, returned time.Time) bool {
	dueDay := time.Date(due.Year(), due.Month(), due.Day(), 0, 0, 0, 0, time.UTC)
	retDay := time.Date(returned.Year(), returned.Month(), returned.Day(), 0, 0, 0, 0, time.UTC)
	return !retDay.After(dueDay)
}
