# Messaging System Architecture

## Overview

The Usual Store application uses a **Kafka-based asynchronous messaging system** for sending email notifications. This architecture ensures reliable, scalable, and non-blocking email delivery for user creation, password resets, and other notifications.

## Architecture Diagram

```
Frontend Creates User
         ↓
Backend API (port 4001) → User saved to database ✓
         ↓
Frontend Sends Email Request
         ↓
Backend API (port 4001) → Publishes to Kafka ✓
         ↓
Kafka Broker (port 9092)
         ↓
Messaging Service → Consumes from Kafka ✓
         ↓
SMTP (Mailtrap) → Email delivered ✓
```

## Components

### 1. Frontend (Support Dashboard)
- **Port:** 3005
- **Technology:** React, Material-UI
- **Responsibility:** 
  - User interface for creating users
  - Calls backend API to create users
  - Calls backend API to send welcome emails

### 2. Backend API
- **Port:** 4001
- **Technology:** Go, Chi Router
- **Endpoints:**
  - `POST /api/users` - Creates new users in database
  - `POST /api/messaging/send` - Publishes email messages to Kafka
- **Responsibility:**
  - User creation and management
  - Kafka message publishing
  - Database operations

### 3. Kafka Broker
- **Port:** 9092
- **Technology:** Confluent Kafka
- **Topic:** `email-notifications`
- **Responsibility:**
  - Message queue for email notifications
  - Ensures reliable message delivery
  - Supports message persistence and replay

### 4. Zookeeper
- **Port:** 2181
- **Technology:** Apache Zookeeper
- **Responsibility:**
  - Kafka cluster coordination
  - Configuration management
  - Leader election for Kafka partitions

### 5. Messaging Service (Consumer)
- **Port:** 6001
- **Technology:** Go, Kafka Consumer
- **Responsibility:**
  - Consumes messages from Kafka topic
  - Processes email messages
  - Sends emails via SMTP

### 6. SMTP Server (Mailtrap)
- **Port:** 2525 (configurable)
- **Technology:** SMTP (Mailtrap for testing)
- **Responsibility:**
  - Final email delivery
  - Email inbox for testing

## Detailed Message Flow

### Step 1: User Creation
```
Frontend → POST /api/users → Backend API
```
- User fills out the "Create User" form
- Frontend sends HTTP POST request to `/api/users`
- Backend validates input
- Backend creates user in PostgreSQL database
- Backend returns user details (ID, email, role, etc.)

**Request Example:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Response Example:**
```json
{
  "error": false,
  "message": "User created successfully",
  "id": 10,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "role": "user"
}
```

### Step 2: Email Request
```
Frontend → POST /api/messaging/send → Backend API
```
- After user creation succeeds, frontend prepares welcome email
- Email content includes:
  - User's name
  - Login credentials
  - Role-specific login URL
  - Welcome message
- Frontend sends HTTP POST request to `/api/messaging/send`

**Request Example:**
```json
{
  "to": "john.doe@example.com",
  "subject": "Welcome to Usual Store - Your User Account",
  "message": "Hello John Doe,\n\nYour account has been successfully created!\n\nAccount Details:\nRole: User\nLogin URL: http://localhost:3000\nEmail: john.doe@example.com\nPassword: securepassword123\n\nBest regards,\nThe Usual Store Team"
}
```

### Step 3: Kafka Publishing
```
Backend API → Kafka Producer → Kafka Broker
```
- Backend API receives email request
- Creates Kafka producer
- Constructs EmailMessage struct:
  ```go
  type EmailMessage struct {
      ID         string
      From       string
      To         string
      Subject    string
      Template   string
      Data       map[string]interface{}
      Priority   string
      Timestamp  time.Time
      RetryCount int
      MaxRetries int
      Status     string
  }
  ```
- Publishes message to `email-notifications` topic
- Returns success response to frontend

**Response Example:**
```json
{
  "success": true,
  "message": "Message queued for delivery"
}
```

