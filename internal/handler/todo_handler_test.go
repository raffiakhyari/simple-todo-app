package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raffi/todo-app/internal/model"
	"github.com/raffi/todo-app/internal/service"
)

type mockTodoService struct {
	createTodoFunc func(
		ctx context.Context,
		title string,
		description *string,
	) (*model.Todo, error)

	getTodosFunc func(
		ctx context.Context,
	) ([]model.Todo, error)

	getTodoFunc func(
		ctx context.Context,
		id int64,
	) (*model.Todo, error)

	updateTodoFunc func(
		ctx context.Context,
		id int64,
		title string,
		description *string,
		completed bool,
	) (*model.Todo, error)

	deleteTodoFunc func(
		ctx context.Context,
		id int64,
	) error
}

func (m *mockTodoService) CreateTodo(
	ctx context.Context,
	title string,
	description *string,
) (*model.Todo, error) {
	return m.createTodoFunc(
		ctx,
		title,
		description,
	)
}

func (m *mockTodoService) GetTodos(
	ctx context.Context,
) ([]model.Todo, error) {
	if m.getTodosFunc == nil {
		return nil, nil
	}

	return m.getTodosFunc(ctx)
}

func (m *mockTodoService) GetTodo(
	ctx context.Context,
	id int64,
) (*model.Todo, error) {
	if m.getTodoFunc == nil {
		return nil, nil
	}

	return m.getTodoFunc(ctx, id)
}

func (m *mockTodoService) UpdateTodo(
	ctx context.Context,
	id int64,
	title string,
	description *string,
	completed bool,
) (*model.Todo, error) {
	if m.updateTodoFunc == nil {
		return nil, nil
	}

	return m.updateTodoFunc(
		ctx,
		id,
		title,
		description,
		completed,
	)
}

func (m *mockTodoService) DeleteTodo(
	ctx context.Context,
	id int64,
) error {
	if m.deleteTodoFunc == nil {
		return nil
	}

	return m.deleteTodoFunc(ctx, id)
}

