package handlers

import (
	"net/http"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/gin-gonic/gin"
)

// DepartmentHandler serves GET /api/departments.
type DepartmentHandler struct {
	repo repositories.StudentDepartmentRepository
}

// NewDepartmentHandler constructs DepartmentHandler.
func NewDepartmentHandler(repo repositories.StudentDepartmentRepository) *DepartmentHandler {
	return &DepartmentHandler{repo: repo}
}

// List GET /api/departments
func (h *DepartmentHandler) List(c *gin.Context) {
	rows, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	out := make([]models.DepartmentResponse, 0, len(rows))
	for _, d := range rows {
		out = append(out, models.DepartmentResponse{ID: d.ID, Name: d.Name, Code: d.Code})
	}
	c.JSON(http.StatusOK, out)
}
