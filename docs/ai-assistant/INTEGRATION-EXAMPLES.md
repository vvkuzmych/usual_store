# 🔌 AI Assistant Integration Examples

Quick examples for integrating the AI chat widget into your application.

---

## 📦 Available Options

1. **Vanilla JavaScript** (works anywhere)
2. **React Component** (for React apps)
3. **Embed Code** (copy-paste into any HTML)

---

## 1️⃣ Vanilla JavaScript (Standalone)

### **Quick Start (3 files)**

**1. Include in your HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/static/chat-widget.css">
</head>
<body>
    <!-- Your page content -->
    <h1>Welcome to Usual Store</h1>
    
    <!-- Chat widget container -->
    <div id="chat-widget-container"></div>

    <!-- Include scripts -->
    <script src="/static/chat-widget.js"></script>
    <script>
        // Initialize chat widget
        const chatWidget = new ChatWidget({
            apiUrl: 'http://localhost:8080',  // Your AI API URL
            position: 'bottom-right',
            theme: 'blue'
        });
    </script>
</body>
</html>
```

**2. Copy files to your project:**
```bash
cp /Users/vkuzm/Projects/usual_store/static/chat-widget.js /your-project/static/
cp /Users/vkuzm/Projects/usual_store/static/chat-widget.css /your-project/static/
```

**3. Done!** Open your page and you'll see the chat button.

---

## 2️⃣ React Component

### **Installation**

```bash
# Copy component to your React project
cp /Users/vkuzm/Projects/usual_store/frontend/src/components/ChatWidget.jsx \
   /your-react-app/src/components/

cp /Users/vkuzm/Projects/usual_store/frontend/src/components/ChatWidget.css \
   /your-react-app/src/components/
```

### **Usage in Your App**

```jsx
// App.jsx
import React from 'react';
import ChatWidget from './components/ChatWidget';
import './App.css';

function App() {
  return (
    <div className="App">
      {/* Your existing app */}
      <header>
        <h1>Usual Store</h1>
      </header>

      <main>
        {/* Your products, pages, etc */}
      </main>

      {/* Add AI Chat Widget */}
      <ChatWidget apiUrl="http://localhost:8080" />
    </div>
  );
}

export default App;
```

**That's it!** The chat widget will appear in the bottom-right corner.

---

## 3️⃣ CDN / Embed Code (Easiest)

### **Single Line Integration**

Add this to any HTML page:

```html
<!-- Before closing </body> tag -->
<script>
(function() {
    // Load CSS
    var css = document.createElement('link');
    css.rel = 'stylesheet';
    css.href = 'https://your-cdn.com/chat-widget.css';
    document.head.appendChild(css);

    // Load JS
    var js = document.createElement('script');
    js.src = 'https://your-cdn.com/chat-widget.js';
    js.onload = function() {
        new ChatWidget({
            apiUrl: 'http://localhost:8080',
            position: 'bottom-right'
        });
    };
    document.body.appendChild(js);
})();
</script>
```

---

## 4️⃣ WordPress Integration

### **Option A: Code Snippet Plugin**

1. Install "Code Snippets" plugin
2. Add new snippet:

```php
<?php
// Add AI Chat Widget to Footer
add_action('wp_footer', 'add_ai_chat_widget');
function add_ai_chat_widget() {
    ?>
    <link rel="stylesheet" href="<?php echo get_template_directory_uri(); ?>/chat-widget.css">
    <div id="chat-widget-container"></div>
    <script src="<?php echo get_template_directory_uri(); ?>/chat-widget.js"></script>
    <script>
        new ChatWidget({
            apiUrl: 'https://your-domain.com/ai-api',
            position: 'bottom-right'
        });
    </script>
    <?php
}
```

### **Option B: Theme Functions**

Add to your theme's `functions.php`:

```php
function enqueue_chat_widget() {
    wp_enqueue_style('chat-widget', get_template_directory_uri() . '/chat-widget.css');
    wp_enqueue_script('chat-widget', get_template_directory_uri() . '/chat-widget.js', array(), '1.0', true);
}
add_action('wp_enqueue_scripts', 'enqueue_chat_widget');
```

---

## 5️⃣ Shopify Integration

### **Add to theme.liquid**

```html
<!-- In theme.liquid, before </body> -->
<link rel="stylesheet" href="{{ 'chat-widget.css' | asset_url }}">
<div id="chat-widget-container"></div>
<script src="{{ 'chat-widget.js' | asset_url }}"></script>
<script>
  new ChatWidget({
    apiUrl: '{{ shop.metafields.ai.api_url }}',
    position: 'bottom-right'
  });
</script>
```

---

## 6️⃣ Next.js Integration

```jsx
// pages/_app.js
import ChatWidget from '../components/ChatWidget';
import '../styles/globals.css';

function MyApp({ Component, pageProps }) {
  return (
    <>
      <Component {...pageProps} />
      <ChatWidget apiUrl={process.env.NEXT_PUBLIC_AI_API_URL} />
    </>
  );
}

export default MyApp;
```

---

## 7️⃣ Vue.js Integration

```vue
<!-- App.vue -->
<template>
  <div id="app">
    <router-view/>
    <ChatWidget :api-url="apiUrl" />
  </div>
</template>

<script>
import ChatWidget from './components/ChatWidget.vue';

