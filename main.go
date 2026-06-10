package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/carloscfgos1980/todo-auth/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// apiConfig holds the dependencies for the API handlers.
type apiConfig struct {
	db        *database.Queries
	jwtSecret string
	port      string
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()
	// Get configuration from environment variables
	dbURL := os.Getenv("DatabaseURL")
	if dbURL == "" {
		log.Fatal("DatabaseURL must be set")
	}
	// Get the port from environment variables, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	// Get the JWT secret from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	// Connect to the database
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer dbConn.Close()

	// database queries variable
	dbQueries := database.New(dbConn)
	// variable for the apiConfig struct
	apiCfg := apiConfig{
		db:        dbQueries,
		port:      port,
		jwtSecret: jwtSecret,
	}
	// Set up the HTTP server and routes
	mux := http.NewServeMux()
	// health check endpoint
	mux.HandleFunc("/health", apiCfg.handlerHealth)
	// user creation endpoint
	mux.HandleFunc("/auth/register", apiCfg.handlerUsersCreate)

	log.Printf("Starting server on port %s", apiCfg.port)
	// Start the HTTP server
	err = http.ListenAndServe(":"+apiCfg.port, mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
