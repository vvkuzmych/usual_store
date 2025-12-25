# 🎨 Frontend Integration Guide

How to add the AI assistant to your frontend (React, Vue, vanilla JS).

---

## 📦 React Component Example

### **Complete ChatWidget.jsx**

```jsx
import React, { useState, useEffect, useRef } from 'react';
import './ChatWidget.css';

const ChatWidget = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId] = useState(`session-${Date.now()}`);
  const messagesEndRef = useRef(null);

  const API_URL = process.env.REACT_APP_AI_API_URL || 'http://localhost:8080';

  useEffect(() => {
    // Welcome message
    if (messages.length === 0) {
      setMessages([{
        id: 0,
        role: 'assistant',
        content: '👋 Hi! I\'m your shopping assistant. How can I help you find the perfect product today?',
        timestamp: new Date()
      }]);
    }
  }, []);

  useEffect(() => {
    // Auto-scroll to bottom
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const sendMessage = async () => {
    const message = inputMessage.trim();
    if (!message || isLoading) return;

    // Add user message
    const userMessage = {
      id: messages.length,
      role: 'user',
      content: message,
      timestamp: new Date()
    };
    setMessages(prev => [...prev, userMessage]);
    setInputMessage('');
    setIsLoading(true);

    try {
      const response = await fetch(`${API_URL}/api/ai/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          session_id: sessionId,
          message: message,
          user_id: null // TODO: Get from auth context
        })
      });

      if (!response.ok) {
        throw new Error('Failed to get response');
      }

      const data = await response.json();

      // Add assistant message
      const assistantMessage = {
        id: messages.length + 1,
        role: 'assistant',
        content: data.message,
        timestamp: new Date(),
        products: data.products,
        suggestions: data.suggestions
      };
      setMessages(prev => [...prev, assistantMessage]);

    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage = {
        id: messages.length + 1,
        role: 'assistant',
        content: 'Sorry, I\'m having trouble connecting. Please try again in a moment.',
        timestamp: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const submitFeedback = async (messageId, helpful) => {
    try {
      await fetch(`${API_URL}/api/ai/feedback`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message_id: messageId,
          conversation_id: 1, // TODO: Get from conversation
          helpful: helpful,
          feedback_type: helpful ? 'helpful' : 'not_helpful'
        })
      });
    } catch (error) {
      console.error('Error submitting feedback:', error);
    }
  };

  return (
    <>
      {/* Chat Toggle Button */}
      {!isOpen && (
        <button className="chat-toggle" onClick={() => setIsOpen(true)}>
          💬 Need Help?
        </button>
      )}

      {/* Chat Widget */}
      {isOpen && (
        <div className="chat-widget">
          {/* Header */}
          <div className="chat-header">
            <div>
              <h3>🤖 Shopping Assistant</h3>
              <p>We're here to help!</p>
            </div>
            <button onClick={() => setIsOpen(false)}>×</button>
          </div>

          {/* Messages */}
          <div className="chat-messages">
            {messages.map((msg) => (
              <div key={msg.id} className={`message ${msg.role}-message`}>
                <div className="message-content">
                  {msg.content}
                </div>

                {/* Product Cards */}
                {msg.products && msg.products.length > 0 && (
                  <div className="product-recommendations">
                    {msg.products.map((product, idx) => (
                      <div key={idx} className="product-card">
                        {product.image && <img src={product.image} alt={product.name} />}
                        <h4>{product.name}</h4>
                        <p className="product-price">${product.price}</p>
                        <p className="product-reason">{product.reason}</p>
                        <button>View Product</button>
                      </div>
                    ))}
                  </div>
                )}

                {/* Suggestions */}
                {msg.suggestions && msg.suggestions.length > 0 && (
                  <div className="suggestions">
                    {msg.suggestions.map((suggestion, idx) => (
                      <button 
                        key={idx}
                        className="suggestion-chip"
                        onClick={() => setInputMessage(suggestion)}
                      >
                        {suggestion}
                      </button>
                    ))}
                  </div>
                )}

                {/* Feedback */}
                {msg.role === 'assistant' && msg.id > 0 && (
                  <div className="message-feedback">
                    <button onClick={() => submitFeedback(msg.id, true)} title="Helpful">
                      👍
                    </button>
                    <button onClick={() => submitFeedback(msg.id, false)} title="Not helpful">
                      👎
                    </button>
                  </div>
                )}
              </div>
            ))}

            {/* Loading indicator */}
            {isLoading && (
              <div className="message assistant-message">
                <div className="typing-indicator">
                  <span></span><span></span><span></span>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <div className="chat-input">
            <input
              type="text"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Ask me anything..."
              disabled={isLoading}
            />
            <button onClick={sendMessage} disabled={isLoading || !inputMessage.trim()}>
              Send
            </button>
          </div>
        </div>
      )}
    </>
  );
};

export default ChatWidget;
```

### **ChatWidget.css**

```css
.chat-toggle {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 15px 25px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 25px;
  cursor: pointer;
  font-size: 16px;
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.4);
  transition: all 0.3s;
  z-index: 1000;
}

