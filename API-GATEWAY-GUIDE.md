# 🚪 API Gateway Setup Guide

**Add API Gateway to your application for routing, authentication, rate limiting, and more**

---

## 🎯 What is API Gateway?

An API Gateway sits between your clients (frontends) and your backend services, providing:

- ✅ **Single Entry Point** - One URL for all services
- ✅ **Authentication** - Validate tokens before reaching backend
- ✅ **Rate Limiting** - Prevent abuse and DDoS
- ✅ **CORS Handling** - Simplified cross-origin requests
- ✅ **Request/Response Transformation** - Modify data in flight
- ✅ **Logging & Monitoring** - Centralized access logs
- ✅ **Load Balancing** - Distribute traffic across multiple backends
- ✅ **SSL Termination** - Handle HTTPS in one place

---

## 🏗️ Architecture

### Without API Gateway (Old Setup)

```
Browser → Frontend:3007 (HTML/CSS/JS)
Browser → Backend:4001/api (JSON - direct)
Browser → AI Service:5001 (JSON - direct)
Browser → Support:6001 (JSON - direct)
```

**Issues:**
- Multiple endpoints to manage
- Different ports for frontend and API
- CORS configuration needed per service
- No centralized authentication
- No rate limiting
- Direct service exposure

### With API Gateway (Current Setup) ✅

```
Browser → API Gateway:8000 → Frontend:3000 (HTML/CSS/JS for /, /products, etc.)
                           → Backend:4001 (JSON for /api/*)
                           → AI Service:5001 (JSON for /ai/*)
                           → Support:6001 (JSON for /support/*)

ONE URL for EVERYTHING!
```

**Benefits:**
- ✅ Single URL: `http://localhost:8000` for frontend AND backend
- ✅ NO CORS issues (same origin!)
- ✅ Centralized auth, rate limiting, logging
- ✅ Services hidden behind gateway
- ✅ Easy to add new services
- ✅ Production-ready architecture

---

## 🔧 Implementation Options

### Option 1: Kong API Gateway (Recommended)

**What:** Open-source API Gateway based on Nginx

**Pros:**
- Production-ready
- Plugin ecosystem
- Easy configuration
- Docker-friendly
- Free and open-source

**Cons:**
- Requires PostgreSQL database
- More complex setup

### Option 2: Nginx as API Gateway

**What:** Configure Nginx as reverse proxy

**Pros:**
- Lightweight
- Simple configuration
- No external database
- Already familiar

**Cons:**
- Manual plugin development
- Less built-in features

### Option 3: Traefik

**What:** Modern cloud-native edge router

**Pros:**
- Auto-discovery
- Let's Encrypt built-in
- Dashboard UI
- Docker labels configuration

**Cons:**
- Learning curve
- Different paradigm

### Option 4: AWS API Gateway (Cloud)

**What:** Managed API Gateway service

**Pros:**
- Fully managed
- Integrates with Lambda
- Auto-scaling

**Cons:**
- AWS only
- Cost per request

---

## 🚀 Quick Start: Kong API Gateway

### Current Working Setup ✅

Your application is **already configured** with Kong API Gateway serving both frontend and backend through **ONE URL**: `http://localhost:8000`

### What's Configured

**Kong Services:**
1. **frontend** → `usualstore-react-frontend:3000`
   - Route: `/` (catch-all for HTML pages)
   - Serves: Home, products, login, cart pages

2. **backend-api** → `back-end:4001`
   - Route: `/api/*` (all API requests)
   - Serves: JSON data for products, users, orders

3. **support-service** → `support-service:6001`
   - Route: `/support/*`
   - Serves: Support ticket system

**Access Points:**
```bash
# Frontend (HTML/UI)
http://localhost:8000/              # Home page
http://localhost:8000/products      # Products page
http://localhost:8000/login         # Login page

# Backend API (JSON)
http://localhost:8000/api/widgets   # Get products
http://localhost:8000/api/users     # Get users
http://localhost:8000/api/orders    # Get orders

# Support Service
http://localhost:8000/support/*     # Support endpoints
```

