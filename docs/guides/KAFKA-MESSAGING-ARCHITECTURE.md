# 📨 Kafka-Based Messaging Architecture

## Overview

This document describes the **Kafka-based asynchronous messaging architecture** that separates email sending logic into a dedicated microservice.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────┐
│                          Main Application                             │
│  ┌────────────────────────────────────────────────────────────────┐  │
│  │  API Handlers (cmd/api/handlers-api.go)                        │  │
│  │  - SendPasswordResetEmail()                                    │  │
│  │  - OrderConfirmation()                                         │  │
│  │  - WelcomeEmail()                                             │  │
│  └─────────────────────────────┬──────────────────────────────────┘  │
│                                │                                      │
│  ┌─────────────────────────────▼──────────────────────────────────┐  │
│  │  Kafka Producer (internal/messaging/producer.go)              │  │
│  │  - Publishes email events to Kafka                            │  │
│  │  - Non-blocking, async operation                              │  │
│  └─────────────────────────────┬──────────────────────────────────┘  │
└────────────────────────────────┼───────────────────────────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │    Kafka Cluster       │
                    │  Topic: email-queue    │
                    └────────────┬───────────┘
                                 │
                                 ▼
┌────────────────────────────────┴───────────────────────────────────────┐
│                      Messaging Microservice                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Kafka Consumer (internal/messaging/consumer.go)               │   │
│  │  - Consumes email events from Kafka                            │   │
│  │  - Processes messages asynchronously                           │   │
│  └─────────────────────────────┬───────────────────────────────────┘   │
│                                 │                                       │
│  ┌─────────────────────────────▼───────────────────────────────────┐   │
│  │  Email Handler (cmd/messaging-service/email_handler.go)        │   │
│  │  - Renders HTML/Plain text templates                           │   │
│  │  - Sends emails via SMTP                                       │   │
│  │  - Handles retries and error handling                          │   │
│  └─────────────────────────────┬───────────────────────────────────┘   │
│                                 │                                       │
│                                 ▼                                       │
│                         ┌──────────────┐                               │
│                         │ SMTP Server  │                               │
│                         │ (Mailtrap)   │                               │
│                         └──────────────┘                               │
└────────────────────────────────────────────────────────────────────────┘
```

## Components

### 1. **Main Application** (`cmd/api`)

**Responsibilities:**
- Handle HTTP requests
- Validate business logic
- Publish email events to Kafka
- Return immediate response to users

**Key Files:**
- `cmd/api/api.go` - Application config with Kafka producer
- `internal/messaging/producer.go` - Kafka message publisher

**Usage Example:**
```go
// In handlers-api.go
func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
    // ... validation logic ...
    
    // Send email asynchronously via Kafka
    err := app.messagingProducer.SendEmail(
        ctx,
        "noreply@usualstore.com",
        payload.Email,
        "Password Reset Request",
        "password-reset",
        map[string]interface{}{
            "Link": signedLink,
        },
        messaging.PriorityHigh,
    )
    
    // Respond immediately (don't wait for email to be sent)
    app.writeJSON(w, http.StatusAccepted, map[string]bool{"sent": true})
}
```

### 2. **Messaging Microservice** (`cmd/messaging-service`)

**Responsibilities:**
- Consume email events from Kafka
- Render email templates (HTML + Plain text)
- Send emails via SMTP
- Handle retries for failed emails
- Log all email operations

**Key Files:**
- `cmd/messaging-service/main.go` - Service entry point
- `cmd/messaging-service/email_handler.go` - SMTP email sender
- `internal/messaging/consumer.go` - Kafka message consumer

**Features:**
- ✅ Async email processing
- ✅ Automatic retries (3 attempts)
- ✅ Template rendering
- ✅ Priority-based processing
- ✅ Dead Letter Queue for failed messages
- ✅ Graceful shutdown

### 3. **Kafka Infrastructure**

**Components:**
- **Zookeeper** - Kafka cluster coordination
- **Kafka Broker** - Message broker
- **Kafka UI** - Web interface for monitoring

**Topics:**
- `email-queue` - Main email queue
- `email-dlq` - Dead Letter Queue for failed emails

## Message Format

### EmailEvent Structure

```go
type EmailEvent struct {
    Type    string       `json:"type"` // password_reset, welcome, etc.
    Message EmailMessage `json:"message"`
}

type EmailMessage struct {
    ID          string                 `json:"id"`
    From        string                 `json:"from"`
    To          string                 `json:"to"`
    Subject     string                 `json:"subject"`
    Template    string                 `json:"template"`
    Data        map[string]interface{} `json:"data"`
    Priority    string                 `json:"priority"` // high, normal, low
    Timestamp   time.Time              `json:"timestamp"`
    RetryCount  int                    `json:"retry_count"`
    MaxRetries  int                    `json:"max_retries"`
    Status      string                 `json:"status"` // pending, sent, failed
}
```

### Example Kafka Message

```json
{
  "type": "password_reset",
  "message": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "from": "noreply@usualstore.com",
    "to": "user@example.com",
    "subject": "Password Reset Request",
    "template": "password-reset",
    "data": {
      "Link": "https://usualstore.com/reset?token=abc123"
    },
    "priority": "high",
    "timestamp": "2025-12-25T20:00:00Z",
    "retry_count": 0,
    "max_retries": 3,
    "status": "pending"
  }
}
```

## Configuration

### Environment Variables

#### Main Application
```bash
# Kafka Configuration
KAFKA_ENABLED=true
KAFKA_BROKERS=localhost:9093
KAFKA_TOPIC=email-queue

