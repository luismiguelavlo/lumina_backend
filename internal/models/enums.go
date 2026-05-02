package models

// UserRole corresponds to enum user_role in the database.
type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleLibrarian UserRole = "librarian"
)

// IsStaffRole reports whether the role may access protected /api staff routes.
// Admin and librarian have the same permissions for /api (see StaffBearerMiddleware).
func IsStaffRole(r UserRole) bool {
	return r == UserRoleAdmin || r == UserRoleLibrarian
}

// LoanStatus corresponds to enum loan_status in the database.
type LoanStatus string

const (
	LoanStatusActive   LoanStatus = "active"
	LoanStatusReturned LoanStatus = "returned"
	LoanStatusOverdue  LoanStatus = "overdue"
)

// FineStatus corresponds to enum fine_status in the database.
type FineStatus string

const (
	FineStatusPending FineStatus = "pending"
	FineStatusPaid    FineStatus = "paid"
	FineStatusWaived  FineStatus = "waived"
)

// SanctionStatus corresponds to enum sanction_status in the database.
type SanctionStatus string

const (
	SanctionStatusActive SanctionStatus = "active"
	SanctionStatusLifted SanctionStatus = "lifted"
)

// ReputationEventType corresponds to enum reputation_event_type in the database.
type ReputationEventType string

const (
	ReputationEventBookReturnedOnTime ReputationEventType = "book_returned_on_time"
	ReputationEventBookReturnedLate   ReputationEventType = "book_returned_late"
	ReputationEventBadgeEarned        ReputationEventType = "badge_earned"
	ReputationEventStreakMilestone    ReputationEventType = "streak_milestone"
	ReputationEventFinePaid           ReputationEventType = "fine_paid"
	ReputationEventSanctionReceived   ReputationEventType = "sanction_received"
)
