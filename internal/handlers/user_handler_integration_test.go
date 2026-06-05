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
	"github.com/carloscfgos1980/todo-auth/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRegisterRoute_Success(t *testing.T) {
	// Set Gin to test mode to avoid unnecessary output during testing
	gin.SetMode(gin.TestMode)
	// Create a new sqlmock database connection and a mock object to set expectations on database interactions
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	// Create a configuration struct with the mocked database connection to be used in the handler
	cfg := &config.Config{
		DB: database.New(db),
	}
	// Create a new Gin router and register the CreateUserHandler for the /auth/register route
	router := gin.New()
	router.POST("/auth/register", CreateUserHandler(cfg))
	// Set up the expected database interactions for the CreateUserHandler. When the handler executes an INSERT INTO users query with the specified email and any password argument, it will return a row with the generated user ID, email, hashed password, and timestamps for created_at and updated_at.
	userID := uuid.New()
	now := time.Now().UTC()
	//	Create a new sqlmock.Rows object with the expected columns and add a row with the generated user ID, email, hashed password, and timestamps for created_at and updated_at.
	rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
		AddRow(userID, "john@example.com", "hashed-password", now, now)
	// Set up the expectation for the INSERT INTO users query with the specified email and any password argument, and specify that it will return the row defined above.
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users")).
		WithArgs("john@example.com", sqlmock.AnyArg()).
		WillReturnRows(rows)

	// Create a new HTTP POST request to the /auth/register route with a JSON payload containing the email and password for the new user. Set the Content-Type header to application/json.
	body := []byte(`{"email":"john@example.com","password":"Str0ngP@ssword!"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	// Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
	require.Equal(t, http.StatusOK, res.Code, "response body: %s", res.Body.String())
	// Define a struct to unmarshal the JSON response payload, which contains a user object with fields for ID, email, created_at, and updated_at.
	var payload struct {
		User struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"user"`
	}
	// Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the email in the response matches the expected email, and that the ID, created_at, and updated_at fields are not empty. Finally, assert that all expectations set on the mock database were met.
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	assert.Equal(t, "john@example.com", payload.User.Email)
	assert.NotEmpty(t, payload.User.ID)
	assert.NotEmpty(t, payload.User.CreatedAt)
	assert.NotEmpty(t, payload.User.UpdatedAt)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthLoginRoute_Success(t *testing.T) {
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
	// Create a new Gin router and register the LoginUserHandler for the /auth/login route
	router := gin.New()
	router.POST("/auth/login", LoginUserHandler(cfg))
	// Set up the expected database interactions for the LoginUserHandler. When the handler executes a SELECT query to retrieve the user by email, it will return a row with the user's ID, email, hashed password, and timestamps for created_at and updated_at.
	userID := uuid.New()
	now := time.Now().UTC()
	hashedPassword, err := utils.HashPassword("Str0ngP@ssword!")
	require.NoError(t, err)
	// Create a new sqlmock.Rows object with the expected columns and add a row with the generated user ID, email, hashed password, and timestamps for created_at and updated_at.
	rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
		AddRow(userID, "john@example.com", hashedPassword, now, now)
	// Set up the expectation for the SELECT query to retrieve the user by email, with the specified email argument, and specify that it will return the row defined above.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password, created_at, updated_at FROM users")).
		WithArgs("john@example.com").
		WillReturnRows(rows)
	// Create a new HTTP POST request to the /auth/login route with a JSON payload containing the email and password for the user. Set the Content-Type header to application/json.
	body := []byte(`{"email":"john@example.com","password":"Str0ngP@ssword!"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	// Assert that the response status code is 200 OK and print the response body for debugging purposes if the assertion fails.
	require.Equal(t, http.StatusOK, res.Code, "response body: %s", res.Body.String())
	// Define a struct to unmarshal the JSON response payload, which contains a token field for the JWT token.
	var payload struct {
		Token string `json:"token"`
	}
	// Unmarshal the JSON response body into the defined struct and assert that there are no errors during unmarshaling. Then, assert that the token in the response is not empty. Finally, validate the JWT token using the provided JWT secret and assert that there are no errors during validation, and that the user ID extracted from the token matches the expected user ID. Finally, assert that all expectations set on the mock database were met.
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	assert.NotEmpty(t, payload.Token)

	parsedUserID, err := utils.ValidateJWT(payload.Token, jwtSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)

	require.NoError(t, mock.ExpectationsWereMet())
}
