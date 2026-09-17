package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/raffi/todo-app/internal/config"
	"github.com/raffi/todo-app/internal/database"
)

func main() {
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
		log.Fatalf("database connection failed: %v", err)
	}

	defer db.Close()

	http.HandleFunc("/health", healthHandler)

	addr := ":" + cfg.AppPort

	log.Printf("Todo API is running on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status":"ok"}`))
}
