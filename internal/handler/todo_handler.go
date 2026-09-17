package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/raffi/todo-app/internal/model"
	"github.com/raffi/todo-app/internal/service"
)

type TodoService interface {
	CreateTodo(
		ctx context.Context,
		title string,
		description *string,
	) (*model.Todo, error)

	GetTodos(
		ctx context.Context,
	) ([]model.Todo, error)

	GetTodo(
		ctx context.Context,
		id int64,
	) (*model.Todo, error)

	UpdateTodo(
		ctx context.Context,
		id int64,
		title string,
		description *string,
		completed bool,
	) (*model.Todo, error)

	DeleteTodo(
		ctx context.Context,
		id int64,
	) error
}

type TodoHandler struct {
	service TodoService
}

func NewTodoHandler(
	todoService TodoService,
) *TodoHandler {
	return &TodoHandler{
		service: todoService,
	}
}

type CreateTodoRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   bool    `json:"completed"`
}

func (h *TodoHandler) CreateTodo(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	todo, err := h.service.CreateTodo(
		r.Context(),
		req.Title,
		req.Description,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		return
	}
}

func (h *TodoHandler) GetTodos(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	todos, err := h.service.GetTodos(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		return
	}
}

func (h *TodoHandler) GetTodo(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idString := r.URL.Path[len("/todos/"):]

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid todo id",
			http.StatusBadRequest,
		)
		return
	}

	todo, err := h.service.GetTodo(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		return
	}
}

func (h *TodoHandler) UpdateTodo(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idString := r.URL.Path[len("/todos/"):]

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid todo id",
			http.StatusBadRequest,
		)
		return
	}

	var req UpdateTodoRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	todo, err := h.service.UpdateTodo(
		r.Context(),
		id,
		req.Title,
		req.Description,
		req.Completed,
	)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		return
	}
}

func (h *TodoHandler) DeleteTodo(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idString := r.URL.Path[len("/todos/"):]

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid todo id",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.DeleteTodo(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
