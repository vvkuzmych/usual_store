# 🧪 Kafka Messaging - Testing & Verification Guide

## Quick Start Testing

### Step 1: Start Kafka Infrastructure

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Start Kafka and Zookeeper
docker-compose -f docker-compose-kafka.yml up -d

# Wait 30 seconds for Kafka to initialize
sleep 30

# Verify services are running
docker ps | grep -E "(zookeeper|kafka)"
```

**Expected Output:**
```
usual_store-kafka       Up 30 seconds
usual_store-zookeeper   Up 30 seconds
```

### Step 2: Verify Kafka is Healthy

```bash
# Check Zookeeper
docker exec usual_store-zookeeper-1 nc -z localhost 2181 && echo "✅ Zookeeper OK" || echo "❌ Zookeeper FAILED"

# Check Kafka broker
docker exec usual_store-kafka-1 kafka-broker-api-versions --bootstrap-server localhost:9092 >/dev/null 2>&1 && echo "✅ Kafka OK" || echo "❌ Kafka FAILED"

# List Kafka topics (should be empty initially)
docker exec usual_store-kafka-1 kafka-topics --list --bootstrap-server localhost:9092
```

### Step 3: Build Messaging Service

```bash
# Build the messaging service Docker image
docker build -f cmd/messaging-service/Dockerfile -t usual_store-messaging-service .

# OR using docker-compose
docker-compose build messaging-service
```

**Expected Output:**
```
Successfully built <image-id>
Successfully tagged usual_store-messaging-service:latest
```

### Step 4: Start Messaging Service

```bash
# Start messaging service
docker run -d \
  --name messaging-service-test \
  --network usual_store_usualstore_network \
  -e KAFKA_BROKERS=kafka:9092 \
  -e KAFKA_TOPIC=email-queue \
  -e KAFKA_GROUP_ID=messaging-service-group \
  -e SMTP_HOST=${SMTP_HOST} \
  -e SMTP_PORT=${SMTP_PORT} \
  -e SMTP_USER=${SMTP_USER} \
  -e SMTP_PASSWORD=${SMTP_PASSWORD} \
  usual_store-messaging-service

# Check if it started
docker ps | grep messaging-service-test
```

### Step 5: Check Messaging Service Logs

```bash
# View logs in real-time
docker logs -f messaging-service-test

# You should see:
# INFO    Starting Messaging Service v1.0.0
# INFO    Environment: development
# INFO    Kafka Brokers: [kafka:9092]
# INFO    Kafka Topic: email-queue
# INFO    Starting Kafka consumer...
# INFO    Messaging service is running. Press Ctrl+C to stop.
```

### Step 6: Test Kafka Producer (Manual)

**Option A: Using Go Code**

Create a test file `test-kafka-producer.go`:

```bash
cat > test-kafka-producer.go << 'EOF'
package main

import (
    "context"
    "log"
    "time"
    "usual_store/internal/messaging"
)

func main() {
    logger := log.New(log.Writer(), "TEST\t", log.Ldate|log.Ltime)
    
    producer := messaging.NewProducer(
        []string{"localhost:9093"},
        "email-queue",
        logger,
    )
    defer producer.Close()
    
    data := map[string]interface{}{
        "Link": "https://example.com/reset-password?token=test123",
    }
    
    err := producer.SendEmail(
        context.Background(),
        "test@usualstore.com",
        "user@example.com",
        "Test Email",
        "password-reset",
        data,
        messaging.PriorityHigh,
    )
    
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }
    
    log.Println("✅ Email event sent to Kafka successfully!")
    time.Sleep(2 * time.Second) // Wait for message to be processed
}
EOF

# Run the test
go run test-kafka-producer.go

# Clean up
rm test-kafka-producer.go
```

**Option B: Using kafka-console-producer**

```bash
# Send a test message directly to Kafka
docker exec -it usual_store-kafka-1 kafka-console-producer \
  --bootstrap-server localhost:9092 \
  --topic email-queue << EOF
{
  "type": "password_reset",
  "message": {
    "id": "test-123",
    "from": "test@usualstore.com",
    "to": "user@example.com",
    "subject": "Test Email",
    "template": "password-reset",
    "data": {"Link": "https://example.com/test"},
    "priority": "high",
    "timestamp": "2025-12-25T21:00:00Z",
    "retry_count": 0,
    "max_retries": 3,
    "status": "pending"
  }
}
EOF
```

### Step 7: Verify Message Processing

```bash
# Check messaging service logs
docker logs messaging-service-test | tail -20

