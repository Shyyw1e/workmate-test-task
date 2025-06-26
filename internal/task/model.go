package task

import (
	"context"
	"time"
)

type Status string

const (
	StatusCreated Status = "created"
	StatusRunning Status = "running"
	StatusDone Status = "done"
	StatusFailed Status = "failed"
)

type Task struct {
	ID 			string				`json:"id"`
	Status 		Status				`json:"status"`
	CreatedAt 	time.Time			`json:"created_at"`
	StartedAt 	*time.Time			`json:"started_at,omitempty"`
	FinishedAt 	*time.Time			`json:"finished_at,omitempty"`
	Duration 	time.Duration		`json:"duration,omitempty"`
	Result 		string				`json:"result,omitempty"`

	Ctx 		context.Context		`json:"-"`
	Cancel 		context.CancelFunc	`json:"-"`
}