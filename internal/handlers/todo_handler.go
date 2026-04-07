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

// UpdateTodoRequest represents the expected JSON payload for updating an existing todo item. It includes optional fields for the title and completion status, allowing for partial updates of the todo item.
type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// CreateTodoHandler returns a Gin handler function that handles the creation of a new todo item. It extracts the user ID from the Gin context, binds the incoming JSON payload to a CreateTodoRequest struct, and calls the CreateTodo repository function to create the todo item in the database. If successful, it returns the created todo item in the response with a 201 Created status. If there are any errors during binding or database operations, it returns an appropriate error response.
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	// Return a Gin handler function that handles the creation of a new todo item.
	return func(c *gin.Context) {
		// Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Bind the incoming JSON payload to a CreateTodoRequest struct. If there is an error during binding (e.g., missing required fields or invalid data types), return a 400 Bad Request response with the error message.
		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Call the CreateTodo repository function to create the todo item in the database, passing the user ID, title, and completion status. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the todo creation failed. If the todo item is created successfully, return the created todo item in the response with a 201 Created status.
		todo, err := repositories.CreateTodo(pool, req.Title, req.Completed, userID.(string))
		if err != nil {
			log.Printf("CreateTodo error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}
		// Return the created todo item in the response with a 201 Created status.
		c.JSON(http.StatusCreated, todo)
	}
}

// GetTodosHandler returns a Gin handler function that retrieves all todo items for the authenticated user. It extracts the user ID from the Gin context, calls the GetTodos repository function to fetch the todos from the database, and returns the list of todos in the response with a 200 OK status. If there are any errors during database operations, it logs the error and returns a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of todos failed.
func GetTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Call the GetTodos repository function to fetch the todos from the database for the authenticated user. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of todos failed. If the todos are retrieved successfully, return the list of todos in the response with a 200 OK status.
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

// GetTodoByIDHandler returns a Gin handler function that retrieves a specific todo item by its ID for the authenticated user. It extracts the user ID from the Gin context, parses the todo ID from the URL parameter, and calls the GetTodoByID repository function to fetch the todo item from the database. If the todo item is found, it returns the item in the response with a 200 OK status. If the todo item is not found, it returns a 404 Not Found response with an appropriate error message. If there are any errors during database operations, it logs the error and returns a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of the todo item failed.
func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
		idStr := c.Param("id")
		// Parse the todo ID from the URL parameter. If the ID is not a valid integer, return a 400 Bad Request response with an appropriate error message indicating that the ID is invalid.
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Call the GetTodoByID repository function to fetch the todo item from the database for the authenticated user and the specified ID. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of the todo item failed. If the todo item is not found, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found. If the todo item is retrieved successfully, return the item in the response with a 200 OK status.
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

// UpdateTodoHandler returns a Gin handler function that updates an existing todo item for the authenticated user. It extracts the user ID from the Gin context, parses the todo ID from the URL parameter, and binds the incoming JSON payload to an UpdateTodoRequest struct. It then calls the UpdateTodo repository function to update the todo item in the database. If the update is successful, it returns the updated todo item in the response with a 200 OK status. If the todo item is not found, it returns a 404 Not Found response with an appropriate error message. If there are any errors during binding or database operations, it logs the error and returns an appropriate error response.
func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	// Return a Gin handler function that updates an existing todo item for the authenticated user.
	return func(c *gin.Context) {
		// Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		// Bind the incoming JSON payload to an UpdateTodoRequest struct. If there is an error during binding (e.g., invalid data types), return a 400 Bad Request response with the error message. If neither the title nor the completed fields are provided in the request, return a 400 Bad Request response with an appropriate error message indicating that at least one field must be provided for the update.
		var req UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// If neither the title nor the completed fields are provided in the request, return a 400 Bad Request response with an appropriate error message indicating that at least one field must be provided for the update.
		if req.Title == nil && req.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field (title or completed) must be provided"})
			return
		}
		// Call the GetTodoByID repository function to fetch the existing todo item from the database for the authenticated user and the specified ID. This is necessary to ensure that the todo item exists before attempting to update it, and to retrieve the current values of the title and completed fields for use in the update operation. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of the existing todo item failed. If the todo item is not found, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found.
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
		// Determine the new values for the title and completed fields based on the incoming request. If a field is not provided in the request, use the existing value from the database. This allows for partial updates of the todo item, where only the fields that are provided in the request will be updated, while the other fields will remain unchanged.
		title := existingTodo.Title
		if req.Title != nil {
			title = *req.Title
		}

		completed := existingTodo.Completed
		if req.Completed != nil {
			completed = *req.Completed
		}
		// Call the UpdateTodo repository function to update the todo item in the database with the new values for the title and completed fields. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the update of the todo item failed. If the todo item is not found during the update operation, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found. If the update is successful
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
		// Return the updated todo item in the response with a 200 OK status.
		c.JSON(http.StatusOK, todo)
	}
}

// DeleteTodoHandler returns a Gin handler function that deletes a specific todo item by its ID for the authenticated user. It extracts the user ID from the Gin context, parses the todo ID from the URL parameter, and calls the DeleteTodo repository function to delete the todo item from the database. If the deletion is successful, it returns a success message in the response with a 200 OK status. If the todo item is not found, it returns a 404 Not Found response with an appropriate error message. If there are any errors during database operations, it logs the error and returns a 500 Internal Server Error response with an appropriate error message indicating that the deletion of the todo item failed.
func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	// Return a Gin handler function that deletes a specific todo item by its ID for the authenticated user.
	return func(c *gin.Context) {
		// Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		//	Parse the todo ID from the URL parameter. If the ID is not a valid integer, return a 400 Bad Request response with an appropriate error message indicating that the ID is invalid.
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Call the DeleteTodo repository function to delete the todo item from the database for the authenticated user and the specified ID. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the deletion of the todo item failed. If the todo item is not found during the deletion operation, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found. If the deletion is successful.
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
		//	Return a success message in the response with a 200 OK status.
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
