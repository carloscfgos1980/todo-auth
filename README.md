# Todo Auth API (Gin + PostgreSQL + sqlc)

A REST API for user authentication and per-user todo management.

## Features

- User registration with email validation and password hashing (Argon2)
- User login with JWT token generation
- Protected todo routes (JWT middleware)
- CRUD operations for todos
- In-memory HTTP metrics endpoint
- sqlc-generated data layer
- Integration tests for auth and todo handlers (sqlmock + testify)

## Tech Stack

- Go
- Gin
- PostgreSQL
- sqlc
- JWT (`github.com/golang-jwt/jwt/v5`)
- Argon2 (`github.com/alexedwards/argon2id`)
- Testing: `sqlmock`, `testify`

## Project Structure

- `cmd/main.go`: app entrypoint and route wiring
- `internal/config`: environment config loader
- `internal/handlers`: HTTP handlers
- `internal/middleware`: auth middleware
- `internal/database`: sqlc-generated queries/models
- `sql/schema`: database schema files
- `sql/queries`: SQL query definitions for sqlc

## Prerequisites

- Go (matching `go.mod`)
- PostgreSQL
- sqlc (for regenerating DB code)

## Environment Variables

Create a `.env` file in project root:

```env
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/todo_auth?sslmode=disable
PORT=8080
JWT_SECRET=your-very-strong-secret
```

Notes:

- `DATABASE_URL` is required (the loader also checks `DatabaseURL`)
- `PORT` is required
- `JWT_SECRET` is required

## Database Setup

Apply schema files in order:

1. `sql/schema/001_todos.sql`
2. `sql/schema/002_users.sql`
3. `sql/schema/003_todos_userId.sql`

Then regenerate sqlc code when SQL changes:

```bash
sqlc generate
```

## Run the API

```bash
go run ./cmd
```

Health check:

```bash
curl http://localhost:8080/
```

Metrics:

```bash
curl http://localhost:8080/metrics
```

## API Endpoints

### Auth

- `POST /auth/register`
- `POST /auth/login`

Register request:

```json
{
  "email": "john@example.com",
  "password": "Str0ngP@ssword!"
}
```

Login request:

```json
{
  "email": "john@example.com",
  "password": "Str0ngP@ssword!"
}
```

Login response:

```json
{
  "token": "<jwt>"
}
```

### Todos (Protected)

Use header:

```http
Authorization: Bearer <jwt>
```

- `POST /todos/`
- `GET /todos/`
- `GET /todos/:id`
- `PUT /todos/:id`
- `DELETE /todos/:id`

Create/Update payload example:

```json
{
  "title": "Buy milk",
  "completed": false
}
```

### Metrics

- `GET /metrics`

Example response:

```json
{
  "total_requests": 12,
  "average_latency_ms": 1.5,
  "requests_by_route": {
    "GET /": 1,
    "GET /metrics": 1,
    "POST /auth/register": 1,
    "GET /todos/": 2
  },
  "requests_by_status": {
    "200": 10,
    "401": 2
  }
}
```

## Tests

Run all handler tests:

```bash
go test ./internal/handlers -v
```

Run focused integration tests:

```bash
go test ./internal/handlers -run TestAuthRegisterRoute_Success -v
go test ./internal/handlers -run TestAuthLoginRoute_Success -v
go test ./internal/handlers -run TestCreateTodoRoute_Success -v
go test ./internal/handlers -run TestGetTodosRoute_Success -v
go test ./internal/handlers -run TestGetTodoByIDRoute_Success -v
go test ./internal/handlers -run TestUpdateTodoRoute_Success -v
go test ./internal/handlers -run TestDeleteTodoRoute_Success -v
```

## Contributing

1. Fork the repository and create a feature branch.
2. Keep changes focused and small per PR.
3. If SQL changes, update files in `sql/queries` or `sql/schema` and run:

```bash
sqlc generate
```

4. Run tests before opening a PR:

```bash
go test ./...
```

5. Use clear commit messages (for example: `feat: add metrics endpoint`).
6. Open a pull request with:

- A summary of what changed
- Test evidence (`go test` output)
- Notes about any env or migration changes

## Common Issues

- `missing database URL`: set `DATABASE_URL` in `.env`
- `missing port`: set `PORT`
- `missing JWT secret`: set `JWT_SECRET`
- `invalid email or password`: check login credentials and stored hash
