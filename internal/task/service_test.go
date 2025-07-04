package task

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewTaskService(logger)

	task := service.CreateTask()

	assert.NotNil(t, task)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, StatusCreated, task.Status)
	assert.WithinDuration(t, time.Now(), task.CreatedAt, time.Second)
}

func TestGetTask(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewTaskService(logger)

	task := service.CreateTask()
	found, err := service.GetTask(task.ID)
	
	assert.NoError(t, err)
	assert.Equal(t, task.ID, found.ID)
}

func TestGetTask_NotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewTaskService(logger)

	_, err := service.GetTask("nonexistent")

	assert.Error(t, err)
}

func TestDeleteTask(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewTaskService(logger)
	
	task := service.CreateTask()
	err := service.DeleteTask(task.ID)

	assert.NoError(t, err)
	_, err = service.GetTask(task.ID)
	assert.Error(t, err)
}

func TestTaskCancellation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewTaskService(logger)

	task := service.CreateTask()
	service.DeleteTask(task.ID)

	time.Sleep(200 * time.Millisecond)

	_, err := service.GetTask(task.ID)
	assert.Error(t, err)
}

func TestListTasks_Filtering(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewTaskService(logger)

	t1 := service.CreateTask()
	t2 := service.CreateTask()

	service.DeleteTask(t2.ID)

	time.Sleep(100 * time.Millisecond)

	all := service.ListTasks(nil)
	assert.GreaterOrEqual(t, len(all), 2)

	status := StatusRunning
	running := service.ListTasks(&status)
	assert.True(t, containsTaskWithID(running, t1.ID))
	assert.False(t, containsTaskWithID(running, t2.ID))
}

func containsTaskWithID(tasks []*Task, id string) bool {
	for _, task := range tasks {
		if task.ID == id {
			return true
		}
	}
	return false
}
