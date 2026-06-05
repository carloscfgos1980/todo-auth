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
go get github.com/DATA-DOG/go-sqlmock
go get github.com/stretchr/testify/assert@v1.11.1

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

## 4. Create Todo - auth

1. Create querie to insert values into todos table
2. run "sqlc generate" to convert SQL to GO
3. Copy middleware folder form previos version
4. responseTodo represents the structure of a todo item that will be sent back in the response. It includes the ID, user ID, title, and completion status of the todo item.
5. CreateTodoHandler handles the creation of a new todo item. It retrieves the user ID from the context, binds the incoming JSON payload to a CreateTodoRequest struct, and calls the CreateTodo method from the database layer to insert the new todo item into the database. If successful, it returns the created todo item in the response; otherwise, it returns an appropriate error message.
5.1 returns a Gin handler function that processes the creation of a new todo item.
5.2 Bind the incoming JSON payload to a CreateTodoRequest struct. If there is an error during binding (e.g., missing required fields), return a bad request error response.
5.3 Call the CreateTodo method from the database layer, passing the user ID, title, and completion status from the request. If there is an error during the database operation, return an internal server error response. If successful, construct a responseTodo struct with the created todo item and return it in the response.
5.4 Construct a responseTodo struct with the created todo item and return it in the response.
5.5 Return the created todo item in the response with a status of 200 OK.
6.Register todo-related routes
 todoRoutes := router.Group("/todos")
 todoRoutes.Use(middleware.AuthMiddleware(cfg))
 todoRoutes.POST("/", handlers.CreateTodoHandler(cfg))

## 5. Get todos - auth

1. Create querie to get todos by user is
2. run "sqlc generate" to convert SQL to GO
3. GetTodosHandler handles the retrieval of all todo items for a specific user. It retrieves the user ID from the context, calls the GetTodosByUserID method from the database layer to fetch the todo items associated with that user, and returns the list of todo items in the response. If there is an error during the process, it returns an appropriate error message.
3.1 returns a Gin handler function that processes the retrieval of all todo items for a specific user.
3.2 Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
3.3 Call the GetTodosByUserID method from the database layer, passing the user ID to fetch the todo items associated with that user. If there is an error during the database operation, return an internal server error response. If successful, construct a slice of responseTodo structs with the retrieved todo items and return it in the response.
3.4 Construct a slice of responseTodo structs with the retrieved todo items and return it in the response.
3.5 Return the list of todo items in the response with a status of 200 OK.
4. Register todo-related routes
 todoRoutes.GET("/", handlers.GetTodosHandler(cfg))

## 6. Get a todo-auth

1. Create querie to fetch todo by id
2. Run sqlc generate
3. GetTodoByIDHandler handles the retrieval of a specific todo item by its ID. It retrieves the user ID from the context, extracts the todo ID from the URL parameters, and calls the GetTodoByID method from the database layer to fetch the todo item. If the todo item is found and belongs to the authenticated user, it returns the todo item in the response; otherwise, it returns an appropriate error message (e.g., not found, unauthorized, or forbidden).
3.1 returns a Gin handler function that processes the retrieval of a specific todo item by its ID.
3.2 Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
3.3 Extract the todo ID from the URL parameters and convert it to an integer. If there is an error during conversion (e.g., invalid ID format), return a bad request error response.
3.4 Call the GetTodoByID method from the database layer, passing the todo ID to fetch the specific todo item. If there is an error during the database operation, return an internal server error response. If the todo item is not found, return a not found error response. If the todo item is found but does not belong to the authenticated user, return a forbidden error response. If successful, construct a responseTodo struct with the retrieved todo item and return it in the response.
3.5 Check if the retrieved todo item belongs to the authenticated user. If not, return a forbidden error response.
3.6 Construct a responseTodo struct with the retrieved todo item and return it in the response.
3.7 Return the retrieved todo item in the response with a status of 200 OK.
4. Register todo-related routes
 todoRoutes.GET("/:id", handlers.GetTodoByIDHandler(cfg))

## 7. Update todo

