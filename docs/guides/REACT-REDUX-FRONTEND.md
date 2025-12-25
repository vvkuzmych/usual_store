# React + Redux Frontend Setup Guide

**Port:** 3002  
**Framework:** React 18 + Redux Toolkit  
**Build Tool:** Vite

---

## 📋 Overview

The React + Redux frontend is a modern single-page application (SPA) with centralized state management using Redux Toolkit. It provides a rich user experience with features like:

- ✅ Centralized state management with Redux Toolkit
- ✅ Async operations with `createAsyncThunk`
- ✅ Type-safe Redux hooks
- ✅ Material-UI components
- ✅ Persistent authentication
- ✅ Light/Dark theme toggle
- ✅ Notification system
- ✅ Client-side routing

---

## 🏗️ Architecture

### Redux Store Structure

```
store/
├── store.js              # Store configuration
└── slices/
    ├── authSlice.js      # Authentication (login, logout, session)
    ├── productsSlice.js  # Products (fetch, filter, search)
    ├── cartSlice.js      # Shopping cart (add, remove, update)
    └── uiSlice.js        # UI state (theme, notifications, drawers)
```

### State Flow

```
Component → Dispatch Action → Redux Store → Update State → Re-render Component
                ↓
         API Call (if async)
                ↓
         Backend Response
```

---

## 🚀 Quick Start

### 1. Development Mode

```bash
cd react-redux-frontend
npm install
npm run dev
```

Access: http://localhost:3002

### 2. Docker Mode

```bash
# Build image
docker build -t usual-store/redux-frontend:latest react-redux-frontend/

# Run with Docker Compose
docker-compose --profile redux-frontend up
```

### 3. Kubernetes Mode

```bash
# Build and load image
docker build -t usual-store/redux-frontend:latest react-redux-frontend/
kubectl apply -f k8s/18-redux-frontend-deployment.yaml

# Port forward
kubectl port-forward service/redux-frontend 3002:3002 -n usual-store
```

---

## 🎯 Redux Features

### 1. Authentication State

**File:** `src/store/slices/authSlice.js`

```javascript
// Actions
dispatch(loginUser({ email, password }));
dispatch(logoutUser());
dispatch(checkAuthStatus());

// Selectors
const user = useAppSelector((state) => state.auth.user);
const isAuthenticated = useAppSelector((state) => state.auth.isAuthenticated);
```

### 2. Products State

**File:** `src/store/slices/productsSlice.js`

```javascript
// Actions
dispatch(fetchProducts());
dispatch(fetchProductById(id));
dispatch(setFilter({ filterType: 'category', value: 'widgets' }));

// Selectors
const products = useAppSelector((state) => state.products.items);
const filteredProducts = useAppSelector(selectFilteredProducts);
```

### 3. Cart State

**File:** `src/store/slices/cartSlice.js`

```javascript
// Actions
dispatch(addToCart({ id, name, price, image }));
dispatch(removeFromCart(id));
dispatch(updateQuantity({ id, quantity }));
dispatch(clearCart());

// Selectors
const cartItems = useAppSelector(selectCartItems);
const totalAmount = useAppSelector(selectCartTotal);
const totalQuantity = useAppSelector(selectCartQuantity);
```

### 4. UI State

**File:** `src/store/slices/uiSlice.js`

```javascript
// Actions
dispatch(toggleTheme());
dispatch(addNotification({ message, severity }));
dispatch(toggleCartDrawer());

// Selectors
const theme = useAppSelector((state) => state.ui.theme);
const notifications = useAppSelector((state) => state.ui.notifications);
```

---

## 📁 Project Structure

```
react-redux-frontend/
├── src/
│   ├── components/          # Reusable UI components
│   │   ├── Header.jsx       # Navigation + cart badge + theme toggle
│   │   ├── Footer.jsx       # Footer component
│   │   └── NotificationBar.jsx  # Toast notifications
│   ├── pages/               # Route pages
│   │   ├── Home.jsx         # Landing page
│   │   ├── Products.jsx     # Product listing
│   │   ├── ProductDetail.jsx # Single product view
│   │   ├── Login.jsx        # Authentication
│   │   └── Cart.jsx         # Shopping cart
│   ├── store/               # Redux configuration
│   │   ├── store.js         # Configure store
│   │   └── slices/          # Feature slices
│   ├── services/            # API layer
│   │   └── api.js           # Axios instance + interceptors
│   ├── hooks/               # Custom React hooks
│   │   └── useRedux.js      # Typed dispatch/selector
│   ├── App.jsx              # Main app with routing
│   └── main.jsx             # Entry point with Redux Provider
├── public/                  # Static assets
├── Dockerfile               # Multi-stage Docker build
├── nginx.conf               # Production server config
├── vite.config.js           # Vite configuration
└── package.json             # Dependencies
```

---

## 🔧 Configuration

### Environment Variables

Create `.env` file in `react-redux-frontend/`:

```env
VITE_API_URL=http://localhost:4000
NODE_ENV=development
```

### Vite Configuration

**File:** `vite.config.js`

```javascript
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 3002,
    proxy: {
      '/api': {
        target: 'http://localhost:4000',
        changeOrigin: true,
      }
    }
  },
});
```

### Nginx Configuration (Production)

**File:** `nginx.conf`

- Serves static files from `/usr/share/nginx/html`
- Proxies `/api/*` to backend
- Handles client-side routing
- Gzip compression enabled

---

## 🐳 Docker Setup

### Dockerfile

Two-stage build:
1. **Build stage**: Install deps, build with Vite
2. **Production stage**: Nginx alpine serving static files

```bash
# Build
docker build -t usual-store/redux-frontend:latest react-redux-frontend/

# Run
docker run -p 3002:3002 usual-store/redux-frontend:latest
```

