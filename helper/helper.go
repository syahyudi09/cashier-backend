package helper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct memeriksa validitas struktur data menggunakan validator
func ValidateStruct(c *gin.Context, data interface{}) bool {
	if err := validate.Struct(data); err != nil {
		validationErrors := FormatValidationError(err)
		response := APIResponse(
			strings.Join(validationErrors, ", "),
			http.StatusBadRequest,
			"error",
			nil,
		)
		c.JSON(http.StatusBadRequest, response)
		return false
	}

	return true
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("%s", e.Translate(nil)))
	}
	return errors
}
