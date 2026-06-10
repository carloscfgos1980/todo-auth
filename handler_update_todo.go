package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
	"github.com/carloscfgos1980/todo-auth/internal/database"
)

// handlerTodoUpdate handles the HTTP request for updating a todo item. It validates the user's authorization, checks the provided parameters, and updates the todo in the database if everything is valid.
func (cfg *apiConfig) handlerTodoUpdate(w http.ResponseWriter, r *http.Request) {
	// Validate the user's authorization to update the todo by checking the provided JWT token
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
	// Get the updated todo information from the request body and validate it
	type parameters struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	// Decode the request body into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// Check if at least one of the parameters is provided
	if params.Title == nil && params.Completed == nil {
		respondWithError(w, http.StatusBadRequest, "At least one parameter (title or completed) must be provided", fmt.Errorf("at least one parameter (title or completed) must be provided"))
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
		respondWithError(w, http.StatusForbidden, "You don't have permission to update this todo", fmt.Errorf("you don't have permission to update this todo"))
		return
	}
	// Update the todo fields with the provided parameters if they are not nil
	title := dbTodo.Title
	if params.Title != nil {
		title = *params.Title
	}
	completed := dbTodo.Completed
	if params.Completed != nil {
		completed = *params.Completed
	}
	// Update the todo in the database using the provided parameters and the authenticated user's ID
	todo, err := cfg.db.UpdateTodo(r.Context(), database.UpdateTodoParams{
		ID:        int32(todoID),
		Title:     title,
		Completed: completed,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Todo not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update todo", err)
		return
	}
	// Prepare the response with the updated todo's information
	response := responseTodo{
		ID:        todo.ID,
		UserID:    todo.UserID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
	// Return the updated todo in the response
	respondWithJSON(w, http.StatusOK, response)
}
