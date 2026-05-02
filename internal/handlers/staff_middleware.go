package handlers

import (
	"net/http"

	"library_back/internal/models"
	"library_back/internal/repositories"
	"library_back/internal/services"

	"github.com/gin-gonic/gin"
)

// StaffBearerMiddleware requires a valid access JWT and an active staff user
// with role admin or librarian (same API permissions for /api routes).
func StaffBearerMiddleware(tokens *services.TokenService, users repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidToken})
			return
		}
		claims, err := tokens.ParseAccess(c.Request.Context(), raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidToken})
			return
		}
		u, err := users.FindByID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
			return
		}
		if u == nil || !u.IsActive || !models.IsStaffRole(u.Role) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: msgInvalidToken})
			return
		}
		c.Set(ContextAuthUserID, claims.UserID)
		c.Next()
	}
}
