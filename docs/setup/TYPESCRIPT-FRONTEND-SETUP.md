# TypeScript Frontend Setup Guide

Complete guide for setting up and using the **TypeScript-based Vite frontend** for Usual Store.

## 📋 Overview

The TypeScript frontend is a modern, type-safe alternative to the JavaScript React frontend, built with:

- **TypeScript 5.3** - Full type safety
- **Vite 5.0** - Lightning-fast builds
- **React 18.2** - Latest React features
- **Material UI** - Enterprise-grade components
- **Strict typing** - Compile-time error detection

## 🆚 Comparison with Other Frontends

| Feature | Go HTML | React JS | TypeScript + Vite |
|---------|---------|----------|-------------------|
| **Port** | 4000 | 3000 | **3001** |
| **Language** | Go templates | JavaScript | TypeScript |
| **Type Safety** | ❌ | ❌ | ✅ |
| **Build Tool** | Go | CRA/Webpack | Vite |
| **Dev Speed** | N/A | Slow | **Very Fast** |
| **Bundle Size** | N/A | Large | **Small** |
| **Hot Reload** | ❌ | Slow | **Instant** |
| **Theme** | Purple | Purple | **Blue** |

## 🏗️ Architecture

```
┌─────────────────────┐
│  TypeScript Frontend │
│    (Port 3001)      │
│  Vite + React + MUI │
└──────────┬──────────┘
           │ HTTP/REST
           ↓
┌─────────────────────┐
│   Backend API (Go)  │
│    (Port 4001)      │
└──────────┬──────────┘
           │
           ↓
┌─────────────────────┐
│  PostgreSQL Database│
│    (Port 5433)      │
└─────────────────────┘
```

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Start TypeScript frontend + backend
cd /Users/vkuzm/Projects/UsualStore/usual_store
./scripts/start-typescript.sh

# Or manually:
docker compose --profile typescript-frontend up -d
```

**Access at**: http://localhost:3001

### Option 2: Local Development

```bash
cd typescript-frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

**Access at**: http://localhost:3001 (with hot reload)

## 📦 Installation

### Prerequisites

- **Node.js** 18+ and npm
- **Docker** and Docker Compose
- **Go** 1.23+ (for backend)

### Install Dependencies

```bash
cd typescript-frontend
npm install
```

This installs:
- React 18.2
- TypeScript 5.3
- Vite 5.0
- Material UI 5.14
- React Router 6.20
- Axios
- Stripe React Elements

## 🛠️ Development

### Commands

```bash
# Development server with hot reload
npm run dev

# Type checking (no build)
npm run type-check

# Linting
npm run lint

# Production build
npm run build

# Preview production build
npm run preview
```

### Configuration Files

| File | Purpose |
|------|---------|
| `tsconfig.json` | TypeScript compiler options |
| `tsconfig.node.json` | TypeScript config for Vite |
| `vite.config.ts` | Vite build configuration |
| `package.json` | Dependencies and scripts |

### Project Structure

```
typescript-frontend/
├── src/
│   ├── components/         # Reusable components
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   ├── pages/              # Route pages
│   │   ├── Home.tsx
│   │   ├── Products.tsx
│   │   ├── ProductDetail.tsx
│   │   ├── Login.tsx
│   │   └── Cart.tsx
│   ├── services/           # API integration
│   │   └── api.ts
│   ├── contexts/           # React contexts
│   │   └── AuthContext.tsx
│   ├── types/              # TypeScript types
│   │   └── index.ts
│   ├── theme.ts            # MUI theme
│   ├── App.tsx             # Root component
│   └── main.tsx            # Entry point
├── public/                 # Static files
├── Dockerfile              # Docker config
├── nginx.conf              # Nginx config
├── vite.config.ts          # Vite config
├── tsconfig.json           # TS config
└── package.json
```

## 🎨 Theming

The TypeScript frontend uses a **blue gradient theme** to differentiate from other frontends:

```typescript
// theme.ts
const theme = createTheme({
  palette: {
    primary: {
      main: '#2196f3',  // Blue
    },
    secondary: {
      main: '#00bcd4',  // Cyan
    },
  },
});
```

### Color Scheme

