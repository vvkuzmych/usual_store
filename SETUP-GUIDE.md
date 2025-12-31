# UsualStore - Complete Setup & Operations Guide

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Getting API Keys & Credentials](#getting-api-keys--credentials)
3. [Initial Setup](#initial-setup)
4. [Configuration](#configuration)
5. [Starting the Application](#starting-the-application)
6. [Verification & Health Checks](#verification--health-checks)
7. [Development Workflow](#development-workflow)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

Before starting, install the following tools:

#### 1. Docker & Docker Compose
```bash
# macOS (using Homebrew)
brew install --cask docker

# Or download Docker Desktop from:
# https://www.docker.com/products/docker-desktop

# Verify installation
docker --version
docker-compose --version
```

#### 2. Terraform
```bash
# macOS
brew tap hashicorp/tap
brew install hashicorp/tap/terraform

# Verify installation
terraform --version
# Should be v1.0.0 or higher
```

#### 3. Kubernetes (Optional - for K8s deployment)
```bash
# Install kubectl
brew install kubectl

# Install minikube for local Kubernetes
brew install minikube

# Verify installation
kubectl version --client
minikube version
```

#### 4. Go (for backend development)
```bash
# macOS
brew install go

# Verify installation
go version
# Should be Go 1.21 or higher
```

#### 5. Node.js & npm (for frontend development)
```bash
# macOS
brew install node

# Verify installation
node --version  # Should be v18.x or higher
npm --version   # Should be 9.x or higher
```

#### 6. Git
```bash
# macOS
brew install git

# Verify
git --version
```

#### 7. PostgreSQL Client (Optional but Recommended)

**Note**: PostgreSQL **database server runs in Docker** (no installation needed), but the `psql` client tool is useful for database access.

```bash
# macOS - Install just the client tools
brew install libpq

# Add to PATH (add to ~/.zshrc or ~/.bashrc)
echo 'export PATH="/opt/homebrew/opt/libpq/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verify
psql --version

# Alternative: Install full PostgreSQL (includes psql)
brew install postgresql@15
```

**What you get**:
- `psql` - PostgreSQL interactive terminal
- Database access for debugging and migrations
- Direct SQL query execution

**Not needed**:
- PostgreSQL server (runs in Docker)
- Database initialization
- Port configuration

### System Requirements

- **OS**: macOS, Linux, or Windows (WSL2)
- **RAM**: Minimum 8GB (16GB recommended)
- **Disk Space**: 10GB free space
- **CPU**: 4 cores or more recommended
- **Database**: PostgreSQL 15 (runs in Docker - included)

---

## Getting API Keys & Credentials

### 1. OpenAI API Key (Required for AI features)

**Purpose**: Speech-to-text, text-to-speech, chat AI

**Steps**:
1. Go to https://platform.openai.com/api-keys
2. Sign up or log in
3. Click "Create new secret key"
4. Name it (e.g., "UsualStore Development")
5. Copy the key (starts with `sk-...`)
6. **⚠️ Save it immediately - you won't see it again!**

**Cost**: Pay-as-you-go (starts at ~$0.006 per 1K tokens)

### 2. Stripe API Keys (Required for payments)

**Purpose**: Payment processing, subscriptions

**Steps**:
1. Go to https://dashboard.stripe.com/register
2. Sign up for a Stripe account
3. Navigate to **Developers → API Keys**
4. Copy **Publishable key** (starts with `pk_test_...`)
5. Copy **Secret key** (starts with `sk_test_...`)

**Note**: Use test keys for development (contain `test`)

### 3. Mailtrap Credentials (Required for email testing)

**Purpose**: Email testing in development (catches all emails)

**Steps**:
1. Go to https://mailtrap.io/register/signup
2. Sign up for free account
3. Create an inbox (e.g., "UsualStore Dev")
4. Go to **SMTP Settings**
5. Copy:
   - **Host**: `smtp.mailtrap.io`
   - **Port**: `2525`
   - **Username**: Your username
   - **Password**: Your password

**Alternative**: Use Gmail SMTP (requires App Password)

### 4. PostgreSQL Database

**Database Type**: PostgreSQL 15

**Installation**: ✅ **NOT REQUIRED** - Runs in Docker container

The PostgreSQL database server runs completely in Docker and is managed by Terraform. You don't need to install PostgreSQL server locally.

**What Terraform provides**:
- PostgreSQL 15 Docker image (automatically pulled)
- Database container: `usualstore-database`
- Persistent storage via Docker volume
- Network configuration
- Health checks
- Port mapping: `localhost:5433` → container `5432`

**Optional**: Install PostgreSQL client tools (`psql`) for direct database access (see Prerequisites section above)

---

## Initial Setup

### Step 1: Clone the Repository

```bash
cd ~/Projects
git clone <your-repo-url> UsualStore
cd UsualStore/usual_store
```

### Step 2: Create Configuration Files

#### 2.1 Copy Terraform Variables Template

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

#### 2.2 Edit `terraform/terraform.tfvars`

```hcl
# Docker Configuration
docker_host = "unix:///var/run/docker.sock"

# Network Configuration (usually no need to change)
ipv4_subnet  = "172.22.0.0/16"
ipv4_gateway = "172.22.0.1"
ipv6_subnet  = "2001:db8:1::/64"
ipv6_gateway = "2001:db8:1::1"

# Database Configuration
postgres_user     = "postgres"
postgres_password = "your_secure_password_here"  # ⚠️ CHANGE THIS
postgres_db       = "usualstore"

# Service Ports (default values - change if conflicts)
api_port        = 4001
support_port    = 5001
react_port      = 3007
typescript_port = 3001
redux_port      = 3002
support_ui_port = 3005
ai_port         = 8080

# ⚠️ ADD YOUR API KEYS HERE
stripe_key    = "pk_test_your_stripe_publishable_key"
stripe_secret = "sk_test_your_stripe_secret_key"
openai_api_key = "sk-your_openai_api_key"

# Mailtrap Configuration (add to environment variables)
# SMTP_HOST=smtp.mailtrap.io
# SMTP_PORT=2525
# SMTP_USER=your_mailtrap_username
# SMTP_PASS=your_mailtrap_password

# Environment
environment = "development"

# Feature Flags
enable_kafka         = true
enable_observability = true
enable_ai_assistant  = true  # Set to true to enable AI features
enable_support       = true
```

#### 2.3 Create Environment File (Optional)

```bash
# In project root
cat > .env << 'EOF'
# Database
DATABASE_URL=postgres://postgres:your_password@localhost:5433/usualstore

# OpenAI
OPENAI_API_KEY=sk-your_openai_api_key

# Stripe
STRIPE_PUBLISHABLE_KEY=pk_test_your_key
STRIPE_SECRET_KEY=sk_test_your_key

# Mailtrap
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your_username
SMTP_PASS=your_password

# Application
PORT=4001
NODE_ENV=development
EOF
```

### Step 3: Initialize Terraform

```bash
cd terraform
terraform init

# Expected output:
# Terraform has been successfully initialized!
```

### Step 4: Validate Configuration

```bash
terraform validate

# Expected output:
# Success! The configuration is valid.
```

### Step 5: Plan Infrastructure

```bash
terraform plan

# This shows what will be created (do NOT apply yet, just review)
```

### Step 6: Install Migration Tool (soda)

**Purpose**: Soda is used to run database migrations

```bash
# macOS
brew install gobuffalo/tap/pop

# Verify installation
soda --version

# Or download directly
go install github.com/gobuffalo/pop/v6/soda@latest
```

**Note**: The `soda` binary path is configured in `Makefile` (line 113)

---

## Database Migrations

### Understanding Migrations

**Location**: `/Users/vkuzm/Projects/UsualStore/usual_store/migrations/`

**Format**: Soda migrations (up/down pairs)
- `*.up.sql` - Creates/adds database objects
- `*.down.sql` - Reverts changes

**Example files**:
- `20240730184028_create_users_table.up.sql`
- `20240730184028_create_users_table.down.sql`

### Run Migrations (Required Before First Start)

**Option 1: Using Makefile (Recommended)**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Start database first
terraform -chdir=terraform apply \
  -target=module.database.docker_container.database \
  -auto-approve

# Wait for database to be healthy (10 seconds)
sleep 10

# Run all migrations
make migrate

# Expected output:
# > Status
# > All migrations are up to date!
```

**Option 2: Using soda directly**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Run migrations
soda migrate up

# Check migration status
soda migrate status

# Rollback last migration (if needed)
soda migrate down
```

**Option 3: Manual SQL execution (not recommended)**

```bash
# Apply each migration manually
for file in migrations/*.up.sql; do
  docker exec -i usualstore-database psql -U postgres -d usualstore < "$file"
done
```

### Verify Migrations

```bash
# Connect to database
docker exec -it usualstore-database psql -U postgres -d usualstore

# Check tables exist
\dt

# Expected tables:
# - users
# - widgets
# - orders
# - transactions
# - ai_conversations
# - ai_messages
# - support_tickets
# (and more...)

# Exit
\q
```

### Create New Migration

```bash
# Using Makefile
make new-migration MIGRATION_NAME=add_user_preferences

# Using soda directly
soda generate migration add_user_preferences

# This creates two files:
# migrations/YYYYMMDDHHMMSS_add_user_preferences.up.sql
# migrations/YYYYMMDDHHMMSS_add_user_preferences.down.sql
```

---

## Starting the Application

### ⚠️ IMPORTANT: Database Migrations Must Run First!

Before starting the application for the first time, you **must** run database migrations.

### Quick Start (First Time)

```bash
# From project root: /Users/vkuzm/Projects/UsualStore/usual_store

# Step 1: Start database first
cd terraform
terraform apply \
  -target=module.database.docker_container.database \
  -auto-approve

# Step 2: Wait for database to be ready
sleep 10

# Step 3: Run migrations (REQUIRED!)
cd ..
make migrate

# Step 4: Now start everything else
cd terraform
terraform apply -auto-approve

# This will:
# - Pull/build Docker images
# - Create network and volumes
# - Start all containers
# - Configure Kong API Gateway
```

**⏱️ Time**: First run takes 10-15 minutes (building images)

**Alternative (All at once - migrations may fail)**:
```bash
# Start everything then run migrations
terraform apply -auto-approve
sleep 30  # Wait for database
make migrate
```

### What Gets Started

| Service | Container Name | Port | URL |
|---------|---------------|------|-----|
| API Gateway | api-gateway | 8000 | http://localhost:8000 |
| Backend API | usualstore-back-end | 4001 | http://localhost:4001 |
| React Frontend | usualstore-react-frontend | 3007 | http://localhost:3007 |
| Support Frontend | usualstore-support-frontend | 3005 | http://localhost:3005 |
| AI Assistant | usualstore-ai-assistant | 8080 | http://localhost:8080 |
| Support Service | usualstore-support-service | 5001 | http://localhost:5001 |
| Database | usualstore-database | 5433 | localhost:5433 |
| Jaeger UI | usualstore-jaeger | 16686 | http://localhost:16686 |
| Kafka UI | usualstore-kafka-ui | 8090 | http://localhost:8090 |
| OPA Server | usualstore-opa-server | 8181 | http://localhost:8181 |

---

## Verification & Health Checks

### Step 1: Check Container Status

```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# All containers should show "Up" and "(healthy)"
```

### Step 2: Test API Gateway

```bash
# Health check
curl http://localhost:8000/health

# Expected: {"status":"ok"}
```

### Step 3: Test Backend API

```bash
# Get products
curl http://localhost:8000/api/products

# Expected: JSON array of products
```

### Step 4: Test Frontends

Open in browser:
- **Main App**: http://localhost:8000
- **React Frontend**: http://localhost:3007
- **Support UI**: http://localhost:3005

### Step 5: Test AI Assistant (if enabled)

```bash
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-123",
    "message": "Hello!",
    "user_id": null
  }'

# Expected: JSON response with AI message
```

### Step 6: Check Database

```bash
# Connect to database
docker exec -it usualstore-database psql -U postgres -d usualstore

# Run a test query
\dt  # List tables
SELECT COUNT(*) FROM widgets;
\q   # Quit
```

### Step 7: View Logs

```bash
# View logs for specific service
docker logs usualstore-react-frontend
docker logs usualstore-back-end
docker logs usualstore-ai-assistant

# Follow logs in real-time
docker logs -f usualstore-back-end
```

---

## Development Workflow

### When You Make Code Changes

Different types of changes require different update procedures:

#### 1. Frontend Changes (React, TypeScript, Redux)

**Files affected**: Anything in `react-frontend/`, `typescript-frontend/`, `support-frontend/`

**What to do**:
```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Option A: Rebuild specific frontend
terraform -chdir=terraform apply \
  -target=module.frontends.docker_image.react_frontend \
  -target=module.frontends.docker_container.react_frontend \
  -auto-approve

# Option B: Use Docker directly (faster for iteration)
cd react-frontend
docker build -t usual_store-react-frontend:latest .
docker stop usualstore-react-frontend
docker rm usualstore-react-frontend
# Then run terraform apply to recreate container

# ⏱️ Time: 5-10 minutes
```

**Clear browser cache**: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows)

#### 2. Backend Changes (Go)

**Files affected**: Anything in `cmd/`, `internal/`, `pkg/`

**What to do**:
```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Rebuild backend
terraform -chdir=terraform apply \
  -target=module.backend_api.docker_image.backend \
  -target=module.backend_api.docker_container.backend \
  -auto-approve

# ⏱️ Time: 3-5 minutes
```

#### 3. AI Assistant Changes

**Files affected**: `internal/ai/`, `cmd/ai-assistant-example/`

**What to do**:
```bash
# Remove old images to force clean rebuild
docker stop usualstore-ai-assistant
docker rm usualstore-ai-assistant
docker rmi usual_store-ai-assistant:latest

# Rebuild from scratch
docker build --no-cache -f Dockerfile.ai-assistant \
  -t usual_store-ai-assistant:latest .

# Recreate container
terraform -chdir=terraform apply \
  -target=module.ai_assistant.docker_container.ai_assistant \
  -auto-approve

# ⏱️ Time: 5-10 minutes
```

#### 4. Database Schema Changes

**Files affected**: `migrations/*.sql`

**What to do**:
```bash
# Migrations run automatically on container start

# Option A: Restart database container
docker restart usualstore-database

# Option B: Run migrations manually
docker exec -it usualstore-database psql -U postgres -d usualstore -f /migrations/your_migration.sql

# Option C: Full recreate (⚠️ deletes data)
terraform -chdir=terraform destroy -target=module.database.docker_container.database
terraform -chdir=terraform apply -target=module.database.docker_container.database -auto-approve
```

#### 5. API Gateway Configuration Changes

**Files affected**: Kong routes, services in `terraform/modules/api_gateway/`

**What to do**:
```bash
# Recreate Kong migration and gateway
terraform -chdir=terraform apply \
  -target=module.api_gateway.docker_container.kong_migration \
  -target=module.api_gateway.docker_container.kong \
  -auto-approve

# ⏱️ Time: 1-2 minutes
```

#### 6. CSS/Style Changes Only

**Fastest update method**:
```bash
cd react-frontend

# Rebuild just the image (uses cache for node_modules)
docker build -t usual_store-react-frontend:latest .

# Recreate container
terraform -chdir=terraform apply \
  -target=module.frontends.docker_container.react_frontend \
  -auto-approve

# ⏱️ Time: 2-3 minutes
```

### Complete Rebuild (Nuclear Option)

**When to use**: Major changes, switching branches, strange behavior

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# 1. Stop and remove ALL containers
docker stop $(docker ps -q --filter "name=usualstore")
docker rm $(docker ps -aq --filter "name=usualstore")

# 2. Remove all images
docker rmi $(docker images 'usual_store*' -q)

# 3. Clean build cache (optional)
docker builder prune -af

# 4. Rebuild everything
cd terraform
terraform destroy -auto-approve  # ⚠️ Careful! Deletes volumes too
terraform apply -auto-approve

# ⏱️ Time: 15-20 minutes
```

### Selective Rebuild (Recommended)

**When to use**: Most code changes

```bash
# Only rebuild what changed
terraform -chdir=terraform apply \
  -target=module.frontends.docker_image.react_frontend \
  -target=module.ai_assistant.docker_image.ai_assistant \
  -auto-approve
```

---

## Common Commands

### Start/Stop Services

```bash
# Stop all services
terraform -chdir=terraform destroy -auto-approve

# Start all services
terraform -chdir=terraform apply -auto-approve

# Restart specific service
docker restart usualstore-react-frontend
```

### View Status

```bash
# All containers
docker ps -a

# Specific service logs
docker logs usualstore-back-end

# Follow logs
docker logs -f usualstore-ai-assistant

# Resource usage
docker stats
```

### Access Services

```bash
# Connect to database
docker exec -it usualstore-database psql -U postgres -d usualstore

# Shell into container
docker exec -it usualstore-back-end sh

# Check environment variables
docker exec usualstore-ai-assistant env | grep OPENAI
```

### Clean Up

```bash
# Remove stopped containers
docker container prune -f

# Remove unused images
docker image prune -a -f

# Remove unused volumes (⚠️ deletes data)
docker volume prune -f

# Remove everything (⚠️ nuclear option)
docker system prune -a --volumes -f
```

---

## Troubleshooting

### Issue: Port Already in Use

**Error**: `Bind for 0.0.0.0:3007 failed: port is already allocated`

**Solution**:
```bash
# Find what's using the port
lsof -i :3007

# Kill the process
kill -9 <PID>

# Or change the port in terraform.tfvars
react_port = 3008  # Use different port
```

### Issue: Container Won't Start (Unhealthy)

**Solution**:
```bash
# Check logs
docker logs usualstore-react-frontend

# Check health check
docker inspect usualstore-react-frontend | grep -A 10 Health

# Restart container
docker restart usualstore-react-frontend
```

### Issue: Database Connection Failed

**Solution**:
```bash
# Verify database is running
docker ps | grep database

# Check database logs
docker logs usualstore-database

# Test connection
docker exec -it usualstore-database psql -U postgres -d usualstore
```

### Issue: OpenAI API Errors

**Error**: `Incorrect API key provided`

**Solution**:
```bash
# Check API key is set
docker exec usualstore-ai-assistant env | grep OPENAI

# Update key in terraform.tfvars
openai_api_key = "sk-your-correct-key"

# Recreate AI assistant
terraform -chdir=terraform apply \
  -target=module.ai_assistant.docker_container.ai_assistant \
  -auto-approve
```

### Issue: Migration Errors

**Error**: `relation "users" does not exist` or similar table errors

**Solution**:
```bash
# Check if migrations have been run
docker exec -it usualstore-database psql -U postgres -d usualstore -c "\dt"

# If no tables exist, run migrations
make migrate

# Or check migration status
soda migrate status
```

**Error**: `soda: command not found`

**Solution**:
```bash
# Install soda
brew install gobuffalo/tap/pop

# Or
go install github.com/gobuffalo/pop/v6/soda@latest

# Verify
which soda
```

**Error**: `migration already applied`

**Solution**:
```bash
# Check status
soda migrate status

# If stuck, manually check schema_migration table
docker exec -it usualstore-database psql -U postgres -d usualstore -c "SELECT * FROM schema_migration;"

# If needed, reset (⚠️ DESTRUCTIVE)
docker exec -it usualstore-database psql -U postgres -d usualstore -c "DROP TABLE IF EXISTS schema_migration CASCADE;"
make migrate
```

### Issue: Frontend Not Updating

**Solution**:
```bash
# Hard refresh browser
# Mac: Cmd+Shift+R
# Windows: Ctrl+Shift+R

# Clear browser cache completely
# Chrome: Settings → Privacy → Clear browsing data

# Verify new build
docker exec usualstore-react-frontend ls -la /usr/share/nginx/html/static/js/

# Force rebuild without cache
cd react-frontend
docker build --no-cache -t usual_store-react-frontend:latest .
```

### Issue: Terraform State Lock

**Error**: `Error acquiring the state lock`

**Solution**:
```bash
cd terraform

# Force unlock (use with caution)
terraform force-unlock <LOCK_ID>

# Or remove state lock file
rm -f .terraform.tfstate.lock.info
```

### Issue: Out of Disk Space

**Solution**:
```bash
# Check disk usage
docker system df

# Clean up
docker system prune -a --volumes -f

# Remove old images
docker images | grep '<none>' | awk '{print $3}' | xargs docker rmi
```

---

## Quick Reference

### Essential Commands

```bash
# Start everything
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
terraform apply -auto-approve

# Stop everything
terraform destroy -auto-approve

# Rebuild frontend after changes
terraform apply \
  -target=module.frontends.docker_image.react_frontend \
  -target=module.frontends.docker_container.react_frontend \
  -auto-approve

# Rebuild backend after changes
terraform apply \
  -target=module.backend_api.docker_image.backend \
  -target=module.backend_api.docker_container.backend \
  -auto-approve

# View logs
docker logs -f usualstore-<service-name>

# Check status
docker ps --format "table {{.Names}}\t{{.Status}}"
```

### URLs Quick Access

- Main App: http://localhost:8000
- React Frontend: http://localhost:3007
- Support UI: http://localhost:3005
- API Direct: http://localhost:4001
- AI Assistant: http://localhost:8080
- Jaeger Tracing: http://localhost:16686
- Kafka UI: http://localhost:8090
- Kong Admin: http://localhost:8001

---

## Production Deployment

For production deployment, see:
- `KUBERNETES-DEPLOYMENT.md` - Kubernetes setup
- `AWS-DEPLOYMENT.md` - AWS ECS/EKS deployment
- `DOCKER-COMPOSE-PROD.md` - Docker Compose production setup

---

## Getting Help

- **Issues**: Check troubleshooting section above
- **Logs**: Always check `docker logs <container-name>`
- **Documentation**: See `/docs` folder
- **API Docs**: http://localhost:8000/api/docs (when running)

---

**Last Updated**: December 31, 2025
**Version**: 1.0.0

