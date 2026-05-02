package handlers

import (
	"context"
	"errors"
	"library_back/internal/services/student"
	"math"
	"net/http"
	"strconv"
	"strings"

	"library_back/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

// studentServiceAPI is implemented by *services.StudentService and by test mocks.
type studentServiceAPI interface {
	Create(ctx context.Context, req models.CreateStudentRequest, adminID string) (*models.StudentResponse, error)
	GetByID(ctx context.Context, id string) (*models.StudentResponse, error)
	List(ctx context.Context, filter models.StudentFilter, limit, offset int) ([]models.StudentResponse, int64, error)
	Update(ctx context.Context, id string, req models.UpdateStudentRequest) (*models.StudentResponse, error)
	Deactivate(ctx context.Context, id string) error
	GetProfile(ctx context.Context, id string, loanLimit int) (*models.StudentProfileResponse, error)
	AwardBadge(ctx context.Context, studentID, badgeID string) (*models.AwardBadgeResponse, error)
	RevokeBadge(ctx context.Context, studentID, badgeID string) error
}

var _ studentServiceAPI = (*student.Service)(nil)

const (
	msgStudentNotFound = "Estudiante no encontrado"
	msgStudentConflict = "El recurso ya existe"
	msgDepartmentBad   = "Departamento no válido"
	msgBadgeNotFound   = "Insignia no encontrada"
	msgBadgeDuplicate  = "El estudiante ya tiene esta insignia"
	msgBadgeNotEarned  = "El estudiante no tiene esta insignia"
)

// StudentHandler exposes /api/students routes.
type StudentHandler struct {
	svc studentServiceAPI
}

// NewStudentHandler constructs StudentHandler.
func NewStudentHandler(svc *student.Service) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// Create POST /api/students
func (h *StudentHandler) Create(c *gin.Context) {
	var req models.CreateStudentRequest
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
		writeStudentServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// List GET /api/students
func (h *StudentHandler) List(c *gin.Context) {
	limit, offset := normalizeStudentListQuery(c.Query("limit"), c.Query("offset"))
	filter := models.StudentFilter{Search: c.Query("search")}
	items, total, err := h.svc.List(c.Request.Context(), filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
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
	c.JSON(http.StatusOK, models.ListStudentsResponse{
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

// Get GET /api/students/:id
func (h *StudentHandler) Get(c *gin.Context) {
	out, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeStudentServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Profile GET /api/students/:id/profile
func (h *StudentHandler) Profile(c *gin.Context) {
	loanLimit := 10
	if v := c.Query("loan_limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			loanLimit = n
		}
	}
	out, err := h.svc.GetProfile(c.Request.Context(), c.Param("id"), loanLimit)
	if err != nil {
		writeStudentServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Update PATCH /api/students/:id
func (h *StudentHandler) Update(c *gin.Context) {
	var req models.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	out, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		writeStudentServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Deactivate DELETE /api/students/:id
func (h *StudentHandler) Deactivate(c *gin.Context) {
	if err := h.svc.Deactivate(c.Request.Context(), c.Param("id")); err != nil {
		writeStudentServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// AwardBadge POST /api/students/:id/badges
func (h *StudentHandler) AwardBadge(c *gin.Context) {
	var req models.AwardBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindError(c)
		return
	}
	out, err := h.svc.AwardBadge(c.Request.Context(), c.Param("id"), req.BadgeID)
	if err != nil {
		writeStudentServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// RevokeBadge DELETE /api/students/:id/badges/:badgeId
func (h *StudentHandler) RevokeBadge(c *gin.Context) {
	badgeID := strings.TrimSpace(c.Param("badgeId"))
	if _, err := uuid.Parse(badgeID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Parámetro badge_id no válido",
			Errors:  map[string]string{"badge_id": "UUID inválido"},
		})
		return
	}
	if err := h.svc.RevokeBadge(c.Request.Context(), c.Param("id"), badgeID); err != nil {
		writeStudentServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeStudentServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, student.ErrStudentNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgStudentNotFound})
	case errors.Is(err, student.ErrDuplicateStudentIDCode), errors.Is(err, student.ErrDuplicateStudentEmail):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgStudentConflict})
	case errors.Is(err, student.ErrDepartmentNotFound):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgDepartmentBad, Errors: map[string]string{"department_id": msgDepartmentBad}})
	case errors.Is(err, student.ErrBadgeNotFound):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgBadgeNotFound, Errors: map[string]string{"badge_id": msgBadgeNotFound}})
	case errors.Is(err, student.ErrBadgeAlreadyEarned):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgBadgeDuplicate})
	case errors.Is(err, student.ErrStudentDoesNotHaveBadge):
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: msgBadgeNotEarned,
			Errors:  map[string]string{"badge_id": msgBadgeNotEarned},
		})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
	}
}

func normalizeStudentListQuery(limitStr, offsetStr string) (limit, offset int) {
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
