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
	GetTodos(ctx context.Context, userID pgtype.UUID) ([]database.Todo, error)
	GetTodoByID(ctx context.Context, id int32) (*database.Todo, error)
	UpdateTodo(ctx context.Context, todo database.UpdateTodoParams) (*database.Todo, error)
	DeleteTodo(ctx context.Context, id int32) error
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

// GetTodos retrieves all todo items for a given user ID from the database
func (s *svc) GetTodos(ctx context.Context, userID pgtype.UUID) ([]database.Todo, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)
	// get the todos for the user
	todos, err := qtx.GetTodosByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return todos, nil
}

// GetTodoByID retrieves a todo item by its ID from the database
func (s *svc) GetTodoByID(ctx context.Context, id int32) (*database.Todo, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)
	// get the todo by ID
	todo, err := qtx.GetTodoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo updates a todo item in the database
func (s *svc) UpdateTodo(ctx context.Context, todo database.UpdateTodoParams) (*database.Todo, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)
	// update the todo item
	updatedTodo, err := qtx.UpdateTodo(ctx, database.UpdateTodoParams{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
	})
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &updatedTodo, nil
}

// DeleteTodo deletes a todo item from the database by its ID
func (s *svc) DeleteTodo(ctx context.Context, id int32) error {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)
	// delete the todo item
	err = qtx.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
