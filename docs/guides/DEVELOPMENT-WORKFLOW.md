# 🔄 Development Workflow Guide

Complete guide for updating your application after making code changes.

---

## 📋 Table of Contents

1. [Quick Reference](#quick-reference)
2. [Terraform + Docker Workflow](#terraform-docker-workflow)
3. [Terraform + Kubernetes Workflow](#terraform-kubernetes-workflow)
4. [Hot Reload vs Manual Rebuild](#hot-reload-vs-manual-rebuild)
5. [Best Practices](#best-practices)
6. [Common Scenarios](#common-scenarios)
7. [Troubleshooting](#troubleshooting)

---

## 🚀 Quick Reference

### For Docker Deployment (Recommended for Development)

```bash
# After changing frontend code:
cd terraform
terraform apply    # Automatically rebuilds changed containers

# Or rebuild specific service:
docker-compose build backend
docker-compose up -d backend

# View logs:
docker-compose logs -f backend
```

### For Kubernetes Deployment

```bash
# After changing any code:
cd terraform-k8s
./MINIKUBE-DEPLOY.sh    # Rebuilds everything and redeploys

# Or manual rebuild:
eval $(minikube docker-env)
docker-compose build
kubectl rollout restart deployment/backend -n usualstore
```

---

## 🐳 Terraform + Docker Workflow

### How It Works

**Docker Compose watches for image changes and automatically rebuilds when you run `terraform apply` or `docker-compose up`.**

### Development Workflow

#### Step 1: Make Your Code Changes

Edit any file in your project:
```bash
# Examples:
vim cmd/api/handlers.go           # Backend change
vim react-frontend/src/App.jsx    # Frontend change
vim typescript-frontend/src/App.tsx
```

#### Step 2: Rebuild and Restart

**Option A: Using Terraform (Recommended)**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
terraform apply
```

Terraform will:
- Detect which containers need rebuilding
- Rebuild only changed services
- Restart affected containers
- Leave unchanged services running

**Option B: Using Docker Compose Directly (Faster for single service)**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Rebuild specific service
docker-compose build backend
docker-compose up -d backend

# Or rebuild all
docker-compose build
docker-compose up -d
```

**Option C: Rebuild Without Cache (if having issues)**

```bash
# Force complete rebuild
docker-compose build --no-cache backend
docker-compose up -d backend
```

#### Step 3: Verify Changes

```bash
# Check if container restarted
docker ps | grep backend

# View logs
docker-compose logs -f backend

# Test the change
curl http://localhost:4001/health
# Or visit in browser
```

---

## ☸️ Terraform + Kubernetes Workflow

### How It Works

**Kubernetes does NOT automatically detect code changes. You must rebuild images and restart pods manually.**

### Development Workflow

#### Method 1: Automatic Script (Easiest)

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform-k8s
./MINIKUBE-DEPLOY.sh
```

This script automatically:
1. Configures Minikube Docker environment
2. Rebuilds all images
3. Redeploys to Kubernetes
4. Shows status

**Use this when:** You want everything updated automatically.

#### Method 2: Rebuild Specific Service (Faster)

```bash
# Step 1: Configure Minikube Docker
eval $(minikube docker-env)

# Step 2: Rebuild specific service
cd /Users/vkuzm/Projects/UsualStore/usual_store
docker-compose build backend

# Step 3: Restart deployment in Kubernetes
kubectl rollout restart deployment/backend -n usualstore

# Step 4: Watch the rollout
kubectl rollout status deployment/backend -n usualstore

# Step 5: Verify
kubectl get pods -n usualstore
```

**Use this when:** You changed only one service and want faster updates.

#### Method 3: Manual Terraform Reapply

```bash
# Step 1: Configure Minikube Docker
eval $(minikube docker-env)

# Step 2: Rebuild images
cd /Users/vkuzm/Projects/UsualStore/usual_store
docker-compose build

# Step 3: Reapply Terraform
cd terraform-k8s/local
terraform apply

# Note: This recreates resources, slower than rollout restart
```

**Use this when:** You changed Kubernetes manifests or Terraform configuration.

---

## 🔥 Hot Reload vs Manual Rebuild

### What Has Hot Reload?

**✅ Frontend (Development Mode Only)**

If you run frontends in development mode with `npm start`:
```bash
cd react-frontend
npm start    # Hot reload enabled on localhost:3000
```

Changes to `.jsx/.tsx` files reload automatically in browser.

**❌ Does NOT work in Docker/Kubernetes** - You must rebuild.

**❌ Backend (Go)**

Go applications require recompilation. No hot reload in production builds.

**For development hot reload in Go:**
```bash
# Install air (Go hot reload tool)
go install github.com/cosmtrek/air@latest

# Run backend with hot reload
cd cmd/api
air
```

### Summary Table

| Component | Hot Reload | Needs Rebuild |
|-----------|------------|---------------|
| React Frontend (npm start) | ✅ Yes | ❌ No |
| React Frontend (Docker) | ❌ No | ✅ Yes |
| TypeScript Frontend (Docker) | ❌ No | ✅ Yes |
| Backend Go (Docker) | ❌ No | ✅ Yes |
| Backend Go (with air) | ✅ Yes | ❌ No |
| Database schema | N/A | Run migration |

---

## 🎯 Best Practices

### For Active Development

**Use Docker (not Kubernetes):**
```bash
# Start services
cd terraform
terraform apply

# After code changes:
docker-compose build <service>
docker-compose up -d <service>
```

**Why?**
- Faster rebuild
- Easier debugging
- Direct access to logs
- No image transfer needed

### For Testing Kubernetes Features

**Use Kubernetes only when you need to test:**
- Scaling (multiple replicas)
- Self-healing
- Service discovery
- Resource limits
- Rolling updates

### Recommended Development Setup

**Terminal 1: Watch Logs**
```bash
# Docker
docker-compose logs -f backend frontend

# Kubernetes
kubectl logs -f deployment/backend -n usualstore
```

**Terminal 2: Run Commands**
```bash
# Make changes, rebuild, test
vim cmd/api/handlers.go
docker-compose build backend
docker-compose up -d backend
curl http://localhost:4001/health
```

**Terminal 3: Interactive Shell (if needed)**
```bash
# Docker
docker exec -it usualstore-back-end sh

# Kubernetes
kubectl exec -it <pod-name> -n usualstore -- sh
```

---

## 📝 Common Scenarios

### Scenario 1: Changed Backend API Handler

**What changed:** Modified `cmd/api/handlers.go`

**Docker:**
```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store
docker-compose build backend
docker-compose up -d backend
docker-compose logs -f backend
```

**Kubernetes:**
```bash
eval $(minikube docker-env)
docker-compose build backend
kubectl rollout restart deployment/backend -n usualstore
kubectl logs -f deployment/backend -n usualstore
```

### Scenario 2: Changed React Frontend Component

**What changed:** Modified `react-frontend/src/components/Header.jsx`

**Docker:**
```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store
docker-compose build react-frontend
docker-compose up -d react-frontend
```

**Kubernetes:**
```bash
eval $(minikube docker-env)
docker-compose build react-frontend
kubectl rollout restart deployment/frontend -n usualstore
```

### Scenario 3: Changed Database Schema

**What changed:** Created new migration file

**Apply migration (works for both Docker and Kubernetes):**
```bash
# Docker
docker exec -it usualstore-database psql -U postgres -d usualstore -f /migrations/your_migration.sql

# Or rebuild with migrations
cd terraform
terraform apply
```

**Kubernetes:**
```bash
# Run migration manually
kubectl exec -it database-0 -n usualstore -- psql -U postgres -d usualstore -f /migrations/your_migration.sql

# Or rebuild database pod (data loss!)
kubectl delete pod database-0 -n usualstore
```

### Scenario 4: Changed Environment Variables

**What changed:** Modified `.env` or ConfigMap

**Docker:**
```bash
# Edit terraform/terraform.tfvars or docker-compose.yml
cd terraform
terraform apply    # Recreates containers with new env vars
```

**Kubernetes:**
```bash
# Edit k8s/02-configmap.yaml or k8s/03-secrets.yaml
cd terraform-k8s/local
terraform apply
# Then restart deployments to pick up changes
kubectl rollout restart deployment --all -n usualstore
```

### Scenario 5: Added New Dependency

**Backend (Go):**
```bash
# Added new import in Go code
go mod tidy                          # Update go.mod
docker-compose build backend         # Rebuild
docker-compose up -d backend         # Restart
```

**Frontend (npm):**
```bash
# Added new package
cd react-frontend
npm install <package>                # Updates package.json & package-lock.json
cd ..
docker-compose build react-frontend  # Rebuild
docker-compose up -d react-frontend  # Restart
```

### Scenario 6: Changed Nginx Configuration

**What changed:** Modified `react-frontend/nginx.conf`

**Docker:**
```bash
docker-compose build react-frontend
docker-compose up -d react-frontend
```

**Kubernetes:**
```bash
eval $(minikube docker-env)
docker-compose build react-frontend
kubectl rollout restart deployment/frontend -n usualstore
```

---

## 🔍 Viewing Changes

### Check if Container/Pod Restarted

**Docker:**
```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.CreatedAt}}"
```

**Kubernetes:**
```bash
kubectl get pods -n usualstore -o wide
# Look at AGE and RESTARTS columns
```

### View Logs

**Docker:**
```bash
# Single service
docker-compose logs -f backend

# Multiple services
docker-compose logs -f backend frontend

# All services
docker-compose logs -f

# Last 50 lines
docker-compose logs --tail=50 backend
```

**Kubernetes:**
```bash
# Single pod
kubectl logs -f <pod-name> -n usualstore

# Deployment (all replicas)
kubectl logs -f deployment/backend -n usualstore

# Last 50 lines
kubectl logs <pod-name> -n usualstore --tail=50

# Previous pod (if crashed)
kubectl logs <pod-name> -n usualstore --previous
```

### Verify Changes Applied

**Docker:**
```bash
# Check environment variable
docker exec usualstore-back-end env | grep API_PORT

# Check file exists
docker exec usualstore-back-end ls -la /app/

# Run command
docker exec usualstore-back-end curl http://localhost:4001/health
```

**Kubernetes:**
```bash
# Check environment variable
kubectl exec <pod-name> -n usualstore -- env | grep API_PORT

# Check file exists
kubectl exec <pod-name> -n usualstore -- ls -la /app/

# Run command
kubectl exec <pod-name> -n usualstore -- curl http://localhost:4001/health
```

---

## ⚡ Fast Development Tips

### 1. Use Makefile Commands

Create a `Makefile` in project root:

```makefile
.PHONY: rebuild-backend rebuild-frontend logs restart-all

rebuild-backend:
	docker-compose build backend
	docker-compose up -d backend
	docker-compose logs -f backend

rebuild-frontend:
	docker-compose build react-frontend
	docker-compose up -d react-frontend

logs:
	docker-compose logs -f backend frontend

restart-all:
	cd terraform && terraform apply

# Kubernetes versions
k8s-rebuild-backend:
	eval $$(minikube docker-env) && docker-compose build backend
	kubectl rollout restart deployment/backend -n usualstore
	kubectl logs -f deployment/backend -n usualstore
```

Then just run:
```bash
make rebuild-backend
make rebuild-frontend
make logs
```

### 2. Use Shell Aliases

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# Docker aliases
alias dps='docker ps --format "table {{.Names}}\t{{.Status}}"'
alias dlogs='docker-compose logs -f'
alias dbuild='docker-compose build'
alias dup='docker-compose up -d'

# Kubernetes aliases
alias kgp='kubectl get pods -n usualstore'
alias kgpw='kubectl get pods -n usualstore --watch'
alias kl='kubectl logs -f'
alias kr='kubectl rollout restart deployment'
alias ke='kubectl exec -it'

# Rebuild shortcuts
alias rebuild-be='docker-compose build backend && docker-compose up -d backend'
alias rebuild-fe='docker-compose build react-frontend && docker-compose up -d react-frontend'
```

Then reload:
```bash
source ~/.zshrc
```

### 3. Watch for Changes (Advanced)

Install `entr` for automatic rebuilds:

```bash
# Install
brew install entr

# Watch Go files and rebuild backend
find cmd internal -name "*.go" | entr -r make rebuild-backend

# Watch React files and rebuild frontend
find react-frontend/src -name "*.jsx" -o -name "*.js" | entr -r make rebuild-frontend
```

---

## 🐛 Troubleshooting

### Changes Not Appearing

**Problem:** Rebuilt but changes don't show

**Solutions:**

1. **Clear Docker build cache:**
   ```bash
   docker-compose build --no-cache backend
   docker-compose up -d backend
   ```

2. **Verify image was rebuilt:**
   ```bash
   docker images | grep usual_store-back-end
   # Check CREATED timestamp
   ```

3. **Check correct container is running:**
   ```bash
   docker ps | grep backend
   # Verify container ID and image ID match
   ```

4. **For Kubernetes, delete old images:**
   ```bash
   eval $(minikube docker-env)
   docker images | grep usual_store
   # Delete old images if multiple exist
   docker rmi <old-image-id>
   ```

### Container Won't Start After Rebuild

**Problem:** Container immediately exits or crashes

**Solutions:**

1. **Check logs:**
   ```bash
   docker-compose logs backend
   # or
   kubectl logs <pod-name> -n usualstore
   ```

2. **Check syntax errors:**
   ```bash
   # For Go
   cd cmd/api
   go build -o /dev/null .
   
   # For React/TypeScript
   cd react-frontend
   npm run build
   ```

3. **Verify dependencies:**
   ```bash
   # Go
   go mod verify
   
   # npm
   cd react-frontend
   npm ci
   ```

### "Permission Denied" Errors

**Problem:** Can't access files in container

**Solution:**
```bash
# Check file permissions in image
docker run --rm usual_store-back-end:latest ls -la /app/

# If needed, add to Dockerfile:
RUN chmod +x /app/myapp
```

### Kubernetes Pod Stuck in "ImagePullBackOff"

**Problem:** Pod can't find image

**Solution:**
```bash
# Verify you're using Minikube's Docker
eval $(minikube docker-env)

# Rebuild image
docker-compose build backend

# Verify image exists
docker images | grep usual_store-back-end

# Restart deployment
kubectl rollout restart deployment/backend -n usualstore
```

---

## 📊 Workflow Comparison

| Aspect | Docker | Kubernetes |
|--------|--------|------------|
| **Rebuild Speed** | ⚡ Fast (5-30s) | 🐢 Slower (1-3 min) |
| **Command** | `docker-compose build` | `docker-compose build` + `kubectl rollout` |
| **Auto-restart** | ✅ Yes (`up -d`) | ❌ Manual rollout |
| **Log Access** | ✅ Easy | ✅ Easy |
| **Best for** | Development | Testing K8s features |
| **Complexity** | ⭐ Low | ⭐⭐⭐ Medium |

---

## 🎯 Recommendations

### For Daily Development

**Use Docker:**
```bash
# Edit code
vim cmd/api/handlers.go

# Rebuild & restart
docker-compose build backend
docker-compose up -d backend

# Check logs
docker-compose logs -f backend
```

**Reason:** Faster, simpler, better for rapid iteration.

### For Testing Before Production

**Use Kubernetes:**
```bash
# After changes are working in Docker
cd terraform-k8s
./MINIKUBE-DEPLOY.sh

# Test scaling, health checks, etc.
```

**Reason:** Verify behavior in production-like environment.

### For Production Deployment

**Use Terraform + Kubernetes (AWS EKS):**
```bash
# When ready to deploy to AWS
cd terraform-k8s/aws-future
# Uncomment resources in main.tf
terraform apply
```

---

## 📚 Quick Command Reference

### Docker Development

```bash
# Build specific service
docker-compose build backend

# Build all services
docker-compose build

# Start/restart services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f backend

# Rebuild without cache
docker-compose build --no-cache backend
```

### Kubernetes Development

```bash
# Configure Minikube Docker
eval $(minikube docker-env)

# Build images
docker-compose build

# Restart deployment
kubectl rollout restart deployment/backend -n usualstore

# Watch rollout
kubectl rollout status deployment/backend -n usualstore

# View logs
kubectl logs -f deployment/backend -n usualstore

# Full redeploy
cd terraform-k8s && ./MINIKUBE-DEPLOY.sh
```

---

## 💡 Pro Tips

1. **Keep Docker deployment running during development** - Use it as your primary dev environment

2. **Test in Kubernetes before deploying to production** - Catch K8s-specific issues early

3. **Use hot reload for rapid UI development:**
   ```bash
   cd react-frontend
   npm start    # Development server with hot reload
   ```

4. **Create shell scripts for common tasks:**
   ```bash
   # rebuild-all.sh
   #!/bin/bash
   docker-compose build
   docker-compose up -d
   docker-compose logs -f
   ```

5. **Use VS Code Dev Containers** - Edit and run code inside containers with instant feedback

6. **Monitor resource usage:**
   ```bash
   docker stats
   # or in Kubernetes:
   kubectl top pods -n usualstore
   ```

---

## 🔗 Related Documentation

- [KUBERNETES-DASHBOARD-GUIDE.md](KUBERNETES-DASHBOARD-GUIDE.md) - View your pods visually
- [../terraform/README.md](../../terraform/README.md) - Docker deployment details
- [../terraform-k8s/README.md](../../terraform-k8s/README.md) - Kubernetes deployment details

---

**Happy coding!** 🚀 Remember: Docker for development, Kubernetes for production testing!

