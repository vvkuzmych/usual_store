# Live Support Chat System

## Overview

The **Live Support Chat** system provides real-time communication between website users and support agents using WebSocket technology. This is separate from the AI Assistant and allows human-to-human interaction.

## Features

### 🎯 **For Users:**
- **Real-time Chat Widget** - Floating chat button on website
- **Instant Messaging** - WebSocket-based real-time communication
- **Session Persistence** - Chat history saved in database
- **Typing Indicators** - See when support agent is typing
- **Queue System** - Automatic ticket creation and priority management
- **Responsive Design** - Works on desktop and mobile

### 🎯 **For Support Agents:**
- **Supporter Dashboard** - Manage multiple chats simultaneously
- **Ticket Management** - View all open support tickets
- **Priority System** - Urgent, high, medium, low priority levels
- **Real-time Notifications** - Get notified of new messages
- **Status Tracking** - Open, assigned, in-progress, resolved, closed
- **Message History** - Access full conversation history

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Support Chat System                      │
└─────────────────────────────────────────────────────────────┘

 ┌──────────────┐        WebSocket         ┌──────────────┐
 │              │◄─────────────────────────►│              │
 │   User       │                           │   Support    │
 │   Widget     │                           │   Service    │
 │  (React)     │                           │  (Go + WS)   │
 │              │      REST API             │              │
 │              │◄─────────────────────────►│              │
 └──────────────┘                           └──────┬───────┘
                                                   │
      ┌────────────────────────────────────────────┼────────┐
      │                                            │        │
      ▼                                            ▼        ▼
┌─────────────┐                           ┌──────────┐  ┌─────┐
│  Supporter  │                           │PostgreSQL│  │ Hub │
│  Dashboard  │◄────WebSocket────────────►│ Database │  │     │
│   (React)   │                           │          │  │(WS) │
└─────────────┘                           └──────────┘  └─────┘
```

## Components

### 1. **Database Tables**

#### `support_tickets`
Stores support ticket information:
- `id` - Unique ticket ID
- `user_id` - User ID (nullable for guests)
- `supporter_id` - Assigned supporter ID (nullable)
- `subject` - Ticket subject
- `status` - open, assigned, in_progress, resolved, closed
- `priority` - low, medium, high, urgent
- `session_id` - Unique session identifier (UUID)
- `user_email` - User email
- `user_name` - User name
- `created_at`, `updated_at`, `closed_at` - Timestamps

#### `support_messages`
Stores chat messages:
- `id` - Message ID
- `ticket_id` - Related ticket
- `sender_id` - User/supporter ID (nullable)
- `sender_type` - user, supporter, system
- `sender_name` - Display name
- `message` - Message text
- `is_read` - Read status
- `created_at` - Timestamp

#### `support_sessions`
Tracks active WebSocket connections:
- `id` - Session ID
- `ticket_id` - Related ticket
- `user_id`, `supporter_id` - Connected users
- `is_active` - Connection status
- `started_at`, `ended_at` - Timestamps
- `user_connected`, `supporter_connected` - Connection flags

### 2. **Backend Service** (`cmd/support-service`)

**Technology:** Go + gorilla/websocket

**Features:**
- WebSocket hub for managing connections
- Real-time message broadcasting
- REST API for ticket management
- Message persistence
- Typing indicators
- Connection management

**Endpoints:**

REST API:
```
POST   /api/support/ticket                 - Create new ticket
GET    /api/support/tickets                - Get all open tickets
GET    /api/support/ticket/{sessionID}     - Get ticket by session
GET    /api/support/ticket/{sessionID}/messages - Get messages
POST   /api/support/ticket/{ticketID}/assign    - Assign to supporter
GET    /api/support/health                 - Health check
```

WebSocket:
```
WS     /ws/support/user/{sessionID}        - User connection
WS     /ws/support/supporter/{sessionID}   - Supporter connection
```

**WebSocket Message Format:**

Incoming:
```json
{
  "type": "message" | "typing" | "read",
  "message": "text content",
  "data": {}
}
```

Outgoing:
```json
{
  "type": "message" | "typing" | "system" | "supporter_joined",
  "message": {
    "id": 123,
    "ticket_id": 45,
    "sender_type": "user",
    "sender_name": "John Doe",
    "message": "Hello!",
    "created_at": "2025-12-26T10:00:00Z"
  }
}
```

### 3. **Frontend Components**

#### **User Chat Widget** (`support-frontend/src/components/SupportChatWidget.jsx`)

A floating chat widget that appears on the website:

**Features:**
- Floating button with unread count badge
- Initial form (name, email, subject)
- Real-time chat interface
- Typing indicators
- Message history
- Session persistence
- Minimizable/closable

**Usage:**
```javascript
import SupportChatWidget from './components/SupportChatWidget';

