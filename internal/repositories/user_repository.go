package repositories

import (
	"context"
	"time"

	"github.com/carloscfgos1980/todo-auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUser inserts a new user into the database and returns the created user with its ID and timestamps.
func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to insert a new user and return the created record
	var query string = `
	INSERT INTO users (email, password) 
	VALUES ($1, $2) 
	RETURNING id, email, created_at, updated_at`
	// Execute the query and scan the returned record into a User model
	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the created user
	return user, nil
}

// GetUserByEmail retrieves a user from the database by their email address. If a user with the specified email is found, it returns the user; otherwise, it returns an error indicating that the user was not found.
func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to select a user by email
	query := `
			SELECT id, email, password, created_at, updated_at
			FROM users
			WHERE email = $1
		`
	// Execute the query and scan the returned record into a User model
	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the found user
	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*models.User, error) {
	// Set a timeout for the database query to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SQL query to select a user by ID
	query := `
			SELECT id, email, password, created_at, updated_at
			FROM users
			WHERE id = $1
		`
	// Execute the query and scan the returned record into a User model
	var user models.User
	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the found user
	return &user, nil
}