func TestCreateTodoHandler(t *testing.T) {
	description := "Learn Go"

	expectedTodo := &model.Todo{
		ID:          1,
		Title:       "Learn Go",
		Description: &description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	serviceCalled := false

	mockService := &mockTodoService{
		createTodoFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			serviceCalled = true

			if title != "Learn Go" {
				t.Fatalf(
					"expected title %q, got %q",
					"Learn Go",
					title,
				)
			}

			if description == nil {
				t.Fatal("expected description, got nil")
			}

			if *description != "Learn Go" {
				t.Fatalf(
					"expected description %q, got %q",
					"Learn Go",
					*description,
				)
			}

			return expectedTodo, nil
		},
	}

	handler := NewTodoHandler(mockService)

	requestBody := `{
		"title": "Learn Go",
		"description": "Learn Go"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		bytes.NewBufferString(requestBody),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	var response model.Todo

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response.ID != expectedTodo.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			expectedTodo.ID,
			response.ID,
		)
	}

	if response.Title != expectedTodo.Title {
		t.Fatalf(
			"expected title %q, got %q",
			expectedTodo.Title,
			response.Title,
		)
	}

	if response.Completed != expectedTodo.Completed {
		t.Fatalf(
			"expected completed %v, got %v",
			expectedTodo.Completed,
			response.Completed,
		)
	}
}

func TestCreateTodoHandlerInvalidJSON(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		createTodoFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			serviceCalled = true
			return nil, nil
		},
	}

	handler := NewTodoHandler(mockService)

	requestBody := `{
		"title": "Learn Go",
		"description":
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		bytes.NewBufferString(requestBody),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateTodo(
		recorder,
		req,
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if serviceCalled {
		t.Fatal("expected service not to be called")
	}
}

func TestCreateTodoHandlerServiceError(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		createTodoFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			serviceCalled = true

			return nil, errors.New("service error")
		},
	}

	handler := NewTodoHandler(mockService)

	requestBody := `{
		"title": "Learn Go",
		"description": "Learn Go"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		bytes.NewBufferString(requestBody),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("service error"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"service error",
			recorder.Body.String(),
		)
	}
}

func TestCreateTodoHandlerMethodNotAllowed(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		createTodoFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			serviceCalled = true
			return nil, nil
		},
	}

	handler := NewTodoHandler(mockService)

	requestBody := `{
		"title": "Learn Go",
		"description": "Learn Go"
	}`

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos",
		bytes.NewBufferString(requestBody),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateTodo(
		recorder,
		req,
	)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}

	if serviceCalled {
		t.Fatal("expected service not to be called")
	}
}

func TestGetTodosHandler(t *testing.T) {
	description1 := "Learn Go"
	description2 := "Learn PostgreSQL"

	expectedTodos := []model.Todo{
		{
			ID:          1,
			Title:       "Learn Go",
			Description: &description1,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Title:       "Learn PostgreSQL",
			Description: &description2,
			Completed:   true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	serviceCalled := false

	mockService := &mockTodoService{
		createTodoFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			return nil, nil
		},

		getTodosFunc: func(
			ctx context.Context,
		) ([]model.Todo, error) {
			serviceCalled = true

			return expectedTodos, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetTodos(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response []model.Todo

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if len(response) != len(expectedTodos) {
		t.Fatalf(
			"expected %d todos, got %d",
			len(expectedTodos),
			len(response),
		)
	}

	if response[0].ID != expectedTodos[0].ID {
		t.Fatalf(
			"expected first todo ID %d, got %d",
			expectedTodos[0].ID,
			response[0].ID,
		)
	}

	if response[0].Title != expectedTodos[0].Title {
		t.Fatalf(
			"expected first todo title %q, got %q",
			expectedTodos[0].Title,
			response[0].Title,
		)
	}

	if response[1].ID != expectedTodos[1].ID {
		t.Fatalf(
			"expected second todo ID %d, got %d",
			expectedTodos[1].ID,
			response[1].ID,
		)
	}
}

func TestGetTodosHandlerMethodNotAllowed(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		getTodosFunc: func(
			ctx context.Context,
		) ([]model.Todo, error) {
			serviceCalled = true

			return nil, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetTodos(
		recorder,
		req,
	)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}

	if serviceCalled {
		t.Fatal("expected service not to be called")
	}
}

func TestGetTodoHandler(t *testing.T) {
	description := "Learn Go"

	expectedTodo := &model.Todo{
		ID:          1,
		Title:       "Learn Go",
		Description: &description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	serviceCalled := false

	mockService := &mockTodoService{
		getTodoFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Todo, error) {
			serviceCalled = true

			if id != 1 {
				t.Fatalf(
					"expected ID %d, got %d",
					1,
					id,
				)
			}

			return expectedTodo, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response model.Todo

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response.ID != expectedTodo.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			expectedTodo.ID,
			response.ID,
		)
	}

	if response.Title != expectedTodo.Title {
		t.Fatalf(
			"expected title %q, got %q",
			expectedTodo.Title,
			response.Title,
		)
	}
}

func TestGetTodoHandlerInvalidID(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		getTodoFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Todo, error) {
			serviceCalled = true

			return nil, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetTodo(
		recorder,
		req,
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if serviceCalled {
		t.Fatal("expected service not to be called")
	}
}

func TestGetTodoHandlerServiceError(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		getTodoFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Todo, error) {
			serviceCalled = true

			return nil, errors.New("service error")
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("service error"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"service error",
			recorder.Body.String(),
		)
	}
}

func TestGetTodoHandlerNotFound(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		getTodoFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Todo, error) {
			serviceCalled = true

			return nil, service.ErrTodoNotFound
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/todos/999",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("todo not found"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"todo not found",
			recorder.Body.String(),
		)
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	description := "Learn advanced Go"

	expectedTodo := &model.Todo{
		ID:          1,
		Title:       "Learn Go",
		Description: &description,
		Completed:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	serviceCalled := false

	mockService := &mockTodoService{
		updateTodoFunc: func(
			ctx context.Context,
			id int64,
			title string,
			description *string,
			completed bool,
		) (*model.Todo, error) {
			serviceCalled = true

			if id != 1 {
				t.Fatalf(
					"expected ID %d, got %d",
					1,
					id,
				)
			}

			if title != "Learn Go" {
				t.Fatalf(
					"expected title %q, got %q",
					"Learn Go",
					title,
				)
			}

			if description == nil {
				t.Fatal("expected description, got nil")
			}

			if *description != "Learn advanced Go" {
				t.Fatalf(
					"expected description %q, got %q",
					"Learn advanced Go",
					*description,
				)
			}

			if !completed {
				t.Fatal("expected completed to be true")
			}

			return expectedTodo, nil
		},
	}

	handler := NewTodoHandler(mockService)

	requestBody := `{
		"title": "Learn Go",
		"description": "Learn advanced Go",
		"completed": true
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/1",
		bytes.NewBufferString(requestBody),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.UpdateTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response model.Todo

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	if err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response.ID != expectedTodo.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			expectedTodo.ID,
			response.ID,
		)
	}

	if response.Title != expectedTodo.Title {
		t.Fatalf(
			"expected title %q, got %q",
			expectedTodo.Title,
			response.Title,
		)
	}

	if response.Completed != expectedTodo.Completed {
		t.Fatalf(
			"expected completed %v, got %v",
			expectedTodo.Completed,
			response.Completed,
		)
	}
}

