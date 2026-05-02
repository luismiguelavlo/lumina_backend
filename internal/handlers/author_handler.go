package handlers

import (
	"context"
	"errors"
	"net/http"

	"library_back/internal/models"
	"library_back/internal/services/author"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authorServiceAPI interface {
	Create(ctx context.Context, req models.CreateAuthorRequest) (*models.AuthorCreatedResponse, error)
}

// AuthorHandler exposes staff catalog routes for authors.
type AuthorHandler struct {
	svc authorServiceAPI
}

// NewAuthorHandler constructs AuthorHandler.
func NewAuthorHandler(svc *author.Service) *AuthorHandler {
	return &AuthorHandler{svc: svc}
}

// Create POST /api/authors
func (h *AuthorHandler) Create(c *gin.Context) {
	var req models.CreateAuthorRequest
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
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusCreated, out)
}

var _ authorServiceAPI = (*author.Service)(nil)
