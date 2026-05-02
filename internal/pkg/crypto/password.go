package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = bcrypt.DefaultCost

// HashPassword hashes a plain password using bcrypt.
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), defaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword compares a plain password with a bcrypt hash.
// Returns true if they match.
func ComparePassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
