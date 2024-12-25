package routes

import (
	"ecommerce-api/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ecommerce-api/utils"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		controllers.Register(c, db) 
	})
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db) 
	})

	authorized := r.Group("/auth")
	authorized.Use(utils.AuthMiddleware(db)) 
	{
		authorized.POST("/create-product", utils.IsAdminMiddleware(), func(c *gin.Context) {
			controllers.CreateProduct(c, db)
		})		
		authorized.GET("/products", utils.IsAdminMiddleware(), func(c *gin.Context) {
			controllers.GetProducts(c, db)
		})
		authorized.PUT("/product/:id", utils.IsAdminMiddleware(), func(c *gin.Context) {
			controllers.UpdateProduct(c, db) 
		})
		authorized.DELETE("/product/:id", utils.IsAdminMiddleware(), func(c *gin.Context) {
			controllers.DeleteProduct(c, db)  
		})

		authorized.POST("/place-order", func(c *gin.Context) {
			controllers.PlaceOrder(c, db) 
		})
		authorized.GET("/orders", func(c *gin.Context) {
			controllers.ListOrders(c, db) 
		})
		authorized.PUT("/order/:id/cancel", func(c *gin.Context) {
			controllers.CancelOrder(c, db) 
		})

		authorized.PUT("/order/:id/status", utils.IsAdminMiddleware(), func(c *gin.Context) {
			controllers.UpdateOrderStatus(c, db) 
		})
	}

	return r
}