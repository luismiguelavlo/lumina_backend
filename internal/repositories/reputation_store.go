package repositories

import (
	"context"

	"library_back/internal/models"
)

// ReputationStore supports loan-related reputation updates (internal use).
type ReputationStore interface {
	ExistsReturnEventForLoan(ctx context.Context, loanID string) (bool, error)
	CreateReputationEvent(ctx context.Context, ev *models.ReputationEvent) error
	IncrementStudentActiveLoans(ctx context.Context, studentID string, delta int) error
	// ApplyStudentReturn updates student_stats after a book return; returns current_streak_days after the update.
	ApplyStudentReturn(ctx context.Context, studentID string, points int, onTime bool) (newStreakDays int, err error)
	// ExistsReputationEventByTypeAndReference is true when an event with this type and reference_id already exists (idempotency).
	ExistsReputationEventByTypeAndReference(ctx context.Context, typ models.ReputationEventType, referenceID string) (bool, error)
	// CreateReputationEventWithPointsIfNew inserts one ledger event and adjusts total_points if no row exists yet for type+reference.
	CreateReputationEventWithPointsIfNew(ctx context.Context, studentID string, points int, typ models.ReputationEventType, referenceID string, description *string) error
}
