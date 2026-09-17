package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/raffi/todo-app/internal/model"
)

type TodoRepository interface {
	Create(
		ctx context.Context,
		title string,
		description *string,
	) (*model.Todo, error)

	FindAll(
		ctx context.Context,
	) ([]model.Todo, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (*model.Todo, error)

	Update(
		ctx context.Context,
		id int64,
		title string,
		description *string,
		completed bool,
	) (*model.Todo, error)

	Delete(
		ctx context.Context,
		id int64,
	) error
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) CreateTodo(
	ctx context.Context,
	title string,
	description *string,
) (*model.Todo, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return s.repo.Create(
		ctx,
		title,
		description,
	)
}

func (s *TodoService) GetTodos(
	ctx context.Context,
) ([]model.Todo, error) {
	return s.repo.FindAll(ctx)
}

func (s *TodoService) GetTodo(
	ctx context.Context,
	id int64,
) (*model.Todo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TodoService) UpdateTodo(
	ctx context.Context,
	id int64,
	title string,
	description *string,
	completed bool,
) (*model.Todo, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return s.repo.Update(
		ctx,
		id,
		title,
		description,
		completed,
	)
}

func (s *TodoService) DeleteTodo(
	ctx context.Context,
	id int64,
) error {
	return s.repo.Delete(ctx, id)
}
