package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/carloscfgos1980/todo-auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTodo inserts a new todo into the database and returns the created todo with its ID and timestamps.
func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to insert a new todo and return the created record
	var query string = `
	INSERT INTO todos (title, completed, user_id) 
	VALUES ($1, $2, $3) 
	RETURNING id, title, completed, created_at, updated_at, user_id`
	// Execute the query and scan the returned record into a Todo model
	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)
	if err != nil {
		return nil, err
	}
	// Return the created todo
	return &todo, nil
}

// GetTodos retrieves all todos from the database, ordered by creation date in descending order.
func GetTodos(pool *pgxpool.Pool, userID string) ([]models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to select all todos, ordered by creation date in descending order
	query := `
			SELECT id, title, completed, created_at, updated_at, user_id
			FROM todos
			WHERE user_id = $1
			ORDER BY created_at DESC
		`
	// Execute the query and get the rows
	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Iterate through the rows and scan the data into a slice of Todo models
	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserID,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	// Check for any errors that occurred during iteration
	return todos, nil
}

// GetTodoByID retrieves a specific todo from the database by its ID. If the todo is not found, it returns nil and an error.
func GetTodoByID(pool *pgxpool.Pool, id int, userID string) (*models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, title, completed, created_at, updated_at, user_id
		FROM todos
		WHERE id = $1 AND user_id = $2
	`
	// Execute the query and scan the returned record into a Todo model
	var todo models.Todo
	err := pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)
	if err != nil {
		return nil, err
	}
	// Return the retrieved todo
	return &todo, nil
}

// UpdateTodo updates an existing todo in the database with the provided title and completion status. It returns the updated todo with its ID and timestamps.
func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, userID string) (*models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		UPDATE todos
		SET title = $1, completed = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, completed, created_at, updated_at, user_id
	`
	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// DeleteTodo deletes a specific todo from the database by its ID. If the todo is not found, it returns an error indicating that the todo with the specified ID was not found. If the deletion is successful, it returns nil.
func DeleteTodo(pool *pgxpool.Pool, id int, userID string) error {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		DELETE FROM todos
		WHERE id = $1 AND user_id = $2
	`
	// Execute the query and check the number of rows affected to determine if the todo was found and deleted
	commandTag, err := pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	// If no rows were affected, it means the todo with the specified ID was not found
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("todo with id %d not found for user %s", id, userID)
	}
	return nil
}
