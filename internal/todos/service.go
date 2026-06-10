package todos

import (
	"context"

	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service defines the interface for the users service
type Service interface {
	CreateTodo(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error)
	GetUserByID(ctx context.Context, id string) (*database.User, error)
}

// svc defines the struct for the users service
type svc struct {
	repo *database.Queries
	db   *pgx.Conn
}

// NewService creates a new service for the users package
func NewService(repo *database.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// GetUserByID retrieves a user by their ID from the database
func (s *svc) GetUserByID(ctx context.Context, id string) (*database.User, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)
	userID := pgtype.UUID{}
	err = userID.Scan(id)
	if err != nil {
		return nil, err
	}
	user, err := qtx.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateTodo creates a new todo item in the database
func (s *svc) CreateTodo(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)

	// create the todo
	createdTodo, err := qtx.CreateTodo(ctx, database.CreateTodoParams{
		Title:     todo.Title,
		Completed: todo.Completed,
		UserID:    todo.UserID,
	})
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	// return the created todo
	return &createdTodo, nil
}
