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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pool, err := database.ConnectPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Todo API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))
	protected.POST("/", handlers.CreateTodoHandler(pool))
	protected.GET("/", handlers.GetTodosHandler(pool))
	protected.GET("/:id", handlers.GetTodoByIDHandler(pool))
	protected.PUT("/:id", handlers.UpdateTodoHandler(pool))
	protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))

	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())
	router.Run(":" + cfg.Port)
}
