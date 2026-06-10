package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/carloscfgos1980/todo-auth/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Define the expected parameters for user login and the response structure
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		Token string `json:"token"`
	}

	// Decode the JSON request body into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// Retrieve the user from the database using the provided email address
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Check if the provided password matches the hashed password stored in the database for the retrieved user
	match, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// If the password is correct, generate a JWT token for the user to authenticate future requests
	token, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		24*7*time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	// Respond with the generated JWT token
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})

}
