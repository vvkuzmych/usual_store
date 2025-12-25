# Usual Store - TypeScript Frontend

Modern **TypeScript** frontend for Usual Store built with **Vite**, **React 18**, and **Material UI**.

## 🚀 Features

- ✨ **Full TypeScript** support with strict type checking
- ⚡ **Vite** for lightning-fast development and builds
- 🎨 **Material UI (MUI)** with custom blue theme
- 🔐 **Authentication** with JWT token management
- 💳 **Stripe integration** ready
- 📱 **Responsive design** optimized for all devices
- 🎯 **Type-safe API** calls with axios
- 🔄 **React Router** for client-side navigation

## 📦 Tech Stack

- **React** 18.2
- **TypeScript** 5.3
- **Vite** 5.0
- **Material UI** 5.14
- **React Router** 6.20
- **Axios** for HTTP requests
- **Emotion** for styled components

## 🏃‍♂️ Running Locally

### Development Mode (with hot reload):

```bash
cd typescript-frontend
npm install
npm run dev
```

Open [http://localhost:3001](http://localhost:3001)

### Build for Production:

```bash
npm run build
npm run preview
```

### Type Checking:

```bash
npm run type-check
```

### Linting:

```bash
npm run lint
```

## 🐳 Docker

### Build the Docker image:

```bash
docker build -t usual_store-typescript-frontend ./typescript-frontend
```

### Run with Docker:

```bash
docker run -p 3001:3001 usual_store-typescript-frontend
```

### Run with Docker Compose:

```bash
# Start TypeScript frontend + backend
docker compose --profile typescript-frontend up -d

# Or use the helper script
./scripts/start-typescript.sh
```

Access the app at:
- **TypeScript Frontend**: http://localhost:3001
- **Backend API**: http://localhost:4001
- **AI Assistant**: http://localhost:8080

## ☸️ Kubernetes

Deploy to Kubernetes:

```bash
# Build and load image
docker build -t usual_store-typescript-frontend:latest ./typescript-frontend

# For Docker Desktop Kubernetes
docker save usual_store-typescript-frontend:latest | docker load

# Deploy to Kubernetes
kubectl apply -f k8s/01-namespace.yaml
kubectl apply -f k8s/17-typescript-frontend-deployment.yaml

# Check status
kubectl get pods -n usual-store
kubectl get svc -n usual-store

# Access the service
kubectl port-forward -n usual-store svc/typescript-frontend 3001:3001
```

Then open http://localhost:3001

## 📁 Project Structure

```
typescript-frontend/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   ├── pages/            # Page components
│   │   ├── Home.tsx
│   │   ├── Products.tsx
│   │   ├── ProductDetail.tsx
│   │   ├── Login.tsx
│   │   └── Cart.tsx
│   ├── services/         # API services
│   │   └── api.ts
│   ├── contexts/         # React contexts
│   │   └── AuthContext.tsx
│   ├── types/            # TypeScript type definitions
│   │   └── index.ts
│   ├── theme.ts          # Material UI theme
│   ├── App.tsx           # Main app component
│   └── main.tsx          # Entry point
├── public/               # Static assets
├── Dockerfile            # Multi-stage Docker build
├── nginx.conf            # Nginx configuration
├── tsconfig.json         # TypeScript config
├── vite.config.ts        # Vite configuration
└── package.json          # Dependencies
```

## 🎨 Theme

The TypeScript frontend uses a custom **blue gradient theme** to differentiate from the purple React frontend:

- **Primary**: Blue (#2196f3)
- **Secondary**: Cyan (#00bcd4)
- **Typography**: Roboto font family
- **Components**: Custom styled MUI components

## 🔐 Authentication

Demo credentials:
- **Email**: admin@example.com
- **Password**: qwerty

The frontend stores JWT tokens in localStorage and includes them in API requests automatically.

## 🌐 API Integration

The TypeScript frontend connects to the same Go backend API:

- **Products API**: `GET /api/products`, `GET /api/product/:id`
- **Auth API**: `POST /api/authenticate`
- **Payment API**: `POST /api/payment-intent`, `POST /api/charge`
- **AI Assistant API**: `POST /api/ai/chat`, `POST /api/ai/feedback`, `GET /api/ai/stats`

All API calls are type-safe using TypeScript interfaces.

## 📊 Comparison with React Frontend

| Feature | React Frontend (Port 3000) | TypeScript Frontend (Port 3001) |
|---------|---------------------------|--------------------------------|
| **Language** | JavaScript | TypeScript |
| **Build Tool** | Create React App | Vite |
| **Theme** | Purple gradient | Blue gradient |
| **Type Safety** | Runtime only | Compile-time |
| **Dev Server** | Webpack | Vite (faster) |
| **Bundle Size** | Larger | Smaller |
| **Hot Reload** | Slower | Faster |

## 🆚 Ports

- **3000**: React frontend
- **3001**: TypeScript frontend ⬅️ This one!
- **4000**: Go HTML frontend
- **4001**: Backend API
- **5000**: Invoice microservice
- **8080**: AI Assistant

## 🚀 Deployment

### Production Build:

```bash
npm run build
```

The `dist/` folder contains optimized static files ready for deployment.

### Environment Variables:

- `VITE_API_URL`: Backend API URL (default: proxied through Nginx)

### Nginx:

The production Docker image uses Nginx to:
- Serve static files
- Proxy API requests to backend
- Handle React Router (SPA) routes

## 📝 License

© 2025 Usual Store. All rights reserved.

## 🤝 Contributing

This is a demo TypeScript frontend showcasing modern React development practices with full type safety.

---

**Happy Coding! 🎉**

