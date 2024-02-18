package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/ent/todo"
	"todo/ent/user"

	"github.com/go-chi/chi/v5"
)

func (handler *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get the userID from the context
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}

	// Parse the request body to get the todo details
	var todoDetails struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&todoDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the new Todo and link it to the User using SetUserID
	newTodo, err := handler.Client.Todo.Create().
		SetTitle(todoDetails.Title).
		SetUserID(userID). // Correctly link the Todo to the User
		Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the newly created Todo as a response
	json.NewEncoder(w).Encode(newTodo)
}

func (handler *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get the userID from the context
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		log.Println("no userID in context!!!!")
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}
	todos, err := handler.Client.Todo.Query().Where(todo.HasUserWith(user.ID(userID))).All(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func (handler *Handler) MarkTodoComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract the Todo ID from the URL path parameters.
	// Assuming you're using chi router and the URL pattern is "/todos/{id}/complete"
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Get the userID from the context to ensure the Todo belongs to the user making the request
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized or invalid user ID", http.StatusUnauthorized)
		return
	}

	// Find the Todo by ID and ensure it belongs to the user
	todoItem, err := handler.Client.Todo.
		Query().
		Where(todo.ID(todoID), todo.HasUserWith(user.ID(userID))).
		Only(ctx)
	if err != nil {
		http.Error(w, "Todo not found or does not belong to user", http.StatusNotFound)
		return
	}

	// Update the Todo's status to "complete"
	_, err = todoItem.Update().SetStatus(todo.StatusComplete).Save(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo marked as complete"})
}
