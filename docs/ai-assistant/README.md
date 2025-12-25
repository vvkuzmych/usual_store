# 🤖 AI Shopping Assistant Documentation

Complete guide to implementing an AI-powered shopping assistant for Usual Store.

---

## 🎯 What is This?

An intelligent chatbot that helps customers:
- **Choose products** based on their needs and budget
- **Answer questions** about products, pricing, and features
- **Get recommendations** personalized to their preferences
- **Support** 24/7 without human intervention

**Built with:** OpenAI GPT, Go, PostgreSQL, React

---

## 📚 Documentation Index

### **🌟 Start Here**

**[QUICK-START.md](QUICK-START.md)** - Get your AI assistant running in 10 minutes
- Prerequisites and setup
- Database migrations
- API testing
- Frontend integration

**Time: 10-15 minutes**

---

### **📖 Understanding AI Assistants**

**[AI-ASSISTANT-OVERVIEW.md](AI-ASSISTANT-OVERVIEW.md)** - Comprehensive overview
- What is an AI shopping assistant?
- Use cases and examples
- Architecture options (OpenAI, Claude, Open-source)
- Cost breakdown
- Feature roadmap
- Security considerations

**Time: 20-30 minutes read**

---

### **🎨 Frontend Integration**

**[FRONTEND-INTEGRATION.md](FRONTEND-INTEGRATION.md)** - Add chat widget to your app
- Complete React component (with code)
- Vue.js example
- Vanilla JavaScript version
- Styling and customization
- Mobile responsive
- Advanced features (voice, files)

**Time: 15-20 minutes**

---

## 🗂️ Quick Reference

### **Architecture**

```
┌──────────────┐
│   Customer   │
└──────┬───────┘
       │ HTTP
       ▼
┌──────────────────────────┐
│  Frontend Chat Widget    │
│  (React/Vue/JS)          │
└──────────┬───────────────┘
           │ POST /api/ai/chat
           ▼
┌──────────────────────────┐
│  Go Backend              │
│  ├── AI Service          │
│  ├── Conversation DB     │
│  └── OpenAI Client       │
└──────────┬───────────────┘
           │ API Call
           ▼
┌──────────────────────────┐
│  OpenAI GPT-3.5/4        │
│  • Natural language      │
│  • Context-aware         │
│  • Product knowledge     │
└──────────────────────────┘
```

---

### **Cost Estimate**

| Conversations | Cost (GPT-3.5) | Cost (GPT-4) |
|--------------|----------------|--------------|
| 100          | ~$0.50         | ~$3.00       |
| 1,000        | ~$5.00         | ~$30.00      |
| 10,000       | ~$50.00        | ~$300.00     |

**Recommendation:** Start with GPT-3.5, upgrade to GPT-4 for better quality if needed.

---

### **Files Created**

```
Usual Store Project
├── docs/ai-assistant/
│   ├── README.md (this file)
│   ├── QUICK-START.md ⭐
│   ├── AI-ASSISTANT-OVERVIEW.md
│   └── FRONTEND-INTEGRATION.md
│
├── internal/ai/
│   ├── models.go (data structures)
│   ├── openai_client.go (OpenAI integration)
│   ├── service.go (business logic)
│   └── handlers.go (HTTP API)
│
├── cmd/ai-assistant-example/
│   └── main.go (standalone service)
│
├── migrations/
│   ├── 20241225000001_create_ai_conversations.up.sql
│   ├── 20241225000002_create_ai_messages.up.sql
│   ├── 20241225000003_create_ai_user_preferences.up.sql
│   ├── 20241225000004_create_ai_feedback.up.sql
│   └── 20241225000005_create_ai_product_cache.up.sql
│
└── AI-ASSISTANT-SUMMARY.txt (quick reference)
```

---

## 🚀 Quick Start Guide

### **1. Get OpenAI API Key**
```bash
# Visit: https://platform.openai.com/api-keys
export OPENAI_API_KEY="sk-your-key-here"
```

### **2. Run Migrations**
```bash
cd /Users/vkuzm/Projects/usual_store

for file in migrations/20241225000*.up.sql; do
  cat "$file" | docker compose exec -T database psql -U postgres -d usualstore
done
```

