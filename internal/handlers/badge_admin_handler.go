package handlers

import (
	"context"
	"errors"
	"net/http"

	"library_back/internal/models"
	"library_back/internal/services/badgedef"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	msgBadgeDefNotFound  = "Insignia no encontrada"
	msgBadgeDefInUse     = "La insignia ya fue otorgada a estudiantes y no se puede eliminar"
	msgBadgeDefSlugTaken = "El slug de la insignia ya existe"
)

type badgeDefServiceAPI interface {
	List(ctx context.Context) ([]models.BadgeDefinitionResponse, error)
	GetByID(ctx context.Context, id string) (*models.BadgeDefinitionResponse, error)
	Create(ctx context.Context, req models.CreateBadgeDefinitionRequest) (*models.BadgeDefinitionResponse, error)
	Patch(ctx context.Context, id string, req models.PatchBadgeDefinitionRequest) (*models.BadgeDefinitionResponse, error)
	Delete(ctx context.Context, id string) error
}

// BadgeAdminHandler exposes CRUD for badge definitions under /api/badges.
type BadgeAdminHandler struct {
	svc badgeDefServiceAPI
}

// NewBadgeAdminHandler constructs BadgeAdminHandler.
func NewBadgeAdminHandler(svc *badgedef.Service) *BadgeAdminHandler {
	return &BadgeAdminHandler{svc: svc}
}

// List GET /api/badges
func (h *BadgeAdminHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, models.ListBadgeDefinitionsResponse{Data: items})
}

// Get GET /api/badges/:id
func (h *BadgeAdminHandler) Get(c *gin.Context) {
	out, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeBadgeDefError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create POST /api/badges
func (h *BadgeAdminHandler) Create(c *gin.Context) {
	var req models.CreateBadgeDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	out, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		writeBadgeDefError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Patch PATCH /api/badges/:id
func (h *BadgeAdminHandler) Patch(c *gin.Context) {
	var req models.PatchBadgeDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	out, err := h.svc.Patch(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		writeBadgeDefError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete DELETE /api/badges/:id
func (h *BadgeAdminHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		writeBadgeDefError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeBadgeDefError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, badgedef.ErrNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgBadgeDefNotFound})
	case errors.Is(err, badgedef.ErrInUse):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgBadgeDefInUse})
	case errors.Is(err, badgedef.ErrSlugTaken):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgBadgeDefSlugTaken, Errors: map[string]string{"slug": msgBadgeDefSlugTaken}})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
	}
}