function App() {
  return (
    <div>
      {/* Your app content */}
      <SupportChatWidget />
    </div>
  );
}
```

#### **Supporter Dashboard** (`support-frontend/src/components/SupporterDashboard.jsx`)

A full-page dashboard for support agents:

**Features:**
- Ticket list with filters
- Priority and status badges
- Real-time chat interface
- Multiple ticket handling
- Message history
- Auto-refresh

**Usage:**
```javascript
import SupporterDashboard from './components/SupporterDashboard';

function App() {
  return <SupporterDashboard />;
}
```

## Installation & Setup

### 1. **Run Database Migrations**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Run migrations
soda migrate up
```

This will create the following tables:
- `support_tickets`
- `support_messages`
- `support_sessions`

### 2. **Start Services with Docker**

```bash
# Start support service + frontend
docker-compose --profile support up -d

# Or start everything including support
docker-compose --profile react-frontend --profile support up -d
```

### 3. **Access the Application**

**User Chat Widget:**
- URL: `http://localhost:3005/support`
- Or: `http://[::1]:3005/support` (IPv6)

**Supporter Dashboard:**
- URL: `http://localhost:3005/support/dashboard`
- Or: `http://[::1]:3005/support/dashboard` (IPv6)

**Backend API:**
- REST: `http://localhost:5000/api/support/`
- WebSocket: `ws://localhost:5000/ws/support/`

### 4. **Environment Variables**

Backend (`.env` or `docker-compose.yml`):
```bash
SUPPORT_SERVICE_PORT=5000
DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
```

Frontend (`support-frontend/.env`):
```bash
REACT_APP_SUPPORT_API_URL=http://localhost:5000
REACT_APP_SUPPORT_WS_URL=ws://localhost:5000
```

## Usage Examples

### As a Website User:

1. **Start a Chat:**
   - Click the floating chat button
   - Enter your name, email (optional), and subject
   - Click "Start Chat"

2. **Send Messages:**
   - Type your message in the input field
   - Press Enter or click Send
   - Messages appear instantly via WebSocket

3. **Wait for Support:**
   - A support agent will join your chat
   - You'll see a notification when they connect
   - Chat history is preserved

4. **End Chat:**
   - Click "End Chat" when done
   - Your session is saved in the database

### As a Support Agent:

1. **Access Dashboard:**
   - Go to `/support/dashboard`
   - View all open tickets

2. **Select a Ticket:**
   - Click on a ticket from the list
   - See priority and status
   - Chat interface opens

3. **Respond to User:**
   - Type your response
   - Press Enter to send
   - User receives message instantly

4. **Manage Multiple Chats:**
   - Switch between tickets in the list
   - Each ticket has its own WebSocket connection
   - Messages are saved to database

## Integration with Existing Frontend

### Add to React Frontend (`react-frontend`):

```javascript
// react-frontend/src/App.js
import SupportChatWidget from 'support-frontend/src/components/SupportChatWidget';

function App() {
  return (
    <div>
      <Routes>
        {/* Your existing routes */}
      </Routes>
      
      {/* Add support chat widget */}
      <SupportChatWidget />
    </div>
  );
}
```

### Add to Go HTML Templates:

```html
<!-- cmd/web/templates/base.layout.gohtml -->
<body>
  {{block "content" .}}{{end}}
  
  <!-- Support Chat Widget (React Component) -->
  <div id="support-chat-widget"></div>
  <script src="http://localhost:3005/static/js/main.js"></script>
</body>
```

## Comparison: AI Assistant vs Live Support

| Feature | AI Assistant | Live Support Chat |
|---------|-------------|-------------------|
| **Response Time** | Instant | Human response time |
| **Availability** | 24/7 | Business hours |
| **Cost** | Per API call | Per agent |
| **Personalization** | ML-based | Human empathy |
| **Complex Issues** | Limited | Full resolution |
| **Integration** | OpenAI API | WebSocket + Database |
| **Use Case** | Quick questions, product recommendations | Technical support, complaints |

## Best Practices

### For Users:
✅ **Be specific** in your subject line  
✅ **Provide details** about your issue  
✅ **Stay in the chat** for faster resolution  
✅ **Check message history** before asking again

