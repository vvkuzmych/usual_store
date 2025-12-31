# UsualStore - Quick Start Guide

## 🚀 First Time Setup (30 minutes)

### 1. Install Prerequisites
```bash
# Core requirements
brew install docker terraform node go

# Migration tool
brew install gobuffalo/tap/pop  # For database migrations (soda)

# PostgreSQL client (optional but useful)
brew install libpq  # Just the client, database runs in Docker
```

**Database**: PostgreSQL 15 (runs in Docker - no separate installation needed)

### 2. Get API Keys
- OpenAI: https://platform.openai.com/api-keys
- Stripe: https://dashboard.stripe.com/apikeys
- Mailtrap: https://mailtrap.io/

### 3. Configure
```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your API keys
```

### 4. Start Database & Run Migrations ⚠️ IMPORTANT
```bash
# Start database first
terraform init
terraform apply -target=module.database.docker_container.database -auto-approve

# Wait for database
sleep 10

# Run migrations (REQUIRED!)
cd ..
make migrate

# Start everything else
cd terraform
terraform apply -auto-approve
```

### 5. Access
- Main App: http://localhost:8000
- React Frontend: http://localhost:3007

---

## 📋 Daily Development Commands

### Start/Stop Application
```bash
# Start everything
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
terraform apply -auto-approve

# Stop everything
terraform destroy -auto-approve
```

### After Code Changes

#### Frontend Changes (React/CSS/JS)
```bash
terraform apply \
  -target=module.frontends.docker_image.react_frontend \
  -target=module.frontends.docker_container.react_frontend \
  -auto-approve

# Clear browser cache: Cmd+Shift+R (Mac)
```

#### Backend Changes (Go)
```bash
terraform apply \
  -target=module.backend_api.docker_image.backend \
  -target=module.backend_api.docker_container.backend \
  -auto-approve
```

#### AI Assistant Changes
```bash
# Clean rebuild
docker stop usualstore-ai-assistant
docker rm usualstore-ai-assistant
docker rmi usual_store-ai-assistant:latest

docker build --no-cache -f Dockerfile.ai-assistant -t usual_store-ai-assistant:latest .

terraform apply \
  -target=module.ai_assistant.docker_container.ai_assistant \
  -auto-approve
```

---

## 🔍 Debugging

### Check Status
```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### View Logs
```bash
docker logs -f usualstore-react-frontend
docker logs -f usualstore-back-end
docker logs -f usualstore-ai-assistant
```

### Test APIs
```bash
# Health check
curl http://localhost:8000/health

# Get products
curl http://localhost:8000/api/products

# Test AI
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"session_id":"test","message":"Hello"}'
```

---

## 🆘 Common Issues

### Port in Use
```bash
lsof -i :3007
kill -9 <PID>
```

### Container Won't Start
```bash
docker logs usualstore-<service-name>
docker restart usualstore-<service-name>
```

### Frontend Not Updating
```bash
# Hard refresh: Cmd+Shift+R
# Or rebuild without cache
cd react-frontend
docker build --no-cache -t usual_store-react-frontend:latest .
```

### Clean Everything
```bash
docker stop $(docker ps -q --filter "name=usualstore")
docker rm $(docker ps -aq --filter "name=usualstore")
docker rmi $(docker images 'usual_store*' -q)
terraform apply -auto-approve
```

---

## 📚 Full Documentation

See `SETUP-GUIDE.md` for complete documentation.

