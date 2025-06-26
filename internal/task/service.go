package task

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	mu sync.Mutex
	tasks map[string]*Task
	logger *slog.Logger
}

func NewTaskService(logger *slog.Logger) *TaskService{
	return &TaskService{
		tasks: make(map[string]*Task),
		logger: logger,
	}
}

func (s *TaskService) CreateTask() *Task {
	id := uuid.New().String()
	now := time.Now()

	ctx, cancel := context.WithCancel(context.Background())

	task := &Task{
		ID: id,
		Status: StatusCreated,
		CreatedAt: now,
		Ctx: ctx,
		Cancel: cancel,
	}

	s.mu.Lock()
	s.tasks[id] = task
	s.mu.Unlock()

	s.logger.Info("Task created", slog.String("task_id", id))

	go s.runTask(task)

	return task
}

func (s *TaskService) runTask(task *Task) {
	s.mu.Lock()
	started := time.Now()
	task.Status = StatusRunning
	task.StartedAt = &started
	s.mu.Unlock()

	sleepDuration := time.Duration(3 + time.Now().UnixNano()%3) * time.Minute
	s.logger.Info("Task started", slog.String("task_id", task.ID),
	 	slog.String("estimated_duration", sleepDuration.String()))
	
	select {
	case <- time.After(sleepDuration):
		s.mu.Lock()
		finished := time.Now()
		task.Status = StatusDone
		task.FinishedAt = &finished
		task.Duration = finished.Sub(*task.StartedAt)
		task.Result = "task completed successfully"
		s.mu.Unlock()

		s.logger.Info("Task completed", slog.String("task_id", task.ID))
	case <- task.Ctx.Done():
		s.mu.Lock()
		task.Status = StatusFailed
		task.Result = "task was cancelled"
		s.mu.Unlock()

		s.logger.Warn("Task cancelled", slog.String("task_id", task.ID))
	}
}

func (s *TaskService) GetTask(id string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		s.logger.Warn("Task not found", slog.String("task_id", id))
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (s *TaskService) DeleteTask(id string) error {
	s.mu.Lock()
	task, ok := s.tasks[id]
	if !ok {
		s.mu.Unlock()
		s.logger.Warn("Task not found for deletion", slog.String("task_id", id))
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	s.mu.Unlock()

	task.Cancel()

	s.logger.Info("Task deleted", slog.String("task_id", id))
	return nil
}

func (s *TaskService) ListTasks(filterStatus *Status) []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []*Task
	for _, task := range s.tasks {
		if filterStatus == nil || filterStatus == &task.Status {
			result = append(result, task)
		}
	}

	return result
}