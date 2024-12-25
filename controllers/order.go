package controllers

import (
	"ecommerce-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

// PlaceOrder - Place a new order
func PlaceOrder(c *gin.Context, db *gorm.DB) {
	var orderRequest struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	// Bind the request data
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	// Check if the product exists
	var product models.Product
	if err := db.Where("id = ?", orderRequest.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if there is enough stock
	if product.Stock < orderRequest.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock"})
		return
	}

	// Create the order
	order := models.Order{
		UserID: uint(c.MustGet("user").(models.User).ID),
		Status: "pending",
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order"})
		return
	}

	// Create order items
	orderItem := models.OrderItem{
		OrderID:  order.ID,
		ProductID: orderRequest.ProductID,
		Quantity:  orderRequest.Quantity,
	}

	if err := db.Create(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order items"})
		return
	}

	// Update product stock
	product.Stock -= orderRequest.Quantity
	db.Save(&product)

	c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully", "order": order})
}

// ListOrders - List all orders for a specific user
func ListOrders(c *gin.Context, db *gorm.DB) {
	userID := uint(c.MustGet("user").(models.User).ID)

	var orders []models.Order
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// CancelOrder - Cancel an order if it is still in 'pending' status
func CancelOrder(c *gin.Context, db *gorm.DB) {
	orderID := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if the order is still in 'pending' status
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel order in current status"})
		return
	}

	// Update order status to 'cancelled'
	order.Status = "cancelled"
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// UpdateOrderStatus - Update the status of an order (admin only)
func UpdateOrderStatus(c *gin.Context, db *gorm.DB) {
	orderID := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if user is admin
	user, _ := c.Get("user")
	if user.(models.User).Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Update order status
	order.Status = statusUpdate.Status
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
