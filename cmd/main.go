package main

import (
	"log"

	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/carloscfgos1980/todo-auth/internal/handlers"
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
	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos", handlers.GetTodosHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))
	router.PUT("/todos/:id", handlers.UpdateTodoHandler(pool))
	router.DELETE("/todos/:id", handlers.DeleteTodoHandler(pool))
	router.Run(":" + cfg.Port)
}
