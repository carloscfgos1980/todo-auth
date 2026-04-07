# TODO-AUTH

Todo App with Gin, PostgreSQL & JWT Auth
[REST API for Beginners](https://www.youtube.com/watch?v=S069igHKUIw)
30-03-2026

## Description

In the tutotorial, the todo endpoint is built first and then the user endpoint wich is a bit contraproducent coz later a new migrawtion is needed to alter the todo table in order to include user_id.

This tutorial does not use a package to convert sql queries into go, instead the whole database code is written from scratch. It is helpful to understand but unnecessary

## 1. set up

Get framework for go
go get -u github.com/gin-gonic/gin

Driver for postgres
go get -u github.com/jackc/pgx/v5

package pool... no fucking idea...
go get -u github.com/jackc/pgx/v5/pgxpool

JWT package
go get -u github.com/golang-jwt/jwt/v5

Migrates postgres
go install --tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

Hash password
go get -u golang.org/x/crypto/bcrypt

Get .env
go get -u github.com/joho/godotenv

Install Air: live reloading tool designed specifically for go
go install github.com/air-verse/air@latest

## 2. Config

1. Create error constants
2. Config struct
3. LoadConfig function that takes no argument and return a pointer to config struct and error
4. Load dotgot end
5. Get DatabaseURL, Port and JWT_Secret
6. Return the configuration struct with the loaded values

## 3. Connect postgres function /database/postgres.go

ConnectPostgres(databaseURL string) (*pgxpool.Pool, error)

ConnectPostgres establishes a connection to the PostgreSQL database using the provided database URL. It returns a connection pool that can be used for executing queries. If there is an error during the connection process, it logs the error and returns it to the caller.

## 4. Start server cmd/main.go

1. Load configuration from environment variables
2. Connect to the PostgreSQL database using the provided URL
3. Initialize the Gin router and set up routes and middleware
4. Set trusted proxies to nil to disable Gin's default behavior of trusting all proxies
5. Define a simple health check endpoint at the root URL to verify that the API is running and the database connection is successful
6. Start the Gin server on the specified port from the configuration and log any errors that occur while running the server
 (to be continue)

## 5. Create sql files to run migrations

Fist the db base has to be created from the CLI

* Here it was a bit difficult coz the way the tutorial goes, todos table is already created so when I ran the migration to add a colmun, I had a error since the user_id can not be null

To do this I use goose. the tutorial use a different method

## 6. Models (user.go)

User represents a user in the system with fields for ID, email, password, and timestamps for creation and updates. The password field is excluded from JSON responses for security reasons.

## 7. Function to run the queries /repository/user_repository

1. CreateUser inserts a new user into the database and returns the created user with its ID and timestamps.
2. GetUserByEmail retrieves a user from the database by their email address. If a user with the specified email is found, it returns the user; otherwise, it returns an error indicating that the user was not found.
3. GetUserByID retrieves a user from the database by their unique ID. If a user with the specified ID is found, it returns the user; otherwise, it returns an error indicating that the user was not found.

## 8. User handler handlers/user_handler.go

1. UserRequest represents the expected JSON payload for user registration and login requests, containing an email and a password with validation rules.
2. LoginResponse represents the JSON response containing a JWT token returned upon successful user login.

3. **CreateUserHandler** returns a Gin handler function that processes requests to create a new user.
3.1 Bind the incoming JSON payload to a UserRequest struct and validate the input. If there is an error during binding or validation, return a 400 Bad Request response with an appropriate error message.
3.2 Validate that the password is at least 6 characters long. If the password does not meet this requirement, return a 400 Bad Request response with an appropriate error message.
3.3 Hash the user's password using bcrypt before storing it in the database. If there is an error during hashing, log the error and return a 500 Internal Server Error response with an appropriate error message.
3.4 Create a new User model with the provided email and the hashed password, then call the CreateUser function from the repositories package to insert the new user into the database. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
3.5 Call the CreateUser function from the repositories package to insert the new user into the database. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
3.6 Return a 200 OK response with the created user (excluding the password) in the response body.

4. **LoginHandler** returns a Gin handler function that processes user login requests.
4.1 Bind the incoming JSON payload to a UserRequest struct and validate the input. If there is an error during binding or validation, return a 400 Bad Request response with an appropriate error message.
4.2 Fetch the user from the database using the provided email. If the user is not found or there is an error during this process, log the error and return a 401 Unauthorized response with an appropriate error message indicating that the email or password is invalid.
4.3 Compare the provided password with the hashed password stored in the database using bcrypt. If the passwords do not match, log the error and return a 401 Unauthorized response with an appropriate error message indicating that the email or password is invalid.
4.4 If the email and password are valid, generate a JWT token containing the user's ID and email as claims. The token should have an expiration time of 24 hours. Sign the token using the secret key specified in the configuration. If there is an error during token generation or signing, log the error and return a 500 Internal Server Error response with an appropriate error message
4.5 Create a new JWT token with the specified claims and sign it using the secret key from the configuration. If there is an error during this process, log the error and return a 500 Internal Server Error response with an appropriate error message.
 4.5 Return a 200 OK response with the generated JWT token in the response body.

## 9. Set the routes for the user cmd/main.go

Define authentication routes for user registration and login, and protect the /todos routes with the AuthMiddleware to ensure that only authenticated users can access them

## 10. Create middleware /middleware/auth

AuthMiddleware returns a Gin middleware function that validates JWT tokens in the Authorization header of incoming requests. It checks for the presence of the token, verifies its format, and validates it using the secret key from the configuration. If the token is valid, it extracts the user ID from the token claims and sets it in the Gin context for use in subsequent handlers. If the token is missing, invalid, or expired, it returns a 401 Unauthorized response with an appropriate error message.

1. Extract the Authorization header from the incoming request. If the header is missing, return a 401 Unauthorized response with an appropriate error message indicating that the Authorization header is required.
2. Check that the Authorization header is in the correct format (i.e., starts with "Bearer " followed by the token). If the header does not match this format, return a 401 Unauthorized response with an appropriate error message indicating that the token format is invalid.
3. Extract the token string from the Authorization header and trim any leading or trailing whitespace. If the token string is empty after trimming, return a 401 Unauthorized response with an appropriate error message indicating that the token format is invalid.
4. Parse and validate the JWT token using the secret key from the configuration. If the token is invalid or expired, return a 401 Unauthorized response with an appropriate error message indicating that the token is invalid or has expired.
5. If there is an error during token parsing or validation, check if the error is due to token expiration and return a specific error message for expired tokens. For other types of errors, return a generic error message indicating that the token is invalid.
6. Check if the token is valid. If the token is not valid, return a 401 Unauthorized response with an appropriate error message indicating that the token is invalid.
7. Extract the user ID from the token claims and set it in the Gin context for use in subsequent handlers. If the user ID is not present in the claims or is of an unexpected type, return a 401 Unauthorized response with an appropriate error message indicating that the token claims are invalid.
8. Set the extracted user ID in the Gin context with the key "userID" for use in subsequent handlers. This allows other handlers to access the user ID associated with the authenticated request.

## Check protected routes /handlers/ user_handler.go and set route

TestProtectedHandler returns a Gin handler function that serves as a test endpoint for verifying the functionality of the authentication middleware. It checks for the presence of a user ID in the Gin context (set by the AuthMiddleware) and returns a JSON response indicating that protected content has been accessed, along with the user ID. If the user ID is not found in the context, it returns a 500 Internal Server Error response with an appropriate error message.

**main.go**
Define a protected test endpoint to verify that the AuthMiddleware is working correctly. This endpoint will only be accessible to requests that include a valid JWT token in the Authorization header. If the token is valid, the endpoint will return a success message; otherwise, it will return an unauthorized error response.

## 11. Models for todo modes/todo.go

Todo represents a task or item in a to-do list, with fields for ID, title, user ID, completion status, and timestamps for creation and updates. The struct tags specify how the fields should be serialized to JSON and mapped to database columns.

## 12. Todos queries in go /repository/todo_repository

1. CreateTodo inserts a new todo into the database and returns the created todo with its ID and timestamps.
2. GetTodos retrieves all todos from the database, ordered by creation date in descending order.
3. GetTodoByID retrieves a specific todo from the database by its ID. If the todo is not found, it returns nil and an error.
4. UpdateTodo updates an existing todo in the database with the provided title and completion status. It returns the updated todo with its ID and timestamps.
5. DeleteTodo deletes a specific todo from the database by its ID. If the todo is not found, it returns an error indicating that the todo with the specified ID was not found. If the deletion is successful, it returns nil.

## 13. todo handlers /handlers/todo_handler.go

1. CreateTodoRequest represents the expected JSON payload for creating a new todo item. It includes the title of the todo and its completion status.
2. UpdateTodoRequest represents the expected JSON payload for updating an existing todo item. It includes optional fields for the title and completion status, allowing for partial updates of the todo item. Here using pointer is import to differenciate between an empty field ("") of not fiel at all (nil).

3. **CreateTodoHandler** returns a Gin handler function that handles the creation of a new todo item. It extracts the user ID from the Gin context, binds the incoming JSON payload to a CreateTodoRequest struct, and calls the CreateTodo repository function to create the todo item in the database. If successful, it returns the created todo item in the response with a 201 Created status. If there are any errors during binding or database operations, it returns an appropriate error response.
3.1 Return a Gin handler function that handles the creation of a new todo item.
3.2 Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
3.3 Bind the incoming JSON payload to a CreateTodoRequest struct. If there is an error during binding (e.g., missing required fields or invalid data types), return a 400 Bad Request response with the error message.
3.3 Call the CreateTodo repository function to create the todo item in the database, passing the user ID, title, and completion status. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the todo creation failed. If the todo item is created successfully, return the created todo item in the response with a 201 Created status.
3.4 Return the created todo item in the response with a 201 Created status.

4. **GetTodosHandler** returns a Gin handler function that retrieves all todo items for the authenticated user. It extracts the user ID from the Gin context, calls the GetTodos repository function to fetch the todos from the database, and returns the list of todos in the response with a 200 OK status. If there are any errors during database operations, it logs the error and returns a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of todos failed.
4.1 Same as 3.1
4.2 Same as 3.2
4.3 Call the GetTodos repository function to fetch the todos from the database for the authenticated user. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of todos failed.
4.4 If the todos are retrieved successfully, return the list of todos in the response with a 200 OK status.

5. **GetTodoByIDHandler** returns a Gin handler function that retrieves a specific todo item by its ID for the authenticated user. It extracts the user ID from the Gin context, parses the todo ID from the URL parameter, and calls the GetTodoByID repository function to fetch the todo item from the database. If the todo item is found, it returns the item in the response with a 200 OK status. If the todo item is not found, it returns a 404 Not Found response with an appropriate error message. If there are any errors during database operations, it logs the error and returns a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of the todo item failed.
5.1 Same as 3.1
5.2 Same as 3.2
5.3 Call the GetTodoByID repository function to fetch the todo item from the database for the authenticated user and the specified ID. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of the todo item failed. If the todo item is not found, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found. If the todo item is retrieved successfully, return the item in the response with a 200 OK status.
5.4 Return the retrieved todo item in the response with a 200 OK status.

6. **UpdateTodoHandler** returns a Gin handler function that updates an existing todo item for the authenticated user. It extracts the user ID from the Gin context, parses the todo ID from the URL parameter, and binds the incoming JSON payload to an UpdateTodoRequest struct. It then calls the UpdateTodo repository function to update the todo item in the database. If the update is successful, it returns the updated todo item in the response with a 200 OK status. If the todo item is not found, it returns a 404 Not Found response with an appropriate error message. If there are any errors during binding or database operations, it logs the error and returns an appropriate error response.
6.1 Return a Gin handler function that updates an existing todo item for the authenticated user.
6.2 Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
6.3 Bind the incoming JSON payload to an UpdateTodoRequest struct. If there is an error during binding (e.g., invalid data types), return a 400 Bad Request response with the error message. If neither the title nor the completed fields are provided in the request, return a 400 Bad Request response with an appropriate error message indicating that at least one field must be provided for the update.
6.4 If neither the title nor the completed fields are provided in the request, return a 400 Bad Request response with an appropriate error message indicating that at least one field must be provided for the update. Here it/s crucial the use of pointers
6.5 Call the GetTodoByID repository function to fetch the existing todo item from the database for the authenticated user and the specified ID. This is necessary to ensure that the todo item exists before attempting to update it, and to retrieve the current values of the title and completed fields for use in the update operation. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the retrieval of the existing todo item failed. If the todo item is not found, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found.
6.6 Determine the new values for the title and completed fields based on the incoming request. If a field is not provided in the request, use the existing value from the database. This allows for partial updates of the todo item, where only the fields that are provided in the request will be updated, while the other fields will remain unchanged.
6.7 Call the UpdateTodo repository function to update the todo item in the database with the new values for the title and completed fields. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the update of the todo item failed. If the todo item is not found during the update operation, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found. If the update is successful
6.8 Return the updated todo item in the response with a 200 OK status.

7. **DeleteTodoHandler** returns a Gin handler function that deletes a specific todo item by its ID for the authenticated user. It extracts the user ID from the Gin context, parses the todo ID from the URL parameter, and calls the DeleteTodo repository function to delete the todo item from the database. If the deletion is successful, it returns a success message in the response with a 200 OK status. If the todo item is not found, it returns a 404 Not Found response with an appropriate error message. If there are any errors during database operations, it logs the error and returns a 500 Internal Server Error response with an appropriate error message indicating that the deletion of the todo item failed.
7.1 Return a Gin handler function that deletes a specific todo item by its ID for the authenticated user.
7.2 Extract the user ID from the Gin context. If the user ID is not found, return a 401 Unauthorized response with an appropriate error message indicating that the user ID is required for authentication.
7.3 Parse the todo ID from the URL parameter. If the ID is not a valid integer, return a 400 Bad Request response with an appropriate error message indicating that the ID is invalid.
7.4 Call the DeleteTodo repository function to delete the todo item from the database for the authenticated user and the specified ID. If there is an error during database operations, log the error and return a 500 Internal Server Error response with an appropriate error message indicating that the deletion of the todo item failed. If the todo item is not found during the deletion operation, return a 404 Not Found response with an appropriate error message indicating that the todo item was not found. If the deletion is successful.
7.5 Return a success message in the response with a 200 OK status.

## 14. Protected routes group /cmd/main.go

Define a protected route group for /todos that uses the AuthMiddleware to ensure that only authenticated users can access the todo-related endpoints. Within this group, define routes for creating, retrieving, updating, and deleting todo items, and associate each route with the corresponding handler function from the handlers package.

## Conclusion

* I learnd how to use go with **gin** framework and how to use the pointer in the struct to check if the value is missing from the request or is empty value.
* It is also interesting how it is fetched the value in the database and update it just if this value is provided in the request
* The use of AIR so I don neeed to constanstly initate the server and all the logs service that come along this tool
* Write databse (postgres) code for go from scratch instead of using a package
