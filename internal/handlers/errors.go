package handlers

import (
	"net/http"
	"strings"

	"library_back/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func writeJSONBindError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Solicitud inválida"})
}

func writeJSONBindErrorWithDetail(c *gin.Context, err error) {
	resp := models.ErrorResponse{Message: "Solicitud inválida"}
	if err != nil {
		resp.Errors = map[string]string{
			"body": err.Error(),
		}
	}
	c.JSON(http.StatusBadRequest, resp)
}

func writeValidationErrors(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Validación fallida"})
		return
	}
	fields := make(map[string]string)
	for _, e := range errs {
		fields[e.Field()] = humanizeValidationTag(e.Tag())
	}
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Message: "Errores de validación",
		Errors:  fields,
	})
}

func humanizeValidationTag(tag string) string {
	switch tag {
	case "required":
		return "campo obligatorio"
	case "email":
		return "email inválido"
	case "min":
		return "longitud o valor mínimo no cumplido"
	case "max":
		return "longitud máxima excedida"
	case "gte":
		return "valor demasiado pequeño"
	case "uuid":
		return "UUID inválido"
	case "oneof":
		return "valor no permitido"
	case "url":
		return "URL inválida"
	default:
		return "valor inválido"
	}
}

// bearerToken returns the token from "Authorization: Bearer <token>".
func bearerToken(c *gin.Context) (string, bool) {
	h := c.GetHeader("Authorization")
	if h == "" {
		return "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	t := strings.TrimSpace(parts[1])
	if t == "" {
		return "", false
	}
	return t, true
}