### **3. Build & Run**
```bash
go build -o ai-assistant ./cmd/ai-assistant-example/
export DATABASE_DSN="postgres://postgres:password@localhost:5433/usualstore?sslmode=disable"
./ai-assistant
```

### **4. Test**
```bash
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"session_id": "test-123", "message": "Hi! Help me choose a product"}'
```

---

## 💡 Features

### **✅ Implemented**
- Natural language chat (OpenAI GPT)
- Product recommendations
- Conversation history
- User preferences learning
- Feedback collection
- Analytics & cost tracking
- React frontend component
- RESTful API

### **🔜 Coming Soon** (you can implement)
- Voice input/output
- Image-based search
- Multi-language support
- Order tracking integration
- A/B testing framework
- Advanced analytics dashboard

---

## 📊 API Endpoints

### **POST /api/ai/chat**
Send a message to the AI assistant.

```bash
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "unique-session-id",
    "message": "I need a gift for my mom, budget $30",
    "user_id": null
  }'
```

**Response:**
```json
{
  "session_id": "unique-session-id",
  "message": "Perfect! The Golden Plan ($30/month) would be ideal...",
  "tokens_used": 156,
  "response_time_ms": 1234,
  "suggestions": ["Tell me more", "What else do you have?"],
  "products": [...]
}
```

### **POST /api/ai/feedback**
Submit feedback on an AI response.

```bash
curl -X POST http://localhost:8080/api/ai/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "message_id": 1,
    "conversation_id": 1,
    "helpful": true,
    "rating": 5
  }'
```

### **GET /api/ai/stats**
Get analytics about AI usage.

```bash
curl http://localhost:8080/api/ai/stats?days=7
```

**Response:**
```json
{
  "total_conversations": 150,
  "total_messages": 890,
  "total_cost": 4.52,
  "purchases": 30,
  "conversion_rate": 20.0
}
```

---

## 🗄️ Database Schema

### **ai_conversations**
Tracks chat sessions.
- `id`, `session_id`, `user_id`
- `started_at`, `ended_at`
- `total_messages`, `total_tokens_used`, `total_cost`
- `resulted_in_purchase`

### **ai_messages**
Stores individual messages.
- `id`, `conversation_id`
- `role` (user/assistant/system)
- `content`, `tokens_used`
- `model`, `response_time_ms`

### **ai_user_preferences**
Learns from user interactions.
- `id`, `user_id`, `session_id`
- `preferred_categories`, `budget_min/max`
- `last_products_viewed/purchased`
- `conversation_style`

### **ai_feedback**
Collects user feedback.
- `id`, `message_id`, `conversation_id`
- `helpful`, `rating` (1-5)
- `feedback_text`, `feedback_type`

### **ai_product_cache**
Optimizes product lookups.
- `id`, `product_id`
- `description_text`, `search_keywords`
- `popularity_score`

---

## 🎨 Frontend Example

```jsx
import ChatWidget from './components/ChatWidget';

function App() {
  return (
    <div>
      <YourStore />
      <ChatWidget /> {/* AI Assistant */}
    </div>
  );
}
```

**Features:**
- Floating chat button (bottom-right)
- Expandable chat window
- Typing indicators
- Product recommendations
- Feedback buttons (👍 👎)
- Mobile responsive

See [FRONTEND-INTEGRATION.md](FRONTEND-INTEGRATION.md) for complete code.

---

## 🔒 Security

- ✅ API keys in environment variables
- ✅ Rate limiting (10 msg/min per user)
- ✅ Input validation and sanitization
- ✅ SQL injection prevention
- ✅ CORS configuration
- ✅ Session management
- ✅ No PII sent to external APIs

---

## 💰 Cost Management

### **Monitor Costs**
```sql
-- Daily cost
SELECT SUM(total_cost) as cost
FROM ai_conversations
WHERE started_at > NOW() - INTERVAL '1 day';

-- Most expensive conversations
SELECT session_id, total_messages, total_cost
FROM ai_conversations
ORDER BY total_cost DESC
LIMIT 10;
```

