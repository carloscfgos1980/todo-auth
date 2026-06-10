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

## 3. Login

1. Service
1.1 GetUserByEmail gets a user from the database by email
1.2 Add GetUserByEmail
 to service interface
2. LoginUser handles the HTTP request for logging in a user
2.1 Parse the JSON request body into a UserRequest struct
2.2 Check if email and password are provided
2.3 Get the user by email from the database
2.4 Check if the provided password matches the stored hashed password
2.5 Generate a JWT token for the authenticated user
2.6 Create a response struct to send back to the client with the access token
2.7 Write the response as JSON with a 200 OK status code
3. users endpoints
 // create the user service and handler
 userService := users.NewService(database.New(app.db), app.db)
 userHandler := users.NewHandler(userService, app.config.JWTSecret)
 // set up the users routes
 r.Route("/auth", func(r chi.Router) {
  r.Post("/register", userHandler.CreateUser)
  r.Post("/login", userHandler.LoginUser)
 })

## 4. Auth middleware

1. HTTP middleware setting a value on the request context
1.1 Return a new http.HandlerFunc that wraps the original handler and adds the authentication logic
1.2 Extract the token from the Authorization header
1.3 Validate the token and extract the user ID
1.4 Create a new context with the user ID value
1.5 Call the next handler with the new context
2. protected routes
 r.Route("/api", func(r chi.Router) {
  // Add authentication middleware here if available
  r.Use(func(next http.Handler) http.Handler {
   return authmiddleware.AuthMiddleware(next, app.config.JWTSecret)
  })

## 5. Create todo

1. Types
1.1 CreateTodoRequest represents the expected JSON payload for creating a new todo item. It includes the title of the todo and its completion status.
1.2 ResponseTodo represents the structure of a todo item that will be sent back in the response. It includes the ID, user ID, title, and completion status of the todo item.
2. service set up
2.1 Service defines the interface for the users service
2.2 svc defines the struct for the users service
2.3 NewService creates a new service for the users package
3. GetUserByID retrieves a user by their ID from the database
4. CreateTodo creates a new todo item in the database
5. Add GetUserByID and CreateTodo to service interface
6. Setup handler
6.1 andler is the HTTP handler for users endpoints
6.2 NewHandler creates a new handler for users endpoints
7. CreateTodo handles the HTTP request for creating a new todo item
7.1 Get the user ID from the request context (set by the authentication middleware)
7.2 The auth middleware stores the JWT subject as a UUID in the request context.
7.3 Check if the user exists in the database
7.4 Parse the JSON request body into a CreateUpdateTodoRequest struct
7.5 validate the request data
7.6 Create a CreateTodoParams struct to pass to the service layer
7.7 Call the service to create the todo item
7.8 Create a ResponseTodo struct to send back to the client
7.9 Write the created todo item as JSON with a 201 Created status code
8. protected routes
 r.Route("/api", func(r chi.Router) {
  // Add authentication middleware here if available
  r.Use(func(next http.Handler) http.Handler {
   return authmiddleware.AuthMiddleware(next, app.config.JWTSecret)
  })
  // create the todo service and handler
  todoService := todos.NewService(database.New(app.db), app.db)
  todoHandler := todos.NewHandler(todoService, app.config.JWTSecret)
  // set up the todos routes
  r.Post("/todos", todoHandler.CreateTodo)

## 6. Get todos

1. GetTodos retrieves all todo items for a given user ID from the database
2. Add GetTodos method to service interface
3. GetTodos handles the HTTP request for retrieving all todo items for the authenticated user
3.1 Get the user ID from the request context (set by the authentication middleware)
3.2 Check if the user ID is present in the context
3.3 The auth middleware stores the JWT subject as a UUID in the request context.
3.4 convert the user ID to pgtype.UUID
3.5 Check if the user exists in the database
3.6 Call the service to get all todo items for the authenticated user
3.7 Create a slice of ResponseTodo to send back to the client
3.8 Write the todos as JSON with a 200 OK status code
4. protected routes
 r.Route("/api", func(r chi.Router) {
  // Add authentication middleware here if available
  r.Use(func(next http.Handler) http.Handler {
   return authmiddleware.AuthMiddleware(next, app.config.JWTSecret)
  })
  // create the todo service and handler
  todoService := todos.NewService(database.New(app.db), app.db)
  todoHandler := todos.NewHandler(todoService, app.config.JWTSecret)
  // set up the todos routes

  r.Get("/todos", todoHandler.GetTodos)

## 7. Get a todo

1. GetTodoByID retrieves a todo item by its ID from the database
2. Add GetTodoByID to service interface
3. GetTodoByID handles the HTTP request for retrieving a specific todo item by its ID for the authenticated user
3.1 Get the user ID from the request context (set by the authentication middleware)
3.2 Check if the user ID is present in the context
3.3 The auth middleware stores the JWT subject as a UUID in the request context.
3.4 Check if the user exists in the database
3.5 Get the todo ID from the URL parameters
3.6 Convert the todo ID from string to int
3.7 Call the service to get the todo item by ID
3.8 convert the user ID to pgtype.UUID
3.9 check if the todo item belongs to the authenticated user
3.10 Create a ResponseTodo struct to send back to the client
3.11 Write the todo item as JSON with a 200 OK status code
4. protected routes
 r.Route("/api", func(r chi.Router) {
  // Add authentication middleware here if available
  r.Use(func(next http.Handler) http.Handler {
   return authmiddleware.AuthMiddleware(next, app.config.JWTSecret)
  })
  // create the todo service and handler
  todoService := todos.NewService(database.New(app.db), app.db)
  todoHandler := todos.NewHandler(todoService, app.config.JWTSecret)

  r.Get("/todos/{todoID}", todoHandler.GetTodoByID)

## 8. Update todo

1. UpdateTodo updates a todo item in the database
2. Add UpdateTodo to service interface
3. UpdateTodo handles the HTTP request for updating a specific todo item by its ID for the authenticated user
3.1 Get the user ID from the request context (set by the authentication middleware)
3.2 Check if the user ID is present in the context
3.3 The auth middleware stores the JWT subject as a UUID in the request context.
3.4 Check if the user exists in the database
3.5 Get the todo ID from the URL parameters
3.6 Convert the todo ID from string to int
3.7 Parse the JSON request body into a CreateUpdateTodoRequest struct
3.8 get the existing todo item from the database to check if it belongs to the authenticated user and to get the current values of the fields
3.9 check if the todo item belongs to the authenticated user
3.10 Update the fields of the todo item based on the request data and the existing values in the database. If a field is not provided in the request, keep the existing value from the database.
3.11 Create an UpdateTodoParams struct to pass to the service layer
3.12 Call the service to update the todo item
3.13 Create a ResponseTodo struct to send back to the client
3.14 Write the updated todo item as JSON with a 200 OK status code
4. create the todo service and handler
  todoService := todos.NewService(database.New(app.db), app.db)
  todoHandler := todos.NewHandler(todoService, app.config.JWTSecret)

  r.Put("/todos/{todoID}", todoHandler.UpdateTodo)
