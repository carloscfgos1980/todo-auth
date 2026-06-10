package todos

import (
	"context"
	stdjson "encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createTodoServiceMock struct {
	createTodoFn  func(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error)
	getUserByIDFn func(ctx context.Context, id string) (*database.User, error)
	getTodosFn    func(ctx context.Context, userID pgtype.UUID) ([]database.Todo, error)
	getTodoByIDFn func(ctx context.Context, id int32) (*database.Todo, error)
	updateTodoFn  func(ctx context.Context, todo database.UpdateTodoParams) (*database.Todo, error)
	deleteTodoFn  func(ctx context.Context, id int32) error
}

func (m *createTodoServiceMock) CreateTodo(ctx context.Context, todo database.CreateTodoParams) (*database.Todo, error) {
	if m.createTodoFn == nil {
		return nil, nil
	}
	return m.createTodoFn(ctx, todo)
}

func (m *createTodoServiceMock) GetUserByID(ctx context.Context, id string) (*database.User, error) {
	if m.getUserByIDFn == nil {
		return nil, nil
	}
	return m.getUserByIDFn(ctx, id)
}

func (m *createTodoServiceMock) GetTodos(ctx context.Context, userID pgtype.UUID) ([]database.Todo, error) {
	if m.getTodosFn == nil {
		return nil, nil
	}
	return m.getTodosFn(ctx, userID)
}

func (m *createTodoServiceMock) GetTodoByID(ctx context.Context, id int32) (*database.Todo, error) {
	if m.getTodoByIDFn == nil {
		return nil, nil
	}
	return m.getTodoByIDFn(ctx, id)
}

func (m *createTodoServiceMock) UpdateTodo(ctx context.Context, todo database.UpdateTodoParams) (*database.Todo, error) {
	if m.updateTodoFn == nil {
		return nil, nil
	}
	return m.updateTodoFn(ctx, todo)
}

func (m *createTodoServiceMock) DeleteTodo(ctx context.Context, id int32) error {
	if m.deleteTodoFn == nil {
		return nil
	}
	return m.deleteTodoFn(ctx, id)
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

	handler := NewHandler(service, "test-secret", nil)

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

func TestUpdateTodo(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	var capturedUpdate database.UpdateTodoParams

	userUUID := pgtype.UUID{}
	require.NoError(t, userUUID.Scan(userID.String()))

	service := &createTodoServiceMock{
		getUserByIDFn: func(ctx context.Context, id string) (*database.User, error) {
			return &database.User{ID: userUUID, Email: "user@example.com"}, nil
		},
		getTodoByIDFn: func(ctx context.Context, id int32) (*database.Todo, error) {
			return &database.Todo{
				ID:        id,
				Title:     "Old title",
				Completed: false,
				UserID:    userUUID,
			}, nil
		},
		updateTodoFn: func(ctx context.Context, todo database.UpdateTodoParams) (*database.Todo, error) {
			capturedUpdate = todo
			return &database.Todo{
				ID:        todo.ID,
				Title:     todo.Title,
				Completed: todo.Completed,
				UserID:    userUUID,
			}, nil
		},
	}

	handler := NewHandler(service, "test-secret", nil)

	reqBody := `{"title":"Updated title","completed":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/todos/1", strings.NewReader(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), "userID", userID))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("todoID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	rec := httptest.NewRecorder()

	handler.UpdateTodo(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int32(1), capturedUpdate.ID)
	assert.Equal(t, "Updated title", capturedUpdate.Title)
	assert.True(t, capturedUpdate.Completed)

	var response map[string]any
	require.NoError(t, stdjson.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Updated title", response["title"])
	assert.Equal(t, true, response["completed"])
	assert.Equal(t, userID.String(), response["user_id"])
}
