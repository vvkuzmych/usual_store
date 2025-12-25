# 🚀 AI Assistant Quick Start

Get your AI shopping assistant running in 10 minutes!

---

## 📋 Prerequisites

1. ✅ OpenAI API key ([Get one here](https://platform.openai.com/api-keys))
2. ✅ PostgreSQL database running (you already have this!)
3. ✅ Go installed

---

## 🎯 Step 1: Get OpenAI API Key

```bash
# 1. Go to https://platform.openai.com/api-keys
# 2. Click "Create new secret key"
# 3. Copy the key (starts with sk-...)
# 4. Add to your environment

export OPENAI_API_KEY="sk-your-key-here"
```

**Cost:** ~$5 for 5000 messages (very cheap!)

---

## 🗄️ Step 2: Run Database Migrations

```bash
cd /Users/vkuzm/Projects/usual_store

# Apply migrations
cat migrations/20241225000001_create_ai_conversations.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore

cat migrations/20241225000002_create_ai_messages.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore

cat migrations/20241225000003_create_ai_user_preferences.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore

cat migrations/20241225000004_create_ai_feedback.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore

cat migrations/20241225000005_create_ai_product_cache.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore
```

Verify:
```bash
docker compose exec database psql -U postgres -d usualstore -c "\dt ai_*"

# Should show:
# ai_conversations
# ai_messages
# ai_user_preferences
# ai_feedback
# ai_product_cache
```

---

## 🏗️ Step 3: Build and Run

### **Option A: Standalone AI Service**

```bash
# Build
go build -o ai-assistant ./cmd/ai-assistant-example/

# Run
export OPENAI_API_KEY="sk-your-key-here"
export DATABASE_DSN="postgres://postgres:password@localhost:5433/usualstore?sslmode=disable"
export PORT="8080"

./ai-assistant

# Server starts on http://localhost:8080
```

### **Option B: Integrate with Existing Backend**

Add to your existing `back-end/main.go`:

```go
import (
    "usualstore/internal/ai"
)

func main() {
    // ... existing code ...
    
    // Create AI client
    openaiClient := ai.NewOpenAIClient(
        os.Getenv("OPENAI_API_KEY"),
        "gpt-3.5-turbo",
        0.7,
    )
    
    // Create AI service
    aiLogger := log.New(os.Stdout, "[AI] ", log.LstdFlags)
    aiService := ai.NewService(db, openaiClient, aiLogger)
    aiHandler := ai.NewHandler(aiService, aiLogger)
    
    // Register routes
    http.HandleFunc("/api/ai/chat", aiHandler.HandleChatRequest)
    http.HandleFunc("/api/ai/feedback", aiHandler.HandleFeedback)
    http.HandleFunc("/api/ai/stats", aiHandler.HandleStats)
    
    // ... rest of your code ...
}
```

---

## 🧪 Step 4: Test the API

### **Test with curl:**

```bash
# Start a conversation
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-session-123",
    "message": "Hi! I need help choosing a product for my mom",
    "user_id": null
  }'

# Response:
{
  "session_id": "test-session-123",
  "message": "Hello! I'd be happy to help you find something special for your mom! To give you the best recommendations, could you tell me:\n\n1. What's your budget?\n2. What are her interests or hobbies?\n3. Is this for a specific occasion?\n\nWe have great options including our Widget ($10) and Golden Plan subscription ($30/month with 30% off for 3+ subscriptions)!",
  "tokens_used": 156,
  "response_time_ms": 1234,
  "suggestions": [
    "Tell me more about this product",
    "What other options do you have?",
    "What's your best deal?"
  ]
}
```

### **Continue conversation:**

```bash
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-session-123",
    "message": "She likes gardening and my budget is around $30",
    "user_id": null
  }'
```

### **Submit feedback:**

```bash
curl -X POST http://localhost:8080/api/ai/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "message_id": 1,
    "conversation_id": 1,
    "helpful": true,
    "rating": 5,
    "feedback_text": "Very helpful!",
    "feedback_type": "helpful"
  }'
```

### **Get statistics:**

```bash
curl http://localhost:8080/api/ai/stats?days=7

# Response:
{
  "total_conversations": 15,
  "total_messages": 87,
  "total_cost": 0.152,
  "purchases": 3,
  "conversion_rate": 20.0
}
```

---

## 🎨 Step 5: Add Frontend Chat Widget

### **Simple HTML/JavaScript Example:**

```html
<!DOCTYPE html>
<html>
<head>
    <title>AI Shopping Assistant</title>
    <style>
        #chat-widget {
            position: fixed;
            bottom: 20px;
            right: 20px;
            width: 350px;
            height: 500px;
            border: 1px solid #ccc;
            border-radius: 10px;
            background: white;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            display: flex;
            flex-direction: column;
        }
        #chat-header {
            background: #007bff;
            color: white;
            padding: 15px;
            border-radius: 10px 10px 0 0;
        }
        #chat-messages {
            flex: 1;
            overflow-y: auto;
            padding: 15px;
        }
        .message {
            margin-bottom: 10px;
            padding: 10px;
            border-radius: 5px;
        }
        .user-message {
            background: #e3f2fd;
            text-align: right;
        }
        .assistant-message {
            background: #f5f5f5;
        }
        #chat-input {
            display: flex;
            padding: 15px;
            border-top: 1px solid #eee;
        }
        #message-input {
            flex: 1;
            padding: 10px;
            border: 1px solid #ccc;
            border-radius: 5px;
        }
        #send-button {
            margin-left: 10px;
            padding: 10px 20px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <div id="chat-widget">
        <div id="chat-header">
            <h3>🤖 Shopping Assistant</h3>
        </div>
        <div id="chat-messages"></div>
        <div id="chat-input">
            <input type="text" id="message-input" placeholder="Type your message...">
            <button id="send-button">Send</button>
        </div>
    </div>

    <script>
        const sessionId = 'session-' + Date.now();
        const messagesDiv = document.getElementById('chat-messages');
        const messageInput = document.getElementById('message-input');
        const sendButton = document.getElementById('send-button');

        function addMessage(content, isUser) {
            const msgDiv = document.createElement('div');
            msgDiv.className = 'message ' + (isUser ? 'user-message' : 'assistant-message');
            msgDiv.textContent = content;
            messagesDiv.appendChild(msgDiv);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }

        async function sendMessage() {
            const message = messageInput.value.trim();
            if (!message) return;

            // Add user message to UI
            addMessage(message, true);
            messageInput.value = '';

            // Send to API
            try {
                const response = await fetch('http://localhost:8080/api/ai/chat', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        session_id: sessionId,
                        message: message,
                        user_id: null
                    })
                });

                const data = await response.json();
                
                // Add assistant response to UI
                addMessage(data.message, false);

            } catch (error) {
                addMessage('Sorry, something went wrong. Please try again.', false);
                console.error('Error:', error);
            }
        }

        sendButton.addEventListener('click', sendMessage);
        messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') sendMessage();
        });

        // Welcome message
        addMessage('Hi! I\'m your shopping assistant. How can I help you today?', false);
    </script>
</body>
</html>
```

Save as `chat-test.html` and open in browser!

---

## ✅ Step 6: Verify Everything Works

### **Checklist:**

- [ ] Database tables created (5 ai_* tables)
- [ ] OpenAI API key configured
- [ ] Server running on port 8080
- [ ] Can send chat messages via curl
- [ ] Chat widget displays in browser
- [ ] Conversation history persists
- [ ] Stats endpoint returns data

---

## 🎉 You're Done!

Your AI shopping assistant is now live! Users can:
- Ask questions about products
- Get personalized recommendations
- Search using natural language
- Get instant support

---

## 📊 Monitor Usage

```bash
# Check logs
tail -f logs/ai-assistant.log

# Query database
docker compose exec database psql -U postgres -d usualstore

SELECT 
    COUNT(*) as conversations,
    SUM(total_messages) as messages,
    SUM(total_cost) as cost,
    COUNT(CASE WHEN resulted_in_purchase THEN 1 END) as purchases
FROM ai_conversations
WHERE started_at > NOW() - INTERVAL '1 day';
```

---

## 💰 Cost Tracking

```bash
# Today's cost
curl http://localhost:8080/api/ai/stats?days=1

# This week's cost
curl http://localhost:8080/api/ai/stats?days=7

# This month's cost
curl http://localhost:8080/api/ai/stats?days=30
```

**Expected costs:**
- 100 conversations: ~$0.50
- 1000 conversations: ~$5.00
- 10000 conversations: ~$50.00

---

## 🐛 Troubleshooting

### **"OPENAI_API_KEY not set"**
```bash
export OPENAI_API_KEY="sk-your-key-here"
```

### **"Connection refused to database"**
```bash
# Make sure Docker is running
docker compose ps

# Check database DSN
echo $DATABASE_DSN
```

### **"Table ai_conversations doesn't exist"**
```bash
# Run migrations again
cat migrations/20241225000001_create_ai_conversations.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore
```

### **"OpenAI API error: Insufficient quota"**
- Add credits to your OpenAI account
- Or switch to a free tier alternative (see AI-ASSISTANT-OVERVIEW.md)

---

## 🚀 Next Steps

1. **Customize the personality**: Edit system prompt in `openai_client.go`
2. **Add more features**: See `AI-ASSISTANT-IMPLEMENTATION.md`
3. **Deploy to production**: See Kubernetes or Docker deployment guides
4. **Monitor performance**: Set up alerts for costs and errors

---

**You now have a working AI shopping assistant!** 🎉

For advanced features and customization, see:
- `AI-ASSISTANT-IMPLEMENTATION.md` - Full implementation guide
- `AI-ASSISTANT-OVERVIEW.md` - Architecture and options
- `FRONTEND-INTEGRATION.md` - React/Vue/Angular examples

