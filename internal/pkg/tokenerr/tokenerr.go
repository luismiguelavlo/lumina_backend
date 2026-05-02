package tokenerr

import "errors"

// Sentinel errors returned by TokenService parse/validate paths.
var (
	ErrTokenRevoked = errors.New("token revoked")
	ErrTokenInvalid = errors.New("token invalid")
)
