package workerpool

import "context"

// Job represents a unit of work to be processed by a worker
type Job interface {
	Execute(ctx context.Context) error
	ID() string
}

// Result represents the outcome of processing a job
type Result struct {
	JobID   string
	Success bool
	Error   error
}
