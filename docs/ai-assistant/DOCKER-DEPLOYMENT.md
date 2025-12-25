# 🐳 AI Assistant - Docker Deployment Guide

Deploy the AI assistant using Docker and Docker Compose.

---

## 📦 What's Included

- **Dockerfile.ai-assistant** - Multi-stage build for optimal image size
- **docker-compose.yml** - Updated with AI assistant service
- **Environment variables** - Secure configuration

---

## 🚀 Quick Start

### **1. Set OpenAI API Key**

```bash
# Add to your .env file
echo "OPENAI_API_KEY=sk-your-actual-key-here" >> .env
echo "AI_ASSISTANT_PORT=8080" >> .env
```

### **2. Build and Run**

```bash
cd /Users/vkuzm/Projects/usual_store

# Build the AI assistant image
docker compose build ai-assistant

# Start all services (including AI assistant)
docker compose up -d

# Or start only AI assistant
docker compose up -d ai-assistant
```

### **3. Verify It's Running**

```bash
# Check logs
docker compose logs -f ai-assistant

# Test health check
curl http://localhost:8080/health
# Expected: OK

# Test chat endpoint
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-123",
    "message": "Hi! Help me choose a product"
  }'
```

---

## 📋 Architecture

```
Docker Compose Services:
┌────────────────────────────────────────────┐
│  front-end (port 4000)                     │
│  ├── Web application                       │
│  └── Connects to: back-end, database      │
├────────────────────────────────────────────┤
│  back-end (port 4001)                      │
│  ├── API server                            │
│  └── Connects to: database                │
├────────────────────────────────────────────┤
│  ai-assistant (port 8080) ⭐ NEW           │
│  ├── AI chat service                       │
│  ├── Connects to: database, OpenAI API    │
│  └── 2 replicas for HA (in K8s)           │
├────────────────────────────────────────────┤
│  invoice (port 5000)                       │
│  └── Invoice microservice                 │
├────────────────────────────────────────────┤
│  database (port 5433)                      │
│  ├── PostgreSQL 15                         │
│  └── Persistent volume                    │
└────────────────────────────────────────────┘
```

---

## 🔧 Configuration

### **Environment Variables**

Create/update `.env` file:

```bash
# AI Assistant
OPENAI_API_KEY=sk-your-actual-openai-key-here
AI_ASSISTANT_PORT=8080

# Existing variables
USUAL_STORE_PORT=4000
API_PORT=4001
INVOICE_PORT=5000
DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
```

### **Docker Compose Service**

```yaml
ai-assistant:
  build:
    context: .
    dockerfile: Dockerfile.ai-assistant
  environment:
    - PORT=${AI_ASSISTANT_PORT:-8080}
    - DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
    - OPENAI_API_KEY=${OPENAI_API_KEY}
  ports:
    - "${AI_ASSISTANT_PORT:-8080}:8080"
  networks:
    - usualstore_network
  depends_on:
    database:
      condition: service_healthy
  healthcheck:
    test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
    interval: 30s
    timeout: 3s
    retries: 3
```

---

## 🏗️ Build Process

### **Multi-Stage Dockerfile**

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-assistant ./cmd/ai-assistant-example/

# Stage 2: Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/ai-assistant .
EXPOSE 8080
CMD ["./ai-assistant"]
```

**Benefits:**
- ✅ Small image size (~20MB vs 300MB+)
- ✅ No source code in final image
- ✅ Only runtime dependencies
- ✅ Security best practices

---

## 📊 Docker Commands

### **Building**

```bash
# Build AI assistant image
docker compose build ai-assistant

# Build with no cache (force rebuild)
docker compose build --no-cache ai-assistant

# Build and tag for registry
docker build -f Dockerfile.ai-assistant \
  -t your-registry/usualstore-ai-assistant:v1.0 .
```

### **Running**

```bash
# Start all services
docker compose up -d

# Start only AI assistant
docker compose up -d ai-assistant

# Scale AI assistant (multiple instances)
docker compose up -d --scale ai-assistant=3

# Run in foreground (see logs)
docker compose up ai-assistant
```

### **Managing**

```bash
# View logs
docker compose logs -f ai-assistant

# View last 100 lines
docker compose logs --tail=100 ai-assistant

# Restart service
docker compose restart ai-assistant

# Stop service
docker compose stop ai-assistant

# Remove service
docker compose down ai-assistant
```

### **Debugging**

```bash
# Shell into running container
docker compose exec ai-assistant sh

# View container stats
docker compose stats ai-assistant

# Inspect container
docker compose ps ai-assistant

# View health status
docker inspect --format='{{.State.Health.Status}}' \
  $(docker compose ps -q ai-assistant)
```

---

## 🔍 Health Checks

### **Built-in Health Check**

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
  interval: 30s  # Check every 30 seconds
  timeout: 3s    # Timeout after 3 seconds
  retries: 3     # Retry 3 times before marking unhealthy
  start_period: 10s  # Grace period on startup
```

### **Manual Health Check**

```bash
# From host
curl http://localhost:8080/health

# From another container
docker compose exec front-end wget -qO- http://ai-assistant:8080/health

# Check Docker health status
docker compose ps | grep ai-assistant
```

---

## 🌐 Networking

### **Internal Communication**

