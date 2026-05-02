package models

import "time"

// RegisterRequest is the JSON body for POST /auth/register.
type RegisterRequest struct {
	FirstName string  `json:"first_name" binding:"required,max=100"`
	LastName  string  `json:"last_name" binding:"required,max=100"`
	Email     string  `json:"email" binding:"required,email,max=255"`
	Password  string  `json:"password" binding:"required,min=8"`
	Role      *string `json:"role" binding:"omitempty,oneof=admin librarian"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty,url"`
}

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest is the JSON body for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest is the optional JSON body for POST /auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"omitempty"`
}

// LoginResponse returns JWT pair and the authenticated user after login or refresh.
type LoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int              `json:"expires_in"`
	User         RegisterResponse `json:"user"`
}

// RegisterResponse is the public user payload after registration (no password).
type RegisterResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse is the structured error payload for 400 and similar.
type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}
