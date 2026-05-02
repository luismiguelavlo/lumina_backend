package loan

import "errors"

var (
	ErrLoanNotFound          = errors.New("loan not found")
	ErrLoanAlreadyReturned   = errors.New("loan already returned")
	ErrStudentNotFound       = errors.New("student not found")
	ErrBookNotFound          = errors.New("book not found")
	ErrDueDateInPast         = errors.New("due_date cannot be in the past")
	ErrNoCopiesAvailable     = errors.New("no copies available for loan")
	ErrInvalidLoanStatus     = errors.New("invalid loan status filter")
	ErrInvalidDueDateFormat  = errors.New("due_date must be YYYY-MM-DD")
	ErrInvalidStudentIDQuery = errors.New("invalid student_id query parameter")
	ErrStudentHasActiveSanction = errors.New("student has an active sanction and cannot borrow")
)
