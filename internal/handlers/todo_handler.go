package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/carloscfgos1980/todo-auth/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTodoRequest represents the expected JSON payload for creating a new todo item. It includes the title of the todo and its completion status.
type CreateTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repositories.CreateTodo(pool, req.Title, req.Completed, userID.(string))
		if err != nil {
			log.Printf("CreateTodo error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}
		c.JSON(http.StatusCreated, todo)
	}
}

func GetTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		todos, err := repositories.GetTodos(pool, userID.(string))
		if err != nil {
			log.Printf("GetTodos error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todos"})
			return
		}
		// Return the list of todos in the response with a 200 OK status.
		c.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {

	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		todo, err := repositories.GetTodoByID(pool, id, userID.(string))
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return the retrieved todo item in the response with a 200 OK status.
		c.JSON(http.StatusOK, todo)
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		var req UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Title == nil && req.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field (title or completed) must be provided"})
			return
		}

		existingTodo, err := repositories.GetTodoByID(pool, id, userID.(string))
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			log.Printf("GetTodoByID error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		title := existingTodo.Title
		if req.Title != nil {
			title = *req.Title
		}

		completed := existingTodo.Completed
		if req.Completed != nil {
			completed = *req.Completed
		}
		todo, err := repositories.UpdateTodo(pool, id, title, completed, userID.(string))
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			log.Printf("UpdateTodo error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = repositories.DeleteTodo(pool, id, userID.(string))
		if err != nil {
			if err.Error() == "todo with id "+idStr+" not found for user "+userID.(string) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			log.Printf("DeleteTodo error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