# Existing SMTP config (for fallback)
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your_user
SMTP_PASSWORD=your_password
```

#### Messaging Service
```bash
# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=email-queue
KAFKA_GROUP_ID=messaging-service-group

# SMTP Configuration
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your_user
SMTP_PASSWORD=your_password
```

## Deployment

### Docker Compose

#### 1. Start Kafka Infrastructure
```bash
docker-compose -f docker-compose-kafka.yml up -d
```

#### 2. Start Main Application
```bash
docker-compose up -d back-end
```

#### 3. Start Messaging Service
```bash
docker-compose up -d messaging-service
```

### Verify Services

```bash
# Check Kafka
docker-compose -f docker-compose-kafka.yml ps

# Check Messaging Service logs
docker logs -f usual_store-messaging-service-1

# Access Kafka UI
open http://localhost:8090
```

## Benefits

### 🚀 Performance
- **Non-blocking**: API responds immediately, email sent asynchronously
- **Scalable**: Add more messaging service instances for higher throughput
- **Resilient**: Failed emails automatically retried

### 🛡️ Reliability
- **Guaranteed Delivery**: Kafka ensures messages aren't lost
- **Retry Logic**: Automatic retries for transient failures
- **Dead Letter Queue**: Failed messages stored for analysis

### 📊 Observability
- **Kafka UI**: Monitor message flow in real-time
- **Logging**: Comprehensive logs for debugging
- **Metrics**: Track email success/failure rates

### 🔧 Maintainability
- **Separation of Concerns**: Email logic isolated from business logic
- **Easy Testing**: Mock Kafka producer in tests
- **Independent Deployment**: Deploy messaging service without affecting main app

## Monitoring

### Kafka UI (http://localhost:8090)
- View topics and partitions
- Monitor consumer lag
- Inspect messages
- Track throughput

### Service Logs
```bash
# Messaging Service
docker logs -f usual_store-messaging-service-1

# Main Application
docker logs -f usual_store-back-end-1
```

### Key Metrics to Monitor
- **Consumer Lag**: How far behind is the consumer?
- **Success Rate**: Percentage of emails sent successfully
- **Retry Rate**: How many messages need retries?
- **Average Processing Time**: Time to send an email

## Troubleshooting

### Issue: Messages not being consumed

**Solution:**
```bash
# Check if messaging service is running
docker ps | grep messaging-service

# Check messaging service logs
docker logs usual_store-messaging-service-1

# Verify Kafka connectivity
docker exec usual_store-kafka-1 kafka-topics --list --bootstrap-server localhost:9092
```

### Issue: Emails not being sent

**Solution:**
1. Check SMTP credentials in messaging service
2. Verify SMTP server is accessible
3. Check messaging service logs for errors
4. Inspect failed messages in Dead Letter Queue

### Issue: High consumer lag

**Solution:**
1. Scale up messaging service instances
2. Increase consumer batch size
3. Optimize SMTP connection pooling

## Migration from Direct SMTP

### Old Code (Direct SMTP)
```go
func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
    // ... validation ...
    
    // Blocking SMTP call
    err = app.SendEmail("info@usualstore.com", payload.Email, "Password Reset", "password-reset", data)
    if err != nil {
        app.errorLog.Println(err)
        return
    }
    
    app.writeJSON(w, http.StatusAccepted, resp)
}
```

### New Code (Kafka Async)
```go
func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
    // ... validation ...
    
    // Non-blocking Kafka publish
    err = app.messagingProducer.SendEmail(
        r.Context(),
        "info@usualstore.com",
        payload.Email,
        "Password Reset",
        "password-reset",
        data,
        messaging.PriorityHigh,
    )
    
    // Respond immediately
    app.writeJSON(w, http.StatusAccepted, resp)
}
```

## Future Enhancements

- [ ] Add email templates for order confirmations
- [ ] Implement priority queues
- [ ] Add email tracking (opens, clicks)
- [ ] Integrate with SendGrid/AWS SES
- [ ] Add email rate limiting
- [ ] Implement scheduled emails
- [ ] Add email analytics dashboard

## Related Documentation

- [Kafka Official Docs](https://kafka.apache.org/documentation/)
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
- [go-simple-mail](https://github.com/xhit/go-simple-mail)
- [Docker Compose](../setup/DOCKER-SETUP.md)
- [Kubernetes Deployment](./KUBERNETES-DEPLOYMENT.md)

---

**Created**: December 25, 2025
**Version**: 1.0.0
**Status**: ✅ Ready for Production

