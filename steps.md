# STEPS

## 1. Set up

go mod init github.com/carloscfgos1980/todo-auth

Copy .env and .gitignore from previous version

### packges

go get github.com/joho/godotenv
go get github.com/lib/pq
go get github.com/gin-gonic/gin
go get github.com/google/uuid
go get github.com/alexedwards/argon2id
go get github.com/golang-jwt/jwt/v5

Copy config directory from previous version

### cmd/main.go

1. Load configuration from environment variables
2. Connect to the database
3. Initialize the Gin router
4. Set trusted proxies to nil to avoid warnings in Gin 1.7+
5. Define a simple health check route
6. Start the server on the specified port

### save to github

```bash
git init
git add .
git commit -m "set up"
git remote add origin https://github.com/carloscfgos1980/todo-auth.git
git checkout -b gin_framework
git push origin gin_framework
```

## 2. Register route

1. update config struct. Add DB field for the queries
2. Add queries to the configuration variable
2.1 Create a new database queries instance
2.2 Assign the database queries instance to the configuration struct for use in handlers
3. Copy utils directory from the taskSpehre api
4. Create querie to insert new user in users table
5. Structs to handle request and response
5.1 structs and handler for creating a new user in the system
5.2 UserRequest is the struct for the request body when creating a new user
6. CreateUserHandler is the handler for creating a new user in the system
6.1 Return a handler function that can be used in the Gin router
6.2 Bind the JSON request body to the UserRequest struct
6.3 Validate email format
6.4 Validate the password strength
6.5 Hash the password before storing it in the database
6.6 Create the user in the database using the provided configuration and request data
6.7 Return the created user as a response, excluding the password
7. Register user-related routes
 router.POST("/auth/register", handlers.CreateUserHandler(cfg))

## 3. Login

1. Create queries to get hte user by email and by id
2. sqlc generate to create Go code
3. LoginUserHandler is the handler for logging in a user and generating a JWT token
3.1 Define a struct for the response that will be sent back to the client after successful login
3.2 Return a handler function that can be used in the Gin router
3.3 Bind the JSON request body to the UserRequest struct
3.4 Validate email format
3.5 Retrieve the user from the database using the provided email
3.6 Check if the provided password matches the stored hashed password
3.7 Generate a JWT token for the authenticated user
3.8 Create the response struct with the generated token
3.9 Send the response back to the client with a 200 OK status
4. Register user-related routes
 router.POST("/auth/login", handlers.LoginUserHandler(cfg))
