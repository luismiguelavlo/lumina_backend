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
	"library_back/internal/services/fine"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	msgFineNotFound            = "Multa no encontrada"
	msgFineNotPending          = "La multa no está pendiente"
	msgFineLoanNotFound        = "Préstamo no encontrado"
	msgLoanStudentMismatch     = "El préstamo no corresponde al estudiante indicado"
	msgFineDuplicatePending    = "Ya existe una multa pendiente para este préstamo"
	msgInvalidFineStatus       = "Estado de multa no válido"
	msgInvalidStudentIDQuery   = "Parámetro student_id no válido"
)

type fineServiceAPI interface {
	List(ctx context.Context, filter repositories.FineListFilter, limit, offset int) ([]models.FineListItem, int64, error)
	Create(ctx context.Context, req models.CreateFineRequest) (*models.FineResponse, error)
	GetByID(ctx context.Context, id string) (*models.FineResponse, error)
	MarkPaid(ctx context.Context, id string) (*models.FineResponse, error)
	MarkWaived(ctx context.Context, id string) (*models.FineResponse, error)
}

// FineHandler exposes /api/fines.
type FineHandler struct {
	svc fineServiceAPI
}

// NewFineHandler constructs FineHandler.
func NewFineHandler(svc *fine.Service) *FineHandler {
	return &FineHandler{svc: svc}
}

// List GET /api/fines
func (h *FineHandler) List(c *gin.Context) {
	limit, offset := normalizeFineListQuery(c.Query("limit"), c.Query("offset"))
	filter := repositories.FineListFilter{
		StudentID: strings.TrimSpace(c.Query("student_id")),
		Status:    strings.TrimSpace(c.Query("status")),
	}
	items, total, err := h.svc.List(c.Request.Context(), filter, limit, offset)
	if err != nil {
		writeFineServiceError(c, err)
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
	c.JSON(http.StatusOK, models.ListFinesResponse{
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

// Create POST /api/fines
func (h *FineHandler) Create(c *gin.Context) {
	var req models.CreateFineRequest
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
		writeFineServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Get GET /api/fines/:id
func (h *FineHandler) Get(c *gin.Context) {
	out, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeFineServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// MarkPaid PATCH /api/fines/:id/paid
func (h *FineHandler) MarkPaid(c *gin.Context) {
	out, err := h.svc.MarkPaid(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeFineServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// MarkWaived PATCH /api/fines/:id/waived
func (h *FineHandler) MarkWaived(c *gin.Context) {
	out, err := h.svc.MarkWaived(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeFineServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func writeFineServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, fine.ErrFineLoanNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgFineLoanNotFound})
	case errors.Is(err, fine.ErrFineNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgFineNotFound})
	case errors.Is(err, fine.ErrFineNotPending):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgFineNotPending})
	case errors.Is(err, fine.ErrLoanStudentMismatch):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgLoanStudentMismatch,
			Errors:  map[string]string{"student_id": msgLoanStudentMismatch},
		})
	case errors.Is(err, fine.ErrInvalidFineStatus):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgInvalidFineStatus,
			Errors:  map[string]string{"status": msgInvalidFineStatus},
		})
	case errors.Is(err, fine.ErrInvalidStudentIDQuery):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgInvalidStudentIDQuery,
			Errors:  map[string]string{"student_id": "UUID inválido"},
		})
	case errors.Is(err, fine.ErrFineDuplicatePending):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgFineDuplicatePending})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
	}
}

func normalizeFineListQuery(limitStr, offsetStr string) (limit, offset int) {
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