1. Create querie to update todo by id
2. Run sqlc generate
3. UpdateTodoRequest represents the expected JSON payload for updating an existing todo item. It includes optional fields for the title and completion status, allowing for partial updates of the todo item.
4. UpdateTodoHandler handles the update of an existing todo item. It retrieves the user ID from the context, extracts the todo ID from the URL parameters, binds the incoming JSON payload to an UpdateTodoRequest struct, and calls the UpdateTodo method from the database layer to update the todo item in the database. If the todo item is found and belongs to the authenticated user, it returns the updated todo item in the response; otherwise, it returns an appropriate error message (e.g., not found, unauthorized, or forbidden).
4.1 returns a Gin handler function that processes the update of an existing todo item.
4.2 Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
4.3 Extract the todo ID from the URL parameters and convert it to an integer. If there is an error during conversion (e.g., invalid ID format), return a bad request error response.
4.4 Bind the incoming JSON payload to an UpdateTodoRequest struct. If there is an error during binding (e.g., invalid JSON format), return a bad request error response.
4.5 Call the GetTodoByID method from the database layer to fetch the existing todo item. If there is an error during the database operation, return an internal server error response. If the todo item is not found, return a not found error response. If the todo item is found but does not belong to the authenticated user, return a forbidden error response.
4.6 Check if the retrieved todo item belongs to the authenticated user. If not, return a forbidden error response.
4.7 Determine the updated title and completion status for the todo item. If the corresponding fields in the UpdateTodoRequest struct are nil, use the existing values from the database; otherwise, use the new values from the request.
4.8 Call the UpdateTodo method from the database layer, passing the todo ID, updated title, and completion status. If there is an error during the database operation, return an internal server error response. If successful, construct a responseTodo struct with the updated todo item and return it in the response.
4.9 Construct a responseTodo struct with the updated todo item and return it in the response.
4.10 Return the updated todo item in the response with a status of 200 OK.
5. Register todo-related routes
 todoRoutes.PUT("/:id", handlers.UpdateTodoHandler(cfg))

## 8. Delete todo

1. Create query to delete todo
2. Run sqlc generate to convert SL to GO
3. DeleteTodoHandler handles the deletion of a specific todo item by its ID. It retrieves the user ID from the context, extracts the todo ID from the URL parameters, and calls the DeleteTodo method from the database layer to delete the todo item. If the todo item is found and belongs to the authenticated user, it deletes the item and returns a success message in the response; otherwise, it returns an appropriate error message (e.g., not found, unauthorized, or forbidden).
3.1 returns a Gin handler function that processes the deletion of a specific todo item by its ID.
3.2 Retrieve the user ID from the context, which is set by the authentication middleware. If the user ID is not found, return an unauthorized error response.
3.3 Extract the todo ID from the URL parameters and convert it to an integer. If there is an error during conversion (e.g., invalid ID format), return a bad request error response.
3.4 Call the GetTodoByID method from the database layer to fetch the existing todo item. If there is an error during the database operation, return an internal server error response. If the todo item is not found, return a not found error response. If the todo item is found but does not belong to the authenticated user, return a forbidden error response.
3.5 Call the DeleteTodo method from the database layer, passing the todo ID. If there is an error during the database operation, return an internal server error response. If successful, return a success message in the response.
3.6 Return a success message in the response with a status of 200 OK.
3.7 Register todo-related routes
 todoRoutes.DELETE("/:id", handlers.DeleteTodoHandler(cfg))

## 9. Integration test to create user

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection to be used in the handler
4. Create a new Gin router and register the CreateUserHandler for the /auth/register route
5. Set up the expected database interactions for the CreateUserHandler. When the handler executes an INSERT INTO users query with the specified email and any password argument, it will return a row with the generated user ID, email, hashed password, and timestamps for created_at and updated_at.
6. Create a new sqlmock.Rows object with the expected columns and add a row with the generated user ID, email, hashed password, and timestamps for created_at and updated_at.
7. Set up the expectation for the INSERT INTO users query with the specified email and any password argument, and specify that it will return the row defined above.
8. Create a new HTTP POST request to the /auth/register route with a JSON payload containing the email and password for the new user. Set the Content-Type header to application/json.
9. Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
10. Define a struct to unmarshal the JSON response payload, which contains a user object with fields for ID, email, created_at, and updated_at.
11. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the email in the response matches the expected email, and that the ID, created_at, and updated_at fields are not empty. Finally, assert that all expectations set on the mock database were met.

```bash
go test ./internal/handlers -run TestAuthRegisterRoute_Success -v
```

