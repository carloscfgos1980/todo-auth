package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

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

// GetTodosHandler handles the retrieval of all todo items for a specific user. It retrieves the user ID from the context, calls the GetTodosByUserID method from the database layer to fetch the todo items associated with that user, and returns the list of todo items in the response. If there is an error during the process, it returns an appropriate error message.
func GetTodosHandler(cfg *config.Config) gin.HandlerFunc {
	// returns a Gin handler function that processes the retrieval of all todo items for a specific user.
	return func(c *gin.Context) {
		// Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Call the GetTodosByUserID method from the database layer, passing the user ID to fetch the todo items associated with that user. If there is an error during the database operation, return an internal server error response. If successful, construct a slice of responseTodo structs with the retrieved todo items and return it in the response.
		todos, err := cfg.DB.GetTodosByUserID(c, userID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Construct a slice of responseTodo structs with the retrieved todo items and return it in the response.
		var response []responseTodo
		for _, todo := range todos {
			response = append(response, responseTodo{
				ID:        todo.ID,
				UserID:    todo.UserID,
				Title:     todo.Title,
				Completed: todo.Completed,
			})
		}
		// Return the list of todo items in the response with a status of 200 OK.
		c.JSON(http.StatusOK, response)
	}
}

// GetTodoByIDHandler handles the retrieval of a specific todo item by its ID. It retrieves the user ID from the context, extracts the todo ID from the URL parameters, and calls the GetTodoByID method from the database layer to fetch the todo item. If the todo item is found and belongs to the authenticated user, it returns the todo item in the response; otherwise, it returns an appropriate error message (e.g., not found, unauthorized, or forbidden).
func GetTodoByIDHandler(cfg *config.Config) gin.HandlerFunc {
	// returns a Gin handler function that processes the retrieval of a specific todo item by its ID.
	return func(c *gin.Context) {
		// Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Extract the todo ID from the URL parameters and convert it to an integer. If there is an error during conversion (e.g., invalid ID format), return a bad request error response.
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Call the GetTodoByID method from the database layer, passing the todo ID to fetch the specific todo item. If there is an error during the database operation, return an internal server error response. If the todo item is not found, return a not found error response. If the todo item is found but does not belong to the authenticated user, return a forbidden error response. If successful, construct a responseTodo struct with the retrieved todo item and return it in the response.
		todo, err := cfg.DB.GetTodoByID(c, int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Check if the retrieved todo item belongs to the authenticated user. If not, return a forbidden error response.
		if todo.UserID != userID.(uuid.UUID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this todo"})
			return
		}
		// Construct a responseTodo struct with the retrieved todo item and return it in the response.
		response := responseTodo{
			ID:        todo.ID,
			UserID:    todo.UserID,
			Title:     todo.Title,
			Completed: todo.Completed,
		}
		// Return the retrieved todo item in the response with a status of 200 OK.
		c.JSON(http.StatusOK, response)
	}
}

// UpdateTodoHandler handles the update of an existing todo item. It retrieves the user ID from the context, extracts the todo ID from the URL parameters, binds the incoming JSON payload to an UpdateTodoRequest struct, and calls the UpdateTodo method from the database layer to update the todo item in the database. If the todo item is found and belongs to the authenticated user, it returns the updated todo item in the response; otherwise, it returns an appropriate error message (e.g., not found, unauthorized, or forbidden).
func UpdateTodoHandler(cfg *config.Config) gin.HandlerFunc {
	// returns a Gin handler function that processes the update of an existing todo item.
	return func(c *gin.Context) {
		// Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Extract the todo ID from the URL parameters and convert it to an integer. If there is an error during conversion (e.g., invalid ID format), return a bad request error response.
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Bind the incoming JSON payload to an UpdateTodoRequest struct. If there is an error during binding (e.g., invalid JSON format), return a bad request error response.
		var req UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Call the GetTodoByID method from the database layer to fetch the existing todo item. If there is an error during the database operation, return an internal server error response. If the todo item is not found, return a not found error response. If the todo item is found but does not belong to the authenticated user, return a forbidden error response.
		dbTodo, err := cfg.DB.GetTodoByID(c, int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Check if the retrieved todo item belongs to the authenticated user. If not, return a forbidden error response.
		if dbTodo.UserID != userID.(uuid.UUID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this todo"})
			return
		}
		// Determine the updated title and completion status for the todo item. If the corresponding fields in the UpdateTodoRequest struct are nil, use the existing values from the database; otherwise, use the new values from the request.
		title := dbTodo.Title
		if req.Title != nil {
			title = *req.Title
		}
		completed := dbTodo.Completed
		if req.Completed != nil {
			completed = *req.Completed
		}
		// Call the UpdateTodo method from the database layer, passing the todo ID, updated title, and completion status. If there is an error during the database operation, return an internal server error response. If successful, construct a responseTodo struct with the updated todo item and return it in the response.
		todo, err := cfg.DB.UpdateTodo(c, database.UpdateTodoParams{
			ID:        int32(id),
			Title:     title,
			Completed: completed,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Construct a responseTodo struct with the updated todo item and return it in the response.
		response := responseTodo{
			ID:        todo.ID,
			UserID:    todo.UserID,
			Title:     todo.Title,
			Completed: todo.Completed,
		}
		// Return the updated todo item in the response with a status of 200 OK.
		c.JSON(http.StatusOK, response)

	}
}

func DeleteTodoHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Extract the todo ID from the URL parameters and convert it to an integer. If there is an error during conversion (e.g., invalid ID format), return a bad request error response.
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Call the GetTodoByID method from the database layer to fetch the existing todo item. If there is an error during the database operation, return an internal server error response. If the todo item is not found, return a not found error response. If the todo item is found but does not belong to the authenticated user, return a forbidden error response.
		dbTodo, err := cfg.DB.GetTodoByID(c, int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if dbTodo.UserID != userID.(uuid.UUID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this todo"})
			return
		}
		// Call the DeleteTodo method from the database layer, passing the todo ID. If there is an error during the database operation, return an internal server error response. If successful, return a success message in the response.
		err = cfg.DB.DeleteTodo(c, int32(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
