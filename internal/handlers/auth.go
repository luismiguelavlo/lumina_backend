package handlers

import (
	"errors"
	"library_back/internal/services/auth"
	"net/http"

	"library_back/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	msgInvalidCredentials = "Credenciales inválidas"
	msgInvalidToken       = "Token inválido"
	msgEmailConflict      = "El recurso ya existe"
	msgInternal           = "Error interno"
	msgUserNotFound       = "Usuario no encontrado"
)

// AuthHandler exposes /auth routes.
type AuthHandler struct {
	auth *auth.AuthService
}

// NewAuthHandler constructs AuthHandler.
func NewAuthHandler(auth *auth.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	user, err := h.auth.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrEmailExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgEmailConflict})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusCreated, models.RegisterResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
		AvatarURL: user.AvatarURL,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	})
}

// Login POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	loggedUser, access, refresh, exp, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidCredentials})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    exp,
		User:         publicUserResponse(loggedUser),
	})
}

// Refresh POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	refreshedUser, access, refresh, exp, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) ||
			errors.Is(err, auth.ErrTokenInvalid) ||
			errors.Is(err, auth.ErrTokenRevoked) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidToken})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    exp,
		User:         publicUserResponse(refreshedUser),
	})
}

func publicUserResponse(u *models.User) models.RegisterResponse {
	return models.RegisterResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      string(u.Role),
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// Logout POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	token, ok := bearerToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidToken})
		return
	}
	var body models.LogoutRequest
	_ = c.ShouldBindJSON(&body)
	if err := h.auth.Logout(c.Request.Context(), token, body.RefreshToken); err != nil {
		if errors.Is(err, auth.ErrTokenInvalid) || errors.Is(err, auth.ErrTokenRevoked) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidToken})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.Status(http.StatusNoContent)
}

// Me GET /auth/me — requires AuthBearerMiddleware; returns public user profile.
func (h *AuthHandler) Me(c *gin.Context) {
	v, ok := c.Get(ContextAuthUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	userID, _ := v.(string)
	if userID == "" {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	u, err := h.auth.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, publicUserResponse(u))
}
