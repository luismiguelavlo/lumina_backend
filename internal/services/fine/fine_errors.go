package fine

import "errors"

var (
	// ErrFineNotFound is returned when a fine ID does not exist.
	ErrFineNotFound = errors.New("fine not found")
	// ErrFineNotPending is returned when marking paid/waived on a non-pending fine.
	ErrFineNotPending = errors.New("fine is not pending")
	// ErrFineLoanNotFound is returned when the loan_id does not exist.
	ErrFineLoanNotFound = errors.New("loan not found")
	// ErrLoanStudentMismatch is returned when loan.student_id != request student_id.
	ErrLoanStudentMismatch = errors.New("loan does not belong to the given student")
	// ErrInvalidFineStatus is returned for an invalid status query param.
	ErrInvalidFineStatus = errors.New("invalid fine status filter")
	// ErrInvalidStudentIDQuery is returned when student_id query is not a valid UUID.
	ErrInvalidStudentIDQuery = errors.New("invalid student_id query")
	// ErrFineDuplicatePending is returned when a pending fine already exists for the loan.
	ErrFineDuplicatePending = errors.New("pending fine already exists for this loan")
)
