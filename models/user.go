package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email string `json:"email" binding:"required,email" gorm:"unique"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}