### For Support Agents:
✅ **Respond promptly** to new tickets  
✅ **Use proper status** updates (assigned → in_progress → resolved)  
✅ **Set priority** correctly (urgent for critical issues)  
✅ **Close tickets** when resolved  
✅ **Be polite** and professional

### For Developers:
✅ **Monitor WebSocket** connections  
✅ **Handle reconnections** gracefully  
✅ **Implement rate limiting** to prevent abuse  
✅ **Log all messages** for audit trail  
✅ **Add authentication** for supporter dashboard  
✅ **Implement load balancing** for high traffic

## Troubleshooting

### User Can't Connect:
1. Check if support-service is running: `docker ps | grep support-service`
2. Check WebSocket URL in browser console
3. Verify firewall allows port 5000
4. Check database connection

### Supporter Dashboard Shows No Tickets:
1. Verify tickets exist: `docker exec -it usual_store-database-1 psql -U postgres -d usualstore -c "SELECT * FROM support_tickets;"`
2. Check API endpoint: `curl http://localhost:5000/api/support/tickets`
3. Restart support-service: `docker-compose restart support-service`

### Messages Not Appearing:
1. Check WebSocket connection in Network tab
2. Verify both user and supporter are connected to same session
3. Check database for message: `SELECT * FROM support_messages WHERE ticket_id = X;`
4. Restart WebSocket hub: `docker-compose restart support-service`

### Database Migration Errors:
```bash
# Check migration status
soda migrate status

# Force migration
soda migrate down
soda migrate up

# Verify tables exist
docker exec -it usual_store-database-1 psql -U postgres -d usualstore -c "\dt support_*"
```

## Future Enhancements

🚀 **Planned Features:**
- [ ] File upload support
- [ ] Video/voice chat
- [ ] Canned responses for agents
- [ ] Chat transcripts via email
- [ ] Customer satisfaction surveys
- [ ] Agent performance analytics
- [ ] Multi-language support
- [ ] Mobile app integration
- [ ] CRM integration (Salesforce, Zendesk)
- [ ] Chatbot handoff from AI Assistant to Live Support

## Security Considerations

### Authentication:
- **Users:** Session-based (no login required for guests)
- **Supporters:** Should implement JWT authentication (TODO)

### Rate Limiting:
```go
// Add to support-service/main.go
router.Use(RateLimitMiddleware) // Limit to 60 requests/min
```

### Input Validation:
- Sanitize all message inputs
- Validate email format
- Limit message length (max 8192 chars)

### CORS:
```go
// Configured in support-service/main.go
cors.Options{
  AllowedOrigins: []string{"https://yourdomain.com"},
  // ...
}
```

## Performance Optimization

### Database Indexes:
✅ Already added in migrations:
- `idx_support_tickets_session_id`
- `idx_support_messages_ticket_id`
- `idx_support_sessions_is_active`

### WebSocket Scaling:
For high traffic, consider:
- Redis Pub/Sub for multi-instance broadcasting
- Load balancer with sticky sessions
- Separate WebSocket server cluster

### Caching:
- Cache open tickets list (60 second TTL)
- Cache recent messages (30 second TTL)

## Testing

### Manual Testing:

1. **Create a ticket:**
```bash
curl -X POST http://localhost:5000/api/support/ticket \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "Test User",
    "user_email": "test@example.com",
    "subject": "Test Ticket",
    "priority": "high"
  }'
```

2. **Get tickets:**
```bash
curl http://localhost:5000/api/support/tickets
```

3. **Test WebSocket:**
```javascript
const ws = new WebSocket('ws://localhost:5000/ws/support/user/<sessionID>?name=TestUser');
ws.onopen = () => ws.send(JSON.stringify({type: 'message', message: 'Hello!'}));
ws.onmessage = (e) => console.log(JSON.parse(e.data));
```

### Automated Testing:
```bash
# Backend tests
go test ./internal/support/...

# Frontend tests
cd support-frontend
npm test
```

## API Reference

See full API documentation: [SUPPORT-API-REFERENCE.md](./SUPPORT-API-REFERENCE.md)

## License

Same as main application (usual_store)

## Support

For issues with the Live Support Chat system:
1. Check logs: `docker logs usual_store-support-service-1`
2. Check database: `docker exec -it usual_store-database-1 psql -U postgres -d usualstore`
3. Open GitHub issue with logs and steps to reproduce

---

**Created:** 2025-12-26  
**Last Updated:** 2025-12-26  
**Version:** 1.0.0

