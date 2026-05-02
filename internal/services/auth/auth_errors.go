package auth

import (
	"errors"

	"library_back/internal/pkg/tokenerr"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserInactive       = errors.New("user inactive")
	ErrTokenRevoked       = tokenerr.ErrTokenRevoked
	ErrTokenInvalid       = tokenerr.ErrTokenInvalid
)
