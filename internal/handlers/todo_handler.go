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

// CreateTodoHandler returns a Gin handler function that processes requests to create a new todo item. It validates the incoming JSON payload, calls the repository function to create the todo in the database, and returns the created todo item in the response.
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	// The returned handler function processes the incoming request to create a new todo item. It first attempts to bind the JSON payload to a CreateTodoRequest struct. If the binding fails (e.g., due to missing required fields), it responds with a 400 Bad Request status and an error message. If the binding is successful, it calls the CreateTodo function from the repositories package to insert the new todo into the database. If there is an error during this process, it logs the error and responds with a 500 Internal Server Error status. If the todo is created successfully, it responds with a 201 Created status and includes the created todo item in the response body.
	return func(c *gin.Context) {
		// Bind the incoming JSON payload to a CreateTodoRequest struct. This will validate that the required fields are present and correctly formatted. If the binding fails, return a 400 Bad Request response with the error message.
		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Call the CreateTodo function from the repositories package to insert the new todo into the database. Pass the database connection pool, title, and completion status as arguments. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message. If the todo is created successfully, return a 201 Created response with the created todo item in the response body.
		todo, err := repositories.CreateTodo(pool, req.Title, req.Completed)
		if err != nil {
			log.Printf("CreateTodo error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}
		// Return the created todo item in the response with a 201 Created status.
		c.JSON(http.StatusCreated, gin.H{"todo": todo})
	}
}

// GetTodosHandler returns a Gin handler function that processes requests to retrieve all todo items. It calls the repository function to fetch the todos from the database and returns them in the response. If there is an error during the retrieval process, it logs the error and responds with a 500 Internal Server Error status.
func GetTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the GetTodos function from the repositories package to fetch all todo items from the database. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message. If the todos are retrieved successfully, return a 200 OK response with the list of todos in the response body.
		todos, err := repositories.GetTodos(pool)
		if err != nil {
			log.Printf("GetTodos error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todos"})
			return
		}
		// Return the list of todos in the response with a 200 OK status.
		c.JSON(http.StatusOK, gin.H{"todos": todos})
	}
}

// GetTodoByIDHandler returns a Gin handler function that processes requests to retrieve a specific todo item by its ID. It extracts the ID from the URL parameters, validates it, and calls the repository function to fetch the todo from the database. If the ID is invalid, it responds with a 400 Bad Request status. If there is an error during the retrieval process, it logs the error and responds with a 500 Internal Server Error status. If the todo is not found, it responds with a 404 Not Found status. If the todo is retrieved successfully, it returns a 200 OK response with the todo item in the response body.
func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	// The returned handler function processes the incoming request to retrieve a specific todo item by its ID. It first extracts the ID from the URL parameters and attempts to convert it to an integer. If the conversion fails (e.g., due to an invalid ID format), it responds with a 400 Bad Request status and an error message. If the ID is valid, it calls the GetTodoByID function from the repositories package to fetch the todo item from the database. If there is an error during this process, it logs the error and responds with a 500 Internal Server Error status. If the todo item is not found (i.e., the repository function returns nil), it responds with a 404 Not Found status. If the todo item is retrieved successfully, it responds with a 200 OK status and includes the todo item in the response body.
	return func(c *gin.Context) {
		// Extract the ID from the URL parameters and attempt to convert it to an integer. If the conversion fails, return a 400 Bad Request response with an error message indicating that the ID is invalid.
		idStr := c.Param("id")
		// Convert the ID string to an integer. If the conversion fails, return a 400 Bad Request response with an error message indicating that the ID is invalid.
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Call the GetTodoByID function from the repositories package to fetch the todo item from the database using the provided ID.
		todo, err := repositories.GetTodoByID(pool, id)

		// If the todo item is not found (i.e., the repository function returns nil), return a 404 Not Found response with an error message indicating that the todo item was not found. If there is an error during the retrieval process, log the error and return a 500 Internal Server Error response with an appropriate error message.
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return the retrieved todo item in the response with a 200 OK status.
		c.JSON(http.StatusOK, gin.H{"todo": todo})
	}
}

// UpdateTodoHandler returns a Gin handler function that processes requests to update an existing todo item by its ID. It extracts the ID from the URL parameters, validates it, and calls the repository function to update the todo in the database. If the ID is invalid, it responds with a 400 Bad Request status. If there is an error during the update process, it logs the error and responds with a 500 Internal Server Error status. If the todo is not found, it responds with a 404 Not Found status. If the todo is updated successfully, it returns a 200 OK response with the updated todo item in the response body.
func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the ID from the URL parameters and attempt to convert it to an integer. If the conversion fails, return a 400 Bad Request response with an error message indicating that the ID is invalid.
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Bind the incoming JSON payload to a CreateTodoRequest struct. This will validate that the required fields are present and correctly formatted. If the binding fails, return a 400 Bad Request response with the error message.
		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the UpdateTodo function from the repositories package to update the existing todo item in the database using the provided ID, title, and completion status. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message. If the todo item is not found (i.e., the repository function returns nil), return a 404 Not Found response with an error message indicating that the todo item was not found. If the todo item is updated successfully, return a 200 OK response with the updated todo item in the response body.
		todo, err := repositories.UpdateTodo(pool, id, req.Title, req.Completed)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			log.Printf("UpdateTodo error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"todo": todo})
	}
}
