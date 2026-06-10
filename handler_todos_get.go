package main

import (
	"net/http"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
)

// createTodosGetHandler handles the GET /todos endpoint to retrieve the list of todos for the authenticated user.
func (cfg *apiConfig) createTodosGetHandler(w http.ResponseWriter, r *http.Request) {
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
