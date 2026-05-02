package sanction

import "errors"

var (
	// ErrSanctionNotFound is returned when the sanction does not exist or cannot be lifted (e.g. already lifted).
	ErrSanctionNotFound = errors.New("sanction not found")
	// ErrStudentNotFound is returned when the target student does not exist or is inactive.
	ErrStudentNotFound = errors.New("student not found")
	// ErrInvalidStudentIDQuery is returned when student_id query is not a valid UUID.
	ErrInvalidStudentIDQuery = errors.New("invalid student_id query")
)
