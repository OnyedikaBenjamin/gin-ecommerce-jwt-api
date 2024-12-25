package routes

import (
	"ecommerce-api/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ecommerce-api/utils"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", func(c *gin.Context) {
		controllers.Register(c, db)  // Pass db to controller
	})
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db)  // Pass db to controller
	})

	// Authenticated routes
	authorized := r.Group("/auth")
	authorized.Use(utils.AuthMiddleware(db))  // Pass db to AuthMiddleware
	{
		// Product routes (admin only)
		authorized.POST("/create-product", func(c *gin.Context) {
			controllers.CreateProduct(c, db)  // Pass db to controller
		})
		authorized.GET("/products", func(c *gin.Context) {
			controllers.GetProducts(c, db)  // Pass db to controller
		})
		authorized.PUT("/product/:id", func(c *gin.Context) {
			controllers.UpdateProduct(c, db)  // Pass db to controller
		})
		authorized.DELETE("/product/:id", func(c *gin.Context) {
			controllers.DeleteProduct(c, db)  // Pass db to controller
		})

		// Order routes
		authorized.POST("/place-order", func(c *gin.Context) {
			controllers.PlaceOrder(c, db)  // Pass db to controller
		})
		authorized.GET("/orders", func(c *gin.Context) {
			controllers.ListOrders(c, db)  // Pass db to controller
		})
		authorized.PUT("/order/:id/cancel", func(c *gin.Context) {
			controllers.CancelOrder(c, db)  // Pass db to controller
		})

		// Admin-only order status update
		authorized.PUT("/order/:id/status", func(c *gin.Context) {
			controllers.UpdateOrderStatus(c, db)  // Pass db to controller
		})
	}

	return r
}
