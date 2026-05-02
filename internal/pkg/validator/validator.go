package validator

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator for use with Gin.
type Validator struct {
	validate *validator.Validate
}

// New returns a Validator instance.
func New() *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})
	return &Validator{validate: v}
}

// BindAndValidate binds JSON from the request into dst and validates it.
// On validation error responds with 400 and error details and returns false.
func (v *Validator) BindAndValidate(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	if err := v.validate.Struct(dst); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			fields := make(map[string]string)
			for _, e := range errs {
				fields[e.Field()] = e.Tag()
			}
			c.JSON(http.StatusUnprocessableEntity, gin.H{"validation": fields})
			return false
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return false
	}
	return true
}

// ValidateStruct runs struct validation only (no JSON bind).
func (v *Validator) ValidateStruct(dst any) error {
	return v.validate.Struct(dst)
}