```bash
# From front-end to AI assistant
http://ai-assistant:8080/api/ai/chat

# From back-end to AI assistant
http://ai-assistant:8080/api/ai/stats

# AI assistant to database
postgres://postgres:password@database:5432/usualstore
```

### **External Access**

```bash
# From host machine
http://localhost:8080/api/ai/chat

# IPv6 (if enabled)
http://[::1]:8080/api/ai/chat

# From other machines on network
http://<your-ip>:8080/api/ai/chat
```

---

## 🔒 Security Best Practices

### **1. Secure API Keys**

```bash
# Never commit .env to git
echo ".env" >> .gitignore

# Use Docker secrets (Swarm mode)
echo "sk-your-key" | docker secret create openai_api_key -

# Use in compose
secrets:
  - openai_api_key
```

### **2. Network Isolation**

```yaml
networks:
  usualstore_network:
    driver: bridge
    internal: false  # Set to true to block external access
```

### **3. Resource Limits**

```yaml
ai-assistant:
  deploy:
    resources:
      limits:
        cpus: '0.5'
        memory: 512M
      reservations:
        cpus: '0.25'
        memory: 256M
```

---

## 📈 Monitoring

### **View Logs**

```bash
# Real-time logs
docker compose logs -f ai-assistant

# Logs with timestamps
docker compose logs -f --timestamps ai-assistant

# Filter logs
docker compose logs ai-assistant | grep ERROR
```

### **Resource Usage**

```bash
# CPU and memory usage
docker compose stats ai-assistant

# Container processes
docker compose top ai-assistant
```

### **Analytics Endpoint**

```bash
# Get AI usage stats
curl http://localhost:8080/api/ai/stats?days=7

# Parse with jq
curl -s http://localhost:8080/api/ai/stats?days=7 | jq .
```

---

## 🐛 Troubleshooting

### **Issue: Container Won't Start**

```bash
# Check logs
docker compose logs ai-assistant

# Common causes:
# 1. Missing OPENAI_API_KEY
echo "OPENAI_API_KEY=sk-..." >> .env

# 2. Database not ready
docker compose up -d database
# Wait 10 seconds, then start AI assistant
docker compose up -d ai-assistant

# 3. Port already in use
# Change port in .env
echo "AI_ASSISTANT_PORT=8081" >> .env
```

### **Issue: Health Check Failing**

```bash
# Check if service is responding
docker compose exec ai-assistant wget -qO- http://localhost:8080/health

# Check logs for errors
docker compose logs --tail=50 ai-assistant

# Restart service
docker compose restart ai-assistant
```

### **Issue: "connection refused to database"**

```bash
# Verify database is running
docker compose ps database

# Check database health
docker compose exec database pg_isready -U postgres

# Verify network connectivity
docker compose exec ai-assistant ping database
```

### **Issue: High Memory Usage**

```bash
# Check memory usage
docker compose stats ai-assistant

# Restart to free memory
docker compose restart ai-assistant

# Set memory limit in docker-compose.yml
deploy:
  resources:
    limits:
      memory: 512M
```

---

## 🚀 Production Deployment

### **1. Build Production Image**

```bash
# Build optimized image
docker build -f Dockerfile.ai-assistant \
  -t your-registry.com/usualstore-ai-assistant:1.0.0 .

# Push to registry
docker push your-registry.com/usualstore-ai-assistant:1.0.0
```

### **2. Use Docker Swarm (Production)**

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml usualstore

# Scale AI assistant
docker service scale usualstore_ai-assistant=3

# Update service (zero downtime)
docker service update \
  --image your-registry.com/usualstore-ai-assistant:1.0.1 \
  usualstore_ai-assistant
```

### **3. Production Checklist**

- [ ] Use specific image tags (not `latest`)
- [ ] Set resource limits
- [ ] Configure log rotation
- [ ] Enable health checks
- [ ] Use secrets for API keys
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure backups
- [ ] Enable HTTPS
- [ ] Set up rate limiting

---

## 📊 Performance Optimization

### **1. Image Size**

```bash
# Check image size
docker images | grep ai-assistant

# Optimize with multi-stage build (already done)
# Result: ~20MB instead of 300MB+
```

### **2. Build Cache**

```bash
# Use BuildKit for faster builds
DOCKER_BUILDKIT=1 docker compose build ai-assistant

# Prune unused images
docker image prune -a
```

### **3. Container Optimization**

```yaml
ai-assistant:
  # Limit restart attempts
  restart: on-failure:3
  
  # Optimize health checks
  healthcheck:
    interval: 60s  # Less frequent checks
    
  # Set resource limits
  deploy:
    resources:
      limits:
        cpus: '0.5'
        memory: 512M
```

---

## 🎯 Next Steps

1. ✅ Docker setup complete
2. Test AI assistant locally
3. Deploy to staging
4. Monitor performance
5. Deploy to production

**For Kubernetes deployment**, see: [KUBERNETES-AI-DEPLOYMENT.md](KUBERNETES-AI-DEPLOYMENT.md)

---

## 📚 Resources

- Docker Compose Docs: https://docs.docker.com/compose/
- Multi-Stage Builds: https://docs.docker.com/build/building/multi-stage/
- Health Checks: https://docs.docker.com/engine/reference/builder/#healthcheck

---

**Your AI assistant is now containerized!** 🎉

