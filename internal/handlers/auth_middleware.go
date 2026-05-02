package handlers

import (
	"net/http"

	"library_back/internal/models"
	"library_back/internal/services"

	"github.com/gin-gonic/gin"
)

// ContextAuthUserID is the Gin context key for the authenticated user ID (JWT sub).
const ContextAuthUserID = "auth_user_id"

// AuthBearerMiddleware validates the access JWT (signature, expiry, type, blacklist)
// and stores the user ID in the context under ContextAuthUserID.
func AuthBearerMiddleware(tokens *services.TokenService) gin.HandlerFunc {
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
		c.Set(ContextAuthUserID, claims.UserID)
		c.Next()
	}
}
