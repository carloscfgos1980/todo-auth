package todos

import (
	"log"
	"net/http"

	"github.com/carloscfgos1980/todo-auth/internal/database"
	"github.com/carloscfgos1980/todo-auth/internal/json"
	"github.com/google/uuid"
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
	//validate the request data
	if todoReq.Title == nil || *todoReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if todoReq.Completed == nil {
		http.Error(w, "Completed is required", http.StatusBadRequest)
		return
	}

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
