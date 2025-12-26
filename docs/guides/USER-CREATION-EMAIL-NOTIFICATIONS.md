# User Creation Email Notifications

## Overview

The Usual Store application automatically sends welcome emails to all newly created users via the **Kafka-based messaging microservice**. This ensures reliable, asynchronous email delivery without blocking the user creation process.

## Architecture

```
┌─────────────────────┐
│  Support Frontend   │
│  (User Management)  │
└──────────┬──────────┘
           │ 1. Create User
           ▼
┌─────────────────────┐
│   Backend API       │
│  (port 4001)        │
└──────────┬──────────┘
           │ 2. User Created
           │    (Return user data)
           ▼
┌─────────────────────┐
│  Support Frontend   │
└──────────┬──────────┘
           │ 3. Send Welcome Email
           ▼
┌─────────────────────┐
│  Messaging Service  │──── 4. Publish to Kafka ───► ┌──────────────┐
│  (port 6001)        │                               │    Kafka     │
└─────────────────────┘                               └──────┬───────┘
                                                             │
                                                             ▼
                                               ┌──────────────────────┐
                                               │  Email Consumer      │
                                               │  (sends via SMTP)    │
                                               └──────────────────────┘
```

## Email Sending Locations

### 1. **UserManagement Component** (Super Admin)
**File:** `support-frontend/src/components/UserManagement.jsx`

**Who can use:** Super Admin only
**What it does:** Creates users with ANY role (user, supporter, admin, super_admin)

**Email Content:**
- Personalized greeting with first name and last name
- Role-specific information (Super Administrator, Administrator, Support Agent, User)
- Login URL (varies by role):
  - Regular users: `http://localhost:3000` (main store)
  - Staff (supporter/admin/super_admin): `http://localhost:3005/support/dashboard`
- Email address
- Temporary password
- Security reminder to change password

**Implementation:**
```javascript
await axios.post(`${MESSAGING_API_URL}/api/messaging/send`, {
  to: newUser.email,
  subject: `Welcome to Usual Store - Your ${roleDisplayName} Account`,
  message: `
    Hello ${newUser.firstName} ${newUser.lastName},
    
    Your account has been successfully created!
    
    Account Details:
    Role: ${roleDisplayName}
    Login URL: ${loginUrl}
    Email: ${newUser.email}
    Password: ${newUser.password}
    
    ...
  `.trim(),
});
```

### 2. **CreateSupporterAccount Component** (Admin)
**File:** `support-frontend/src/components/CreateSupporterAccount.jsx`

**Who can use:** Admin and Super Admin
**What it does:** Creates users with 'supporter' role only

**Email Content:**
- Personalized greeting
- Support dashboard login URL
- Email address
- Temporary password
- Instructions to change password after first login

## Email Templates by Role

### Super Admin
```
Subject: Welcome to Usual Store - Your Super Administrator Account

Hello [First Name] [Last Name],

Your account has been successfully created!

Account Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Role: Super Administrator
Login URL: http://localhost:3005/support/dashboard
Email: [email@example.com]
Password: [generated_password]

You can now access the support dashboard and manage tickets.

For security reasons, we recommend changing your password after your first login.

Best regards,
The Usual Store Team
```

### Admin
```
Subject: Welcome to Usual Store - Your Administrator Account

Hello [First Name] [Last Name],

Your account has been successfully created!

Account Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Role: Administrator
Login URL: http://localhost:3005/support/dashboard
Email: [email@example.com]
Password: [generated_password]

You can now access the support dashboard and manage tickets.

For security reasons, we recommend changing your password after your first login.

Best regards,
The Usual Store Team
```

### Supporter
```
Subject: Welcome to Usual Store - Your Support Agent Account

Hello [First Name] [Last Name],

Your account has been successfully created!

Account Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Role: Support Agent
Login URL: http://localhost:3005/support/dashboard
Email: [email@example.com]
Password: [generated_password]

You can now access the support dashboard and manage tickets.

For security reasons, we recommend changing your password after your first login.

Best regards,
The Usual Store Team
```

### Regular User
```
Subject: Welcome to Usual Store - Your User Account

Hello [First Name] [Last Name],

Your account has been successfully created!

Account Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Role: User
Login URL: http://localhost:3000
Email: [email@example.com]
Password: [generated_password]

You can now log in and start shopping!

For security reasons, we recommend changing your password after your first login.

Best regards,
The Usual Store Team
```

## Kafka Integration

### Configuration

The messaging service URL is configured via environment variables:

**Support Frontend (.env or docker-compose.yml):**
```bash
REACT_APP_MESSAGING_API_URL=http://localhost:6001
```

### Message Flow

1. **User Created:** Frontend calls backend API to create user
2. **User Data Returned:** Backend returns user ID and details
3. **Email Request:** Frontend sends email request to messaging service
4. **Kafka Publish:** Messaging service publishes to `email-notifications` topic
5. **Consumer Processes:** Email consumer reads from Kafka
6. **SMTP Send:** Consumer sends email via configured SMTP server (e.g., Mailtrap)

### Error Handling

Email sending is **non-blocking**:
- If email fails, user creation still succeeds
- Error is logged to console: `⚠️ Failed to send welcome email`
- Success is logged: `✅ Welcome email sent successfully via Kafka`

This ensures that email delivery issues don't prevent user account creation.

## Testing Email Delivery

### 1. **Check Mailtrap Inbox**

