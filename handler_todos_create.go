package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
	"github.com/carloscfgos1980/todo-auth/internal/database"
)

// responseTodo is a struct that represents the JSON response for a created task. It includes the task's ID, user ID, title, and completion status.
type responseTodo struct {
	ID        int32     `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}

// handlerCreateTodo handles requests to create a new task. It validates the user's authorization using a JWT token, checks the provided parameters for creating a new task, and creates the task in the database if everything is valid.
func (cfg *apiConfig) handlerCreateTodo(w http.ResponseWriter, r *http.Request) {
	// Define the expected parameters for creating a new task
	type parameters struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	// Validate the user's authorization to create a new task by checking the provided JWT token
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
	// Decode the JSON request body into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	// Validate the provided parameters for creating a new task (e.g., check if priority, state, tag, and date formats are valid)
	if params.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Title is required", fmt.Errorf("title is required"))
		return
	}
	// check if the user is register in database
	_, err = cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}
	// If the parameters are valid and the user is authorized, create a new task in the database associated with the user's ID and the provided parameters
	todo, err := cfg.db.CreateTodo(r.Context(), database.CreateTodoParams{
		UserID:    userID,
		Title:     params.Title,
		Completed: params.Completed,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create todo", err)
		return
	}
	// Respond with the created task's information
	response := responseTodo{
		ID:        todo.ID,
		UserID:    todo.UserID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
	// Respond with the created task's information in JSON format
	respondWithJSON(w, http.StatusOK, response)
}
