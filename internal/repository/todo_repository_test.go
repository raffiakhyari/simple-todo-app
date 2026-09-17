package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTodoRepositoryCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := pgxpool.New(
		ctx,
		"postgres://todo-app:testpassword@localhost:5432/todo-app?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}

	defer db.Close()

	repo := NewTodoRepository(db)

	description := "Learn Go"

	todo, err := repo.Create(
		ctx,
		"Build Todo API",
		&description,
	)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	if todo.ID == 0 {
		t.Fatal("expected todo ID to be generated")
	}

	if todo.Title != "Build Todo API" {
		t.Fatalf(
			"expected title %q, got %q",
			"Build Todo API",
			todo.Title,
		)
	}

	if todo.Description == nil {
		t.Fatal("expected description to be set")
	}

	if *todo.Description != "Learn Go" {
		t.Fatalf(
			"expected description %q, got %q",
			"Learn Go",
			*todo.Description,
		)
	}

	if todo.Completed {
		t.Fatal("expected todo to be incomplete")
	}
}

func TestTodoRepositoryFindAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := pgxpool.New(
		ctx,
		"postgres://todo-app:testpassword@localhost:5432/todo-app?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}

	defer db.Close()

	repo := NewTodoRepository(db)

	description := "Repository test"

	_, err = repo.Create(
		ctx,
		"Find All Test",
		&description,
	)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	todos, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("failed to find todos: %v", err)
	}

	if len(todos) == 0 {
		t.Fatal("expected at least one todo")
	}
}

func TestTodoRepositoryFindByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := pgxpool.New(
		ctx,
		"postgres://todo-app:testpassword@localhost:5432/todo-app?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	description := "Test description"

	createdTodo, err := repo.Create(
		ctx,
		"Test FindByID",
		&description,
	)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	todo, err := repo.FindByID(
		ctx,
		createdTodo.ID,
	)
	if err != nil {
		t.Fatalf("failed to find todo by id: %v", err)
	}

	if todo.ID != createdTodo.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			createdTodo.ID,
			todo.ID,
		)
	}

	if todo.Title != createdTodo.Title {
		t.Fatalf(
			"expected title %q, got %q",
			createdTodo.Title,
			todo.Title,
		)
	}

	if todo.Description == nil {
		t.Fatal("expected description, got nil")
	}

	if *todo.Description != description {
		t.Fatalf(
			"expected description %q, got %q",
			description,
			*todo.Description,
		)
	}

	if todo.Completed != createdTodo.Completed {
		t.Fatalf(
			"expected completed %v, got %v",
			createdTodo.Completed,
			todo.Completed,
		)
	}
}

func TestTodoRepositoryUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := pgxpool.New(
		ctx,
		"postgres://todo-app:testpassword@localhost:5432/todo-app?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	description := "Original description"

	createdTodo, err := repo.Create(
		ctx,
		"Original Title",
		&description,
	)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	updatedDescription := "Updated description"

	updatedTodo, err := repo.Update(
		ctx,
		createdTodo.ID,
		"Updated Title",
		&updatedDescription,
		true,
	)
	if err != nil {
		t.Fatalf("failed to update todo: %v", err)
	}

	if updatedTodo.ID != createdTodo.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			createdTodo.ID,
			updatedTodo.ID,
		)
	}

	if updatedTodo.Title != "Updated Title" {
		t.Fatalf(
			"expected title %q, got %q",
			"Updated Title",
			updatedTodo.Title,
		)
	}

	if updatedTodo.Description == nil {
		t.Fatal("expected description, got nil")
	}

	if *updatedTodo.Description != "Updated description" {
		t.Fatalf(
			"expected description %q, got %q",
			"Updated description",
			*updatedTodo.Description,
		)
	}

	if !updatedTodo.Completed {
		t.Fatal("expected todo to be completed")
	}

	if !updatedTodo.UpdatedAt.After(createdTodo.UpdatedAt) {
		t.Fatalf(
			"expected updated_at to change, created_at=%v updated_at=%v",
			createdTodo.UpdatedAt,
			updatedTodo.UpdatedAt,
		)
	}
}

func TestTodoRepositoryDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := pgxpool.New(
		ctx,
		"postgres://todo-app:testpassword@localhost:5432/todo-app?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	description := "Delete test"

	createdTodo, err := repo.Create(
		ctx,
		"Delete Test",
		&description,
	)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	err = repo.Delete(
		ctx,
		createdTodo.ID,
	)
	if err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}

	deletedTodo, err := repo.FindByID(
		ctx,
		createdTodo.ID,
	)

	if err == nil {
		t.Fatal("expected error when finding deleted todo")
	}

	if deletedTodo != nil {
		t.Fatal("expected deleted todo to be nil")
	}
}
