# Todo Auth API (Go + PostgreSQL)

Simple REST API for user authentication and todo management.

## Features

- User registration with password hashing (Argon2)
- User login with JWT token generation
- CRUD operations for todos (scoped to authenticated user)
- Health endpoint
- Lightweight metrics endpoint (uptime and total requests)

## Tech Stack

- Go
- PostgreSQL
- sqlc-generated data access layer
- JWT (github.com/golang-jwt/jwt/v5)

## Requirements

- Go 1.25+
- PostgreSQL

## Environment Variables

Create a .env file in the project root with:

~~~env
DatabaseURL=postgres://postgres:postgres@localhost:5432/todo_auth?sslmode=disable
PORT=8080
JWT_SECRET=your-super-secret-key
~~~

## Database Setup

1. Create a PostgreSQL database.
2. Enable pgcrypto extension (needed for gen_random_uuid()).
3. Run schema files in order:

~~~sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
~~~

~~~bash
psql "$DatabaseURL" -f sql/schema/001_todos.sql
psql "$DatabaseURL" -f sql/schema/002_users.sql
psql "$DatabaseURL" -f sql/schema/003_todos_userId.sql
~~~

## Run

~~~bash
go mod tidy
go run .
~~~

Server starts on the value of PORT.

## API Endpoints

### Public

- GET /health
- GET /metrics
- POST /auth/register
- POST /auth/login

### Authenticated (Bearer token required)

- POST /todos
- GET /todos
- GET /todos/{todoID}
- PUT /todos/{todoID}
- DELETE /todos/{todoID}

## Quick Usage

### 1) Register

~~~bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"StrongP@ssw0rd"}'
~~~

### 2) Login

~~~bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"StrongP@ssw0rd"}'
~~~

Save the returned token.

### 3) Create Todo

~~~bash
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","completed":false}'
~~~

### 4) List Todos

~~~bash
curl http://localhost:8080/todos \
  -H "Authorization: Bearer YOUR_TOKEN"
~~~

### 5) Metrics

~~~bash
curl http://localhost:8080/metrics
~~~

Example response:

~~~json
{
  "uptime_seconds": 123,
  "requests_total": 45
}
~~~

## Project Structure

- main server bootstrap and route registration in main.go
- HTTP handlers in handler_*.go
- auth utilities in internal/auth
- sqlc output in internal/database
- SQL schema and queries in sql

## Notes

- This project uses Go method-aware routes for todos and metrics.
- requests_total includes all incoming requests, including /metrics.
