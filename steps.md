# STEPS

## 1. Setup

1. Start the api
go mod init github.com/carloscfgos1980/todo-auth

2. Packages
github.com/joho/godotenv
github.com/lib/pq
github.com/google/uuid

3. Copy files from previos version:

- sql directory
- yaml file
- json.go
- .env
- .gitignore
- health endpoint method (handler_health.go)

**main**
4. apiConfig holds the dependencies for the API handlers. 
5. main func
5.1 Load environment variables from .env file
5.2 Get configuration from environment variables
5.3 Get the port from environment variables, default to 8080 if not set
5.4 Get the JWT secret from environment variables
5.5 Connect to the database
5.6 database queries variable
5.7 variable for the apiConfig struct
5.8 Set up the HTTP server and routes
5.9 health check endpoint
5.10 Start the HTTP server