### How to Test

```bash
# Test frontend (returns HTML)
curl http://localhost:8000/

# Test API (returns JSON)
curl http://localhost:8000/api/widgets

# Test in browser
# Open: http://localhost:8000
# You'll see the full React application!
```

### Adding to a New Project

If you need to add this to a new project:

**Step 1: Add to Terraform**

Update your `terraform/main.tf`:

```hcl
# Add API Gateway module
module "api_gateway" {
  source = "./modules/api-gateway"
  
  network_name   = docker_network.usualstore_network.name
  gateway_port   = 8000
  admin_port     = 8001
  backend_port   = var.backend_port
  
  depends_on = [
    module.backend_api,
    module.frontends
  ]
}

# Output gateway URL
output "api_gateway_url" {
  value       = module.api_gateway.gateway_url
  description = "API Gateway URL - use this for all requests"
}
```

**Step 2: Deploy with Terraform**

```bash
cd terraform
terraform init
terraform apply
```

**Step 3: Configure Routes**

```bash
# Add frontend service
curl -i -X POST http://localhost:8001/services/ \
  --data name=frontend \
  --data url='http://react-frontend:3000'

# Add frontend route (catch-all)
curl -i -X POST http://localhost:8001/services/frontend/routes \
  --data 'paths[]=/' \
  --data 'strip_path=false'

# Backend route is auto-configured in Terraform
```

**Step 4: Update Frontend**

Update your frontend to use relative URLs:

```javascript
// In src/services/api.js
const API_BASE_URL = '';  // Empty string = relative URLs!

// API calls now go to same domain
fetch('/api/widgets')  // → http://localhost:8000/api/widgets
```

**Step 5: Rebuild Frontend**

```bash
cd react-frontend
npm run build

# Rebuild Docker image
docker build -t react-frontend:latest .

# Restart container
terraform apply
```

---

## 📋 Kong Configuration

### View Current Configuration

```bash
# List all services
curl http://localhost:8001/services

# Expected output:
# - frontend (usualstore-react-frontend:3000)
# - backend-api (back-end:4001)
# - support-service (support-service:6001)

# List all routes
curl http://localhost:8001/routes

# Expected output:
# - /api/* → backend-api
# - /support/* → support-service
# - / → frontend (catch-all)

# List all plugins
curl http://localhost:8001/plugins

# Expected output:
# - rate-limiting (100/min, 1000/hour)
# - cors
# - file-log
```

### Test Your Routes

```bash
# Test frontend route (should return HTML)
curl http://localhost:8000/

# Test API route (should return JSON)
curl http://localhost:8000/api/widgets

# Check Kong headers (proves it's working)
curl -I http://localhost:8000/api/widgets | grep -i kong
# Should see: Via: kong/3.4.2
```

### Add New Service

```bash
# Add AI Assistant service
curl -i -X POST http://localhost:8001/services/ \
  --data name=ai-assistant \
  --data url='http://ai-assistant:5001'

# Add route
curl -i -X POST http://localhost:8001/services/ai-assistant/routes \
  --data 'paths[]=/ai' \
  --data 'strip_path=false'
```

### Add Rate Limiting

```bash
# Limit to 100 requests per minute
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=rate-limiting \
  --data config.minute=100 \
  --data config.hour=1000
```

### Add JWT Authentication

```bash
# Enable JWT plugin
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=jwt
```

### Add Request Logging

```bash
# Log all requests
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=file-log \
  --data config.path=/tmp/kong-access.log
```

---

## 🔒 Authentication Flow

### With API Gateway

```
1. User logs in → Backend generates JWT token
2. Frontend stores token
3. Frontend makes request with token:
   GET http://localhost:8000/api/orders
   Authorization: Bearer <token>
   
4. API Gateway validates token (JWT plugin)
5. If valid → forwards to backend
6. If invalid → returns 401 Unauthorized
```

### Configure JWT in Kong

