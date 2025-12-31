# Development Workflow - Decision Tree

## 🔄 When to Rebuild What?

```
┌─────────────────────────────────────────────────────────────────────┐
│                    What Did You Change?                              │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────────────┐
        │                                                   │
        ▼                                                   ▼
┌──────────────────┐                           ┌──────────────────────┐
│ Frontend Files?  │                           │   Backend Files?     │
│ (.jsx, .css)     │                           │   (.go files)        │
└──────────────────┘                           └──────────────────────┘
        │                                                   │
        ▼                                                   ▼
┌──────────────────────────────────────┐      ┌──────────────────────────┐
│ Rebuild Frontend Container:           │      │ Rebuild Backend Container:│
│                                        │      │                          │
│ terraform apply \                      │      │ terraform apply \        │
│   -target=module.frontends.\           │      │   -target=module.\       │
│   docker_image.react_frontend \        │      │   backend_api.\          │
│   -target=module.frontends.\           │      │   docker_image.backend \ │
│   docker_container.react_frontend \    │      │   -auto-approve          │
│   -auto-approve                        │      │                          │
│                                        │      │ ⏱️  Time: 3-5 minutes    │
│ ⏱️  Time: 5-10 minutes                │      └──────────────────────────┘
│                                        │
│ ⚠️  IMPORTANT:                         │
│ Clear browser cache: Cmd+Shift+R      │
└──────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────────────┐
        │                                                   │
        ▼                                                   ▼
┌──────────────────────┐                      ┌──────────────────────┐
│ AI Assistant Files?  │                      │ Database Migration?   │
│ (internal/ai/*)      │                      │ (.sql files)         │
└──────────────────────┘                      └──────────────────────┘
        │                                                   │
        ▼                                                   ▼
┌────────────────────────────────────┐       ┌──────────────────────────┐
│ Clean Rebuild AI Assistant:        │       │ Restart Database:        │
│                                    │       │                          │
│ docker stop usualstore-ai-assistant│       │ docker restart \         │
│ docker rm usualstore-ai-assistant  │       │   usualstore-database    │
│ docker rmi usual_store-ai-assistant│       │                          │
│                                    │       │ OR recreate:             │
│ docker build --no-cache \          │       │ terraform destroy \      │
│   -f Dockerfile.ai-assistant \     │       │   -target=module.\       │
│   -t usual_store-ai-assistant .    │       │   database.\             │
│                                    │       │   docker_container.\     │
│ terraform apply \                  │       │   database               │
│   -target=module.ai_assistant.\    │       │ terraform apply \        │
│   docker_container.ai_assistant \  │       │   -target=module.\       │
│   -auto-approve                    │       │   database.\             │
│                                    │       │   docker_container.\     │
│ ⏱️  Time: 5-10 minutes             │       │   database -auto-approve │
└────────────────────────────────────┘       │                          │
                                             │ ⚠️  WARNING: May lose data│
                                             └──────────────────────────┘
```

---

## 📊 Change Impact Matrix

| Changed Files | Rebuild Required | Command | Time | Cache Clear |
|--------------|------------------|---------|------|-------------|
| `react-frontend/src/*.jsx` | React Frontend | `terraform apply -target=module.frontends.docker_image.react_frontend ...` | 5-10 min | ✅ Yes |
| `react-frontend/src/*.css` | React Frontend | Same as above | 5-10 min | ✅ Yes |
| `support-frontend/src/*` | Support Frontend | `terraform apply -target=module.frontends.docker_image.support_frontend ...` | 5-10 min | ✅ Yes |
| `cmd/`, `internal/` (Go) | Backend | `terraform apply -target=module.backend_api.docker_image.backend ...` | 3-5 min | ❌ No |
| `internal/ai/` | AI Assistant | Clean rebuild (see above) | 5-10 min | ❌ No |
| `migrations/*.sql` | Database | `docker restart usualstore-database` | 30 sec | ❌ No |
| `terraform/modules/api_gateway/` | Kong Gateway | `terraform apply -target=module.api_gateway ...` | 1-2 min | ❌ No |
| `terraform.tfvars` | Full Apply | `terraform apply -auto-approve` | 2-5 min | ❌ No |
| Everything | Nuclear Rebuild | `terraform destroy && apply` | 15-20 min | ✅ Yes |

---

## 🎯 Optimization Tips

### Fast Iteration (CSS/Frontend Only)

```bash
# 1. Make changes in react-frontend/src/
# 2. Rebuild just the image (uses cache for node_modules)
cd react-frontend
docker build -t usual_store-react-frontend:latest .

# 3. Recreate container (fast)
terraform -chdir=../terraform apply \
  -target=module.frontends.docker_container.react_frontend \
  -auto-approve

# 4. Hard refresh browser: Cmd+Shift+R
```

