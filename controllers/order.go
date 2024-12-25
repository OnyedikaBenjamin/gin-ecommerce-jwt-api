package controllers

import (
	"ecommerce-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

func PlaceOrder(c *gin.Context, db *gorm.DB) {
    var orderRequest struct {
        ProductIDs []uint `json:"product_ids"`
        Quantities []int  `json:"quantity"`
    }

    if err := c.ShouldBindJSON(&orderRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
        return
    }

    if len(orderRequest.ProductIDs) != len(orderRequest.Quantities) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatched product IDs and quantities"})
        return
    }

    order := models.Order{
        UserID: uint(c.MustGet("user").(models.User).ID),
        Status: "pending",
    }

    if err := db.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order"})
        return
    }

    for i, productID := range orderRequest.ProductIDs {
        var product models.Product

        if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "product_id": productID})
            return
        }

        if product.Stock < orderRequest.Quantities[i] {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":      "Not enough stock",
                "product_id": productID,
            })
            return
        }

        orderItem := models.OrderItem{
            OrderID:   order.ID,
            ProductID: productID,
            Quantity:  orderRequest.Quantities[i],
        }

        if err := db.Create(&orderItem).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order items"})
            return
        }

        product.Stock -= orderRequest.Quantities[i]
        db.Save(&product)
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully", "order": order})
}

func ListOrders(c *gin.Context, db *gorm.DB) {
	user := c.MustGet("user").(models.User)

	var orders []models.Order
	if err := db.Preload("Items").Where("user_id = ?", user.ID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}



func CancelOrder(c *gin.Context, db *gorm.DB) {
	orderID := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel order in current status"})
		return
	}

	order.Status = "cancelled"
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func UpdateOrderStatus(c *gin.Context, db *gorm.DB) {
	orderID := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	user := c.MustGet("user").(models.User)
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status", "details": err.Error()})
		return
	}

	order.Status = statusUpdate.Status
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}