package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

func CreateProduct(c *gin.Context, db *gorm.DB) {
	var product models.Product
	validator := utils.NewValidator()

	if !validator.ValidateRequest(c, &product) {
		return
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func GetProducts(c *gin.Context, db *gorm.DB) {
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context, db *gorm.DB) {
	var product models.Product
	id := c.Param("id")
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	db.Save(&product)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}