func TestUpdateTodoHandlerInvalidID(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		updateTodoFunc: func(
			ctx context.Context,
			id int64,
			title string,
			description *string,
			completed bool,
		) (*model.Todo, error) {
			serviceCalled = true
			return nil, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/abc",
		bytes.NewBufferString(`{
			"title": "Learn Go",
			"description": "Learn advanced Go",
			"completed": true
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.UpdateTodo(
		recorder,
		req,
	)

	if serviceCalled {
		t.Fatal("expected service not to be called")
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("invalid todo id"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"invalid todo id",
			recorder.Body.String(),
		)
	}
}

func TestUpdateTodoHandlerInvalidJSON(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		updateTodoFunc: func(
			ctx context.Context,
			id int64,
			title string,
			description *string,
			completed bool,
		) (*model.Todo, error) {
			serviceCalled = true
			return nil, nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/1",
		bytes.NewBufferString(`{
			"title": "Learn Go",
			"description":
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.UpdateTodo(
		recorder,
		req,
	)

	if serviceCalled {
		t.Fatal("expected service not to be called")
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("invalid request body"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"invalid request body",
			recorder.Body.String(),
		)
	}
}

func TestUpdateTodoHandlerServiceError(t *testing.T) {
	mockService := &mockTodoService{
		updateTodoFunc: func(
			ctx context.Context,
			id int64,
			title string,
			description *string,
			completed bool,
		) (*model.Todo, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/1",
		bytes.NewBufferString(`{
			"title": "Learn Go",
			"description": "Learn advanced Go",
			"completed": true
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.UpdateTodo(
		recorder,
		req,
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("database error"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"database error",
			recorder.Body.String(),
		)
	}
}

func TestUpdateTodoHandlerNotFound(t *testing.T) {
	mockService := &mockTodoService{
		updateTodoFunc: func(
			ctx context.Context,
			id int64,
			title string,
			description *string,
			completed bool,
		) (*model.Todo, error) {
			return nil, service.ErrTodoNotFound
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/todos/999",
		bytes.NewBufferString(`{
			"title": "Learn Go",
			"description": "Learn advanced Go",
			"completed": true
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.UpdateTodo(
		recorder,
		req,
	)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}

	if !bytes.Contains(
		recorder.Body.Bytes(),
		[]byte("todo not found"),
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			"todo not found",
			recorder.Body.String(),
		)
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	serviceCalled := false

	mockService := &mockTodoService{
		deleteTodoFunc: func(
			ctx context.Context,
			id int64,
		) error {
			serviceCalled = true

			if id != 1 {
				t.Fatalf("expected id 1, got %d", id)
			}

			return nil
		},
	}

	handler := NewTodoHandler(mockService)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/todos/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.DeleteTodo(
		recorder,
		req,
	)

	if !serviceCalled {
		t.Fatal("expected service to be called")
	}

	if recorder.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			recorder.Code,
		)
	}
}
