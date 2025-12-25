# React + Redux Frontend - Usual Store

A modern React frontend application with Redux Toolkit for state management, running on port **3002**.

## 🚀 Features

- **React 18** - Latest React features with hooks
- **Redux Toolkit** - Powerful state management with minimal boilerplate
- **Material-UI** - Beautiful, responsive UI components
- **React Router** - Client-side routing
- **Vite** - Lightning-fast build tool
- **Axios** - HTTP client with interceptors
- **TypeScript-ready** - Easy to migrate to TypeScript

## 📦 Redux Store Structure

### Slices

1. **authSlice** - Authentication state
   - Login/Logout
   - User session management
   - Token handling

2. **productsSlice** - Product management
   - Fetch all products
   - Fetch product by ID
   - Filtering and search

3. **cartSlice** - Shopping cart
   - Add/remove items
   - Update quantities
   - Calculate totals

4. **uiSlice** - UI state
   - Sidebar toggle
   - Cart drawer
   - Notifications
   - Theme (light/dark)

## 🛠️ Development

### Prerequisites

- Node.js 18+
- npm or yarn

### Install Dependencies

```bash
cd react-redux-frontend
npm install
```

### Run Development Server

```bash
npm run dev
```

Application will be available at `http://localhost:3002`

### Build for Production

```bash
npm run build
```

## 🐳 Docker

### Build Docker Image

```bash
docker build -t usual-store/redux-frontend:latest .
```

### Run Docker Container

```bash
docker run -p 3002:3002 usual-store/redux-frontend:latest
```

### Using Docker Compose

```bash
# Start with Redux frontend
docker-compose --profile redux-frontend up

# Or start all services
docker-compose --profile go-frontend --profile react-frontend --profile typescript-frontend --profile redux-frontend up
```

## ☸️ Kubernetes

### Deploy to Kubernetes

```bash
kubectl apply -f k8s/18-redux-frontend-deployment.yaml
```

### Access the Application

```bash
# Get service URL
kubectl get service redux-frontend -n usual-store

# Port forward for local access
kubectl port-forward service/redux-frontend 3002:3002 -n usual-store
```

Then visit `http://localhost:3002`

## 📂 Project Structure

```
react-redux-frontend/
├── src/
│   ├── components/          # Reusable components
│   │   ├── Header.jsx
│   │   ├── Footer.jsx
│   │   └── NotificationBar.jsx
│   ├── pages/               # Page components
│   │   ├── Home.jsx
│   │   ├── Products.jsx
│   │   ├── ProductDetail.jsx
│   │   ├── Login.jsx
│   │   └── Cart.jsx
│   ├── store/               # Redux store
│   │   ├── store.js         # Store configuration
│   │   └── slices/          # Redux slices
│   │       ├── authSlice.js
│   │       ├── productsSlice.js
│   │       ├── cartSlice.js
│   │       └── uiSlice.js
│   ├── services/            # API services
│   │   └── api.js           # Axios instance
│   ├── hooks/               # Custom hooks
│   │   └── useRedux.js      # Redux hooks
│   ├── utils/               # Utility functions
│   ├── App.jsx              # Main App component
│   ├── main.jsx             # Entry point
│   └── index.css            # Global styles
├── public/                  # Static assets
├── Dockerfile               # Docker configuration
├── nginx.conf               # Nginx configuration
├── vite.config.js           # Vite configuration
└── package.json             # Dependencies
```

## 🎯 Redux Features

### Async Operations

Using `createAsyncThunk` for async actions:

```javascript
export const fetchProducts = createAsyncThunk(
  'products/fetchAll',
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.get('/api/widget/all');
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response?.data?.message);
    }
  }
);
```

### Selectors

Memoized selectors for derived state:

```javascript
export const selectFilteredProducts = (state) => {
  const { items, filters } = state.products;
  return items.filter(/* filtering logic */);
};
```

### Custom Hooks

Type-safe Redux hooks:

```javascript
import { useAppDispatch, useAppSelector } from './hooks/useRedux';

function Component() {
  const dispatch = useAppDispatch();
  const products = useAppSelector((state) => state.products.items);
  
  // Use dispatch and state...
}
```

## 🔧 Configuration

### Environment Variables

Create `.env` file:

```env
VITE_API_URL=http://localhost:4000
NODE_ENV=development
```

### API Configuration

API base URL and interceptors in `src/services/api.js`:

```javascript
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:4000',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});
```

## 🎨 Theming

Toggle between light and dark themes:

```javascript
import { toggleTheme } from './store/slices/uiSlice';

const handleThemeToggle = () => {
  dispatch(toggleTheme());
};
```

## 📡 API Integration

### Authentication

```javascript
const result = await dispatch(loginUser({ email, password }));

if (loginUser.fulfilled.match(result)) {
  // Login successful
  navigate('/');
}
```

### Products

```javascript
// Fetch all products
dispatch(fetchProducts());

// Fetch single product
dispatch(fetchProductById(id));
```

### Cart Operations

```javascript
// Add to cart
dispatch(addToCart({ id, name, price, image }));

// Remove from cart
dispatch(removeFromCart(id));

// Update quantity
dispatch(updateQuantity({ id, quantity }));

// Clear cart
dispatch(clearCart());
```

## 🔔 Notifications

Show notifications using the UI slice:

```javascript
dispatch(addNotification({
  message: 'Product added to cart!',
  severity: 'success', // 'success', 'error', 'warning', 'info'
}));
```

## 🧪 Testing

```bash
npm run test
```

## 📊 Redux DevTools

Redux DevTools are enabled in development mode. Install the browser extension:

- [Chrome Extension](https://chrome.google.com/webstore/detail/redux-devtools/lmhkpmbekcpmknklioeibfkpmmfibljd)
- [Firefox Extension](https://addons.mozilla.org/en-US/firefox/addon/reduxdevtools/)

## 🚀 Quick Start Commands

```bash
# Development
npm run dev              # Start dev server

# Production
npm run build            # Build for production
npm run preview          # Preview production build

# Docker
docker-compose --profile redux-frontend up    # Start with Docker
docker-compose --profile redux-frontend down  # Stop

# Kubernetes
kubectl apply -f k8s/18-redux-frontend-deployment.yaml  # Deploy
kubectl delete -f k8s/18-redux-frontend-deployment.yaml # Remove
```

## 🌐 Access URLs

- **Development:** http://localhost:3002
- **Docker:** http://localhost:3002
- **Kubernetes:** http://localhost:3002 (with port-forward)

## 🔗 Backend API

The frontend connects to the Go backend API at:
- Default: `http://localhost:4000`
- Docker: `http://backend:4000`

## 📚 Key Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| React | 18.2.0 | UI framework |
| Redux Toolkit | 2.0.1 | State management |
| React Redux | 9.0.4 | React bindings for Redux |
| React Router | 6.20.1 | Routing |
| Material-UI | 5.14.20 | UI components |
| Axios | 1.6.2 | HTTP client |
| Vite | 5.0.8 | Build tool |

## 🎓 Learning Resources

- [Redux Toolkit Documentation](https://redux-toolkit.js.org/)
- [React Redux Hooks](https://react-redux.js.org/api/hooks)
- [Material-UI Components](https://mui.com/material-ui/getting-started/)
- [Vite Guide](https://vitejs.dev/guide/)

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Kill process on port 3002
lsof -ti:3002 | xargs kill -9
```

### API Connection Issues

Check that the backend is running on port 4000:

```bash
curl http://localhost:4000/api/widget/all
```

### Redux DevTools Not Working

1. Install browser extension
2. Check that `devTools: true` in store configuration
3. Open browser DevTools and look for Redux tab

## 📝 License

MIT License - see LICENSE file for details

---

Built with ❤️ using React + Redux Toolkit