.chat-toggle:hover {
  background: #0056b3;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 123, 255, 0.5);
}

.chat-widget {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 380px;
  height: 600px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  z-index: 1000;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.chat-header {
  background: linear-gradient(135deg, #007bff, #0056b3);
  color: white;
  padding: 20px;
  border-radius: 12px 12px 0 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header h3 {
  margin: 0;
  font-size: 18px;
}

.chat-header p {
  margin: 5px 0 0 0;
  font-size: 12px;
  opacity: 0.9;
}

.chat-header button {
  background: transparent;
  border: none;
  color: white;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #f8f9fa;
}

.message {
  margin-bottom: 15px;
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.message-content {
  padding: 12px 16px;
  border-radius: 12px;
  max-width: 85%;
  word-wrap: break-word;
}

.user-message .message-content {
  background: #007bff;
  color: white;
  margin-left: auto;
  border-bottom-right-radius: 4px;
}

.assistant-message .message-content {
  background: white;
  color: #333;
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.product-recommendations {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  overflow-x: auto;
}

.product-card {
  min-width: 200px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 12px;
  text-align: center;
}

.product-card img {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 8px;
}

.product-card h4 {
  font-size: 14px;
  margin: 0 0 5px 0;
}

.product-price {
  font-size: 18px;
  font-weight: bold;
  color: #007bff;
  margin: 5px 0;
}

.product-reason {
  font-size: 12px;
  color: #666;
  margin: 5px 0;
}

.product-card button {
  width: 100%;
  padding: 8px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.suggestion-chip {
  padding: 6px 12px;
  background: #e3f2fd;
  color: #007bff;
  border: 1px solid #90caf9;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.suggestion-chip:hover {
  background: #bbdefb;
}

.message-feedback {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.message-feedback button {
  background: transparent;
  border: none;
  font-size: 16px;
  cursor: pointer;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.message-feedback button:hover {
  opacity: 1;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #999;
  border-radius: 50%;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.7;
  }
  30% {
    transform: translateY(-10px);
    opacity: 1;
  }
}

.chat-input {
  display: flex;
  padding: 15px;
  background: white;
  border-top: 1px solid #dee2e6;
  border-radius: 0 0 12px 12px;
}

.chat-input input {
  flex: 1;
  padding: 12px;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.chat-input input:focus {
  border-color: #007bff;
}

.chat-input button {
  margin-left: 10px;
  padding: 12px 24px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.2s;
}

.chat-input button:hover:not(:disabled) {
  background: #0056b3;
}

.chat-input button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Mobile Responsive */
@media (max-width: 480px) {
  .chat-widget {
    width: 100%;
    height: 100%;
    bottom: 0;
    right: 0;
    border-radius: 0;
  }
  
  .chat-toggle {
    bottom: 15px;
    right: 15px;
  }
}
```

---

## 🔌 Usage in Your App

### **Add to App.jsx:**

```jsx
import ChatWidget from './components/ChatWidget';

function App() {
  return (
    <div className="App">
      {/* Your existing app code */}
      
      {/* Add AI Assistant */}
      <ChatWidget />
    </div>
  );
}
```

---

## 📱 Vue.js Example

```vue
<template>
  <div>
    <!-- Toggle Button -->
    <button v-if="!isOpen" class="chat-toggle" @click="isOpen = true">
      💬 Need Help?
    </button>

    <!-- Chat Widget -->
    <div v-if="isOpen" class="chat-widget">
      <div class="chat-header">
        <div>
          <h3>🤖 Shopping Assistant</h3>
        </div>
        <button @click="isOpen = false">×</button>
      </div>

      <div class="chat-messages" ref="messagesContainer">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message', `${msg.role}-message`]"
        >
          <div class="message-content">{{ msg.content }}</div>
        </div>

        <div v-if="isLoading" class="message assistant-message">
          <div class="typing-indicator">
            <span></span><span></span><span></span>
          </div>
        </div>
      </div>

      <div class="chat-input">
        <input
          v-model="inputMessage"
          @keyup.enter="sendMessage"
          placeholder="Ask me anything..."
          :disabled="isLoading"
        />
        <button @click="sendMessage" :disabled="isLoading || !inputMessage.trim()">
          Send
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChatWidget',
  data() {
    return {
      isOpen: false,
      messages: [],
      inputMessage: '',
      isLoading: false,
      sessionId: `session-${Date.now()}`,
      apiUrl: process.env.VUE_APP_AI_API_URL || 'http://localhost:8080'
    };
  },
  mounted() {
    this.messages.push({
      id: 0,
      role: 'assistant',
      content: '👋 Hi! I\'m your shopping assistant. How can I help you today?'
    });
  },
  methods: {
    async sendMessage() {
      const message = this.inputMessage.trim();
      if (!message || this.isLoading) return;

      // Add user message
      this.messages.push({
        id: this.messages.length,
        role: 'user',
        content: message
      });
      this.inputMessage = '';
      this.isLoading = true;

      try {
        const response = await fetch(`${this.apiUrl}/api/ai/chat`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            session_id: this.sessionId,
            message: message,
            user_id: null
          })
        });

        const data = await response.json();

        this.messages.push({
          id: this.messages.length,
          role: 'assistant',
          content: data.message
        });

      } catch (error) {
        console.error('Error:', error);
        this.messages.push({
          id: this.messages.length,
          role: 'assistant',
          content: 'Sorry, something went wrong.'
        });
      } finally {
        this.isLoading = false;
        this.$nextTick(() => {
          const container = this.$refs.messagesContainer;
          container.scrollTop = container.scrollHeight;
        });
      }
    }
  }
};
</script>

<style scoped>
/* Use the same CSS from React example */
</style>
```

---

## 🎨 Customization Options

### **1. Change Colors**

```css
/* Primary color */
.chat-header {
  background: linear-gradient(135deg, #6a11cb, #2575fc); /* Purple/Blue */
}

.chat-input button,
.user-message .message-content {
  background: #6a11cb;
}
```

### **2. Change Position**

```css
/* Top-right corner */
.chat-widget {
  top: 20px;
  bottom: auto;
}

/* Left side */
.chat-widget {
  left: 20px;
  right: auto;
}
```

### **3. Change Size**

```css
.chat-widget {
  width: 450px;   /* Wider */
  height: 700px;  /* Taller */
}
```

### **4. Add Avatar**

```jsx
<div className="message assistant-message">
  <img src="/bot-avatar.png" className="avatar" alt="Bot" />
  <div className="message-content">{msg.content}</div>
</div>
```

---

## 🚀 Advanced Features

### **1. Voice Input**

```jsx
const handleVoiceInput = () => {
  const recognition = new window.webkitSpeechRecognition();
  recognition.onresult = (event) => {
    const transcript = event.results[0][0].transcript;
    setInputMessage(transcript);
  };
  recognition.start();
};
```

### **2. File Upload**

```jsx
const handleFileUpload = async (file) => {
  // Upload to server, get URL
  const imageUrl = await uploadImage(file);
  
  // Send to AI
  sendMessage(`Can you help me find products similar to this image? ${imageUrl}`);
};
```

### **3. Typing Indicator (Real-time)**

```jsx
useEffect(() => {
  if (isTyping) {
    // Show "AI is typing..."
    const timeout = setTimeout(() => setIsTyping(false), 3000);
    return () => clearTimeout(timeout);
  }
}, [isTyping]);
```

---

## 📚 Complete Examples

See `/examples/` folder for complete working examples:
- `examples/react-chat-widget/` - Full React implementation
- `examples/vue-chat-widget/` - Full Vue implementation
- `examples/vanilla-js-chat/` - Pure JavaScript
- `examples/mobile-react-native/` - React Native (mobile)

---

**Your frontend is ready to chat!** 🎉