### Docker Compose Profile

```yaml
redux-frontend:
  build:
    context: ./react-redux-frontend
    dockerfile: Dockerfile
  ports:
    - "3002:3002"
  networks:
    - usualstore_network
  depends_on:
    - back-end
  profiles:
    - redux-frontend
```

**Start:**
```bash
docker-compose --profile redux-frontend up
```

---

## ☸️ Kubernetes Deployment

### Deployment Configuration

**File:** `k8s/18-redux-frontend-deployment.yaml`

- 2 replicas for high availability
- Resource limits: 512Mi memory, 500m CPU
- Liveness and readiness probes
- LoadBalancer service on port 3002

### Deploy

```bash
# Apply configuration
kubectl apply -f k8s/18-redux-frontend-deployment.yaml

# Check status
kubectl get pods -n usual-store -l app=redux-frontend
kubectl get service redux-frontend -n usual-store

# Port forward
kubectl port-forward service/redux-frontend 3002:3002 -n usual-store
```

### Scale

```bash
# Scale replicas
kubectl scale deployment redux-frontend --replicas=3 -n usual-store
```

---

## 🎨 UI Components

### Material-UI Integration

All components use Material-UI for consistent design:

```javascript
import { Button, TextField, Card, Typography } from '@mui/material';
import { ShoppingCart } from '@mui/icons-material';
```

### Theme Provider

```javascript
const theme = createTheme({
  palette: {
    mode: themeMode, // 'light' or 'dark'
    primary: {
      main: '#1976d2',
    },
  },
});

<ThemeProvider theme={theme}>
  <CssBaseline />
  <App />
</ThemeProvider>
```

---

## 📡 API Integration

### Axios Configuration

**File:** `src/services/api.js`

```javascript
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:4000',
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
});

// Request interceptor - add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle 401
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear auth and redirect
      localStorage.removeItem('authToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

---

## 🔔 Notifications System

### Show Notification

```javascript
import { addNotification } from './store/slices/uiSlice';

dispatch(addNotification({
  message: 'Product added to cart!',
  severity: 'success', // 'success', 'error', 'warning', 'info'
}));
```

### Auto-dismiss

Notifications auto-dismiss after 6 seconds.

---

## 🧪 Testing

```bash
# Run tests (when added)
npm run test

# Coverage
npm run test:coverage
```

---

## 🔍 Redux DevTools

### Enable DevTools

DevTools are automatically enabled in development:

```javascript
export const store = configureStore({
  reducer: { /* ... */ },
  devTools: process.env.NODE_ENV !== 'production',
});
```

### Browser Extensions

- [Chrome](https://chrome.google.com/webstore/detail/redux-devtools/lmhkpmbekcpmknklioeibfkpmmfibljd)
- [Firefox](https://addons.mozilla.org/en-US/firefox/addon/reduxdevtools/)

### Features

- Time-travel debugging
- Action logging
- State inspection
- Action replay

---

## 📊 Comparison with Other Frontends

| Feature | Go (3000) | React (3000) | TypeScript (3001) | **Redux (3002)** |
|---------|-----------|--------------|-------------------|------------------|
| Framework | Go templates | React | React + TypeScript | React + Redux |
| State Mgmt | N/A | Context API | Context API | **Redux Toolkit** |
| Language | Go | JavaScript | TypeScript | JavaScript |
| Build Tool | N/A | Vite | Vite | Vite |
| UI Library | Bootstrap | Material-UI | Material-UI | Material-UI |
| Dev Server | Hot reload | Vite HMR | Vite HMR | Vite HMR |
| Production | Go binary | Nginx | Nginx | Nginx |

### When to Use Redux Frontend?

- ✅ Complex state management needs
- ✅ Multiple components sharing state
- ✅ Need for time-travel debugging
- ✅ Large-scale applications
- ✅ Team familiar with Redux
- ✅ Predictable state updates
- ✅ Advanced dev tools

---

## 🚀 Performance Optimization

### Code Splitting

Vite automatically splits code for optimal loading.

### Lazy Loading

```javascript
const ProductDetail = lazy(() => import('./pages/ProductDetail'));

<Suspense fallback={<Loading />}>
  <Route path="/products/:id" element={<ProductDetail />} />
</Suspense>
```

### Memoization

Use `createSelector` from Redux Toolkit for memoized selectors:

```javascript
export const selectFilteredProducts = createSelector(
  [(state) => state.products.items, (state) => state.products.filters],
  (items, filters) => items.filter(/* ... */)
);
```

---

## 🐛 Troubleshooting

### Issue: Redux state not updating

**Solution:** Check that you're using `createSlice` and returning state from reducers.

### Issue: API calls failing

**Solution:** Verify backend is running and check CORS settings.

### Issue: Theme not switching

**Solution:** Ensure `ThemeProvider` wraps the entire app.

### Issue: Port 3002 in use

```bash
lsof -ti:3002 | xargs kill -9
```

---

## 📚 Resources

- [Redux Toolkit Official Docs](https://redux-toolkit.js.org/)
- [React Redux Hooks](https://react-redux.js.org/api/hooks)
- [Material-UI Documentation](https://mui.com/)
- [Vite Guide](https://vitejs.dev/guide/)

---

## ✅ Checklist

- [x] Redux store configured with slices
- [x] Authentication flow implemented
- [x] Product listing and details
- [x] Shopping cart functionality
- [x] Material-UI theme integration
- [x] Docker configuration
- [x] Kubernetes deployment
- [x] API integration with backend
- [x] Notification system
- [x] Responsive design

---

**Status:** ✅ Production Ready  
**Port:** 3002  
**Technology:** React + Redux Toolkit + Material-UI + Vite

---

*Created: December 25, 2025*  
*Last Updated: December 25, 2025*

