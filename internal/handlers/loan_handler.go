package handlers

import (
	"context"
	"errors"
	"library_back/internal/services/loan"
	"math"
	"net/http"
	"strconv"
	"strings"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// loanServiceAPI is implemented by *loan.Service.
type loanServiceAPI interface {
	List(ctx context.Context, filter repositories.LoanListFilter, limit, offset int) ([]models.LoanListItem, int64, error)
	Create(ctx context.Context, req models.CreateLoanRequest, adminID string) (*models.LoanDetailResponse, error)
	Return(ctx context.Context, loanID string, actorUserID string) (*models.LoanDetailResponse, error)
}

var _ loanServiceAPI = (*loan.Service)(nil)

const (
	msgLoanNotFound        = "Préstamo no encontrado"
	msgLoanAlreadyReturned = "El préstamo ya fue devuelto"
	msgBorrowerNotFound    = "Estudiante no encontrado"
	msgLoanBookNotFound    = "Libro no encontrado"
	msgDueDateBad          = "La fecha de vencimiento no es válida"
	msgNoCopies            = "No hay ejemplares disponibles"
	msgStudentSanctioned   = "El estudiante tiene una sanción activa; no se puede registrar el préstamo"
	msgLoanStatusBad       = "Estado de préstamo no válido"
)

// LoanHandler exposes /api/loans.
type LoanHandler struct {
	svc loanServiceAPI
}

// NewLoanHandler constructs LoanHandler.
func NewLoanHandler(svc *loan.Service) *LoanHandler {
	return &LoanHandler{svc: svc}
}

// List GET /api/loans
func (h *LoanHandler) List(c *gin.Context) {
	limit, offset := normalizeLoanListQuery(c.Query("limit"), c.Query("offset"))
	filter := repositories.LoanListFilter{
		Status:    c.Query("status"),
		StudentID: strings.TrimSpace(c.Query("student_id")),
	}
	items, total, err := h.svc.List(c.Request.Context(), filter, limit, offset)
	if err != nil {
		writeLoanServiceError(c, err)
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
	c.JSON(http.StatusOK, models.ListLoansResponse{
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

// Create POST /api/loans
func (h *LoanHandler) Create(c *gin.Context) {
	var req models.CreateLoanRequest
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
		writeLoanServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Return PATCH /api/loans/:id/return
func (h *LoanHandler) Return(c *gin.Context) {
	adminID, ok := c.Get(ContextAuthUserID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	aid, _ := adminID.(string)
	out, err := h.svc.Return(c.Request.Context(), c.Param("id"), aid)
	if err != nil {
		writeLoanServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func writeLoanServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, loan.ErrStudentNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgBorrowerNotFound})
	case errors.Is(err, loan.ErrBookNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgLoanBookNotFound})
	case errors.Is(err, loan.ErrLoanNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgLoanNotFound})
	case errors.Is(err, loan.ErrDueDateInPast), errors.Is(err, loan.ErrInvalidDueDateFormat):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgDueDateBad, Errors: map[string]string{"due_date": msgDueDateBad}})
	case errors.Is(err, loan.ErrInvalidLoanStatus):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgLoanStatusBad, Errors: map[string]string{"status": msgLoanStatusBad}})
	case errors.Is(err, loan.ErrInvalidStudentIDQuery):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Parámetro student_id no válido", Errors: map[string]string{"student_id": "UUID inválido"}})
	case errors.Is(err, loan.ErrNoCopiesAvailable):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgNoCopies})
	case errors.Is(err, loan.ErrStudentHasActiveSanction):
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: msgStudentSanctioned,
			Errors:  map[string]string{"student_id": msgStudentSanctioned},
		})
	case errors.Is(err, loan.ErrLoanAlreadyReturned):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgLoanAlreadyReturned})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
	}
}

func normalizeLoanListQuery(limitStr, offsetStr string) (limit, offset int) {
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
