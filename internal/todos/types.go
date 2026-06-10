package todos

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateUpdateTodoRequest represents the expected JSON payload for creating or updating a todo item. It includes the title of the todo and its completion status.
type CreateUpdateTodoRequest struct {
	Title     *string `json:"title" binding:"required"`
	Completed *bool   `json:"completed"`
}

// ResponseTodo represents the structure of a todo item that will be sent back in the response. It includes the ID, user ID, title, and completion status of the todo item.
type ResponseTodo struct {
	ID        int32       `json:"id"`
	UserID    pgtype.UUID `json:"user_id"`
	Title     string      `json:"title"`
	Completed bool        `json:"completed"`
}