- **Primary**: Blue gradient (#2196f3 → #21cbf3)
- **Secondary**: Cyan (#00bcd4)
- **Success**: Green (#4caf50)
- **Error**: Red (#f44336)

## 🔐 Authentication

### Login Flow

1. User enters credentials on `/login`
2. Frontend calls `POST /api/authenticate`
3. Backend returns JWT token + user data
4. Token stored in localStorage
5. Token included in all subsequent API requests

### Implementation

```typescript
// AuthContext.tsx
const login = async (credentials: LoginCredentials) => {
  const response = await apiService.login(credentials);
  setUser(response.user);
  localStorage.setItem('user', JSON.stringify(response.user));
  localStorage.setItem('authToken', response.authentication_token.token);
};
```

### Demo Credentials

```
Email:    admin@example.com
Password: qwerty
```

## 🌐 API Integration

### Type-Safe API Calls

```typescript
// types/index.ts
export interface Product {
  id: number;
  name: string;
  price: number;
  description: string;
  image: string;
  is_recurring: boolean;
  inventory_level: number;
}

// services/api.ts
async getProducts(): Promise<Product[]> {
  const response = await this.api.get<Product[]>('/api/products');
  return response.data;
}
```

### Endpoints Used

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/api/products` | List all products |
| GET | `/api/product/:id` | Get single product |
| POST | `/api/authenticate` | User login |
| POST | `/api/payment-intent` | Create Stripe intent |
| POST | `/api/charge` | Process payment |
| POST | `/api/ai/chat` | AI assistant chat |
| GET | `/api/ai/stats` | AI statistics |

## 🐳 Docker Setup

### Dockerfile

Multi-stage build for optimized production image:

```dockerfile
# Stage 1: Build
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Nginx
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 3001
```

### Build and Run

```bash
# Build image
docker build -t usual_store-typescript-frontend ./typescript-frontend

# Run container
docker run -p 3001:3001 usual_store-typescript-frontend

# With Docker Compose
docker compose --profile typescript-frontend up -d
```

### Nginx Configuration

```nginx
server {
    listen 3001;
    root /usr/share/nginx/html;
    
    # API proxy
    location /api/ {
        proxy_pass http://back-end:4001/api/;
    }
    
    # SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes

```bash
# Build image
docker build -t usual_store-typescript-frontend:latest ./typescript-frontend

# Load to Docker Desktop Kubernetes
docker save usual_store-typescript-frontend:latest | docker load

# Deploy
kubectl apply -f k8s/01-namespace.yaml
kubectl apply -f k8s/17-typescript-frontend-deployment.yaml

# Check status
kubectl get pods -n usual-store | grep typescript
kubectl get svc -n usual-store | grep typescript

# Port forward
kubectl port-forward -n usual-store svc/typescript-frontend 3001:3001
```

### Kubernetes Resources

| Resource | File | Replicas | Resources |
|----------|------|----------|-----------|
| Deployment | `17-typescript-frontend-deployment.yaml` | 2 | 256Mi-512Mi RAM |
| Service | Same file | LoadBalancer | Port 3001 |

### Resource Limits

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

## 🧪 Testing

### Type Checking

```bash
npm run type-check
```

This runs TypeScript compiler without emitting files, catching type errors.

### Linting

```bash
npm run lint
```

Uses ESLint with TypeScript rules.

### Manual Testing

1. **Start services**:
   ```bash
   ./scripts/start-typescript.sh
   ```

2. **Open browser**: http://localhost:3001

3. **Test features**:
   - [ ] Homepage loads
   - [ ] Products list displays
   - [ ] Product details work
   - [ ] Login with demo credentials
   - [ ] Payment form shows (after login)
   - [ ] Navigation works
   - [ ] Responsive design

## 📊 Performance

### Build Metrics

- **Build Time**: ~10-20 seconds
- **Bundle Size**: ~200-300 KB (gzipped)
- **First Load**: ~1-2 seconds
- **Time to Interactive**: ~2-3 seconds

### Optimization

- ✅ Code splitting
- ✅ Tree shaking
- ✅ Lazy loading
- ✅ Asset optimization
- ✅ Gzip compression

## 🔧 Troubleshooting

### Issue: "Cannot find module"

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Issue: TypeScript errors

```bash
# Check TypeScript version
npm list typescript

# Reinstall TypeScript
npm install -D typescript@latest
```

### Issue: Vite port already in use

```bash
# Find and kill process on port 3001
lsof -ti:3001 | xargs kill -9

# Or use different port
npm run dev -- --port 3002
```

### Issue: API calls failing

Check that backend is running:

```bash
curl http://localhost:4001/api/products
```

If not running:

```bash
docker compose ps
docker compose logs back-end
```

## 🆚 When to Use TypeScript Frontend

**Use TypeScript Frontend when**:
- ✅ You need type safety
- ✅ Building large applications
- ✅ Working in a team
- ✅ Want best dev experience
- ✅ Need fast builds

**Use React Frontend when**:
- Simple prototype
- Don't need types
- Prefer JavaScript

**Use Go HTML Frontend when**:
- Server-side rendering required
- No JavaScript needed

## 📝 Scripts Reference

| Script | Command | Purpose |
|--------|---------|---------|
| **dev** | `npm run dev` | Start dev server |
| **build** | `npm run build` | Production build |
| **preview** | `npm run preview` | Preview build |
| **lint** | `npm run lint` | Run ESLint |
| **type-check** | `npm run type-check` | Check types |

## 🔗 Related Documentation

- [Material UI Setup](./MATERIAL-UI-SETUP.md)
- [Docker Deployment](../guides/DOCKER-DEPLOYMENT.md)
- [Kubernetes Deployment](../kubernetes/GETTING-STARTED.md)
- [Authentication Setup](./AUTHENTICATION-SETUP.md)

## ✅ Checklist

After setup, verify:

- [ ] TypeScript compiles without errors
- [ ] Development server starts on port 3001
- [ ] Can build for production
- [ ] Docker image builds successfully
- [ ] Can access frontend at http://localhost:3001
- [ ] API calls work
- [ ] Authentication works
- [ ] Products display correctly
- [ ] Theme looks good
- [ ] Navigation works

## 🎉 Success!

Your TypeScript frontend is ready! Access it at:

**http://localhost:3001**

Login with:
- Email: `admin@example.com`
- Password: `qwerty`

---

**Happy TypeScript Development! 🚀**