**Kafka Message Example:**
```json
{
  "type": "notification",
  "message": {
    "id": "b7ae53a1-0942-40ef-a76a-250cdb340327",
    "from": "noreply@usualstore.com",
    "to": "john.doe@example.com",
    "subject": "Welcome to Usual Store - Your User Account",
    "template": "plain",
    "data": {
      "body": "Hello John Doe,\n\nYour account has been successfully created!..."
    },
    "priority": "normal",
    "timestamp": "2025-12-26T13:25:35Z",
    "retry_count": 0,
    "max_retries": 3,
    "status": "pending"
  }
}
```

### Step 4: Kafka Consumption
```
Kafka Broker → Messaging Service Consumer
```
- Messaging service continuously polls Kafka for new messages
- Consumer group: `messaging-service-group`
- When new message arrives:
  - Deserializes JSON message
  - Validates message structure
  - Checks retry count
  - Processes email

**Consumer Logs:**
```
INFO Processing message: offset=1, partition=0
```

### Step 5: SMTP Email Sending
```
Messaging Service → SMTP Server → Email Delivered
```
- Messaging service retrieves SMTP configuration
- Renders email templates (HTML and plain text)
- Connects to SMTP server (Mailtrap)
- Sends email with:
  - From: `noreply@usualstore.com`
  - To: User's email address
  - Subject: From message
  - Body: Rendered from template
- Logs success/failure

**Consumer Logs:**
```
INFO Email sent successfully: From=noreply@usualstore.com, To=john.doe@example.com, Subject=Welcome to Usual Store - Your User Account
INFO Email sent successfully: ID=b7ae53a1-0942-40ef-a76a-250cdb340327, To=john.doe@example.com, Type=notification
```

## Configuration

### Environment Variables

#### Backend API (`.env`)
```bash
# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=email-notifications

# SMTP Configuration (for direct sending, optional)
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=your_mailtrap_username
SMTP_PASSWORD=your_mailtrap_password
SMTP_USER=noreply@usualstore.com
```

#### Messaging Service (`.env`)
```bash
# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=email-notifications
KAFKA_GROUP_ID=messaging-service-group

# SMTP Configuration (required)
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=your_mailtrap_username
SMTP_PASSWORD=your_mailtrap_password
SMTP_USER=noreply@usualstore.com
```

#### Frontend (docker-compose.yml or .env)
```bash
# Backend API URL
REACT_APP_SUPPORT_API_URL=http://localhost:4001

# Messaging API URL (same as backend)
REACT_APP_MESSAGING_API_URL=http://localhost:4001
```

### Docker Compose Services

```yaml
services:
  # Zookeeper (Kafka dependency)
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  # Kafka Broker
  kafka:
    image: confluentinc/cp-kafka:7.5.0
    ports:
      - "9092:9092"
      - "9093:9093"
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:9093
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"

  # Backend API
  back-end:
    build: .
    ports:
      - "4001:4001"
    environment:
      KAFKA_BROKERS: kafka:9092
      KAFKA_TOPIC: email-notifications

  # Messaging Service
  messaging-service:
    build:
      context: .
      dockerfile: cmd/messaging-service/Dockerfile
    ports:
      - "6001:6001"
    env_file:
      - .env
    environment:
      KAFKA_BROKERS: kafka:9092
      KAFKA_TOPIC: email-notifications
```

## Email Templates

The messaging service uses Go templates for email rendering.

### Template Structure

```
cmd/messaging-service/templates/
├── plain.html.tmpl       # HTML version of plain text emails
├── plain.plain.tmpl      # Plain text version
├── password-reset.html.tmpl
└── password-reset.plain.tmpl
```

### Plain Text Template

**File:** `plain.html.tmpl`
```html
{{define "body"}}
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }
    </style>
</head>
<body>
    <div style="white-space: pre-wrap;">{{.body}}</div>
</body>
</html>
{{end}}
```