If using Mailtrap for testing:
1. Go to https://mailtrap.io
2. Sign in to your account
3. Navigate to "Email Testing" → "Inboxes" → "My Inbox"
4. Check for welcome emails

### 2. **Verify Kafka Messages**

Check if messages are being published to Kafka:

```bash
# Enter Kafka container
docker exec -it usual_store-kafka-1 /bin/bash

# List topics
kafka-topics.sh --list --bootstrap-server localhost:9092

# Read messages from email-notifications topic
kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic email-notifications \
  --from-beginning
```

### 3. **Check Messaging Service Logs**

```bash
# View messaging service logs
docker logs -f usual_store-messaging-service-1

# Look for:
# ✅ Message sent to email-notifications topic
# ✅ Email sent successfully
```

### 4. **Check Consumer Logs**

```bash
# View messaging consumer logs (if separate)
docker logs -f usual_store-messaging-consumer-1
```

## Manual Testing

### Test Case 1: Create Super Admin User
1. Login as super admin: `admin@example.com` / `qwerty12`
2. Navigate to http://localhost:3005/support/users
3. Click "Create User"
4. Fill in:
   - First Name: Test
   - Last Name: SuperAdmin
   - Email: test.superadmin@example.com
   - Password: testpass123
   - Role: Super Admin
5. Click "Create User"
6. ✅ Check Mailtrap for welcome email

### Test Case 2: Create Supporter
1. Login as admin: `admin@example.com` / `qwerty12`
2. Navigate to http://localhost:3005/support/dashboard
3. Click "Create Supporter"
4. Fill in user details
5. Click "Create Account"
6. ✅ Check Mailtrap for credentials email

### Test Case 3: Create Regular User
1. Login as super admin
2. Navigate to User Management
3. Create user with "User" role
4. ✅ Verify email contains correct login URL (localhost:3000)

## Environment Variables

### Required

```bash
# Backend API URL (for user creation)
REACT_APP_SUPPORT_API_URL=http://localhost:4001

# Messaging Service URL (for email sending)
REACT_APP_MESSAGING_API_URL=http://localhost:6001
```

### Messaging Service Configuration

```bash
# SMTP Settings (in backend .env)
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USERNAME=your_mailtrap_username
SMTP_PASSWORD=your_mailtrap_password
SMTP_USER=your_email@example.com

# Kafka Settings
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=email-notifications
```

## Troubleshooting

### Email Not Sent

**Symptom:** User created but no email received

**Possible Causes:**
1. Messaging service not running
2. Kafka not running
3. SMTP credentials incorrect
4. Email consumer not running

**Solution:**
```bash
# Check all services are running
docker-compose ps

# Check messaging service logs
docker logs usual_store-messaging-service-1

# Verify Kafka connectivity
docker exec usual_store-messaging-service-1 nc -zv kafka 9092

# Test SMTP connection
docker exec usual_store-messaging-service-1 nc -zv sandbox.smtp.mailtrap.io 2525
```

### Kafka Connection Error

**Symptom:** `Failed to connect to Kafka broker`

**Solution:**
```bash
# Restart Kafka and Zookeeper
docker-compose restart zookeeper kafka

# Wait 30 seconds for Kafka to fully start
sleep 30

# Restart messaging service
docker-compose restart messaging-service
```

### SMTP Authentication Error

**Symptom:** `535 5.7.0 Invalid credentials`

**Solution:**
1. Verify SMTP credentials in `.env` file
2. For Mailtrap:
   - Login to https://mailtrap.io
   - Go to "Email Testing" → "Inboxes" → "My Inbox"
   - Copy credentials from "Integrations" → "SMTP Settings"
3. Update `.env` file with correct credentials
4. Restart services:
   ```bash
   docker-compose restart messaging-service
   ```

## Security Considerations

1. **Password in Email:** The temporary password is sent via email. Users should change it immediately after first login.

2. **Email Privacy:** Emails are sent via Kafka, which provides:
   - Message persistence
   - Guaranteed delivery
   - Audit trail

3. **SMTP Security:** Use TLS/SSL for SMTP connections in production:
   ```bash
   SMTP_ENCRYPTION=tls  # or 'ssl'
   ```

4. **Production URLs:** Update login URLs for production:
   ```javascript
   const loginUrl = newUser.role === 'user' 
     ? 'https://usualstore.com' 
     : 'https://support.usualstore.com/dashboard';
   ```

## Production Checklist

- [ ] Update login URLs to production domains
- [ ] Configure production SMTP server (not Mailtrap)
- [ ] Enable TLS/SSL for SMTP
- [ ] Set up email templates with HTML formatting
- [ ] Add email rate limiting
- [ ] Configure Kafka message retention policy
- [ ] Set up monitoring for email delivery failures
- [ ] Add retry mechanism for failed emails
- [ ] Implement unsubscribe functionality (if needed)
- [ ] Add email delivery status tracking

## Related Documentation

- [Kafka Messaging Service](./KAFKA-MESSAGING-SERVICE.md)
- [Support Dashboard](./SUPPORT-DASHBOARD.md)
- [User Management](./USER-MANAGEMENT.md)
- [RBAC (Role-Based Access Control)](./RBAC.md)

## API Reference

### Send Email via Kafka

**Endpoint:** `POST /api/messaging/send`

**Request Body:**
```json
{
  "to": "user@example.com",
  "subject": "Welcome to Usual Store",
  "message": "Email body content"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Message sent to email-notifications topic"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Failed to send message to Kafka"
}
```

