# Live Support Chat - Quick Start Guide

Get the Live Support Chat system up and running in 5 minutes! 🚀

## Prerequisites

- Docker and Docker Compose installed
- Database running (`usual_store-database-1`)
- Ports 5000 and 3005 available

## Quick Start

### 1. Run Database Migrations

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store
soda migrate up
```

**Expected Output:**
```
✅ support_tickets table created
✅ support_messages table created
✅ support_sessions table created
```

### 2. Start Support Services

```bash
docker-compose --profile support up -d
```

**Services Started:**
- `support-service` (Backend) - Port 5000
- `support-frontend` (React UI) - Port 3005

### 3. Verify Services

```bash
# Check if containers are running
docker ps | grep support

# Test backend health
curl http://localhost:5000/api/support/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "active_sessions": 0,
  "timestamp": 1735200000
}
```

### 4. Open in Browser

**User Chat Widget:**
```
http://localhost:3005/support
```

**Supporter Dashboard:**
```
http://localhost:3005/support/dashboard
```

## Test It Out!

### As a User:

1. Go to: `http://localhost:3005/support`
2. Fill in the form:
   - **Name:** Test User
   - **Email:** test@example.com (optional)
   - **Subject:** Testing support chat
3. Click **"Start Chat"**
4. Send a message: "Hello, I need help!"
5. ✅ Message sent via WebSocket!

### As a Supporter:

1. Open **new browser tab** (or incognito window)
2. Go to: `http://localhost:3005/support/dashboard`
3. You should see the ticket from "Test User"
4. Click on the ticket to open chat
5. Send a response: "Hi! How can I help you?"
6. ✅ User receives your message instantly!

## Integration with Existing Frontend

### Add Chat Widget to React Frontend:

1. **Copy component:**
```bash
cp support-frontend/src/components/SupportChatWidget.jsx react-frontend/src/components/
```

2. **Add to App.js:**
```javascript
import SupportChatWidget from './components/SupportChatWidget';

function App() {
  return (
    <div>
      {/* Your existing app */}
      <SupportChatWidget />
    </div>
  );
}
```

3. **Rebuild frontend:**
```bash
docker-compose --profile react-frontend up -d --build
```

4. **Test:**
- Go to: `http://localhost:3000`
- You should see the floating support button in the bottom-right corner! 🎉

## Common Issues

### ❌ Port Already in Use

```bash
# Check what's using port 5000
lsof -i :5000

# Kill the process
kill -9 <PID>

# Or use different port in docker-compose.yml
```

### ❌ Database Tables Not Found

```bash
# Verify migrations
soda migrate status

# Re-run migrations
soda migrate down
soda migrate up
```

### ❌ WebSocket Connection Failed

```bash
# Check support-service logs
docker logs usual_store-support-service-1

# Restart service
docker-compose restart support-service
```

### ❌ Frontend Shows Blank Page

```bash
# Check build logs
docker logs usual_store-support-frontend-1

# Rebuild with no cache
docker-compose build --no-cache support-frontend
docker-compose up -d support-frontend
```

## Next Steps

1. **Customize the UI** - Edit `support-frontend/src/components/`
2. **Add Authentication** - Implement JWT for supporter dashboard
3. **Set Up Monitoring** - Add Prometheus metrics
4. **Enable HTTPS** - Configure SSL certificates
5. **Read Full Documentation** - See [LIVE-SUPPORT-CHAT.md](./LIVE-SUPPORT-CHAT.md)

## Docker Commands Cheat Sheet

```bash
# Start support services
docker-compose --profile support up -d

# Stop support services
docker-compose --profile support down

# View logs
docker logs -f usual_store-support-service-1
docker logs -f usual_store-support-frontend-1

# Rebuild services
docker-compose build support-service support-frontend
docker-compose --profile support up -d

# Restart services
docker-compose restart support-service support-frontend

# Check status
docker ps | grep support

# Access database
docker exec -it usual_store-database-1 psql -U postgres -d usualstore
```

## Database Queries

```sql
-- View all tickets
SELECT * FROM support_tickets;

-- View all messages for a ticket
SELECT * FROM support_messages WHERE ticket_id = 1 ORDER BY created_at;

-- View open tickets
SELECT * FROM support_tickets WHERE status IN ('open', 'assigned', 'in_progress');

-- Count messages per ticket
SELECT ticket_id, COUNT(*) FROM support_messages GROUP BY ticket_id;
```

## Success! 🎉

You now have a fully functional live support chat system!

**What You Get:**
- ✅ Real-time WebSocket communication
- ✅ Floating chat widget for users
- ✅ Supporter dashboard for agents
- ✅ Message persistence in PostgreSQL
- ✅ Typing indicators
- ✅ Session management
- ✅ Priority and status tracking

**Next:** Integrate with your existing frontend or customize the UI!

---

Need help? Check the full documentation: [LIVE-SUPPORT-CHAT.md](./LIVE-SUPPORT-CHAT.md)

