# STEPS

## 1. Setup

1. Start the api
go mod init github.com/carloscfgos1980/todo-auth

2. Packages
github.com/joho/godotenv
github.com/lib/pq
github.com/google/uuid
github.com/alexedwards/argon2id
github.com/golang-jwt/jwt/v5
github.com/stretchr/testify/assert

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
6. Load to github

```bash
git init
git add .
git commit -m "set up"
git remote add origin https://github.com/carloscfgos1980/todo-auth.git
git checkout -b no_framework
git push origin no_framework
```

## 2. User register

1. Copy auth folder from previous version
2. structs and handler for creating a new user in the system
3. handlerUsersCreate handles the user registration endpoint, allowing clients to create new user accounts by providing an email and password. It validates the input, hashes the password, and stores the new user in the database.
3.1 Define the expected parameters for creating a new user and the response structure
3.2 Decode the JSON request body into the parameters struct
3.3 Validate the provided parameters (e.g., check for valid email format, strong password, etc.)
3.4 strong password validation can be added here before hashing the password and creating the user in the database
3.5 Hash the user's password before storing it in the database
3.6 Create a new user in the database using the provided parameters and the hashed password
3.7 Prepare the response with the created user's information (excluding the password)
3.8 Respond with the created user's information (excluding the password)

4. user creation endpoint
 mux.HandleFunc("/auth/register", apiCfg.handlerUsersCreate)

## 3. Login

1. handlerLogin handles user login requests. It expects a JSON body with the user's email and password, verifies the credentials, and returns a JWT token if the login is successful.
1.1 Define the expected parameters for user login
1.2 Define the response structure
1.3 Decode the JSON request body into the parameters struct
1.4 Retrieve the user from the database using the provided email address
1.5 Check if the provided password matches the hashed password stored in the database for the retrieved user
1.6 If the password is correct, generate a JWT token for the user to authenticate future requests
1.7 Respond with the generated JWT token

2. user login endpoint
 mux.HandleFunc("/auth/login", apiCfg.handlerLogin)

## Create todo

1. responseTodo is a struct that represents the JSON response for a created task. It includes the task's ID, user ID, title, and completion status.
2. createTodoHandler handles requests to create a new task. It validates the user's authorization using a JWT token, checks the provided parameters for creating a new task, and creates the task in the database if everything is valid.
2.1 Define the expected parameters for creating a new task
2.2 Validate the user's authorization to create a new task by checking the provided JWT token
2.3 Validate the JWT token and extract the user ID from it
2.4 Decode the JSON request body into the parameters struct
2.5 Validate the provided parameters for creating a new task (e.g., check if priority, state, tag, and date formats are valid)
2.6 heck if the user is register in database
2.7 If the parameters are valid and the user is authorized, create a new task in the database associated with the user's ID and the provided parameters
2.8 Respond with the created task's information
2.9 Respond with the created task's information in JSON format

3. todo creation endpoint
 mux.HandleFunc("POST /todos", apiCfg.createTodoHandler)
