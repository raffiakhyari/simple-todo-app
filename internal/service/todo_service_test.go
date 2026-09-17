package service

import (
	"context"
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
}

func (m *mockTodoRepository) Create(
	ctx context.Context,
	title string,
	description *string,
) (*model.Todo, error) {
	return m.createFunc(ctx, title, description)
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
