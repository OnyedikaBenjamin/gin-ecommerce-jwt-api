package models

import "github.com/jinzhu/gorm"

// User represents a user in the system
type User struct {
	gorm.Model
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}