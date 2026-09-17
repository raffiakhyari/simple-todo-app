package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raffi/todo-app/internal/model"
)

type mockTodoRepository struct {
	createFunc func(
		ctx context.Context,
		title string,
		description *string,
	) (*model.Todo, error)

	findAllFunc func(
		ctx context.Context,
	) ([]model.Todo, error)

	findByIDFunc func(
		ctx context.Context,
		id int64,
	) (*model.Todo, error)

	updateFunc func(
		ctx context.Context,
		id int64,
		title string,
		description *string,
		completed bool,
	) (*model.Todo, error)
}

func (m *mockTodoRepository) Create(
	ctx context.Context,
	title string,
	description *string,
) (*model.Todo, error) {
	return m.createFunc(ctx, title, description)
}

func (m *mockTodoRepository) FindAll(
	ctx context.Context,
) ([]model.Todo, error) {
	return m.findAllFunc(ctx)
}

func (m *mockTodoRepository) FindByID(
	ctx context.Context,
	id int64,
) (*model.Todo, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockTodoRepository) Update(
	ctx context.Context,
	id int64,
	title string,
	description *string,
	completed bool,
) (*model.Todo, error) {
	return m.updateFunc(
		ctx,
		id,
		title,
		description,
		completed,
	)
}

func TestCreateTodo(t *testing.T) {
	repo := &mockTodoRepository{
		createFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			return &model.Todo{
				ID:          1,
				Title:       title,
				Description: description,
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},

		findAllFunc: func(
			ctx context.Context,
		) ([]model.Todo, error) {
			return nil, nil
		},
	}

	svc := NewTodoService(repo)

	description := "Learn Go"

	todo, err := svc.CreateTodo(
		context.Background(),
		"  Build Todo API  ",
		&description,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if todo.Title != "Build Todo API" {
		t.Fatalf(
			"expected title %q, got %q",
			"Build Todo API",
			todo.Title,
		)
	}
}

func TestCreateTodoRequiresTitle(t *testing.T) {
	repo := &mockTodoRepository{
		createFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			t.Fatal("repository should not be called")

			return nil, nil
		},

		findAllFunc: func(
			ctx context.Context,
		) ([]model.Todo, error) {
			return nil, nil
		},
	}

	svc := NewTodoService(repo)

	_, err := svc.CreateTodo(
		context.Background(),
		"   ",
		nil,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "title is required" {
		t.Fatalf(
			"expected error %q, got %q",
			"title is required",
			err.Error(),
		)
	}
}

func TestGetTodos(t *testing.T) {
	expectedTodos := []model.Todo{
		{
			ID:        1,
			Title:     "Learn Go",
			Completed: false,
		},
		{
			ID:        2,
			Title:     "Build Todo API",
			Completed: true,
		},
	}

	repo := &mockTodoRepository{
		createFunc: func(
			ctx context.Context,
			title string,
			description *string,
		) (*model.Todo, error) {
			return nil, nil
		},

		findAllFunc: func(
			ctx context.Context,
		) ([]model.Todo, error) {
			return expectedTodos, nil
		},
	}

	svc := NewTodoService(repo)

	todos, err := svc.GetTodos(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(todos) != 2 {
		t.Fatalf(
			"expected 2 todos, got %d",
			len(todos),
		)
	}

	if todos[0].Title != "Learn Go" {
		t.Fatalf(
			"expected first todo title %q, got %q",
			"Learn Go",
			todos[0].Title,
		)
	}

	if todos[1].Title != "Build Todo API" {
		t.Fatalf(
			"expected second todo title %q, got %q",
			"Build Todo API",
			todos[1].Title,
		)
	}
}

func TestGetTodo(t *testing.T) {
	ctx := context.Background()

	expectedTodo := &model.Todo{
		ID:        1,
		Title:     "Learn Go",
		Completed: false,
	}

	repo := &mockTodoRepository{
		findByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Todo, error) {
			if id != 1 {
				t.Fatalf("expected ID 1, got %d", id)
			}

			return expectedTodo, nil
		},
	}

	svc := NewTodoService(repo)

	todo, err := svc.GetTodo(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if todo.ID != expectedTodo.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			expectedTodo.ID,
			todo.ID,
		)
	}

	if todo.Title != expectedTodo.Title {
		t.Fatalf(
			"expected title %q, got %q",
			expectedTodo.Title,
			todo.Title,
		)
	}

	if todo.Completed != expectedTodo.Completed {
		t.Fatalf(
			"expected completed %v, got %v",
			expectedTodo.Completed,
			todo.Completed,
		)
	}
}

func TestGetTodoRepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("database error")

	repo := &mockTodoRepository{
		findByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Todo, error) {
			return nil, expectedErr
		},
	}

	svc := NewTodoService(repo)

	todo, err := svc.GetTodo(ctx, 1)

	if todo != nil {
		t.Fatal("expected todo to be nil")
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}
