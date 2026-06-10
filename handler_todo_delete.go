package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
)

// handlerTodoDelete handles the endpoint for deleting a specific todo item. It validates the user's authorization using JWT tokens, checks if the authenticated user is the owner of the todo, and deletes the todo from the database if all checks pass.
func (cfg *apiConfig) handlerTodoDelete(w http.ResponseWriter, r *http.Request) {
	// Validate the user's authorization to delete the todo by checking the provided JWT token
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
	// Get the existing todo from the database to check if it exists and if the authenticated user is the owner of the todo
	dbTodo, err := cfg.db.GetTodoByID(r.Context(), int32(todoID))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Todo not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error getting todo", err)
		return
	}
	// Check if the authenticated user is the owner of the todo
	if dbTodo.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You don't have permission to delete this todo", fmt.Errorf("you don't have permission to delete this todo"))
		return
	}
	// Delete the todo from the database using the authenticated user's ID and the todo ID
	err = cfg.db.DeleteTodo(r.Context(), int32(todoID))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Todo not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete todo", err)
		return
	}
	// Return a success response indicating that the todo was deleted successfully
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Todo deleted successfully"})
}
