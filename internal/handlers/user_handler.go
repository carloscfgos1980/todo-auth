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

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// CreateUserHandler returns a Gin handler function that processes requests to create a new user.
func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long!"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repositories.CreateUser(pool, user)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, createdUser)

	}
}

// LoginHandler returns a Gin handler function that processes user login requests.
func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
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
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		// Create a new JWT token with the specified claims and sign it using the secret key from the configuration. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Printf("Error signing token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})

	}
}
