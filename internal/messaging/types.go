package messaging

import "time"

// EmailMessage represents an email to be sent via Kafka
type EmailMessage struct {
	ID         string                 `json:"id"`
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	Subject    string                 `json:"subject"`
	Template   string                 `json:"template"`
	Data       map[string]interface{} `json:"data"`
	Priority   string                 `json:"priority"` // high, normal, low
	Timestamp  time.Time              `json:"timestamp"`
	RetryCount int                    `json:"retry_count"`
	MaxRetries int                    `json:"max_retries"`
	Status     string                 `json:"status"` // pending, sent, failed
	ErrorMsg   string                 `json:"error_msg,omitempty"`
}

// EmailEvent represents different email event types
type EmailEvent struct {
	Type    string       `json:"type"` // password_reset, welcome, notification, etc.
	Message EmailMessage `json:"message"`
}

const (
	// Topic names
	TopicEmailQueue = "email-queue"
	TopicEmailDLQ   = "email-dlq" // Dead Letter Queue for failed emails

	// Priority levels
	PriorityHigh   = "high"
	PriorityNormal = "normal"
	PriorityLow    = "low"

	// Status
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"

	// Email Types
	TypePasswordReset = "password_reset"
	TypeWelcome       = "welcome"
	TypeNotification  = "notification"
	TypeOrderConfirm  = "order_confirmation"
)