### **Optimize Costs**
1. Use GPT-3.5-turbo instead of GPT-4 (10x cheaper)
2. Set `max_tokens` to limit response length
3. Cache product information
4. Use embeddings for simple queries
5. Implement rate limiting

---

## 🐛 Troubleshooting

### **"OPENAI_API_KEY not set"**
```bash
export OPENAI_API_KEY="sk-your-key-here"
```

### **"Connection refused"**
Check database is running:
```bash
docker compose ps
```

### **"Table doesn't exist"**
Run migrations:
```bash
cat migrations/20241225000001_create_ai_conversations.up.sql | \
  docker compose exec -T database psql -U postgres -d usualstore
```

### **"Insufficient quota" (OpenAI)**
- Add credits: https://platform.openai.com/account/billing
- Or use free alternatives (see OVERVIEW.md)

---

## 📈 Success Metrics

Track these KPIs:

**Engagement:**
- Conversations started
- Messages per conversation
- Return users

**Business:**
- Conversion rate (chat → purchase)
- Average order value (with AI vs without)
- Customer satisfaction (feedback ratings)

**Technical:**
- Response time (target: < 2 seconds)
- API cost per conversation (target: < $0.01)
- Error rate (target: < 1%)

---

## 🎓 Customization

### **Change AI Personality**
Edit `internal/ai/openai_client.go`:
```go
func (c *OpenAIClient) buildSystemPrompt(productContext string) string {
    return fmt.Sprintf(`You are a friendly, enthusiastic shopping assistant...`)
}
```

### **Add More Context**
Edit `internal/ai/service.go`:
```go
func (s *Service) getProductContext() (string, error) {
    // Add user's order history, preferences, etc.
}
```

### **Use Different Model**
```go
openaiClient := ai.NewOpenAIClient(
    apiKey,
    "gpt-4", // or "gpt-3.5-turbo"
    0.7,
)
```

---

## 🚀 Next Steps

### **Today:**
1. Read [QUICK-START.md](QUICK-START.md)
2. Get OpenAI API key
3. Run migrations
4. Test with curl

### **This Week:**
1. Add frontend widget
2. Customize AI personality
3. Test with real users
4. Monitor costs

### **This Month:**
1. Analyze metrics
2. Optimize prompts
3. Add advanced features
4. Scale to production

---

## 📖 Resources

**Official Docs:**
- OpenAI API: https://platform.openai.com/docs
- Go PostgreSQL: https://github.com/lib/pq
- React: https://react.dev

**Learning:**
- Prompt Engineering: https://platform.openai.com/docs/guides/prompt-engineering
- GPT Best Practices: https://platform.openai.com/docs/guides/gpt-best-practices

**Alternatives:**
- Claude (Anthropic): https://www.anthropic.com
- Ollama (Local LLMs): https://ollama.com
- Hugging Face: https://huggingface.co

---

## 💬 Example Conversations

### **Product Selection**
```
Customer: "I need a gift for my mom's birthday"
AI: "What's your budget and what does she like?"
Customer: "Around $30, she likes gardening"
AI: "Perfect! The Golden Plan at $30/month is ideal..."
```

### **Product Comparison**
```
Customer: "What's the difference between Widget and Golden Plan?"
AI: "Widget is $10 one-time, great for basic needs.
     Golden Plan is $30/month with recurring benefits..."
```

### **Budget Constraint**
```
Customer: "Show me your cheapest option"
AI: "Our Widget is just $10! It's perfect for..."
```

---

## 🎉 Summary

You now have:
- ✅ Complete AI assistant implementation
- ✅ Production-ready Go backend
- ✅ Database schema with 5 tables
- ✅ Frontend React component
- ✅ Comprehensive documentation
- ✅ Cost tracking and analytics

**Expected Benefits:**
- 24/7 customer support
- Instant product recommendations
- Higher conversion rates
- Reduced support costs
- Better customer experience

---

**Ready to build your AI assistant?**

👉 **Start here:** [QUICK-START.md](QUICK-START.md)

**Questions or issues?** Review the documentation or check the troubleshooting sections.

---

**🚀 Happy building!**

