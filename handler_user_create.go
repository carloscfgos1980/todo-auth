package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/google/uuid"
)

// structs and handler for creating a new user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// handlerUsersCreate handles the user registration endpoint, allowing clients to create new user accounts by providing an email and password. It validates the input, hashes the password, and stores the new user in the database.
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	// Define the expected parameters for creating a new user and the response structure
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Decode the JSON request body into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// Validate the provided parameters (e.g., check for valid email format, strong password, etc.)
	err = auth.IsValidEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	// strong password validation can be added here before hashing the password and creating the user in the database
	err = auth.IsStrongPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Hash the user's password before storing it in the database
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	// Create a new user in the database using the provided parameters and the hashed password
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			respondWithError(w, http.StatusConflict, "User with this email already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	// Prepare the response with the created user's information (excluding the password)
	response := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	}
	// Respond with the created user's information (excluding the password)
	respondWithJSON(w, http.StatusCreated, response)

}
