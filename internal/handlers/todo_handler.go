package handlers

import (
	"net/http"

	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTodoRequest represents the expected JSON payload for creating a new todo item. It includes the title of the todo and its completion status.
type CreateTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

// UpdateTodoRequest represents the expected JSON payload for updating an existing todo item. It includes optional fields for the title and completion status, allowing for partial updates of the todo item.
type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// responseTodo represents the structure of a todo item that will be sent back in the response. It includes the ID, user ID, title, and completion status of the todo item.
type responseTodo struct {
	ID        int32     `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}

// CreateTodoHandler handles the creation of a new todo item. It retrieves the user ID from the context, binds the incoming JSON payload to a CreateTodoRequest struct, and calls the CreateTodo method from the database layer to insert the new todo item into the database. If successful, it returns the created todo item in the response; otherwise, it returns an appropriate error message.
func CreateTodoHandler(cfg *config.Config) gin.HandlerFunc {
	// returns a Gin handler function that processes the creation of a new todo item.
	return func(c *gin.Context) {
		// Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Bind the incoming JSON payload to a CreateTodoRequest struct. If there is an error during binding (e.g., missing required fields), return a bad request error response.
		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Call the CreateTodo method from the database layer, passing the user ID, title, and completion status from the request. If there is an error during the database operation, return an internal server error response. If successful, construct a responseTodo struct with the created todo item and return it in the response.
		todo, err := cfg.DB.CreateTodo(c, database.CreateTodoParams{
			UserID:    userID.(uuid.UUID),
			Title:     req.Title,
			Completed: req.Completed,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Construct a responseTodo struct with the created todo item and return it in the response.
		response := responseTodo{
			ID:        todo.ID,
			UserID:    todo.UserID,
			Title:     todo.Title,
			Completed: todo.Completed,
		}
		// Return the created todo item in the response with a status of 200 OK.
		c.JSON(http.StatusOK, response)
	}
}
