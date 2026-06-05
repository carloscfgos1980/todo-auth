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
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	jwtSecret := "test-secret"
	cfg := &config.Config{
		DB:        database.New(db),
		JWTSecret: jwtSecret,
	}

	router := gin.New()
	todoRoutes := router.Group("/todos")
	todoRoutes.Use(middleware.AuthMiddleware(cfg))
	todoRoutes.POST("/", CreateTodoHandler(cfg))

	userID := uuid.New()
	token, err := utils.MakeJWT(userID, jwtSecret, time.Hour)
	require.NoError(t, err)

	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at", "completed", "user_id"}).
		AddRow(int32(1), "Buy milk", now, now, false, userID)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO todos (created_at, updated_at, title, completed, user_id)")).
		WithArgs("Buy milk", false, userID).
		WillReturnRows(rows)

	body := []byte(`{"title":"Buy milk","completed":false}`)
	req := httptest.NewRequest(http.MethodPost, "/todos/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code, "response body: %s", res.Body.String())

	var payload struct {
		ID        int32  `json:"id"`
		UserID    string `json:"user_id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	assert.Equal(t, int32(1), payload.ID)
	assert.Equal(t, userID.String(), payload.UserID)
	assert.Equal(t, "Buy milk", payload.Title)
	assert.False(t, payload.Completed)

	require.NoError(t, mock.ExpectationsWereMet())
}