```bash
# Create a consumer (user)
curl -i -X POST http://localhost:8001/consumers/ \
  --data username=testuser

# Add JWT credentials for consumer
curl -i -X POST http://localhost:8001/consumers/testuser/jwt \
  --data key=user123 \
  --data secret=secret123

# Enable JWT plugin on service
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=jwt
```

---

## 🎨 Alternative: Simple Nginx Gateway

If Kong is too complex, use Nginx:

### nginx.conf

```nginx
events {
    worker_connections 1024;
}

http {
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    
    # Upstream services
    upstream backend {
        server backend:4001;
    }
    
    upstream ai_service {
        server ai-assistant:5001;
    }
    
    server {
        listen 8000;
        server_name localhost;
        
        # Enable CORS
        add_header 'Access-Control-Allow-Origin' '*' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type' always;
        
        # Handle preflight
        if ($request_method = 'OPTIONS') {
            return 204;
        }
        
        # Backend API routes
        location /api {
            limit_req zone=api_limit burst=20 nodelay;
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
        
        # AI Assistant routes
        location /ai {
            limit_req zone=api_limit burst=5 nodelay;
            proxy_pass http://ai_service;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        # Health check
        location /health {
            access_log off;
            return 200 "API Gateway is healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
```

### Deploy Nginx Gateway

```yaml
# docker-compose-gateway.yml
version: '3.8'

services:
  gateway:
    image: nginx:alpine
    container_name: api-gateway
    ports:
      - "8000:8000"
    volumes:
      - ./nginx-gateway.conf:/etc/nginx/nginx.conf:ro
    networks:
      - usualstore_network
    depends_on:
      - backend
      - ai-assistant

networks:
  usualstore_network:
    external: true
```

```bash
docker-compose -f docker-compose-gateway.yml up -d
```

---

## ☁️ AWS API Gateway

### For Lambda Deployment

When deploying to AWS Lambda, use AWS API Gateway:

```hcl
# terraform-aws/api-gateway.tf

resource "aws_apigatewayv2_api" "usualstore" {
  name          = "usualstore-api"
  protocol_type = "HTTP"
  
  cors_configuration {
    allow_origins = ["https://usualstore.com"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["content-type", "authorization"]
    max_age       = 3600
  }
}

# Backend integration
resource "aws_apigatewayv2_integration" "backend" {
  api_id           = aws_apigatewayv2_api.usualstore.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.backend.invoke_arn
}

# Routes
resource "aws_apigatewayv2_route" "api_routes" {
  api_id    = aws_apigatewayv2_api.usualstore.id
  route_key = "ANY /api/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.backend.id}"
}

# Stage
resource "aws_apigatewayv2_stage" "prod" {
  api_id      = aws_apigatewayv2_api.usualstore.id
  name        = "prod"
  auto_deploy = true
  
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
      protocol       = "$context.protocol"
      responseLength = "$context.responseLength"
    })
  }
}

# Custom domain (optional)
resource "aws_apigatewayv2_domain_name" "api" {
  domain_name = "api.usualstore.com"
  
  domain_name_configuration {
    certificate_arn = aws_acm_certificate.api.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

output "api_gateway_url" {
  value = aws_apigatewayv2_stage.prod.invoke_url
}
```

---

## 📊 Monitoring & Logging

### Kong Logging

```bash
# View Kong logs
docker logs api-gateway

# View access logs
docker exec api-gateway tail -f /tmp/kong-access.log

# View metrics
curl http://localhost:8001/status
```

### Prometheus Integration

Kong supports Prometheus metrics:

```bash
# Enable Prometheus plugin
curl -i -X POST http://localhost:8001/plugins/ \
  --data name=prometheus
```

Access metrics at: `http://localhost:8001/metrics`

### Custom Logging

```bash
# Add custom logging plugin
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=http-log \
  --data config.http_endpoint=http://logserver:8080/logs
```

---

## 🔄 Multi-Tenant with API Gateway

### Per-Tenant Routing

Route requests based on subdomain or header:

```bash
# Building Shop route
curl -i -X POST http://localhost:8001/services/ \
  --data name=building-shop-backend \
  --data url='http://building-shop-backend:4001'

curl -i -X POST http://localhost:8001/services/building-shop-backend/routes \
  --data 'hosts[]=building-shop.usualstore.com' \
  --data 'paths[]=/api'

# Hardware Store route
curl -i -X POST http://localhost:8001/services/ \
  --data name=hardware-store-backend \
  --data url='http://hardware-store-backend:4002'

curl -i -X POST http://localhost:8001/services/hardware-store-backend/routes \
  --data 'hosts[]=hardware-store.usualstore.com' \
  --data 'paths[]=/api'
```

**Now:**
- `http://building-shop.usualstore.com/api` → Building Shop backend
- `http://hardware-store.usualstore.com/api` → Hardware Store backend

### Header-Based Routing

```bash
# Route based on X-Tenant-ID header
curl -i -X POST http://localhost:8001/services/backend-api/routes \
  --data 'headers.x-tenant-id=building-shop' \
  --data 'paths[]=/api'
```

---

## ⚡ Performance & Caching

### Enable Response Caching

```bash
# Cache GET requests for 5 minutes
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=proxy-cache \
  --data config.strategy=memory \
  --data config.content_type=application/json \
  --data config.cache_ttl=300
```

### Load Balancing

```bash
# Add multiple backend targets
curl -i -X POST http://localhost:8001/upstreams \
  --data name=backend-cluster

curl -i -X POST http://localhost:8001/upstreams/backend-cluster/targets \
  --data target=backend-1:4001 \
  --data weight=100

curl -i -X POST http://localhost:8001/upstreams/backend-cluster/targets \
  --data target=backend-2:4001 \
  --data weight=100

# Update service to use upstream
curl -i -X PATCH http://localhost:8001/services/backend-api \
  --data host=backend-cluster
```

---

## 🎯 Best Practices

### 1. Always Use API Gateway in Production

```
✅ DO: Frontend → API Gateway → Backend
❌ DON'T: Frontend → Backend (direct)
```

### 2. Enable Rate Limiting

```bash
# Prevent abuse
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=rate-limiting \
  --data config.minute=100
```

### 3. Use Authentication

```bash
# Validate JWT tokens
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=jwt
```

### 4. Enable CORS Properly

```bash
# Allow specific origins only
curl -i -X POST http://localhost:8001/services/backend-api/plugins \
  --data name=cors \
  --data 'config.origins=https://usualstore.com'
```

### 5. Monitor Everything

- Enable logging plugins
- Set up alerts for errors
- Monitor rate limit hits
- Track response times

---

## 📁 File Structure

```
usual_store/
├── terraform/
│   ├── main.tf                    ← Add API Gateway module here
│   └── modules/
│       └── api-gateway/           ← New module
│           ├── main.tf
│           ├── variables.tf
│           └── outputs.tf
│
├── nginx-gateway.conf             ← Alternative: Nginx config
├── docker-compose-gateway.yml     ← Alternative: Docker Compose
└── API-GATEWAY-GUIDE.md          ← This guide
```

---

## ✅ Deployment Checklist

### Local Setup

- [ ] Add API Gateway module to Terraform
- [ ] Run `terraform apply`
- [ ] Verify gateway is running: `curl http://localhost:8000/health`
- [ ] Update frontend to use gateway URL
- [ ] Test API requests through gateway
- [ ] Configure rate limiting
- [ ] Enable CORS
- [ ] Add logging

### Production Setup

- [ ] Use AWS API Gateway (for Lambda)
- [ ] Configure custom domain
- [ ] Enable authentication (JWT/OAuth)
- [ ] Set up CloudWatch logging
- [ ] Configure API keys (if needed)
- [ ] Enable caching for GET requests
- [ ] Set up alerts for errors
- [ ] Document API endpoints

---

## 🌐 Unified URL Architecture (Current Setup)

### How It Works

When you access `http://localhost:8000`, Kong API Gateway intelligently routes your request:

```
User Request → http://localhost:8000/products
           ↓
   Kong API Gateway checks path:
           ↓
   Does it match /api/* ? NO
   Does it match /support/* ? NO
   Default: Route to frontend
           ↓
   Frontend Container → Returns HTML page
```

```
User Request → http://localhost:8000/api/widgets
           ↓
   Kong API Gateway checks path:
           ↓
   Does it match /api/* ? YES!
           ↓
   Backend Container → Returns JSON data
```

### Benefits of Unified URL

**1. No CORS Issues**
```javascript
// Frontend at http://localhost:8000/
// API at http://localhost:8000/api/*
// Same origin = No CORS headers needed!

fetch('/api/widgets')  // Just works!
```

**2. Simpler Configuration**
```
Before: Configure CORS on every service
After:  Configure CORS once in Kong
```

**3. Production-Ready**
```
Development: http://localhost:8000
Production:  https://yourdomain.com
Just change the domain, architecture stays the same!
```

**4. Easy Deployment**
```
Single URL = Single domain name needed
No need for api.yourdomain.com, app.yourdomain.com, etc.
```

### Route Priority

Kong checks routes in order:

1. **High Priority**: `/api/*` → Backend API
2. **Medium Priority**: `/support/*` → Support Service
3. **Low Priority**: `/` → Frontend (catch-all)

This ensures API requests are handled first, then frontend pages.

---

## 🎯 Summary

### What You Get (Current Setup) ✅

✅ **Unified URL** - `http://localhost:8000` for frontend AND backend  
✅ **No CORS Issues** - Same origin = no cross-origin problems  
✅ **Rate Limiting** - 100 requests/min, 1000/hour automatically  
✅ **CORS Handling** - Configured once at gateway  
✅ **Request Logging** - All requests logged centrally  
✅ **Production Ready** - Kong 3.4 is battle-tested  
✅ **Easy to Scale** - Add services anytime  

### Quick Commands

```bash
# View gateway status
curl http://localhost:8001/status

# Test frontend (returns HTML)
curl http://localhost:8000/

# Test API (returns JSON)
curl http://localhost:8000/api/widgets

# View configuration
curl http://localhost:8001/services
curl http://localhost:8001/routes
curl http://localhost:8001/plugins

# Check if request went through Kong
curl -I http://localhost:8000/api/widgets | grep -i kong
```

### Current Configuration

**Services Configured:**
- ✅ frontend → usualstore-react-frontend:3000
- ✅ backend-api → back-end:4001
- ✅ support-service → support-service:6001

**Routes Configured:**
- ✅ `/` → Frontend (HTML pages)
- ✅ `/api/*` → Backend API (JSON)
- ✅ `/support/*` → Support Service

**Plugins Enabled:**
- ✅ rate-limiting (100/min, 1000/hour)
- ✅ cors (cross-origin enabled)
- ✅ file-log (request logging)

### Access Your Application

**Open in browser:**
```
http://localhost:8000
```

**You'll see:**
- Beautiful product catalog
- Working shopping cart
- User authentication
- All features working through ONE URL!

### Optional Enhancements

1. **Add JWT Authentication**
   ```bash
   curl -X POST http://localhost:8001/services/backend-api/plugins \
     --data name=jwt
   ```

2. **Add More Services**
   ```bash
   # Add AI Assistant
   curl -X POST http://localhost:8001/services/ \
     --data name=ai-assistant \
     --data url='http://ai-assistant:5001'
   
   curl -X POST http://localhost:8001/services/ai-assistant/routes \
     --data 'paths[]=/ai'
   ```

3. **Enable Response Caching**
   ```bash
   curl -X POST http://localhost:8001/services/backend-api/plugins \
     --data name=proxy-cache \
     --data config.cache_ttl=300
   ```

---

**Your application has a production-grade API Gateway with unified URL architecture!** 🚀

**Key Achievement:** Frontend and backend now run on ONE URL - just like GitHub, Amazon, and other major websites!

