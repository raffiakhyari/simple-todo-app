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