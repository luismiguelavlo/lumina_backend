package handlers

import (
	"log"
	"net/http"
	"os"

	"library_back/internal/models"
	"library_back/internal/seeds"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SeedHandler runs optional reference data seed (no JWT).
type SeedHandler struct {
	db *gorm.DB
}

// NewSeedHandler constructs SeedHandler.
func NewSeedHandler(db *gorm.DB) *SeedHandler {
	return &SeedHandler{db: db}
}

// Run POST /setup/seed — fills departments, genres, authors, demo books, ensures badges.
// Requires environment ENABLE_PUBLIC_SEED=true (no JWT; protect production by leaving unset).
func (h *SeedHandler) Run(c *gin.Context) {
	if os.Getenv("ENABLE_PUBLIC_SEED") != "true" {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Semilla pública desactivada. Define ENABLE_PUBLIC_SEED=true en el entorno solo en desarrollo.",
		})
		return
	}
	sum, err := seeds.Run(c.Request.Context(), h.db)
	if err != nil {
		log.Printf("seed: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Semilla aplicada (idempotente).",
		"summary": sum,
	})
}
