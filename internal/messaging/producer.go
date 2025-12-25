package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Producer handles sending messages to Kafka
type Producer struct {
	writer *kafka.Writer
	logger *log.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string, logger *log.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Compression:  kafka.Snappy,
		MaxAttempts:  3,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

// SendEmail sends an email message to Kafka
func (p *Producer) SendEmail(ctx context.Context, from, to, subject, template string, data map[string]interface{}, priority string) error {
	if priority == "" {
		priority = PriorityNormal
	}

	msg := EmailMessage{
		ID:         uuid.New().String(),
		From:       from,
		To:         to,
		Subject:    subject,
		Template:   template,
		Data:       data,
		Priority:   priority,
		Timestamp:  time.Now(),
		RetryCount: 0,
		MaxRetries: 3,
		Status:     StatusPending,
	}

	event := EmailEvent{
		Type:    determineEmailType(template),
		Message: msg,
	}

	return p.SendEvent(ctx, event)
}

// SendEvent sends a generic email event to Kafka
func (p *Producer) SendEvent(ctx context.Context, event EmailEvent) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		p.logger.Printf("Failed to marshal email event: %v", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(event.Message.ID),
		Value: msgBytes,
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(event.Type)},
			{Key: "priority", Value: []byte(event.Message.Priority)},
			{Key: "timestamp", Value: []byte(event.Message.Timestamp.Format(time.RFC3339))},
		},
	}

	err = p.writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		p.logger.Printf("Failed to write message to Kafka: %v", err)
		return fmt.Errorf("failed to write to kafka: %w", err)
	}

	p.logger.Printf("Email event sent to Kafka: ID=%s, Type=%s, To=%s",
		event.Message.ID, event.Type, event.Message.To)

	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// determineEmailType determines the email type based on the template name
func determineEmailType(template string) string {
	switch template {
	case "password-reset":
		return TypePasswordReset
	case "welcome":
		return TypeWelcome
	case "order-confirmation":
		return TypeOrderConfirm
	default:
		return TypeNotification
	}
}
