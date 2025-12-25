# Makefile Commands Guide

Complete reference for all `make` commands available in the Usual Store project.

## 🚀 Quick Start

### Start Different Frontends

```bash
# Start React Frontend (port 3000) + Backend
make react-start

# Start TypeScript Frontend (port 3001) + Backend
make typescript-start

# Start Go Frontend (port 4000) + Backend
make go-start

# Start ALL frontends simultaneously
make all-frontends-start
```

### Access URLs

| Frontend | Command | URL |
|----------|---------|-----|
| **React** (JavaScript) | `make react-start` | http://localhost:3000 |
| **TypeScript** (Vite) | `make typescript-start` | http://localhost:3001 |
| **Go** (HTML Templates) | `make go-start` | http://localhost:4000 |
| Backend API | (included with all) | http://localhost:4001 |

---

## 📋 Command Reference

### Frontend Management

#### Start Commands
```bash
# Start individual frontends
make react-start          # React frontend + backend
make typescript-start     # TypeScript frontend + backend
make go-start            # Go frontend + backend
make all-frontends-start # All frontends + backend
```

#### Stop Commands
```bash
make react-stop          # Stop React frontend
make typescript-stop     # Stop TypeScript frontend
make go-stop            # Stop Go frontend
make docker-down        # Stop all services
```

#### Build Commands
```bash
make react-build             # Build React Docker image
make typescript-build        # Build TypeScript Docker image
make go-build-docker         # Build Go Docker image
make build-all-frontends     # Build all frontend images
```

#### Logs Commands
```bash
make react-logs          # View React frontend logs
make typescript-logs     # View TypeScript frontend logs
make go-logs            # View Go frontend logs
make docker-logs        # View all service logs
```

#### Restart Commands
```bash
make react-restart       # Restart React frontend
make typescript-restart  # Restart TypeScript frontend
make go-restart         # Restart Go frontend
make docker-restart     # Restart all services
```

#### Status Commands
```bash
make frontend-status     # Show status of all frontends
make docker-ps          # Show all Docker containers
```

---

### Docker Compose Commands

```bash
make docker-up           # Start all services
make docker-down         # Stop all services
make docker-restart      # Restart all services
make docker-logs         # Follow all logs
make docker-ps           # List running containers
```

---

### Build Commands

```bash
make build               # Build all Go binaries
make build_front         # Build Go frontend binary
make build_back          # Build backend API binary
make build_invoice       # Build invoice microservice binary
make clean               # Clean all binaries
```

---

### Database Commands

#### Database Connection
```bash
make db-shell-ipv4       # Connect via IPv4 (127.0.0.1:5433)
make db-shell-ipv6       # Connect via IPv6 ([::1]:5433)
make db-shell-docker     # Connect inside Docker container
```

#### Migrations
```bash
make migrate             # Run database migrations
make rollback            # Rollback migrations
make new-migration       # Create new migration (set MIGRATION_NAME)
make create-db           # Create database
make drop-db             # Drop database
```

Example:
```bash
# Create a new migration
make new-migration MIGRATION_NAME=add_users_table
```

---

### IPv6 Testing Commands

```bash
make test-ipv6           # Comprehensive IPv6 tests
make check-ipv6-network  # Check Docker IPv6 config
make verify-db-ipv6      # Verify PostgreSQL IPv6
make test-db-ipv6-host   # Test DB connection via IPv6
make test-db-ipv4-host   # Test DB connection via IPv4
make show-container-ips  # Show all container IPs
```

---

### Mock Generation

```bash
make mock                # Generate all mocks
make mock_token_repository   # Generate TokenRepository mock
make mock_user_repository    # Generate UserRepository mock
```

---

## 🎯 Common Workflows

### Development Workflow

**1. Start React Frontend for Development:**
```bash
make react-start
# Access at http://localhost:3000
```

**2. View Logs:**
```bash
make react-logs
```

**3. Make Changes & Restart:**
```bash
make react-restart
```

**4. Stop When Done:**
```bash
make react-stop
```

---

### Testing Different Frontends

**Test All Frontends:**
```bash
# Start all frontends
make all-frontends-start

# Check status
make frontend-status

# Test each:
# - React:      http://localhost:3000
# - TypeScript: http://localhost:3001
# - Go:         http://localhost:4000

# Stop all
make docker-down
```

