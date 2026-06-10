package users

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// structs and handler for creating a new user in the system
type User struct {
	ID        pgtype.UUID `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Email     string      `json:"email"`
}

// UserRequest is the struct for the request body when creating a new user
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the response body when logging in a user.
type LoginResponse struct {
	Token string `json:"token"`
}