## 10. integration test login

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
4. Create a new Gin router and register the LoginUserHandler for the /auth/login route
5. Set up the expected database interactions for the LoginUserHandler. When the handler executes a SELECT query to retrieve the user by email, it will return a row with the user's ID, email, hashed password, and timestamps for created_at and updated_at.
6. Create a new sqlmock.Rows object with the expected columns and add a row with the generated user ID, email, hashed password, and timestamps for created_at and updated_at.
7. Set up the expectation for the SELECT query to retrieve the user by email, with the specified email argument, and specify that it will return the row defined above.
8. Create a new HTTP POST request to the /auth/login route with a JSON payload containing the email and password for the user. Set the Content-Type header to application/json.
9. Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
10. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the token in the response is not empty. Finally, validate the JWT token using the provided JWT secret and assert that there are no errors during validation, and that the user ID extracted from the token matches the expected user ID. Finally, assert that all expectations set on the mock database were met.

```bash
go test ./internal/handlers -run TestAuthLoginRoute_Success -v
```

## 11. Integration test create todo

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
4. Create a new Gin router and register the CreateTodoHandler for the /todos/ route, applying the AuthMiddleware to protect the route and require authentication.
5. Set up the expected database interactions for the CreateTodoHandler. When the handler executes an INSERT INTO todos query with the specified title, completed status, and user ID arguments, it will return a row with the generated todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
6. Create a new HTTP POST request to the /todos/ route with a JSON payload containing the title and completed status for the new todo. Set the Content-Type header to application/json and include the Authorization header with the Bearer token for authentication.
7. Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
8. Define a struct to unmarshal the JSON response payload, which contains fields for ID, user_id, title, and completed status of the created todo.
9. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the created todo, including the ID, user ID, title, and completed status. Finally, assert that all expectations set on the mock database were met.

```bash
go test ./internal/handlers -run TestCreateTodoRoute_Success -v
```

## 12. Integration test get TODOS

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
4. Create a new Gin router and register the GetTodosHandler for the /todos/ route, applying AuthMiddleware
5. Set up the expected database interactions for the GetTodosHandler. When the handler executes a SELECT query to retrieve todos by user ID, it will return multiple rows with the todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
6. Create a new HTTP GET request to the /todos/ route. Set the Authorization header with the Bearer token for authentication.
7. Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
8. Define a struct to unmarshal the JSON response payload, which contains a slice of todo objects with fields for ID, user_id, title, and completed status.
9. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the retrieved todos, including the ID, user ID, title, and completed status for each todo. Finally, assert that all expectations set on the mock database were met.
10. Assert that the length of the payload is 2, indicating that two todos were returned in the response.
11. Assert the fields of the first todo in the response match the expected values for the first todo in the database, including the ID, user ID, title, and completed status.

```bash
go test ./internal/handlers -run TestGetTodosRoute_Success -v
```

## 13. Integration test get todo by id

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
4. Create a new Gin router and register the GetTodoByIDHandler for the /todos/:id route, applying AuthMiddleware
5. Set up the expected database interactions for the GetTodoByIDHandler. When the handler executes a SELECT query to retrieve a todo by ID, it will return a row with the todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
6. Create a new HTTP GET request to the /todos/1 route. Set the Authorization header with the Bearer token for authentication.
7. Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
8. Define a struct to unmarshal the JSON response payload, which contains fields for ID, user_id, title, and completed status of the retrieved todo.
9. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the retrieved todo, including the ID, user ID, title, and completed status. Finally, assert that all expectations set on the mock database were met.

```bash
go test ./internal/handlers -run TestGetTodoByIDRoute_Success -v
```

## 14. integration test update todo

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
4. Create a new Gin router and register the UpdateTodoHandler for the /todos/:id route, applying AuthMiddleware
5. Set up the expected database interactions for the UpdateTodoHandler. When the handler executes a SELECT query to retrieve the current todo by ID, it will return a row with the current todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. Then, when the handler executes an UPDATE query to update the todo with the specified title and completed status, it will return a row with the updated todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
6. Create a new HTTP PUT request to the /todos/1 route with a JSON payload containing the updated title and completed status for the todo. Set the Content-Type header to application/json and include the Authorization header with the Bearer token for authentication.
7. Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
8. Define a struct to unmarshal the JSON response payload, which contains fields for ID, user_id, title, and completed status of the updated todo.
9. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the updated todo, including the ID, user ID, title, and completed status. Finally, assert that all expectations set on the mock database were met.

