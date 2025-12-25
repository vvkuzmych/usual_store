package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer handles consuming messages from Kafka
type Consumer struct {
	reader  *kafka.Reader
	logger  *log.Logger
	handler EmailHandler
}

// EmailHandler interface for processing email messages
type EmailHandler interface {
	SendEmail(from, to, subject, template string, data map[string]interface{}) error
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic, groupID string, handler EmailHandler, logger *log.Logger) *Consumer {
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

	return &Consumer{
		reader:  reader,
		logger:  logger,
		handler: handler,
	}
}

// Start starts consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Println("Starting Kafka consumer...")

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

			if err := c.processMessage(ctx, msg); err != nil {
				c.logger.Printf("Error processing message: %v", err)
				// Don't commit on error - message will be retried
				continue
			}

			// Commit the message after successful processing
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Printf("Error committing message: %v", err)
			}
		}
	}
}

// processMessage processes a single Kafka message
func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	c.logger.Printf("Processing message: offset=%d, partition=%d", msg.Offset, msg.Partition)

	var event EmailEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Check if message has exceeded retry limit
	if event.Message.RetryCount >= event.Message.MaxRetries {
		c.logger.Printf("Message exceeded max retries, sending to DLQ: ID=%s", event.Message.ID)
		// In production, you would send this to a Dead Letter Queue (DLQ)
		return fmt.Errorf("message exceeded max retries")
	}

	// Process the email
	err := c.handler.SendEmail(
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
		c.logger.Printf("Failed to send email (retry %d/%d): ID=%s, Error=%v",
			event.Message.RetryCount, event.Message.MaxRetries, event.Message.ID, err)
		return err
	}

	event.Message.Status = StatusSent
	c.logger.Printf("Email sent successfully: ID=%s, To=%s, Type=%s",
		event.Message.ID, event.Message.To, event.Type)

	return nil
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
