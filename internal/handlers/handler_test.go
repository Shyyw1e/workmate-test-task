package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	
	"github.com/Shyyw1e/workmate-test-task/internal/task"
)

func TestCreateTaskHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := task.NewTaskService(logger)
	handler := NewHandler(service, logger)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var task task.Task
	err := json.NewDecoder(rr.Body).Decode(&task)
	assert.NoError(t, err)
	assert.NotEmpty(t, task.ID)
}

func TestGetTaskHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := task.NewTaskService(logger)
	handler := NewHandler(service, logger)

	newTask := service.CreateTask()

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/tasks/" + newTask.ID, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	
	var result task.Task
	err := json.NewDecoder(rr.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, newTask.ID, result.ID)

}

func TestGetTaskHandler_NotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := task.NewTaskService(logger)
	handler := NewHandler(service, logger)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/tasks/invalid-id", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