**File:** `plain.plain.tmpl`
```
{{define "body"}}{{.body}}{{end}}
```

## API Reference

### Create User

**Endpoint:** `POST /api/users`

**Request:**
```json
{
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "password": "string",
  "role": "string"  // user, supporter, admin, super_admin
}
```

**Response:**
```json
{
  "error": false,
  "message": "User created successfully",
  "id": 10,
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "role": "string"
}
```

### Send Email via Kafka

**Endpoint:** `POST /api/messaging/send`

**Request:**
```json
{
  "to": "string",
  "subject": "string",
  "message": "string"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message queued for delivery"
}
```

## Error Handling

### Non-Blocking Email Sending

Email sending is **non-blocking**, meaning:
- ✅ User creation **always succeeds** even if email fails
- ⚠️ Failed emails are logged but don't prevent user creation
- 🔄 Kafka provides retry mechanism for transient failures

### Frontend Error Handling

```javascript
try {
  // 1. Create user
  const userResponse = await axios.post(`${API_URL}/api/users`, userData);
  
  // 2. Send welcome email (non-blocking)
  try {
    await axios.post(`${MESSAGING_API_URL}/api/messaging/send`, emailData);
    console.log('✅ Welcome email sent successfully via Kafka');
  } catch (emailError) {
    console.error('⚠️ Failed to send welcome email:', emailError);
    // Don't fail the whole process if email fails
  }
  
  // 3. Continue with success
  setSuccess(true);
} catch (err) {
  setError(err.message);
}
```

### Retry Logic

The messaging service implements automatic retry for failed emails:
- **Max Retries:** 3
- **Retry Strategy:** Exponential backoff
- **Dead Letter Queue:** Failed messages moved to `email-dlq` topic after max retries

## Monitoring & Debugging

### Check Service Health

```bash
# Check all services
docker-compose ps

# Check Kafka health
docker exec usual_store-kafka kafka-broker-api-versions \
  --bootstrap-server localhost:9092

# Check Zookeeper health
docker exec usual_store-zookeeper nc -z localhost 2181
```

### View Logs

```bash
# Backend API logs
docker logs -f usual_store-back-end-1

# Messaging service logs
docker logs -f usual_store-messaging-service-1

# Kafka logs
docker logs -f usual_store-kafka
```

### Monitor Kafka Messages

```bash
# View all messages in topic
docker exec -it usual_store-kafka kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic email-notifications \
  --from-beginning

# View consumer group status
docker exec -it usual_store-kafka kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group messaging-service-group \
  --describe
```

### Kafka UI

Access Kafka UI at: http://localhost:8090
- View topics
- Monitor consumer groups
- Inspect messages
- Check cluster health

## Testing

### Manual Testing via UI

1. **Open Support Dashboard**
   ```
   http://localhost:3005/support/login
   ```

2. **Login**
   ```
   Email: admin@example.com
   Password: qwerty12
   ```

3. **Create User**
   - Click "Manage Users"
   - Click "Create User"
   - Fill in user details
   - Click "Create User"

4. **Verify Email**
   - Check Mailtrap inbox: https://mailtrap.io
   - Verify email received with correct content

### Testing via API

```bash
# 1. Create user
curl -X POST http://localhost:4001/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "password": "testpass123",
    "role": "user"
  }'

# 2. Send welcome email
curl -X POST http://localhost:4001/api/messaging/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Welcome to Usual Store",
    "message": "Hello! Your account has been created."
  }'

# 3. Check messaging service logs
docker logs usual_store-messaging-service-1 | tail -20
```

### Automated Testing

```bash
# Run full test suite
make test

# Run messaging tests only
go test ./internal/messaging/... -v

# Run with coverage
go test ./internal/messaging/... -cover
```

## Troubleshooting

### Issue: Email not sent

**Symptoms:** User created but no email received

**Possible Causes:**
1. Messaging service not running
2. Kafka not running
3. SMTP credentials incorrect
4. Email consumer not consuming

