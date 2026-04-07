package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware returns a Gin middleware function that validates JWT tokens in the Authorization header of incoming requests. It checks for the presence of the token, verifies its format, and validates it using the secret key from the configuration. If the token is valid, it extracts the user ID from the token claims and sets it in the Gin context for use in subsequent handlers. If the token is missing, invalid, or expired, it returns a 401 Unauthorized response with an appropriate error message.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header from the incoming request. If the header is missing, return a 401 Unauthorized response with an appropriate error message indicating that the Authorization header is required.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// Check that the Authorization header is in the correct format (i.e., starts with "Bearer " followed by the token). If the header does not match this format, return a 401 Unauthorized response with an appropriate error message indicating that the token format is invalid.
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		// Extract the token string from the Authorization header and trim any leading or trailing whitespace. If the token string is empty after trimming, return a 401 Unauthorized response with an appropriate error message indicating that the token format is invalid.
		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		// Parse and validate the JWT token using the secret key from the configuration. If the token is invalid or expired, return a 401 Unauthorized response with an appropriate error message indicating that the token is invalid or has expired.
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		// If there is an error during token parsing or validation, check if the error is due to token expiration and return a specific error message for expired tokens. For other types of errors, return a generic error message indicating that the token is invalid.
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}
		// Check if the token is valid. If the token is not valid, return a 401 Unauthorized response with an appropriate error message indicating that the token is invalid.
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Extract the user ID from the token claims and set it in the Gin context for use in subsequent handlers. If the user ID is not present in the claims or is of an unexpected type, return a 401 Unauthorized response with an appropriate error message indicating that the token claims are invalid.
		var userID string
		switch v := claims["userID"].(type) {
		case string:
			userID = v
		case float64:
			userID = fmt.Sprintf("%.0f", v)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token claims"})
			c.Abort()
			return
		}
		// Set the extracted user ID in the Gin context with the key "userID" for use in subsequent handlers. This allows other handlers to access the user ID associated with the authenticated request.
		c.Set("userID", userID)
		c.Next()
	}
}
