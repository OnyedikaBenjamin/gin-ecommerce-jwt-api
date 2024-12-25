package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func (cv *CustomValidator) ValidateRequest(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return false
	}

	if err := cv.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return false
	}

	return true
}