# STEPS

## 1. Set up

1. Start the app
go mod init github.com/carloscfgos1980/todo-auth

2. packages
github.com/joho/godotenv
github.com/go-chi/chi/v5
github.com/jackc/pgx/v5
github.com/alexedwards/argon2id
github.com/golang-jwt/jwt/v5
github.com/google/uuid

3. Copy files from previous version: .env, .gitignore and yaml 
4. Copy "sql" directory from previous version
5. Run command "sqlc generate" to convert SQL to Go
6. Copy auxialiary directories: utils, json, env to "internal" directory
7. Api set up
7.1 application is the main application struct that holds the configuration and database connection
7.2 config holds the configuration for the application
7.3 dbConfig holds the database configuration for the application
8. mount sets up the routes and middleware for the application
8.1 create a new router
8.2 set up middleware
8.3 health check endpoint
8.4 return the router
9. run starts the HTTP server
9.1 create the HTTP server
9.2 start the server
10. Main
10.1 Load environment variables from .env file
10.2 Get the port from environment variables, default to 8080 if not set
10.3 Get the JWT secret from environment variables
10.4 create a context
10.5 load env variables
10.6 initialize logger
10.7 database connection
10.8 reate the application
10.9 run the application
