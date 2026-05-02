package student

import "errors"

var (
	ErrStudentNotFound           = errors.New("student not found")
	ErrDepartmentNotFound        = errors.New("department not found")
	ErrDuplicateStudentIDCode    = errors.New("student_id_code already exists")
	ErrStudentIDGenerationFailed = errors.New("could not generate unique student_id_code")
	ErrDuplicateStudentEmail     = errors.New("email already exists")
	ErrBadgeNotFound             = errors.New("badge not found")
	ErrBadgeAlreadyEarned        = errors.New("student already has this badge")
	ErrStudentDoesNotHaveBadge   = errors.New("student does not have this badge")
)
