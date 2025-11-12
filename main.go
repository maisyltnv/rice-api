package main

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/handlers"
	"example.com/go-xampp-api/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	// Health check
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// AUTH routes (Admin/User)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// CUSTOMER AUTH routes
	r.POST("/customers/register", handlers.CustomerRegister)
	r.POST("/customers/login", handlers.CustomerLogin)

	// CATEGORY routes
	r.GET("/categories", handlers.GetCategories)
	r.GET("/categories/:id", handlers.GetCategory)
	r.POST("/categories", middleware.AuthMiddleware(), handlers.CreateCategory)
	r.PUT("/categories/:id", middleware.AuthMiddleware(), handlers.UpdateCategory)
	r.DELETE("/categories/:id", middleware.AuthMiddleware(), handlers.DeleteCategory)

	// PRODUCT routes
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProduct)
	r.POST("/products", middleware.AuthMiddleware(), handlers.CreateProduct)
	r.PUT("/products/:id", middleware.AuthMiddleware(), handlers.UpdateProduct)
	r.DELETE("/products/:id", middleware.AuthMiddleware(), handlers.DeleteProduct)

	// CUSTOMER routes
	r.GET("/customers", handlers.GetCustomers)
	r.GET("/customers/:id", handlers.GetCustomer)
	r.POST("/customers", middleware.AuthMiddleware(), handlers.CreateCustomer)
	r.PUT("/customers/:id", middleware.AuthMiddleware(), handlers.UpdateCustomer)
	r.DELETE("/customers/:id", middleware.AuthMiddleware(), handlers.DeleteCustomer)

	// ORDER routes
	r.GET("/orders", middleware.AuthMiddleware(), handlers.GetOrders)
	r.GET("/orders/:id", middleware.AuthMiddleware(), handlers.GetOrder)
	r.POST("/orders", middleware.AuthMiddleware(), handlers.CreateOrder)
	r.PUT("/orders/:id/status", middleware.AuthMiddleware(), handlers.UpdateOrderStatus)
	r.DELETE("/orders/:id", middleware.AuthMiddleware(), handlers.DeleteOrder)

	// CART routes (Customer)
	cartRoutes := r.Group("/cart")
	cartRoutes.Use(middleware.AuthMiddleware())
	{
		cartRoutes.GET("", handlers.GetCart)
		cartRoutes.POST("/items", handlers.AddCartItem)
		cartRoutes.PUT("/items/:item_id", handlers.UpdateCartItem)
		cartRoutes.DELETE("/items/:item_id", handlers.DeleteCartItem)
		cartRoutes.DELETE("", handlers.ClearCart)
	}

	r.Run(":8081")
}
