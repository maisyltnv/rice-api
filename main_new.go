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

	// Health check
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// AUTH routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

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
	r.GET("/orders", handlers.GetOrders)
	r.GET("/orders/:id", handlers.GetOrder)
	r.POST("/orders", middleware.AuthMiddleware(), handlers.CreateOrder)
	r.PUT("/orders/:id/status", middleware.AuthMiddleware(), handlers.UpdateOrderStatus)
	r.DELETE("/orders/:id", middleware.AuthMiddleware(), handlers.DeleteOrder)

	r.Run(":8080")
}
