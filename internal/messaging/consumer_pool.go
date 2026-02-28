package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"usual_store/internal/workerpool"

	"github.com/segmentio/kafka-go"
)

// PooledConsumer handles consuming messages from Kafka with worker pool
type PooledConsumer struct {
	reader     *kafka.Reader
	logger     *log.Logger
	handler    EmailHandler
	workerPool *workerpool.Pool
}

// EmailJob implements workerpool.Job interface for email processing
type EmailJob struct {
	message kafka.Message
	handler EmailHandler
	logger  *log.Logger
	jobID   string
}

// NewEmailJob creates a new email job
func NewEmailJob(msg kafka.Message, handler EmailHandler, logger *log.Logger) workerpool.Job {
	return &EmailJob{
		message: msg,
		handler: handler,
		logger:  logger,
		jobID:   fmt.Sprintf("email-%d-%d", msg.Partition, msg.Offset),
	}
}

// Execute processes the email job
func (e *EmailJob) Execute(ctx context.Context) error {
	e.logger.Printf("Processing email job: %s", e.jobID)

	var event EmailEvent
	if err := json.Unmarshal(e.message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Check if message has exceeded retry limit
	if event.Message.RetryCount >= event.Message.MaxRetries {
		e.logger.Printf("Message exceeded max retries: ID=%s", event.Message.ID)
		return fmt.Errorf("message exceeded max retries")
	}

	// Send the email
	err := e.handler.SendEmail(
		event.Message.From,
		event.Message.To,
		event.Message.Subject,
		event.Message.Template,
		event.Message.Data,
	)

	if err != nil {
		event.Message.RetryCount++
		event.Message.Status = StatusFailed
		event.Message.ErrorMsg = err.Error()
		e.logger.Printf("Failed to send email (retry %d/%d): ID=%s, Error=%v",
			event.Message.RetryCount, event.Message.MaxRetries, event.Message.ID, err)
		return err
	}

	event.Message.Status = StatusSent
	e.logger.Printf("Email sent successfully: ID=%s, To=%s", event.Message.ID, event.Message.To)

	return nil
}

// ID returns the job identifier
func (e *EmailJob) ID() string {
	return e.jobID
}

// NewPooledConsumer creates a new Kafka consumer with worker pool
func NewPooledConsumer(
	brokers []string,
	topic, groupID string,
	handler EmailHandler,
	workerPool *workerpool.Pool,
	logger *log.Logger,
) *PooledConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		MaxAttempts:    3,
	})

	return &PooledConsumer{
		reader:     reader,
		logger:     logger,
		handler:    handler,
		workerPool: workerPool,
	}
}

// Start starts consuming messages from Kafka and submitting to worker pool
func (c *PooledConsumer) Start(ctx context.Context) error {
	c.logger.Println("Starting pooled Kafka consumer...")

	// Start goroutine to handle results
	go c.handleResults(ctx)

	for {
		select {
		case <-ctx.Done():
			c.logger.Println("Consumer context cancelled, shutting down...")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				c.logger.Printf("Error fetching message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Create job and submit to worker pool
			job := NewEmailJob(msg, c.handler, c.logger)

			// Try non-blocking submit first
			if err := c.workerPool.Submit(job); err != nil {
				// If queue is full, use blocking submit
				c.logger.Printf("Job queue full, using blocking submit...")
				if err := c.workerPool.SubmitBlocking(job); err != nil {
					c.logger.Printf("Failed to submit job: %v", err)
					continue
				}
			}

			// Commit message immediately after successful submission to pool
			// Worker pool will handle the actual processing
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Printf("Error committing message: %v", err)
			}
		}
	}
}

// handleResults processes worker pool results
func (c *PooledConsumer) handleResults(ctx context.Context) {
	results := c.workerPool.Results()
	for {
		select {
		case <-ctx.Done():
			return
		case result, ok := <-results:
			if !ok {
				return
			}

			// Type assertion to access fields
			if result.Success {
				c.logger.Printf("Job completed successfully: %s", result.JobID)
			} else {
				c.logger.Printf("Job failed: %s, Error: %v", result.JobID, result.Error)
				// In production, you could implement retry logic or DLQ here
			}
		}
	}
}

// Close closes the Kafka consumer
func (c *PooledConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
