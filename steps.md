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
