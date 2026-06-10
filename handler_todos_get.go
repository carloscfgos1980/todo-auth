package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
)

// handlerTodosGet handles the GET /todos endpoint to retrieve the list of todos for the authenticated user.
func (cfg *apiConfig) handlerTodosGet(w http.ResponseWriter, r *http.Request) {
	// Validate the user's authorization to get the list of todos by checking the provided JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authorization token included", err)
		return
	}
	// Validate the JWT token and extract the user ID from it
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization token", err)
		return
	}
	// Get the list of todos associated with the user's ID from the database and return it in the response
	todos, err := cfg.db.GetTodosByUserID(r.Context(), userID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithJSON(w, http.StatusOK, []responseTodo{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error getting todos", err)
		return
	}
	// Convert the list of todos to the response format and return it in the response
	response := []responseTodo{}
	for _, todo := range todos {
		response = append(response, responseTodo{
			ID:        todo.ID,
			UserID:    todo.UserID,
			Title:     todo.Title,
			Completed: todo.Completed,
		})
	}
	// Return the list of todos in the response
	respondWithJSON(w, http.StatusOK, response)
}

// handlerTodoByID handles the GET /todos/{todoID} endpoint to retrieve a specific todo by its ID for the authenticated user.
func (cfg *apiConfig) handlerTodoByID(w http.ResponseWriter, r *http.Request) {
	// Validate the user's authorization to get the todo by ID by checking the provided JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authorization token included", err)
		return
	}
	// Validate the JWT token and extract the user ID from it
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization token", err)
		return
	}
	// Get the todo ID from the URL parameters and validate it
	todoIDStr := r.PathValue("todoID")
	if todoIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Todo ID is required", fmt.Errorf("todo ID is required"))
		return
	}
	// Convert the todo ID from string to integer
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID format", err)
		return
	}
	// Get the todo associated with the user's ID and the provided todo ID from the database and return it in the response
	todo, err := cfg.db.GetTodoByID(r.Context(), int32(todoID))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Todo not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error getting todo", err)
		return
	}
	// Check if the todo belongs to the authenticated user
	if todo.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You don't have permission to access this todo", fmt.Errorf("you don't have permission to access this todo"))
		return
	}
	// Convert the todo to the response format and return it in the response
	response := responseTodo{
		ID:        todo.ID,
		UserID:    todo.UserID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
	// Return the todo in the response
	respondWithJSON(w, http.StatusOK, response)
}
