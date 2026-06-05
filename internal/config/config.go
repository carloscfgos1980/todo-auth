package config

import (
	"errors"
	"os"

	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/joho/godotenv"
)

var (
	ErrMissingDatabaseURL = errors.New("missing database URL")
	ErrMissingPort        = errors.New("missing port")
	ErrMissingJWT         = errors.New("missing JWT secret")
)

type Config struct {
	DB          *database.Queries
	DatabaseURL string
	Port        string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	// Support running from project root and from cmd/.
	_ = godotenv.Load(".env", "../.env")

	DatabaseURL := getEnv("DatabaseURL", "DATABASE_URL")
	if DatabaseURL == "" {
		return nil, ErrMissingDatabaseURL
	}

	Port := os.Getenv("PORT")
	if Port == "" {
		return nil, ErrMissingPort
	}

	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		return nil, ErrMissingJWT
	}

	// Return the configuration struct with the loaded values
	return &Config{
		DatabaseURL: DatabaseURL,
		Port:        Port,
		JWTSecret:   JWTSecret,
	}, nil
}

func getEnv(keys ...string) string {
	for _, key := range keys {
		value := os.Getenv(key)
		if value != "" {
			return value
		}
	}

	return ""
}
