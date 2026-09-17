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
