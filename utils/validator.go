package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// CustomValidator is used to wrap Gin context and handle validation
type CustomValidator struct {
	Validator *validator.Validate
}

// NewValidator initializes a new CustomValidator
func NewValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

// ValidateRequest validates the request body data using the Gin context
func (cv *CustomValidator) ValidateRequest(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return false
	}

	// Perform validation on the request data
	if err := cv.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return false
	}

	return true
}
