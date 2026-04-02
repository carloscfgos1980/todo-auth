package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPostgres establishes a connection to the PostgreSQL database using the provided database URL. It returns a connection pool that can be used for executing queries. If there is an error during the connection process, it logs the error and returns it to the caller.
func ConnectPostgres(databaseURL string) (*pgxpool.Pool, error) {
	// Create a context for the database connection. This can be used to set timeouts or cancel the connection if needed.
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Error parsing database URL: %v", err)
		return nil, err
	}
	// Create a new connection pool using the parsed configuration. If there is an error creating the pool, log the error and return it.
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Printf("Error creating database pool: %v", err)
		return nil, err
	}
	// Ping the database to ensure that the connection is established successfully. If there is an error pinging the database, log the error, close the pool, and return the error.
	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		pool.Close()
		return nil, err
	}

	// If the connection is successful, return the connection pool to the caller.
	return pool, nil
}