---

### Build & Deploy Workflow

**1. Build All Images:**
```bash
make build-all-frontends
```

**2. Start Services:**
```bash
make typescript-start
```

**3. Verify:**
```bash
make frontend-status
```

---

### Database Management Workflow

**1. Connect to Database:**
```bash
make db-shell-ipv4
```

**2. Run Migrations:**
```bash
make migrate
```

**3. Verify IPv6 Support:**
```bash
make test-ipv6
```

---

## 🔧 Troubleshooting

### Check Service Status
```bash
make docker-ps
make frontend-status
```

### View Logs
```bash
# Specific frontend
make react-logs
make typescript-logs
make go-logs

# All services
make docker-logs
```

### Restart Services
```bash
# Specific frontend
make react-restart

# All services
make docker-restart
```

### Rebuild Images
```bash
# Specific frontend
make react-build

# All frontends
make build-all-frontends
```

### Clean State
```bash
# Stop everything
make docker-down

# Rebuild
make build-all-frontends

# Start fresh
make typescript-start
```

---

## 📊 Comparison Table

| Feature | React | TypeScript | Go |
|---------|-------|------------|-----|
| **Start Command** | `make react-start` | `make typescript-start` | `make go-start` |
| **Port** | 3000 | 3001 | 4000 |
| **Language** | JavaScript | TypeScript | Go |
| **Build Tool** | Create React App | Vite | Go compiler |
| **Theme** | Purple | Blue | Purple |
| **Logs** | `make react-logs` | `make typescript-logs` | `make go-logs` |
| **Stop** | `make react-stop` | `make typescript-stop` | `make go-stop` |

---

## 💡 Tips & Best Practices

### 1. Use Specific Frontend Commands
Instead of starting all services with `make docker-up`, use frontend-specific commands:
```bash
# Better
make typescript-start

# Instead of
docker compose --profile typescript-frontend up -d
```

### 2. Check Status Before Starting
```bash
make frontend-status
```

### 3. View Logs in Real-Time
```bash
make typescript-logs
```

### 4. Rebuild After Code Changes
```bash
make typescript-build
make typescript-restart
```

### 5. Use All-Frontends for Testing
```bash
make all-frontends-start
make frontend-status
```

---

## 🎓 Examples

### Example 1: Start React Development
```bash
# Start React frontend
make react-start

# View logs
make react-logs

# Open browser: http://localhost:3000

# Stop when done
make react-stop
```

### Example 2: Compare All Frontends
```bash
# Start all
make all-frontends-start

# Check status
make frontend-status

# Compare:
# - React:      http://localhost:3000 (Purple, JavaScript)
# - TypeScript: http://localhost:3001 (Blue, TypeScript)
# - Go:         http://localhost:4000 (Purple, Go Templates)

# Stop all
make docker-down
```

### Example 3: Build and Deploy TypeScript
```bash
# Build image
make typescript-build

# Start service
make typescript-start

# Check logs
make typescript-logs

# Verify
curl http://localhost:3001
```

### Example 4: Database Operations
```bash
# Connect to database
make db-shell-ipv4

# In psql:
# \dt  (list tables)
# \q   (quit)

# Run migrations
make migrate

# Test IPv6
make test-ipv6
```

---

## 📝 Help Command

To see all available commands:

```bash
make help
```

This displays a comprehensive list of all available targets with descriptions.

---

## 🔗 Related Documentation

- [TypeScript Frontend Setup](../setup/TYPESCRIPT-FRONTEND-SETUP.md)
- [Material UI Setup](../setup/MATERIAL-UI-SETUP.md)
- [Docker Deployment](./DOCKER-DEPLOYMENT.md)
- [Kubernetes Deployment](../kubernetes/GETTING-STARTED.md)

---

## ✅ Quick Reference

Most commonly used commands:

```bash
# Frontend
make typescript-start      # Start TypeScript frontend
make react-start          # Start React frontend
make frontend-status      # Check status

# Logs
make typescript-logs      # View TypeScript logs
make react-logs          # View React logs

# Stop
make docker-down         # Stop everything

# Build
make typescript-build    # Build TypeScript image

# Database
make db-shell-ipv4      # Connect to database

# Help
make help               # Show all commands
```

---

**Happy Development! 🚀**