# Look for:
# INFO    Processing message: offset=0, partition=0
# INFO    Email sent successfully: ID=test-123, To=user@example.com, Type=password_reset
```

### Step 8: Monitor with Kafka UI (Optional)

```bash
# Start Kafka UI
docker run -d \
  --name kafka-ui-test \
  --network usual_store_usualstore_network \
  -p 8090:8080 \
  -e KAFKA_CLUSTERS_0_NAME=usual-store-cluster \
  -e KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092 \
  -e KAFKA_CLUSTERS_0_ZOOKEEPER=zookeeper:2181 \
  provectuslabs/kafka-ui:latest

# Access Kafka UI
open http://localhost:8090
```

**In Kafka UI, check:**
- Topics → `email-queue` → Should have messages
- Consumers → `messaging-service-group` → Should show consumer lag
- Messages → View actual message content

## Full End-to-End Test

### Prerequisites

1. **Update Main Application** (if not done yet)

Add to `cmd/api/handlers-api.go` (after line 633):

```go
// Test if Kafka producer is available
if app.messagingProducer != nil {
    ctx := r.Context()
    err = app.messagingProducer.SendEmail(
        ctx,
        "noreply@usualstore.com",
        payload.Email,
        "Password Reset Request",
        "password-reset",
        map[string]interface{}{"Link": signedLink},
        messaging.PriorityHigh,
    )
    if err != nil {
        app.errorLog.Printf("Kafka failed, using direct email: %v", err)
        // Fallback to direct email
        err = app.SendEmail("info@usualstore.com", payload.Email, "Password Reset Request", "password-reset", data)
        if err != nil {
            app.errorLog.Println(err)
            return
        }
    }
} else {
    // Direct email (Kafka disabled)
    err = app.SendEmail("info@usualstore.com", payload.Email, "Password Reset Request", "password-reset", data)
    if err != nil {
        app.errorLog.Println(err)
        return
    }
}
```

2. **Add Kafka initialization** in `cmd/api/api.go` main():

```go
// After line 157, add:
import "strings"

// Load Kafka config
kafkaEnabled := os.Getenv("KAFKA_ENABLED") == "true"
if kafkaEnabled {
    kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9093"), ",")
    cfg.kafka.enabled = true
    cfg.kafka.brokers = kafkaBrokers
    cfg.kafka.topic = getEnv("KAFKA_TOPIC", "email-queue")
}

// After initializing app, add:
var messagingProducer *messaging.Producer
if cfg.kafka.enabled {
    messagingProducer = messaging.NewProducer(cfg.kafka.brokers, cfg.kafka.topic, infoLog)
    defer messagingProducer.Close()
    infoLog.Println("Kafka producer initialized")
}

// Update app initialization:
app := &application{
    // ... existing fields ...
    messagingProducer: messagingProducer,
    // ... rest of fields ...
}
```

3. **Update .env file:**

```bash
# Add to .env
KAFKA_ENABLED=true
KAFKA_BROKERS=localhost:9093
KAFKA_TOPIC=email-queue
```

### E2E Test Steps

```bash
# 1. Ensure Kafka and messaging service are running
docker ps | grep -E "(kafka|zookeeper|messaging)"

# 2. Rebuild and restart backend
docker-compose build back-end
docker-compose up -d back-end

# 3. Wait for backend to start
sleep 5

# 4. Check backend logs for Kafka initialization
docker logs usual_store-back-end-1 | grep -i kafka

# Expected: "Kafka producer initialized"

# 5. Send password reset request
curl -X POST http://localhost:4001/api/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com"}'

# Expected response:
# {"error":false,"message":""}

# 6. Check messaging service logs
docker logs messaging-service-test | tail -10

# Expected:
# INFO    Email event sent to Kafka: ID=<uuid>, Type=password_reset, To=admin@example.com
# INFO    Processing message: offset=0, partition=0
# INFO    Email sent successfully: ID=<uuid>, To=admin@example.com
```

## Verification Checklist

### ✅ Infrastructure Health

```bash
# Check all services
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Expected services running:
# ✅ usual_store-zookeeper-1
# ✅ usual_store-kafka-1
# ✅ messaging-service-test (or usual_store-messaging-service-1)
# ✅ usual_store-back-end-1
# ✅ usual_store-database-1
```

### ✅ Kafka Topics

```bash
# List topics
docker exec usual_store-kafka-1 kafka-topics --list --bootstrap-server localhost:9092

# Expected output:
# email-queue

# Get topic details
docker exec usual_store-kafka-1 kafka-topics \
  --describe \
  --topic email-queue \
  --bootstrap-server localhost:9092
```

### ✅ Consumer Groups

```bash
# List consumer groups
docker exec usual_store-kafka-1 kafka-consumer-groups \
  --list \
  --bootstrap-server localhost:9092

