package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrMissingDatabaseURL = errors.New("missing database URL")
	ErrMissingPort        = errors.New("missing port")
)

type Config struct {
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
		return nil, errors.New("missing JWT secret")
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
