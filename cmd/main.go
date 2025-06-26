package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/Shyyw1e/workmate-test-task/internal/task"
	"github.com/Shyyw1e/workmate-test-task/internal/handlers"
	"github.com/Shyyw1e/workmate-test-task/pkg/logger"
)

func main() {
	log := logger.New(slog.LevelDebug)

	service := task.NewTaskService(log)
	handler := handlers.NewHandler(service, log)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler.RegisterRoutes(r)

	log.Info("Server started", slog.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("Server failed", slog.Any("err", err))
		os.Exit(1)
	}
}