**⏱️ Time: 2-3 minutes (vs 5-10 for full rebuild)**

### Backend Hot Reload (Development)

For rapid backend development, consider:

```bash
# Run backend directly (outside Docker)
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Set environment variables
export DATABASE_URL="postgres://postgres:password@localhost:5433/usualstore"
export OPENAI_API_KEY="sk-your-key"
export PORT=4001

# Run with auto-reload (using air or similar)
go run cmd/backend/main.go

# Make changes → automatic reload
```

---

## 🚦 Build Strategy Flowchart

```
┌─────────────────────────────────────────┐
│     Are you actively developing?        │
└─────────────────────────────────────────┘
              │
      ┌───────┴───────┐
      │               │
      ▼               ▼
┌──────────┐    ┌──────────┐
│   YES    │    │    NO    │
└──────────┘    └──────────┘
      │               │
      ▼               ▼
┌─────────────────────────────┐    ┌──────────────────────┐
│ Use Selective Rebuild:       │    │ Use Full Rebuild:    │
│ - Only rebuild changed       │    │ - Everything fresh   │
│ - Faster iteration           │    │ - Known good state   │
│ - May have stale deps        │    │ - Clean slate        │
└─────────────────────────────┘    └──────────────────────┘
              │                              │
              ▼                              ▼
┌─────────────────────────────┐    ┌──────────────────────┐
│ terraform apply -target=...  │    │ terraform destroy    │
│                              │    │ terraform apply      │
└─────────────────────────────┘    └──────────────────────┘
```

---

## 🔧 Troubleshooting Flowchart

```
┌─────────────────────────────────────────┐
│     Is something not working?            │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│ Step 1: Check container status           │
│ docker ps --filter "name=usualstore"    │
└─────────────────────────────────────────┘
              │
      ┌───────┴───────┐
      │               │
      ▼               ▼
┌──────────┐    ┌──────────────┐
│ Running? │    │  Not Running │
└──────────┘    └──────────────┘
      │               │
      NO              ▼
      │         ┌─────────────────────┐
      │         │ Check logs:          │
      │         │ docker logs <name>   │
      │         └─────────────────────┘
      │               │
      ▼               ▼
┌─────────────────────────────────┐
│ Step 2: Check logs for errors   │
│ docker logs <container-name>     │
└─────────────────────────────────┘
              │
      ┌───────┴────────┐
      │                │
      ▼                ▼
┌──────────┐     ┌──────────┐
│ API Key  │     │  Other   │
│  Error?  │     │  Error?  │
└──────────┘     └──────────┘
      │                │
      ▼                ▼
┌────────────────┐   ┌──────────────────┐
│ Fix keys in:   │   │ Try restart:     │
│ terraform.     │   │ docker restart   │
│ tfvars         │   │ <container>      │
└────────────────┘   └──────────────────┘
      │                │
      ▼                ▼
┌────────────────────────────────┐
│ Step 3: Still broken?          │
│ Nuclear option:                │
│ terraform destroy              │
│ terraform apply -auto-approve  │
└────────────────────────────────┘
```

---

## 💡 Best Practices

### 1. Before Starting Work
```bash
# Pull latest code
git pull origin main

# Check for changes
git status

# Start services
terraform apply -auto-approve
```

### 2. During Development
```bash
# Make small, testable changes
# Commit frequently
# Test after each change

# Only rebuild what changed
terraform apply -target=<specific-module>
```

### 3. Before Committing
```bash
# Ensure everything builds
terraform validate

# Run linting (if configured)
go fmt ./...
npm run lint

# Test the build
terraform plan
```

### 4. End of Day
```bash
# Stop services to save resources
terraform destroy -auto-approve

# OR keep running for next day
# (Docker Desktop can be paused)
```

---

## 📈 Performance Tips

### Speed Up Builds

1. **Use Docker Build Cache**
   - Don't use `--no-cache` unless necessary
   - Order Dockerfile commands: dependencies first, code last

2. **Selective Targeting**
   - Use `-target=` flag
   - Only rebuild changed modules

3. **Parallel Builds**
   ```bash
   # Build multiple services in parallel
   docker build -t frontend:latest react-frontend &
   docker build -t backend:latest . &
   wait
   ```

4. **Keep Docker Desktop Resources High**
   - Docker Desktop → Settings → Resources
   - CPUs: 4+
   - Memory: 8GB+

---

**For complete setup instructions, see `SETUP-GUIDE.md`**

