package todos

import (
	"context"
	stdjson "encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createTodoServiceMock struct {
	createTodoFn func(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error)
}

func (m *createTodoServiceMock) CreateTodo(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error) {
	return m.createTodoFn(ctx, todo)
}

func (m *createTodoServiceMock) GetUserByID(ctx context.Context, id string) (*database.User, error) {
	return nil, nil
}

func (m *createTodoServiceMock) GetTodos(ctx context.Context, userID pgtype.UUID) ([]database.Todo, error) {
	return nil, nil
}

func (m *createTodoServiceMock) GetTodoByID(ctx context.Context, id int32) (*database.Todo, error) {
	return nil, nil
}

func (m *createTodoServiceMock) UpdateTodo(ctx context.Context, todo database.UpdateTodoParams) (*database.Todo, error) {
	return nil, nil
}

func (m *createTodoServiceMock) DeleteTodo(ctx context.Context, id int32) error {
	return nil
}

func TestCreateTodo(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	var capturedTodo database.CreateTodoParams
	service := &createTodoServiceMock{
		createTodoFn: func(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error) {
			capturedTodo = todo

			todoID := pgtype.UUID{}
			require.NoError(t, todoID.Scan(userID.String()))

			return &database.Todo{
				ID:        1,
				Title:     todo.Title,
				Completed: todo.Completed,
				UserID:    todoID,
			}, nil
		},
	}

	handler := NewHandler(service, "test-secret")

	reqBody := `{"title":"Buy milk","completed":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "userID", userID))
	rec := httptest.NewRecorder()

	handler.CreateTodo(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Buy milk", capturedTodo.Title)
	assert.True(t, capturedTodo.Completed)

	var response map[string]any
	require.NoError(t, stdjson.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Buy milk", response["title"])
	assert.Equal(t, true, response["completed"])
	assert.Equal(t, userID.String(), response["user_id"])
}
