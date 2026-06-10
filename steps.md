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
11. GitHub

```bash
git init
git add .
git commit -m "set up"
git remote add origin https://github.com/carloscfgos1980/todo-auth.git
git checkout -b chi_framework
git push origin chi_framework
```

## 2. Register user

1. types
1.1 structs and handler for creating a new user in the system
1.2 UserRequest is the struct for the request body when creating a new user
1.3 LoginResponse is the response body when logging in a user.
2. service set up
2.1 Service defines the interface for the users service
2.2 svc defines the struct for the users service
2.3 NewService creates a new service for the users package
3. CreateUser creates a new user in the database
4. Add create user method to service interface
5. Handler setup
5.1 handler is the HTTP handler for users endpoints
5.2 NewHandler creates a new handler for users endpoints
6. CreateUser handles the HTTP request for creating a new user
6.1 Parse the JSON request body into a UserRequest struct
6.2 Check if any field is empty
6.3 Validate email format
6.4 Validate the password strength
6.5 Hash the password before storing it in the database
6.6 Update the user request with the hashed password
6.7 Call the service to create the user
6.7.1 Check if the error is a unique constraint violation (duplicate email)
6.8 Create a response struct to send back to the client, excluding the password
6.9 Write the response as JSON with a 201 Created status code
7. users endpoints
 userService := users.NewService(database.New(app.db), app.db)
 userHandler := users.NewHandler(userService, app.config.JWTSecret)
 // set up the users routes
 r.Route("/auth", func(r chi.Router) {
  r.Post("/register", userHandler.CreateUser)
  // r.Post("/login", userHandler.LoginUser)
 })
