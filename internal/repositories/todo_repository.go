package repositories

import (
	"context"
	"time"

	"github.com/carloscfgos1980/todo-auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTodo inserts a new todo into the database and returns the created todo with its ID and timestamps.
func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to insert a new todo and return the created record
	var query string = `
	INSERT INTO todos (title, completed, created_at, updated_at) 
	VALUES ($1, $2, NOW(), NOW()) 
	RETURNING id, title, completed, created_at, updated_at`
	// Execute the query and scan the returned record into a Todo model
	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the created todo
	return &todo, nil
}

// GetTodos retrieves all todos from the database, ordered by creation date in descending order.
func GetTodos(pool *pgxpool.Pool) ([]models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to select all todos, ordered by creation date in descending order
	query := `
			SELECT id, title, completed, created_at, updated_at
			FROM todos
			ORDER BY created_at DESC
		`
	// Execute the query and get the rows
	rows, err := pool.Query(ctx, query)
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
func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, title, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`
	// Execute the query and scan the returned record into a Todo model
	var todo models.Todo
	err := pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the retrieved todo
	return &todo, nil
}

// UpdateTodo updates an existing todo in the database with the provided title and completion status. It returns the updated todo with its ID and timestamps.
func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool) (*models.Todo, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		UPDATE todos
		SET title = $1, completed = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, completed, created_at, updated_at
	`
	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int) error {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		DELETE FROM todos
		WHERE id = $1
	`
	_, err := pool.Exec(ctx, query, id)
	return err
}