export default {
  name: 'App',
  components: {
    ChatWidget
  },
  data() {
    return {
      apiUrl: process.env.VUE_APP_AI_API_URL || 'http://localhost:8080'
    };
  }
};
</script>
```

---

## ⚙️ Configuration Options

All integration methods support these options:

```javascript
new ChatWidget({
  // Required
  apiUrl: 'http://localhost:8080',      // Your AI API endpoint
  
  // Optional
  position: 'bottom-right',              // 'bottom-right', 'bottom-left', 'top-right', 'top-left'
  theme: 'blue',                         // 'blue', 'green', 'purple', 'dark'
  welcomeMessage: 'Hi! How can I help?', // Custom welcome message
  
  // Advanced
  autoOpen: false,                       // Auto-open on page load
  openDelay: 3000,                       // Delay before auto-open (ms)
  userId: null,                          // User ID for authenticated users
  locale: 'en',                          // Language: 'en', 'es', 'fr', etc.
});
```

---

## 🎨 Customization Examples

### **Change Colors**

```css
/* Add to your CSS file */
.chat-toggle,
.chat-header,
.chat-send-btn {
  background: linear-gradient(135deg, #6a11cb, #2575fc) !important;
}

.user-message .message-content {
  background: linear-gradient(135deg, #6a11cb, #2575fc) !important;
}
```

### **Change Position**

```css
.chat-widget {
  bottom: 20px;
  left: 20px;  /* Changed from right */
  right: auto;
}

.chat-toggle {
  bottom: 20px;
  left: 20px;  /* Changed from right */
  right: auto;
}
```

### **Change Size**

```css
.chat-widget {
  width: 450px;   /* Default: 380px */
  height: 700px;  /* Default: 600px */
}
```

---

## 🔒 Security Best Practices

### **1. Use Environment Variables**

```javascript
// Don't hardcode API URLs
const chatWidget = new ChatWidget({
  apiUrl: process.env.REACT_APP_AI_API_URL  // React
  apiUrl: import.meta.env.VITE_AI_API_URL   // Vite
  apiUrl: process.env.NEXT_PUBLIC_AI_API_URL // Next.js
});
```

### **2. Enable CORS on Backend**

```go
// In your Go backend
w.Header().Set("Access-Control-Allow-Origin", "https://your-domain.com")
w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

### **3. Rate Limiting**

The backend already includes rate limiting, but you can add client-side throttling:

```javascript
// Limit messages to 10 per minute
let messageCount = 0;
let resetTime = Date.now() + 60000;

function canSendMessage() {
  if (Date.now() > resetTime) {
    messageCount = 0;
    resetTime = Date.now() + 60000;
  }
  if (messageCount >= 10) {
    alert('Please wait a moment before sending more messages.');
    return false;
  }
  messageCount++;
  return true;
}
```

---

## 📊 Analytics Integration

### **Google Analytics**

```javascript
// Track chat interactions
function trackChatEvent(action, label) {
  if (window.gtag) {
    gtag('event', action, {
      event_category: 'AI Chat',
      event_label: label
    });
  }
}

// Usage in ChatWidget
sendMessage: function() {
  // ... send message code ...
  trackChatEvent('message_sent', this.inputField.value);
}
```

### **Custom Analytics**

```javascript
// Track conversation metrics
const chatAnalytics = {
  sessionStart: Date.now(),
  messagesSent: 0,
  productsViewed: 0,
  
  track(event, data) {
    console.log('Chat Analytics:', event, data);
    // Send to your analytics service
    fetch('/api/analytics', {
      method: 'POST',
      body: JSON.stringify({ event, data })
    });
  }
};
```

---

## 🧪 Testing

### **Test in Browser Console**

```javascript
// Open chat programmatically
chatWidget.open();

// Send test message
chatWidget.inputField.value = 'Test message';
chatWidget.sendMessage();

// Close chat
chatWidget.close();

// Check session ID
console.log(chatWidget.sessionId);
```

---

## 🐛 Troubleshooting

### **Widget Not Showing**

```javascript
// Check if widget initialized
console.log(window.ChatWidget);

// Check if container exists
console.log(document.getElementById('chat-widget-container'));

// Check for errors
console.log('Check browser console for errors');
```

### **API Connection Issues**

```javascript
// Test API connectivity
fetch('http://localhost:8080/health')
  .then(r => r.text())
  .then(console.log)
  .catch(console.error);
```

### **CORS Errors**

Add to backend:
```go
w.Header().Set("Access-Control-Allow-Origin", "*")  // Development only
// Production: Use specific domain
w.Header().Set("Access-Control-Allow-Origin", "https://your-domain.com")
```

---

## 📚 Files Reference

All frontend files are located in:

```
/Users/vkuzm/Projects/usual_store/

📁 Vanilla JavaScript:
  static/
    ├── chat-widget.html  (Demo page)
    ├── chat-widget.css   (Styles)
    └── chat-widget.js    (Logic)

📁 React Component:
  frontend/src/components/
    ├── ChatWidget.jsx    (Component)
    └── ChatWidget.css    (Styles)
```

---

## 🚀 Next Steps

1. **Choose your integration method** (Vanilla JS, React, etc.)
2. **Copy the files** to your project
3. **Configure API URL** (point to your backend)
4. **Test locally** first
5. **Deploy to production**

---

**Your chat widget is ready to integrate!** 🎉

For more examples, see:
- `static/chat-widget.html` - Working demo
- `docs/ai-assistant/FRONTEND-INTEGRATION.md` - Detailed guide

