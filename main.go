package main

import (
	"ecommerce-api/routes"
	"ecommerce-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func init() {
	// Initialize the database connection
	dsn := "root:black0111@tcp(localhost:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
}

func main() {
	// Set up the Gin router
	r := routes.SetupRouter(db) // Pass db to routes

	// Start the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
