package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/raffi/todo-app/internal/model"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (r *TodoRepository) Create(
	ctx context.Context,
	title string,
	description *string,
) (*model.Todo, error) {
	todo := &model.Todo{}

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO todos (
			title,
			description
		)
		VALUES ($1, $2)
		RETURNING
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
		`,
		title,
		description,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return todo, nil
}