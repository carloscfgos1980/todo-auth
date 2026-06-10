# Todo Auth API

A RESTful API for managing todos with user authentication built with Go, Chi framework, and PostgreSQL.

## Features

- **User Authentication**: Register and login with JWT-based authentication
- **Todo Management**: Create, read, update, and delete todo items
- **Authorization**: Users can only access their own todos
- **Password Security**: Passwords are hashed using Argon2id
- **Database**: PostgreSQL with sqlc for type-safe database access
- **Testing**: Integration tests with mocked services

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- `sqlc` CLI tool for code generation

## Installation

1. Clone the repository:
```bash
git clone https://github.com/carloscfgos1980/todo-auth.git
cd todo-auth-v4
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

## Environment Variables

```env
PORT=8080
DB_DSN=postgres://user:password@localhost:5432/todo_auth_db
JWT_SECRET=your-secret-key-here
```

## Database Setup

1. Create the database:
```sql
CREATE DATABASE todo_auth_db;
```

2. Run migrations:
```bash
# Migrations are in sql/schema/ directory
psql -U postgres -d todo_auth_db -f sql/schema/001_todos.sql
psql -U postgres -d todo_auth_db -f sql/schema/002_users.sql
psql -U postgres -d todo_auth_db -f sql/schema/003_todos_userId.sql
```

3. Generate SQL code:
```bash
sqlc generate
```

## Running the Application

```bash
go run ./cmd/*.go
```

The server will start on `http://localhost:8080` (or the port specified in `PORT` env var).

Health check endpoint: `GET /health`

## API Endpoints

### Authentication

#### Register a new user
```
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}

Response: 201 Created
{
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Login
```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}

Response: 200 OK
{
  "access_token": "jwt_token_here"
}
```

### Todos (Protected Routes)

All todo endpoints require the `Authorization: Bearer <token>` header.

#### Create a todo
```
POST /api/todos
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Buy milk",
  "completed": false
}

Response: 201 Created
{
  "id": 1,
  "title": "Buy milk",
  "completed": false,
  "user_id": "uuid"
}
```

#### Get all todos for authenticated user
```
GET /api/todos
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": 1,
    "title": "Buy milk",
    "completed": false,
    "user_id": "uuid"
  }
]
```

#### Get a specific todo
```
GET /api/todos/{todoID}
Authorization: Bearer <token>

Response: 200 OK
{
  "id": 1,
  "title": "Buy milk",
  "completed": false,
  "user_id": "uuid"
}
```

#### Update a todo
```
PUT /api/todos/{todoID}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Buy milk (updated)",
  "completed": true
}

Response: 200 OK
{
  "id": 1,
  "title": "Buy milk (updated)",
  "completed": true,
  "user_id": "uuid"
}
```

#### Delete a todo
```
DELETE /api/todos/{todoID}
Authorization: Bearer <token>

Response: 200 OK
```

## Project Structure

```
.
├── cmd/
│   ├── api.go          # API setup and routing
│   └── main.go         # Entry point
├── internal/
│   ├── authmiddleware/ # JWT authentication middleware
│   ├── database/       # Database models and queries (sqlc generated)
│   ├── env/            # Environment variable handling
│   ├── json/           # JSON utilities
│   ├── todos/          # Todo handlers, services, and types
│   ├── users/          # User handlers, services, and types
│   └── utils/          # Utility functions
├── sql/
│   ├── queries/        # SQL queries for sqlc
│   └── schema/         # Database schema migrations
├── go.mod
├── go.sum
└── .env
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests for a specific package with verbose output:
```bash
go test ./internal/todos -v
```

Run a specific test:
```bash
go test -run TestCreateTodo ./internal/todos -v
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: User not authorized to access the resource
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Password Requirements

Passwords must meet these criteria:
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit
- At least one special character

## Future Enhancements

- [ ] Add refresh token mechanism
- [ ] Add rate limiting
- [ ] Add request logging and monitoring
- [ ] Add API documentation with Swagger/OpenAPI
- [ ] Add pagination for list endpoints
- [ ] Add todo categories/tags
- [ ] Add due dates for todos

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please create a new branch for your feature and submit a pull request.
