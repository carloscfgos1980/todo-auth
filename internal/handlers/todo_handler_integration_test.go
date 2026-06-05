package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/carloscfgos1980/todo-auth/internal/config"
	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/carloscfgos1980/todo-auth/internal/middleware"
	"github.com/carloscfgos1980/todo-auth/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTodoRoute_Success(t *testing.T) {
	// Set Gin to test mode to avoid unnecessary output during testing
	gin.SetMode(gin.TestMode)
	// Create a new sqlmock database connection and a mock object to set expectations on database interactions
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	// Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
	jwtSecret := "test-secret"
	cfg := &config.Config{
		DB:        database.New(db),
		JWTSecret: jwtSecret,
	}
	// Create a new Gin router and register the CreateTodoHandler for the /todos/ route, applying the AuthMiddleware to protect the route and require authentication.
	router := gin.New()
	todoRoutes := router.Group("/todos")
	todoRoutes.Use(middleware.AuthMiddleware(cfg))
	todoRoutes.POST("/", CreateTodoHandler(cfg))
	// Set up the expected database interactions for the CreateTodoHandler. When the handler executes an INSERT INTO todos query with the specified title, completed status, and user ID arguments, it will return a row with the generated todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
	userID := uuid.New()
	token, err := utils.MakeJWT(userID, jwtSecret, time.Hour)
	require.NoError(t, err)

	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "completed", "user_id"}).
		AddRow(int32(1), "Buy milk", now, now, false, userID)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO todos (created_at, updated_at, title, completed, user_id)")).
		WithArgs("Buy milk", false, userID).
		WillReturnRows(rows)
	// Create a new HTTP POST request to the /todos/ route with a JSON payload containing the title and completed status for the new todo. Set the Content-Type header to application/json and include the Authorization header with the Bearer token for authentication.
	body := []byte(`{"title":"Buy milk","completed":false}`)
	req := httptest.NewRequest(http.MethodPost, "/todos/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	// Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
	require.Equal(t, http.StatusOK, res.Code, "response body: %s", res.Body.String())
	// Define a struct to unmarshal the JSON response payload, which contains fields for ID, user_id, title, and completed status of the created todo.
	var payload struct {
		ID        int32  `json:"id"`
		UserID    string `json:"user_id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	// Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the created todo, including the ID, user ID, title, and completed status. Finally, assert that all expectations set on the mock database were met.
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	assert.Equal(t, int32(1), payload.ID)
	assert.Equal(t, userID.String(), payload.UserID)
	assert.Equal(t, "Buy milk", payload.Title)
	assert.False(t, payload.Completed)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTodosRoute_Success(t *testing.T) {
	// Set Gin to test mode to avoid unnecessary output during testing
	gin.SetMode(gin.TestMode)
	// Create a new sqlmock database connection and a mock object to set expectations on database interactions
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	// Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
	jwtSecret := "test-secret"
	cfg := &config.Config{
		DB:        database.New(db),
		JWTSecret: jwtSecret,
	}
	// Create a new Gin router and register the GetTodosHandler for the /todos/ route, applying AuthMiddleware
	router := gin.New()
	todoRoutes := router.Group("/todos")
	todoRoutes.Use(middleware.AuthMiddleware(cfg))
	todoRoutes.GET("/", GetTodosHandler(cfg))
	// Set up the expected database interactions for the GetTodosHandler. When the handler executes a SELECT query to retrieve todos by user ID, it will return multiple rows with the todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
	userID := uuid.New()
	token, err := utils.MakeJWT(userID, jwtSecret, time.Hour)
	require.NoError(t, err)

	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "completed", "user_id"}).
		AddRow(int32(2), "Buy milk", now, now, false, userID).
		AddRow(int32(1), "Read book", now, now, true, userID)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, created_at, updated_at, completed, user_id FROM todos")).
		WithArgs(userID).
		WillReturnRows(rows)
	// Create a new HTTP GET request to the /todos/ route. Set the Authorization header with the Bearer token for authentication.
	req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	// Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
	require.Equal(t, http.StatusOK, res.Code, "response body: %s", res.Body.String())
	// Define a struct to unmarshal the JSON response payload, which contains a slice of todo objects with fields for ID, user_id, title, and completed status.
	var payload []struct {
		ID        int32  `json:"id"`
		UserID    string `json:"user_id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	// Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the retrieved todos, including the ID, user ID, title, and completed status for each todo. Finally, assert that all expectations set on the mock database were met.
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	// Assert that the length of the payload is 2, indicating that two todos were returned in the response.
	require.Len(t, payload, 2)
	// Assert the fields of the first todo in the response match the expected values for the first todo in the database, including the ID, user ID, title, and completed status.
	assert.Equal(t, int32(2), payload[0].ID)
	assert.Equal(t, userID.String(), payload[0].UserID)
	assert.Equal(t, "Buy milk", payload[0].Title)
	assert.False(t, payload[0].Completed)

	assert.Equal(t, int32(1), payload[1].ID)
	assert.Equal(t, userID.String(), payload[1].UserID)
	assert.Equal(t, "Read book", payload[1].Title)
	assert.True(t, payload[1].Completed)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTodoByIDRoute_Success(t *testing.T) {
	// Set Gin to test mode to avoid unnecessary output during testing
	gin.SetMode(gin.TestMode)
	// Create a new sqlmock database connection and a mock object to set expectations on database interactions
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	// Create a configuration struct with the mocked database connection and a JWT secret to be used in the handler
	jwtSecret := "test-secret"
	cfg := &config.Config{
		DB:        database.New(db),
		JWTSecret: jwtSecret,
	}
	// Create a new Gin router and register the GetTodoByIDHandler for the /todos/:id route, applying AuthMiddleware
	router := gin.New()
	todoRoutes := router.Group("/todos")
	todoRoutes.Use(middleware.AuthMiddleware(cfg))
	todoRoutes.GET("/:id", GetTodoByIDHandler(cfg))
	// Set up the expected database interactions for the GetTodoByIDHandler. When the handler executes a SELECT query to retrieve a todo by ID, it will return a row with the todo ID, title, timestamps for created_at and updated_at, completed status, and user ID. To simulate an authenticated request, we generate a JWT token for a test user ID using the provided JWT secret and a short expiration time.
	userID := uuid.New()
	token, err := utils.MakeJWT(userID, jwtSecret, time.Hour)
	require.NoError(t, err)

	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "completed", "user_id"}).
		AddRow(int32(1), "Buy milk", now, now, false, userID)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, created_at, updated_at, completed, user_id FROM todos WHERE id = $1")).
		WithArgs(int32(1)).
		WillReturnRows(rows)
	// Create a new HTTP GET request to the /todos/1 route. Set the Authorization header with the Bearer token for authentication.
	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	// Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
	require.Equal(t, http.StatusOK, res.Code, "response body: %s", res.Body.String())
	// Define a struct to unmarshal the JSON response payload, which contains fields for ID, user_id, title, and completed status of the retrieved todo.
	var payload struct {
		ID        int32  `json:"id"`
		UserID    string `json:"user_id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	// Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the fields in the response match the expected values for the retrieved todo, including the ID, user ID, title, and completed status. Finally, assert that all expectations set on the mock database were met.
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	assert.Equal(t, int32(1), payload.ID)
	assert.Equal(t, userID.String(), payload.UserID)
	assert.Equal(t, "Buy milk", payload.Title)
	assert.False(t, payload.Completed)

	require.NoError(t, mock.ExpectationsWereMet())
}
