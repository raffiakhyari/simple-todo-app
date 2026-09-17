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
