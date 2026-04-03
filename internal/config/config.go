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
	godotenv.Load()
	DatabaseURL := os.Getenv("DatabaseURL")
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

	return &Config{
		DatabaseURL: DatabaseURL,
		Port:        Port,
		JWTSecret:   JWTSecret,
	}, nil
}
