package models

import "time"

// Todo represents a task or item in a to-do list, with fields for ID, title, user ID, completion status, and timestamps for creation and updates. The struct tags specify how the fields should be serialized to JSON and mapped to database columns.
type Todo struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	UserID    string    `json:"user_id" db:"user_id"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
