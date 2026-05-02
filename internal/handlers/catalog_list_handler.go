package handlers

import (
	"net/http"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/gin-gonic/gin"
)

// CatalogListHandler serves GET /api/authors and GET /api/genres.
type CatalogListHandler struct {
	authors repositories.AuthorCatalogRepository
	genres  repositories.GenreCatalogRepository
}

// NewCatalogListHandler constructs CatalogListHandler.
func NewCatalogListHandler(authors repositories.AuthorCatalogRepository, genres repositories.GenreCatalogRepository) *CatalogListHandler {
	return &CatalogListHandler{authors: authors, genres: genres}
}

// ListAuthors GET /api/authors
func (h *CatalogListHandler) ListAuthors(c *gin.Context) {
	rows, err := h.authors.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	out := make([]models.AuthorResponse, 0, len(rows))
	for _, a := range rows {
		out = append(out, models.AuthorResponse{ID: a.ID, Name: a.Name})
	}
	c.JSON(http.StatusOK, out)
}

// ListGenres GET /api/genres
func (h *CatalogListHandler) ListGenres(c *gin.Context) {
	rows, err := h.genres.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	out := make([]models.GenreResponse, 0, len(rows))
	for _, g := range rows {
		out = append(out, models.GenreResponse{ID: g.ID, Name: g.Name, Code: g.Code})
	}
	c.JSON(http.StatusOK, out)
}