# Expected:
# messaging-service-group

# Check consumer lag
docker exec usual_store-kafka-1 kafka-consumer-groups \
  --describe \
  --group messaging-service-group \
  --bootstrap-server localhost:9092

# Expected:
# GROUP                     TOPIC        PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG
# messaging-service-group   email-queue  0          5               5               0
```

### ✅ Message Flow

```bash
# Check how many messages were sent
docker exec usual_store-kafka-1 kafka-run-class kafka.tools.GetOffsetShell \
  --broker-list localhost:9092 \
  --topic email-queue

# Check messaging service processed messages
docker logs messaging-service-test | grep "Email sent successfully" | wc -l

# Check for any errors
docker logs messaging-service-test | grep -i error
```

## Troubleshooting

### Issue: Kafka not starting

```bash
# Check Zookeeper logs
docker logs usual_store-zookeeper-1

# Check Kafka logs
docker logs usual_store-kafka-1

# Restart Kafka infrastructure
docker-compose -f docker-compose-kafka.yml down
docker-compose -f docker-compose-kafka.yml up -d
```

### Issue: Messaging service can't connect to Kafka

```bash
# Check network
docker network inspect usual_store_usualstore_network | grep -A 5 messaging

# Check if Kafka is accessible from messaging service
docker exec messaging-service-test nc -zv kafka 9092

# Check environment variables
docker exec messaging-service-test env | grep KAFKA
```

### Issue: Messages not being consumed

```bash
# Check if consumer is running
docker logs messaging-service-test | grep "Starting Kafka consumer"

# Check consumer group
docker exec usual_store-kafka-1 kafka-consumer-groups \
  --describe \
  --group messaging-service-group \
  --bootstrap-server localhost:9092

# Manually consume messages to see if they're there
docker exec usual_store-kafka-1 kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic email-queue \
  --from-beginning \
  --max-messages 5
```

### Issue: Backend not sending to Kafka

```bash
# Check backend logs
docker logs usual_store-back-end-1 | grep -i kafka

# Check if Kafka is enabled
docker exec usual_store-back-end-1 env | grep KAFKA_ENABLED

# Test Kafka connection from backend
docker exec usual_store-back-end-1 nc -zv kafka 9092
```

## Performance Testing

### Load Test

```bash
# Install hey (HTTP load testing tool)
# macOS: brew install hey
# Linux: go install github.com/rakyll/hey@latest

# Send 100 concurrent password reset requests
hey -n 100 -c 10 \
  -m POST \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}' \
  http://localhost:4001/api/forgot-password

# Check results:
# - All requests should return 202 (Accepted)
# - Average response time should be <50ms
# - No errors
```

### Monitor Consumer Lag During Load

```bash
# Run this while load testing
watch -n 1 'docker exec usual_store-kafka-1 kafka-consumer-groups \
  --describe \
  --group messaging-service-group \
  --bootstrap-server localhost:9092'

# Consumer lag should stay low (< 10)
# If lag is high, scale up messaging service instances
```

## Clean Up Test Resources

```bash
# Stop and remove test containers
docker stop messaging-service-test kafka-ui-test
docker rm messaging-service-test kafka-ui-test

# Remove test files
rm -f test-kafka-producer.go

# Optional: Stop Kafka infrastructure
docker-compose -f docker-compose-kafka.yml down

# Optional: Remove all data
docker-compose -f docker-compose-kafka.yml down -v
```

## Success Criteria

Your Kafka messaging system is working correctly if:

✅ All services are running (Kafka, Zookeeper, Messaging Service, Backend)  
✅ Kafka topics are created automatically  
✅ Messages are produced to `email-queue` topic  
✅ Messages are consumed by messaging service  
✅ Emails are sent via SMTP  
✅ Consumer lag is low (< 10)  
✅ No errors in any service logs  
✅ API response time is < 50ms  
✅ End-to-end email flow completes in < 5 seconds  

## Monitoring Commands (Quick Reference)

```bash
# Services status
docker ps --format "table {{.Names}}\t{{.Status}}"

# Kafka topics
docker exec usual_store-kafka-1 kafka-topics --list --bootstrap-server localhost:9092

# Consumer lag
docker exec usual_store-kafka-1 kafka-consumer-groups --describe --group messaging-service-group --bootstrap-server localhost:9092

# Messaging service logs
docker logs -f messaging-service-test

# Backend logs
docker logs -f usual_store-back-end-1

# Message count
docker exec usual_store-kafka-1 kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic email-queue

# Kafka UI
open http://localhost:8090
```

---

**Happy Testing! 🧪**

If you encounter any issues, check the troubleshooting section or review logs for detailed error messages.

