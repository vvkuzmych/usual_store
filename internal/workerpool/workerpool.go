package workerpool

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Pool manages a fixed number of workers processing jobs
type Pool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *log.Logger
}

// New creates a new worker pool
func New(numWorkers int, bufferSize int, logger *log.Logger) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	if logger == nil {
		logger = log.Default()
	}

	return &Pool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, bufferSize),
		results:    make(chan Result, bufferSize),
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
	}
}

// Start launches all workers
func (p *Pool) Start() {
	p.logger.Printf("Starting worker pool with %d workers", p.numWorkers)

	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker processes jobs from the jobs channel
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	p.logger.Printf("Worker %d started", id)

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Printf("Worker %d shutting down", id)
			return

		case job, ok := <-p.jobs:
			if !ok {
				p.logger.Printf("Worker %d: jobs channel closed", id)
				return
			}

			// Process job
			err := job.Execute(p.ctx)

			// Send result
			result := Result{
				JobID:   job.ID(),
				Success: err == nil,
				Error:   err,
			}

			select {
			case p.results <- result:
				if err != nil {
					p.logger.Printf("Worker %d: job %s failed: %v", id, job.ID(), err)
				}
			case <-p.ctx.Done():
				return
			}
		}
	}
}

// Submit sends a job to the worker pool (non-blocking if buffer available)
func (p *Pool) Submit(job Job) error {
	select {
	case p.jobs <- job:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	default:
		return fmt.Errorf("worker pool job queue is full")
	}
}

// SubmitBlocking sends a job to the worker pool (blocks if buffer full)
func (p *Pool) SubmitBlocking(job Job) error {
	select {
	case p.jobs <- job:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	}
}

// Results returns the results channel for reading processed job results
func (p *Pool) Results() <-chan Result {
	return p.results
}

// Stop gracefully shuts down the worker pool
func (p *Pool) Stop() {
	p.logger.Println("Stopping worker pool...")

	// Signal workers to stop accepting new jobs
	close(p.jobs)

	// Wait for all workers to finish current jobs
	p.wg.Wait()

	// Cancel context
	p.cancel()

	// Close results channel
	close(p.results)

	p.logger.Println("Worker pool stopped")
}

// StopWithContext stops the pool with a context for timeout
func (p *Pool) StopWithContext(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		p.Stop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("worker pool shutdown timeout: %w", ctx.Err())
	}
}
