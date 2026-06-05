package handlers

import (
	"net/http"
	"time"

	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/carloscfgos1980/todo-auth/internal/utils"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

// structs and handler for creating a new user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// UserRequest is the struct for the request body when creating a new user
type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateUserHandler is the handler for creating a new user in the system
func CreateUserHandler(cfg *config.Config) gin.HandlerFunc {
	// Return a handler function that can be used in the Gin router
	return func(c *gin.Context) {
		// Bind the JSON request body to the UserRequest struct
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate email format
		err := utils.IsValidEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate the password strength
		err = utils.IsStrongPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Hash the password before storing it in the database
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Create the user in the database using the provided configuration and request data
		user, err := cfg.DB.CreateUser(c, database.CreateUserParams{
			Email:    req.Email,
			Password: hashedPassword,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return the created user as a response, excluding the password
		response := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			Email:     user.Email,
		}
		// Send the response back to the client with a 200 OK status
		c.JSON(http.StatusOK, gin.H{"user": response})
	}
}

// LoginUserHandler is the handler for logging in a user and generating a JWT token
func LoginUserHandler(cfg *config.Config) gin.HandlerFunc {
	// Define a struct for the response that will be sent back to the client after successful login
	type response struct {
		Token string `json:"token"`
	}
	// Return a handler function that can be used in the Gin router
	return func(c *gin.Context) {
		// Bind the JSON request body to the UserRequest struct
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate email format
		if err := utils.IsValidEmail(req.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Retrieve the user from the database using the provided email
		user, err := cfg.DB.GetUserByEmail(c, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		// Check if the provided password matches the stored hashed password
		match, err := utils.CheckPasswordHash(req.Password, user.Password)
		if err != nil || !match {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		// Generate a JWT token for the authenticated user
		token, err := utils.MakeJWT(
			user.ID,
			cfg.JWTSecret,
			24*7*time.Hour,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
		// Create the response struct with the generated token
		response := response{
			Token: token,
		}
		// Send the response back to the client with a 200 OK status
		c.JSON(http.StatusOK, response)
	}
}
