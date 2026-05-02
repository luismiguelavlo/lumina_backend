package models

import "time"

// RevokedToken represents a revoked JWT (jti blacklist).
type RevokedToken struct {
	JTI       string    `json:"jti" gorm:"type:varchar(36);primaryKey;comment:Identificador único del JWT jti usado en blacklist de logout"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index;comment:Usuario dueño del token revocado"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index;comment:Instante en que el JWT dejaría de ser válido"`
	RevokedAt time.Time `json:"revoked_at" gorm:"not null;index;comment:Instante en que se registró la revocación"`
}

// TableName overrides GORM pluralization for revoked_tokens.
func (RevokedToken) TableName() string {
	return "revoked_tokens"
}
