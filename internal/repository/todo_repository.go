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

func (r *TodoRepository) FindAll(
	ctx context.Context,
) ([]model.Todo, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
		FROM todos
		ORDER BY created_at DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []model.Todo

	for rows.Next() {
		var todo model.Todo

		err := rows.Scan(
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

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) FindByID(
	ctx context.Context,
	id int64,
) (*model.Todo, error) {
	todo := &model.Todo{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
		FROM todos
		WHERE id = $1
		`,
		id,
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

func (r *TodoRepository) Update(
	ctx context.Context,
	id int64,
	title string,
	description *string,
	completed bool,
) (*model.Todo, error) {
	todo := &model.Todo{}

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE todos
		SET
			title = $1,
			description = $2,
			completed = $3,
			updated_at = NOW()
		WHERE id = $4
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
		completed,
		id,
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
