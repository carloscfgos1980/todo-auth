package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/carloscfgos1980/todo-auth/internal/models"
	"github.com/carloscfgos1980/todo-auth/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// UserRequest represents the expected JSON payload for user registration and login requests, containing an email and a password with validation rules.
type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the JSON response containing a JWT token returned upon successful user login.
type LoginResponse struct {
	Token string `json:"token"`
}

// CreateUserHandler returns a Gin handler function that processes requests to create a new user.
func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the incoming JSON payload to a UserRequest struct and validate the input. If there is an error during binding or validation, return a 400 Bad Request response with an appropriate error message.
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate that the password is at least 6 characters long. If the password does not meet this requirement, return a 400 Bad Request response with an appropriate error message.
		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long!"})
			return
		}
		// Hash the user's password using bcrypt before storing it in the database. If there is an error during hashing, log the error and return a 500 Internal Server Error response with an appropriate error message.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Create a new User model with the provided email and the hashed password, then call the CreateUser function from the repositories package to insert the new user into the database. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
		user := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		}
		// Call the CreateUser function from the repositories package to insert the new user into the database. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
		createdUser, err := repositories.CreateUser(pool, user)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return a 200 OK response with the created user (excluding the password) in the response body.
		c.JSON(http.StatusOK, createdUser)

	}
}

// LoginHandler returns a Gin handler function that processes user login requests.
func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the incoming JSON payload to a UserRequest struct and validate the input. If there is an error during binding or validation, return a 400 Bad Request response with an appropriate error message.
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Fetch the user from the database using the provided email. If the user is not found or there is an error during this process, log the error and return a 401 Unauthorized response with an appropriate error message indicating that the email or password is invalid.
		user, err := repositories.GetUserByEmail(pool, req.Email)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Compare the provided password with the hashed password stored in the database using bcrypt. If the passwords do not match, log the error and return a 401 Unauthorized response with an appropriate error message indicating that the email or password is invalid.
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			log.Printf("Error comparing passwords: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// If the email and password are valid, generate a JWT token containing the user's ID and email as claims. The token should have an expiration time of 24 hours. Sign the token using the secret key specified in the configuration. If there is an error during token generation or signing, log the error and return a 500 Internal Server Error response with an appropriate error message.
		claims := jwt.MapClaims{
			"userID": user.ID,
			"email":  user.Email,
			"exp":    time.Now().Add(24 * time.Hour).Unix(),
		}

		// Create a new JWT token with the specified claims and sign it using the secret key from the configuration. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Printf("Error signing token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return a 200 OK response with the generated JWT token in the response body.
		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})

	}
}

// TestProtectedHandler returns a Gin handler function that serves as a test endpoint for verifying the functionality of the authentication middleware. It checks for the presence of a user ID in the Gin context (set by the AuthMiddleware) and returns a JSON response indicating that protected content has been accessed, along with the user ID. If the user ID is not found in the context, it returns a 500 Internal Server Error response with an appropriate error message.
func TestProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Protected content accessed", "userID": userID})
	}
}