**Solution:**
```bash
# Check messaging service
docker logs usual_store-messaging-service-1

# Check Kafka connectivity
docker exec usual_store-messaging-service-1 nc -zv kafka 9092

# Check SMTP credentials in .env
cat .env | grep SMTP

# Restart messaging service
docker-compose restart messaging-service
```

### Issue: 404 on /api/messaging/send

**Symptoms:** `404 page not found` error

**Solution:**
```bash
# Rebuild backend API
docker-compose build --no-cache back-end

# Restart backend
docker-compose restart back-end

# Verify route is registered
docker logs usual_store-back-end-1 | grep "Starting Back end"
```

### Issue: Kafka connection refused

**Symptoms:** `failed to connect to kafka` error

**Solution:**
```bash
# Check Kafka is running
docker ps | grep kafka

# Restart Kafka and Zookeeper
docker-compose restart zookeeper kafka

# Wait 30 seconds for Kafka to start
sleep 30

# Check Kafka health
docker exec usual_store-kafka kafka-broker-api-versions \
  --bootstrap-server localhost:9092
```

### Issue: SMTP authentication error

**Symptoms:** `535 5.7.0 Invalid credentials` error

**Solution:**
1. Go to https://mailtrap.io
2. Sign in to your account
3. Navigate to "Email Testing" → "Inboxes" → "My Inbox"
4. Copy SMTP credentials
5. Update `.env` file:
   ```bash
   SMTP_HOST=sandbox.smtp.mailtrap.io
   SMTP_PORT=2525
   SMTP_USERNAME=your_actual_username
   SMTP_PASSWORD=your_actual_password
   SMTP_USER=noreply@usualstore.com
   ```
6. Restart messaging service:
   ```bash
   docker-compose restart messaging-service
   ```

## Performance Considerations

### Scalability

- **Kafka Partitions:** Increase partitions for higher throughput
- **Consumer Groups:** Add more messaging service instances
- **Batch Processing:** Process multiple emails in batches

### Optimization

- **Connection Pooling:** Reuse SMTP connections
- **Async Processing:** Non-blocking email sending
- **Message Compression:** Kafka message compression (Snappy)

### Monitoring

- **Kafka Lag:** Monitor consumer lag
- **Email Send Time:** Track SMTP response time
- **Retry Rate:** Monitor failed email retry rate
- **Queue Depth:** Track Kafka queue depth

## Security

### Message Security

- **TLS/SSL:** Enable encryption for SMTP
- **Authentication:** Use strong SMTP credentials
- **Message Signing:** Sign messages for authenticity

### Access Control

- **Kafka ACLs:** Restrict topic access
- **API Authentication:** Require authentication for messaging endpoint
- **Network Isolation:** Use Docker networks for service isolation

## Production Checklist

- [ ] Update SMTP server from Mailtrap to production server
- [ ] Enable TLS/SSL for SMTP connections
- [ ] Configure Kafka retention policy
- [ ] Set up monitoring and alerting
- [ ] Implement email rate limiting
- [ ] Add email delivery status tracking
- [ ] Configure dead letter queue handling
- [ ] Set up log aggregation
- [ ] Enable Kafka replication
- [ ] Configure backup and disaster recovery

## Related Documentation

- [User Creation Email Notifications](./USER-CREATION-EMAIL-NOTIFICATIONS.md)
- [Kafka Testing Guide](../../KAFKA-TESTING-GUIDE.md)
- [Email Quick Reference](./EMAIL-QUICK-REFERENCE.md)
- [Support Dashboard](./SUPPORT-DASHBOARD.md)

## Summary

The messaging system provides:
- ✅ Reliable email delivery via Kafka
- ✅ Non-blocking user creation
- ✅ Automatic retry mechanism
- ✅ Scalable architecture
- ✅ Easy monitoring and debugging
- ✅ Role-specific email content
- ✅ Production-ready design

**All components work together to ensure users receive welcome emails immediately after account creation!**

