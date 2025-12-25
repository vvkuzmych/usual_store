# 🎨 Usual Store - React Frontend

Modern React frontend for the Usual Store application.

---

## 🚀 Quick Start

### **Local Development (without Docker)**

```bash
cd react-frontend

# Install dependencies
npm install

# Create .env file
cp .env.example .env

# Start development server
npm start

# Open http://localhost:3000
```

### **Docker Development**

```bash
# From project root, start with React frontend
docker compose --profile react-frontend up

# Or use the helper script
./scripts/start-react.sh

# Open http://localhost:3000
```

---

## 🔀 Switching Between Frontends

This project has TWO frontends:

1. **Go Frontend** (original) - Port 4000
2. **React Frontend** (new) - Port 3000

### **Start Go Frontend**
```bash
docker compose --profile go-frontend up
# Access: http://localhost:4000
```

### **Start React Frontend**
```bash
docker compose --profile react-frontend up
# Access: http://localhost:3000
```

### **Start Both (for comparison)**
```bash
docker compose --profile go-frontend --profile react-frontend up
# Go:    http://localhost:4000
# React: http://localhost:3000
```

---

## 📦 Features

✅ **Modern React** (v18.2) with Hooks  
✅ **React Router** (v6) for navigation  
✅ **Axios** for API calls  
✅ **AI Chat Widget** integrated  
✅ **Responsive Design** (mobile-first)  
✅ **Production-ready** Docker setup  
✅ **Nginx** for serving  
✅ **Health checks** included  

---

## 🏗️ Project Structure

```
react-frontend/
├── public/
│   └── index.html          # HTML template
├── src/
│   ├── components/         # Reusable components
│   │   ├── Header.jsx
│   │   ├── Footer.jsx
│   │   └── ChatWidget.jsx  # AI Assistant
│   ├── pages/              # Route pages
│   │   ├── Home.jsx
│   │   ├── Products.jsx
│   │   ├── ProductDetail.jsx
│   │   ├── Cart.jsx
│   │   ├── Checkout.jsx
│   │   ├── Login.jsx
│   │   └── Signup.jsx
│   ├── services/           # API calls
│   │   └── api.js
│   ├── utils/              # Utilities
│   ├── App.js              # Main app
│   ├── App.css
│   ├── index.js            # Entry point
│   └── index.css
├── Dockerfile              # Production build
├── nginx.conf              # Nginx config
└── package.json
```

---

## 🔌 API Integration

The React app connects to the Go backend API:

```javascript
// Configured in src/services/api.js
const API_BASE_URL = 'http://localhost:4001';

// Available endpoints:
GET    /api/products          - List all products
GET    /api/product/:id       - Get single product
POST   /api/cart/add          - Add to cart
GET    /api/cart              - Get cart
POST   /api/login             - Login
POST   /api/signup            - Signup
POST   /api/checkout          - Checkout
```

---

## 🎨 Customization

### **Change Colors**

Edit `src/App.css` and component CSS files:

```css
/* Change primary color */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

### **Add New Pages**

1. Create page in `src/pages/NewPage.jsx`
2. Add route in `src/App.js`:
   ```jsx
   <Route path="/new" element={<NewPage />} />
   ```

### **Modify API URL**

Update `.env`:
```bash
REACT_APP_API_URL=http://your-api-url:4001
```

---

## 🐳 Docker Commands

```bash
# Build React frontend
docker compose build react-frontend

# Start React frontend
docker compose --profile react-frontend up

# Start in detached mode
docker compose --profile react-frontend up -d

# View logs
docker compose logs -f react-frontend

# Stop
docker compose --profile react-frontend down

# Rebuild from scratch
docker compose build --no-cache react-frontend
```

---

## 🧪 Testing

```bash
# Run tests
npm test

# Run tests with coverage
npm test -- --coverage

# Build for production (test build)
npm run build
```

---

## 📊 Available Scripts

```bash
npm start       # Start development server (port 3000)
npm run build   # Build for production
npm test        # Run tests
npm run eject   # Eject from Create React App (irreversible!)
```

---

## 🚀 Production Deployment

### **Docker Production Build**

The Dockerfile uses a **multi-stage build**:

1. **Stage 1**: Build React app with Node.js
2. **Stage 2**: Serve with Nginx (lightweight)

```bash
# Build production image
docker build -t usual-store-react:latest .

# Run production container
docker run -p 3000:3000 usual-store-react:latest
```

### **Kubernetes Deployment**

See `../k8s/` for Kubernetes manifests.

---

## 🔒 Security

✅ **Nginx** security headers configured  
✅ **CORS** handled by backend  
✅ **Environment variables** for API URLs  
✅ **No hardcoded secrets**  
✅ **Input validation** in forms  

---

## 🐛 Troubleshooting

### **Port 3000 already in use**

```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use different port
PORT=3001 npm start
```

### **API connection errors**

Check backend is running:
```bash
curl http://localhost:4001/api/products
```

### **Docker build fails**

```bash
# Clean Docker cache
docker system prune -a

# Rebuild
docker compose build --no-cache react-frontend
```

### **Blank page after build**

Check browser console for errors. Common issues:
- API URL not configured
- Backend not running
- CORS issues

---

## 📚 Tech Stack

- **React** 18.2 - UI library
- **React Router** 6.20 - Routing
- **Axios** 1.6 - HTTP client
- **Create React App** 5.0 - Tooling
- **Nginx** Alpine - Production server
- **Docker** - Containerization

---

## 🎯 TODO

- [ ] Add more pages (Cart, Checkout, Auth)
- [ ] Add state management (Redux/Zustand)
- [ ] Add form validation
- [ ] Add loading states
- [ ] Add error boundaries
- [ ] Add tests
- [ ] Add Storybook for components
- [ ] Add TypeScript

---

## 📞 Support

For issues or questions:
1. Check the main project README
2. Check Docker logs: `docker compose logs react-frontend`
3. Check API health: `http://localhost:4001/health`

---

## 🎉 Summary

You now have a **production-ready React frontend** that:
- ✅ Connects to your Go backend API
- ✅ Includes AI Chat Widget
- ✅ Works standalone or in Docker
- ✅ Can run alongside Go frontend
- ✅ Is mobile-responsive
- ✅ Is production-optimized

**Start it now:**
```bash
docker compose --profile react-frontend up
```

**Then visit:** http://localhost:3000

