package models

import "time"

// User represents a user in the system with fields for ID, email, password, and timestamps for creation and updates. The password field is excluded from JSON responses for security reasons.
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
