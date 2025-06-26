package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Shyyw1e/workmate-test-task/internal/task"
)

type Handler struct {
	service *task.TaskService
	logger  *slog.Logger
}

func NewHandler(service *task.TaskService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/tasks", h.CreateTask)
	r.Get("/tasks/{id}", h.GetTask)
	r.Delete("/tasks/{id}", h.DeleteTask)
	r.Get("/tasks", h.ListTasks)

}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := h.service.CreateTask()
	h.logger.Info("Handled CreateTask", slog.String("task_id", task.ID))
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.service.GetTask(id)
	if err != nil {
		h.logger.Warn("Task not found", slog.String("task_id", id))
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteTask(id)
	if err != nil {
		h.logger.Warn("DeleteTask failed", slog.String("task_id", id))
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	h.logger.Info("Handled DeleteTask", slog.String("task_id", id))
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("status")

	var status *task.Status
	if query != "" {
		s := task.Status(query)
		status = &s
	}

	tasks := h.service.ListTasks(status)
	writeJSON(w, http.StatusOK, tasks)
}
