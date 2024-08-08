package main

import (
	"ecommerce-api/internal/db"
	"ecommerce-api/internal/handlers"
	"ecommerce-api/internal/middleware"
	"ecommerce-api/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database
	if err := db.Init(); err != nil {
		panic("failed to connect to database")
	}

	// migrate the models
	db.DB.AutoMigrate(&models.User{}, &models.Product{})

	r := gin.Default()

	// routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the E-commerce API!",
		})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/refresh-token", handlers.RefreshToken)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/", handlers.CreateUser)
		userGroup.GET("/:id", handlers.GetUser)
	}

	// Public product routes (available to all authenticated users)
	productGroup := r.Group("/products")
	productGroup.Use(middleware.AuthMiddleware())
	{
		productGroup.GET("/:id", handlers.GetProduct)
	}

	// Admin-only product routes
	adminProductGroup := productGroup.Group("/")
	adminProductGroup.Use(middleware.RoleMiddleware("admin"))
	{
		adminProductGroup.POST("/", handlers.CreateProduct)
		adminProductGroup.PUT("/:id", handlers.UpdateProduct)
		adminProductGroup.DELETE("/:id", handlers.DeleteProduct)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome to the protected route!"})
		})

		protected.PUT("/profile", handlers.UpdateProfile)

		adminGroup := protected.Group("/")
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.DELETE("/users/:id", handlers.DeleteUser) // Admin-only route
		}
	}

	// Start the server
	r.Run(":8080")
}
