package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/raffi/todo-app/internal/config"
	"github.com/raffi/todo-app/internal/database"
	"github.com/raffi/todo-app/internal/handler"
	"github.com/raffi/todo-app/internal/logger"
	"github.com/raffi/todo-app/internal/middleware"
	"github.com/raffi/todo-app/internal/repository"
	"github.com/raffi/todo-app/internal/service"
)

func main() {
	appLogger := logger.New()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := database.NewPostgresPool(
		ctx,
		database.Config{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			Name:     cfg.DBName,
			SSLMode:  cfg.DBSSLMode,
		},
	)
	if err != nil {
		log.Fatalf(
			"database connection failed: %v",
			err,
		)
	}
	defer db.Close()

	appLogger.Info(
		"database connection established",
	)

	// Repository
	todoRepository := repository.NewTodoRepository(db)

	// Service
	todoService := service.NewTodoService(todoRepository)

	// Handler
	todoHandler := handler.NewTodoHandler(
		todoService,
		appLogger,
	)

	// Health check
	http.HandleFunc(
		"/health",
		healthHandler,
	)

	// Todo collection:
	// GET  /todos
	// POST /todos
	http.HandleFunc(
		"/todos",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				todoHandler.GetTodos(w, r)

			case http.MethodPost:
				todoHandler.CreateTodo(w, r)

			default:
				w.WriteHeader(
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	// Todo item:
	// GET    /todos/{id}
	// PUT    /todos/{id}
	// DELETE /todos/{id}
	http.HandleFunc(
		"/todos/",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				todoHandler.GetTodo(w, r)

			case http.MethodPut:
				todoHandler.UpdateTodo(w, r)

			case http.MethodDelete:
				todoHandler.DeleteTodo(w, r)

			default:
				w.WriteHeader(
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	addr := ":" + cfg.AppPort

	appLogger.Info(
		"todo API is running",
		"address", addr,
	)

	// CORS middleware
	appHandler := middleware.CORS(
		cfg.CORSAllowedOrigin,
	)(http.DefaultServeMux)

	if err := http.ListenAndServe(
		addr,
		appHandler,
	); err != nil {
		log.Fatal(err)
	}

}

func healthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	w.Write(
		[]byte(`{"status":"ok"}`),
	)

}
