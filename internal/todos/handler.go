package todos

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/carloscfgos1980/todo-auth/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// handler is the HTTP handler for users endpoints
type handler struct {
	service   Service
	jwtSecret string
}

// NewHandler creates a new handler for users endpoints
func NewHandler(service Service, jwtSecret string) *handler {
	return &handler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// CreateTodo handles the HTTP request for creating a new todo item
func (h *handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context (set by the authentication middleware)
	userIDValue := r.Context().Value("userID")
	// Check if the user ID is present in the context
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// The auth middleware stores the JWT subject as a UUID in the request context.
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	// Check if the user exists in the database
	_, err := h.service.GetUserByID(r.Context(), userID.String())
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Parse the JSON request body into a CreateUpdateTodoRequest struct
	var todoReq CreateUpdateTodoRequest
	if err := json.ReadJSON(r, &todoReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate the request data
	if todoReq.Title == nil || *todoReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if todoReq.Completed == nil {
		http.Error(w, "Completed is required", http.StatusBadRequest)
		return
	}
	// Create a CreateTodoParams struct to pass to the service layer
	pgUserID := pgtype.UUID{}
	if err := pgUserID.Scan(userID.String()); err != nil {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	todo := database.CreateTodoParams{
		Title:     *todoReq.Title,
		Completed: *todoReq.Completed,
		UserID:    pgUserID,
	}
	// Call the service to create the todo item
	createdTodo, err := h.service.CreateTodo(r.Context(), todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create a ResponseTodo struct to send back to the client
	response := ResponseTodo{
		ID:        createdTodo.ID,
		UserID:    createdTodo.UserID,
		Title:     createdTodo.Title,
		Completed: createdTodo.Completed,
	}
	// Write the created todo item as JSON with a 201 Created status code
	if err := json.WriteJSON(w, http.StatusCreated, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTodos handles the HTTP request for retrieving all todo items for the authenticated user
func (h *handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context (set by the authentication middleware)
	userIDValue := r.Context().Value("userID")
	// Check if the user ID is present in the context
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// The auth middleware stores the JWT subject as a UUID in the request context.
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	// convert the user ID to pgtype.UUID
	pgUserID := pgtype.UUID{}
	if err := pgUserID.Scan(userID.String()); err != nil {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	// Check if the user exists in the database
	_, err := h.service.GetUserByID(r.Context(), userID.String())
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Call the service to get all todo items for the authenticated user
	todos, err := h.service.GetTodos(r.Context(), pgUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create a slice of ResponseTodo to send back to the client
	response := make([]ResponseTodo, len(todos))
	for i, todo := range todos {
		response[i] = ResponseTodo{
			ID:        todo.ID,
			UserID:    todo.UserID,
			Title:     todo.Title,
			Completed: todo.Completed,
		}
	}
	// Write the todos as JSON with a 200 OK status code
	if err := json.WriteJSON(w, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTodoByID handles the HTTP request for retrieving a specific todo item by its ID for the authenticated user
func (h *handler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context (set by the authentication middleware)
	userIDValue := r.Context().Value("userID")
	// Check if the user ID is present in the context
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// The auth middleware stores the JWT subject as a UUID in the request context.
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	// Check if the user exists in the database
	_, err := h.service.GetUserByID(r.Context(), userID.String())
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Get the todo ID from the URL parameters
	todoIDStr := chi.URLParam(r, "todoID")
	// Convert the todo ID from string to int
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	// Call the service to get the todo item by ID
	todo, err := h.service.GetTodoByID(r.Context(), int32(todoID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Todo item not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// convert the user ID to pgtype.UUID
	authUserID := pgtype.UUID{}
	if err := authUserID.Scan(userID.String()); err != nil {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	//check if the todo item belongs to the authenticated user
	if todo.UserID != authUserID {
		http.Error(w, "You are not authorized to access this todo item", http.StatusForbidden)
		return
	}
	// Create a ResponseTodo struct to send back to the client
	response := ResponseTodo{
		ID:        todo.ID,
		UserID:    todo.UserID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
	// Write the todo item as JSON with a 200 OK status code
	if err := json.WriteJSON(w, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context (set by the authentication middleware)
	userIDValue := r.Context().Value("userID")
	// Check if the user ID is present in the context
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// The auth middleware stores the JWT subject as a UUID in the request context.
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	// Check if the user exists in the database
	_, err := h.service.GetUserByID(r.Context(), userID.String())
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Get the todo ID from the URL parameters
	todoIDStr := chi.URLParam(r, "todoID")
	// Convert the todo ID from string to int
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	// Parse the JSON request body into a CreateUpdateTodoRequest struct
	var todoReq CreateUpdateTodoRequest
	if err := json.ReadJSON(r, &todoReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dbTodo, err := h.service.GetTodoByID(r.Context(), int32(todoID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Todo item not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if the todo item belongs to the authenticated user
	authUserID := pgtype.UUID{}
	if err := authUserID.Scan(userID.String()); err != nil {
		http.Error(w, "Invalid user ID in context", http.StatusInternalServerError)
		return
	}
	if dbTodo.UserID != authUserID {
		http.Error(w, "You are not authorized to update this todo item", http.StatusForbidden)
		return
	}
	// Update the fields of the todo item based on the request data

	title := dbTodo.Title
	if todoReq.Title != nil {
		title = *todoReq.Title
	}
	completed := dbTodo.Completed
	if todoReq.Completed != nil {
		completed = *todoReq.Completed
	}
	// Create an UpdateTodoParams struct to pass to the service layer
	updateParams := database.UpdateTodoParams{
		ID:        int32(todoID),
		Title:     title,
		Completed: completed,
	}
	// Call the service to update the todo item
	updatedTodo, err := h.service.UpdateTodo(r.Context(), updateParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create a ResponseTodo struct to send back to the client
	response := ResponseTodo{
		ID:        updatedTodo.ID,
		UserID:    updatedTodo.UserID,
		Title:     updatedTodo.Title,
		Completed: updatedTodo.Completed,
	}
	// Write the updated todo item as JSON with a 200 OK status code
	if err := json.WriteJSON(w, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