```bash
go test ./internal/handlers -run TestUpdateTodoRoute_Success -v
```

## 15. integration test delete todo

1. Set Gin to test mode to avoid unnecessary output during testing
2. Create a new sqlmock database connection and a mock object to set expectations on database interactions
3. Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
4. Create a new Gin router and register the DeleteTodoHandler for the /todos/:id route, applying AuthMiddleware
5. Set up the expected database interactions for the DeleteTodoHandler. When the handler executes a SELECT query to retrieve the current todo by ID, it will return a row with the current todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. Then, when the handler executes a DELETE query to delete the todo by ID, it will return a result indicating that one row was affected. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
6. Create a new HTTP DELETE request to the /todos/1 route. Set the Authorization header with the Bearer token for authentication.
7. Define a struct to unmarshal the JSON response payload, which contains a message field indicating the result of the delete operation.
8. Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the message in the response indicates that the todo was deleted successfully. Finally, assert that all expectations set on the mock database were met.

```bash
go test ./internal/handlers -run TestDeleteTodoRoute_Success -v
```

## 16. Add README

* Features
* Tech Stack
* Project Structure
* Prerequisites
* Environment Variables
* Database Setup
* Run the API
* API Endpoints
* Tests
* Contributing
* Common Issues

## 17. Add metrics

metrics_middleware

1. MetricsCollector is a struct that collects metrics for HTTP requests, including total requests, total latency in milliseconds, requests by route, and requests by status code. It uses a mutex to ensure thread-safe access to the metrics data.
2. NewMetricsCollector creates and returns a new instance of MetricsCollector with initialized fields.
3. Middleware is a Gin middleware function that collects metrics for each incoming HTTP request. It records the start time of the request, processes the request, and then calculates the latency and updates the metrics for total requests, total latency, requests by route, and requests by status code. It uses the request method and route as the key for tracking requests by route, and the response status code for tracking requests by status code. If the route is not matched, it uses "unmatched" as the key for tracking requests by route.
3.1 Return a Gin handler function that will be called for each incoming HTTP request. This function will record the start time of the request, process the request, and then calculate the latency and update the metrics accordingly.
3.2 Record the start time of the request to calculate latency later.
3.3 Process the request by calling the next handler in the chain.
3.4 Get the full path of the route that was matched for the request. If no route was matched, use "unmatched" as the key for tracking requests by route.
3.5 Create a key for tracking requests by route using the request method and route. Get the response status code and calculate the latency in milliseconds. Then, update the total requests, total latency, requests by route, and requests by status code metrics using atomic operations and a helper function to increment the counters in the sync.Map.
4. Handler is a Gin handler function that returns the collected metrics as a JSON response. It calculates the total requests, average latency, requests by route, and requests by status code, and sends this information in the response.
4.1 Return a Gin handler function that will be called when the /metrics endpoint is accessed. This function will gather the collected metrics and return them as a JSON response.
4.2 Load the total requests and total latency metrics using atomic operations. Then, iterate over the byRoute and byStatus sync.Maps to gather the counts for each route and status code. Finally, calculate the average latency and return all the metrics in a JSON response.
4.3 Create maps to hold the counts for requests by route and requests by status code. Use the Range method of sync.Map to iterate over the entries in byRoute and byStatus, and load the counts using atomic operations.
4.4 Create a map to hold the counts for requests by status code. Use the Range method of sync.Map to iterate over the entries in byStatus, and load the counts using atomic operations.
4.5 Calculate the average latency in milliseconds by dividing the total latency by the total number of requests. If there are no requests, set the average latency to 0 to avoid division by zero.
4.6 Return the collected metrics in a JSON response with fields for total requests, average latency, requests by route, and requests by status code.
5. incrementCounter is a helper function that increments the counter for a given key in a sync.Map. It uses the LoadOrStore method to get or create an atomic.Uint64 counter for the key, and then increments the counter using the Add method.
main.go
6. Create a new MetricsCollector instance and register the metrics middleware to collect metrics for each incoming HTTP request. The metrics will include total requests, total latency, requests by route, and requests by status code. The metrics will be collected in a thread-safe manner using atomic operations and a sync.Map to store the counts for each route and status code.
7. Register the /metrics endpoint to return the collected metrics as a JSON response when accessed. This endpoint will allow monitoring tools to scrape the metrics for analysis and visualization.
