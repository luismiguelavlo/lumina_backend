package handlers

import (
	"errors"
	"library_back/internal/services/book"
	"math"
	"net/http"
	"strconv"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)

const (
	msgCatalogConflict = "El recurso ya existe"
	msgBookNotFound    = "Libro no encontrado"
	msgAuthorMissing   = "Autor no encontrado"
	msgGenreMissing    = "Género no encontrado"
)

// BookHandler exposes /api/books routes.
type BookHandler struct {
	svc *book.BookService
}

// NewBookHandler constructs BookHandler.
func NewBookHandler(svc *book.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// Create POST /api/books
func (h *BookHandler) Create(c *gin.Context) {
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindErrorWithDetail(c, err)
		return
	}
	out, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		writeBookServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// List GET /api/books
func (h *BookHandler) List(c *gin.Context) {
	limit, offset := normalizeBookListParams(c.Query("limit"), c.Query("offset"))
	genreID := c.Query("genre_id")
	if genreID != "" {
		if _, err := uuid.Parse(genreID); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Message: "Parámetros inválidos",
				Errors:  map[string]string{"genre_id": "Debe ser un UUID válido"},
			})
			return
		}
	}
	filter := repositories.BookListFilter{Search: c.Query("search"), GenreID: genreID}
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
	c.JSON(http.StatusOK, models.ListBooksResponse{
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

// Get GET /api/books/:id
func (h *BookHandler) Get(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		writeBookServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// UpdatePut PUT /api/books/:id
func (h *BookHandler) UpdatePut(c *gin.Context) {
	h.update(c)
}

// UpdatePatch PATCH /api/books/:id
func (h *BookHandler) UpdatePatch(c *gin.Context) {
	h.update(c)
}

func (h *BookHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			writeValidationErrors(c, errs)
			return
		}
		writeJSONBindErrorWithDetail(c, err)
		return
	}
	out, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		writeBookServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete DELETE /api/books/:id
func (h *BookHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeBookServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeBookServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, book.ErrBookNotFound):
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgBookNotFound})
	case errors.Is(err, book.ErrDuplicateISBN), errors.Is(err, book.ErrDuplicateCatalogCode):
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: msgCatalogConflict})
	case errors.Is(err, book.ErrAuthorNotFound):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgAuthorMissing, Errors: map[string]string{"author_ids": msgAuthorMissing}})
	case errors.Is(err, book.ErrGenreNotFound):
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: msgGenreMissing, Errors: map[string]string{"genre_ids": msgGenreMissing}})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
	}
}

func normalizeBookListParams(limitStr, offsetStr string) (limit, offset int) {
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
