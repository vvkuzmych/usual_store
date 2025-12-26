# ✅ Final Setup - Usual Store (IPv4 Network)

## 🎯 What You Have Now

A **working, simplified Docker setup** with IPv4 networking that actually works on macOS.

---

## 📊 Services Running

```
✅ database   - PostgreSQL 15 on port 5433 (healthy)
✅ back-end   - API server on port 4001
✅ front-end  - Web server on port 4000
✅ invoice    - Invoice service on port 5000
```

---

## 🌐 Network Configuration

```
Network: usual_store_usualstore_network
├── Type: Bridge
├── IPv6: Disabled (doesn't work properly on macOS Docker Desktop)
├── IPv4 Subnet: 172.22.0.0/16
└── IPv4 Gateway: 172.22.0.1
```

**Why IPv4-only?**
- Docker Desktop on macOS doesn't support IPv6 host-to-container properly
- IPv4 is simpler, clearer, and actually works
- If you need IPv6 in production (Linux), it's easy to add later

---

## 🔌 How to Connect

### **From Your Mac (Development)**
```bash
# Database
psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable"

# Web Application
http://localhost:4000

# API
http://localhost:4001

# Invoice Service
http://localhost:5000
```

### **Inside Docker (Container-to-Container)**
```bash
# Your Go applications use this
DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
```
**Note**: Uses port **5432** inside Docker (not 5433)

---

## 🚀 Quick Commands

```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs

# Check status
docker compose ps

# Connect to database
make db-shell-ipv4
# or
psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable"
```

---

## 📝 Important Notes

### **Port Configuration**
- **Host Port 5433** → Database (to avoid conflict with your local PostgreSQL on 5432)
- **Container Port 5432** → Database internal port
- Your **Go applications inside Docker** use `database:5432`
- Your **Mac** uses `127.0.0.1:5433`

### **Network Subnet**
- Using **172.22.0.0/16** (changed from 172.20 due to conflict with headlamp extension)
- Containers get IPs like: 172.22.0.2, 172.22.0.3, etc.

---

## ✅ Verification

```bash
# Test database connection
psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable" -c "SELECT 1;"

# Check services
docker compose ps

# View network details
docker network inspect usual_store_usualstore_network

# Test web app
curl http://localhost:4000
```

---

## 🔧 What Was Fixed

1. ✅ **Network subnet conflict** - Changed from 172.20 to 172.22
2. ✅ **Port conflict** - PostgreSQL on 5433 (not 5432)
3. ✅ **Simplified to IPv4** - Removed non-functional IPv6 complexity

---

## 📚 Documentation

The IPv6 documentation files are kept for reference, but the **current setup is IPv4-only**:
- `FINAL-SETUP.md` ← **You are here (use this!)**
- `SETUP-COMPLETE.md` - Previous setup notes
- IPv6-*.md files - Reference only (IPv6 doesn't work on macOS Docker Desktop)

---

## 🎉 Summary

Your Usual Store is now running with:
- ✅ **Simple IPv4 networking** that actually works
- ✅ **All services healthy** and accessible
- ✅ **No port conflicts**
- ✅ **Clear and understandable** configuration
- ✅ **Ready for development** on macOS

**Everything is working correctly!** 🚀

---

## 💬 Need Help?

```bash
# Services won't start?
make docker-logs

# Can't connect to database?
psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable"

# Port still in use?
lsof -nP -iTCP:5433 -sTCP:LISTEN

# Restart everything
make docker-down && make docker-up
```

---

**Your setup is complete and working!** ✅

