package repositories

import "errors"

// ErrDuplicateEmail is returned when a user insert violates the unique email constraint.
var ErrDuplicateEmail = errors.New("duplicate email")
