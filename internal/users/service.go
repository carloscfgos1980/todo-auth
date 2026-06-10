package users

import (
	"context"

	"github.com/carloscfgos1980/todo-auth/internal/database"

	"github.com/jackc/pgx/v5"
)

// Service defines the interface for the users service
type Service interface {
	CreateUser(ctx context.Context, user UserRequest) (*database.User, error)
	GetUserByEmail(ctx context.Context, email string) (*database.User, error)
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

// CreateUser creates a new user in the database
func (s *svc) CreateUser(ctx context.Context, user UserRequest) (*database.User, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)

	// create the user
	createdUser, err := qtx.CreateUser(ctx, database.CreateUserParams{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	// return the created user
	return &createdUser, nil
}

// GetUserByEmail gets a user from the database by email
func (s *svc) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	// start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// create a new Queries instance with the transaction
	qtx := s.repo.WithTx(tx)
	// get the user by email
	user, err := qtx.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	// return the user
	return &user, nil

}
