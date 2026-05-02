package handlers

import (
	"net/http"

	"library_back/internal/models"
	"library_back/internal/pkg/ratelimit"

	"github.com/gin-gonic/gin"
)

const msgTooManyRequests = "Demasiadas solicitudes. Espere un momento e intente de nuevo."

// IPRateLimit enforces per-IP request limits using ratelimit.Store (always skips /health).
func IPRateLimit(store *ratelimit.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if store == nil || store.ShouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}
		if !store.Allow(c.Request.URL.Path, c.ClientIP()) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorResponse{Message: msgTooManyRequests})
			return
		}
		c.Next()
	}
}
