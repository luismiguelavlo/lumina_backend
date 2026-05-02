package auth

import (
	"context"
	"errors"
	"library_back/internal/models"
	"library_back/internal/pkg/crypto"
	"library_back/internal/repositories"
	"library_back/internal/services"

	"github.com/google/uuid"
)

// AuthService handles registration, login, and refresh.
type AuthService struct {
	users  repositories.UserRepository
	tokens *services.TokenService
}

// NewAuthService builds AuthService.
func NewAuthService(users repositories.UserRepository, tokens *services.TokenService) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

// Register creates an inactive user with hashed password.
func (a *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	existing, err := a.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}
	hash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	role := models.UserRoleAdmin
	if req.Role != nil && *req.Role != "" {
		role = models.UserRole(*req.Role)
	}
	u := &models.User{
		ID:           uuid.NewString(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         role,
		AvatarURL:    req.AvatarURL,
		IsActive:     false,
	}
	if err := a.users.Create(ctx, u); err != nil {
		if errors.Is(err, repositories.ErrDuplicateEmail) {
			return nil, ErrEmailExists
		}
		return nil, err
	}
	created, err := a.users.FindByID(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	if created != nil {
		return created, nil
	}
	return u, nil
}

// Login returns the user and JWT pair for active users with valid password.
func (a *AuthService) Login(ctx context.Context, email, password string) (u *models.User, access, refresh string, expiresIn int, err error) {
	u, err = a.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", 0, err
	}
	if u == nil || !u.IsActive || !crypto.ComparePassword(password, u.PasswordHash) {
		return nil, "", "", 0, ErrInvalidCredentials
	}
	access, refresh, expiresIn, err = a.tokens.GeneratePair(ctx, u.ID)
	if err != nil {
		return nil, "", "", 0, err
	}
	return u, access, refresh, expiresIn, nil
}

// Refresh issues a new token pair and returns the user for a valid refresh token.
func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (u *models.User, access, refresh string, expiresIn int, err error) {
	p, err := a.tokens.ParseRefresh(ctx, refreshToken)
	if err != nil {
		return nil, "", "", 0, err
	}
	u, err = a.users.FindByID(ctx, p.UserID)
	if err != nil {
		return nil, "", "", 0, err
	}
	if u == nil || !u.IsActive {
		return nil, "", "", 0, ErrInvalidCredentials
	}
	access, refresh, expiresIn, err = a.tokens.GeneratePair(ctx, u.ID)
	if err != nil {
		return nil, "", "", 0, err
	}
	return u, access, refresh, expiresIn, nil
}

// GetProfile returns the user by ID (e.g. for GET /auth/me).
func (a *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

// Logout revokes access token (and optional refresh) JTIs.
func (a *AuthService) Logout(ctx context.Context, accessToken, refreshTokenOptional string) error {
	acc, err := a.tokens.ParseAccess(ctx, accessToken)
	if err != nil {
		return err
	}
	if err := a.tokens.Revoke(ctx, acc.UserID, acc.JTI, acc.ExpiresAt); err != nil {
		return err
	}
	if refreshTokenOptional == "" {
		return nil
	}
	ref, err := a.tokens.ParseRefreshIgnoreExpiry(ctx, refreshTokenOptional)
	if err != nil {
		return nil
	}
	if ref.UserID != acc.UserID {
		return nil
	}
	return a.tokens.Revoke(ctx, ref.UserID, ref.JTI, ref.ExpiresAt)
}
