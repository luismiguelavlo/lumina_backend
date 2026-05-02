package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health checks.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler returns a new HealthHandler. Pass nil for db if DATABASE_URL is not used.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check responds with service status and optional database connectivity.
func (h *HealthHandler) Check(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "not_configured",
		})
		return
	}
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"database": "unreachable",
			"error":    err.Error(),
		})
		return
	}
	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"database": "unreachable",
			"error":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "ok",
	})
}
