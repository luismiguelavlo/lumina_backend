package services

import (
	"context"
	"fmt"
	"time"

	"library_back/internal/pkg/jwths256"
	"library_back/internal/pkg/tokenerr"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
	tokenTTL         = 10 * time.Hour
)

// TokenService issues and validates JWT access/refresh pairs (HS256).
type TokenService struct {
	secret  []byte
	revoked repositories.RevokedTokenRepository
}

// NewTokenService builds a TokenService. secret must be at least 32 bytes.
func NewTokenService(secret string, revoked repositories.RevokedTokenRepository) (*TokenService, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	return &TokenService{secret: []byte(secret), revoked: revoked}, nil
}

// ParsedToken is the validated JWT payload used by handlers and AuthService.
type ParsedToken struct {
	UserID    string
	JTI       string
	ExpiresAt time.Time
	TokenType string
}

// GeneratePair returns access and refresh JWTs and expires_in seconds.
func (s *TokenService) GeneratePair(ctx context.Context, userID string) (access, refresh string, expiresIn int, err error) {
	_ = ctx
	expiresIn = int(tokenTTL.Seconds())
	now := time.Now().UTC()
	expUnix := now.Add(tokenTTL).Unix()
	iatUnix := now.Unix()

	makeTok := func(typ string) (string, error) {
		c := jwths256.Claims{
			Typ: typ,
			Sub: userID,
			Jti: uuid.NewString(),
			Exp: expUnix,
			Iat: iatUnix,
			Nbf: iatUnix,
		}
		return jwths256.Sign(s.secret, c)
	}
	access, err = makeTok(tokenTypeAccess)
	if err != nil {
		return "", "", 0, err
	}
	refresh, err = makeTok(tokenTypeRefresh)
	if err != nil {
		return "", "", 0, err
	}
	return access, refresh, expiresIn, nil
}

// ParseAccess validates an access token (signature, expiry, type, blacklist).
func (s *TokenService) ParseAccess(ctx context.Context, tokenStr string) (*ParsedToken, error) {
	return s.parse(ctx, tokenStr, tokenTypeAccess, true)
}

// ParseRefresh validates a refresh token (signature, expiry, type, blacklist).
func (s *TokenService) ParseRefresh(ctx context.Context, tokenStr string) (*ParsedToken, error) {
	return s.parse(ctx, tokenStr, tokenTypeRefresh, true)
}

// ParseRefreshIgnoreExpiry validates signature and type only (logout of expired refresh).
func (s *TokenService) ParseRefreshIgnoreExpiry(ctx context.Context, tokenStr string) (*ParsedToken, error) {
	return s.parse(ctx, tokenStr, tokenTypeRefresh, false)
}

func (s *TokenService) parse(ctx context.Context, tokenStr, wantType string, validateExp bool) (*ParsedToken, error) {
	c, err := jwths256.Parse(s.secret, tokenStr, validateExp)
	if err != nil {
		return nil, tokenerr.ErrTokenInvalid
	}
	if c.Typ != wantType || c.Jti == "" || c.Sub == "" {
		return nil, tokenerr.ErrTokenInvalid
	}
	revoked, err := s.revoked.IsRevoked(ctx, c.Jti)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, tokenerr.ErrTokenRevoked
	}
	return &ParsedToken{
		UserID:    c.Sub,
		JTI:       c.Jti,
		ExpiresAt: time.Unix(c.Exp, 0).UTC(),
		TokenType: c.Typ,
	}, nil
}

// Revoke inserts a single JTI into the blacklist.
func (s *TokenService) Revoke(ctx context.Context, userID, jti string, expiresAt time.Time) error {
	return s.revoked.Revoke(ctx, jti, userID, expiresAt)
}
