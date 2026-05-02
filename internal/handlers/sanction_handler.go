package handlers

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"library_back/internal/models"
	"library_back/internal/repositories"
	"library_back/internal/services/sanction"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const msgSanctionNotFound = "Sanción no encontrada"

type sanctionServiceAPI interface {
	ListActive(ctx context.Context, filter repositories.SanctionListFilter, limit, offset int) ([]models.SanctionListItem, int64, error)
	Create(ctx context.Context, req models.CreateSanctionRequest, adminID string) (*models.SanctionResponse, error)
	Lift(ctx context.Context, sanctionID, adminID string) (*models.SanctionResponse, error)
}

// SanctionHandler exposes /api/sanctions.
type SanctionHandler struct {
	svc sanctionServiceAPI
}

// NewSanctionHandler constructs SanctionHandler.
func NewSanctionHandler(svc *sanction.Service) *SanctionHandler {
	return &SanctionHandler{svc: svc}
}

// List GET /api/sanctions — active sanctions only.
func (h *SanctionHandler) List(c *gin.Context) {
	limit, offset := normalizeSanctionListQuery(c.Query("limit"), c.Query("offset"))
	filter := repositories.SanctionListFilter{
		StudentID: strings.TrimSpace(c.Query("student_id")),
	}
	items, total, err := h.svc.ListActive(c.Request.Context(), filter, limit, offset)
	if err != nil {
		writeSanctionServiceError(c, err)
		return
	}
	count := len(items)
	currentPage := 1
	if limit > 0 {
		currentPage = (offset / limit) + 1
	}
	totalPages := 0
	if limit > 0 && total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
	c.JSON(http.StatusOK, models.ListSanctionsResponse{
		Data:  items,
		Total: total,
		Pagination: models.PaginationMeta{
			Limit:       limit,
			Offset:      offset,
			Count:       count,
			Total:       total,
			CurrentPage: currentPage,
			TotalPages:  totalPages,
			HasNext:     int64(offset+count) < total,
			HasPrev:     offset > 0,
		},
	})
}

// Create POST /api/sanctions
func (h *SanctionHandler) Create(c *gin.Context) {
	var req models.CreateSanctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	adminID, ok := c.Get(ContextAuthUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	aid, _ := adminID.(string)
	out, err := h.svc.Create(c.Request.Context(), req, aid)
	if err != nil {
		writeSanctionServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Lift PATCH /api/sanctions/:id — optional body `{"status":"lifted"}`.
func (h *SanctionHandler) Lift(c *gin.Context) {
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		var body models.LiftSanctionRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				writeValidationErrors(c, errs)
				return
			}
			writeJSONBindError(c)
			return
		}
	}
	adminID, ok := c.Get(ContextAuthUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	aid, _ := adminID.(string)
	out, err := h.svc.Lift(c.Request.Context(), c.Param("id"), aid)
	if err != nil {
		writeSanctionServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func writeSanctionServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, sanction.ErrStudentNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgStudentNotFound})
	case errors.Is(err, sanction.ErrSanctionNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgSanctionNotFound})
	case errors.Is(err, sanction.ErrInvalidStudentIDQuery):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Parámetro student_id no válido",
			Errors:  map[string]string{"student_id": "UUID inválido"},
		})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
	}
}

func normalizeSanctionListQuery(limitStr, offsetStr string) (limit, offset int) {
	limit = 20
	offset = 0
	if limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil {
			if n > 100 {
				limit = 100
			} else if n > 0 {
				limit = n
			}
		}
	}
	if offsetStr != "" {
		if n, err := strconv.Atoi(offsetStr); err == nil && n > 0 {
			offset = n
		}
	}
	return limit, offset
}
