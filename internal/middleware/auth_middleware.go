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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var userID string
		switch v := claims["user_id"].(type) {
		case string:
			userID = v
		case float64:
			userID = fmt.Sprintf("%.0f", v)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token claims"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
