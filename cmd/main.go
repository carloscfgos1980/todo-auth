package main

import (
	"log"

	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/carloscfgos1980/todo-auth/internal/handlers"
	"github.com/carloscfgos1980/todo-auth/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// Connect to the PostgreSQL database using the provided URL
	pool, err := database.ConnectPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer pool.Close()
	// Initialize the Gin router and set up routes and middleware
	var router *gin.Engine = gin.Default()
	// Set trusted proxies to nil to disable Gin's default behavior of trusting all proxies
	router.SetTrustedProxies(nil)
	// Define a simple health check endpoint at the root URL to verify that the API is running and the database connection is successful
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Todo API is running",
			"status":   "success",
			"database": "connected",
		})
	})
	// Define authentication routes for user registration and login, and protect the /todos routes with the AuthMiddleware to ensure that only authenticated users can access them
	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	// Define a protected route group for /todos that uses the AuthMiddleware to ensure that only authenticated users can access the todo-related endpoints. Within this group, define routes for creating, retrieving, updating, and deleting todo items, and associate each route with the corresponding handler function from the handlers package.
	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))
	protected.POST("/", handlers.CreateTodoHandler(pool))
	protected.GET("/", handlers.GetTodosHandler(pool))
	protected.GET("/:id", handlers.GetTodoByIDHandler(pool))
	protected.PUT("/:id", handlers.UpdateTodoHandler(pool))
	protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))

	// Define a protected test endpoint to verify that the AuthMiddleware is working correctly. This endpoint will only be accessible to requests that include a valid JWT token in the Authorization header. If the token is valid, the endpoint will return a success message; otherwise, it will return an unauthorized error response.
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())
	// Start the Gin server on the specified port from the configuration and log any errors that occur while running the server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
