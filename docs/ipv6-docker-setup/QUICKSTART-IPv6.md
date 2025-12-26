# 🚀 Quick Start: IPv6 Setup

This is a quick guide to get your Usual Store project running with IPv6 support.

---

## ⚡ Quick Commands

### 1. **Start Everything**

```bash
make docker-up
```

### 2. **Test IPv6 Connectivity**

```bash
make test-ipv6
```

### 3. **Verify Database IPv6**

```bash
make verify-db-ipv6
```

### 4. **Test Database Connection from Host**

```bash
# IPv6
make test-db-ipv6-host

# IPv4 (backward compatible)
make test-db-ipv4-host
```

### 5. **Show Container IP Addresses**

```bash
make show-container-ips
```

### 6. **View Logs**

```bash
make docker-logs
```

### 7. **Stop Everything**

```bash
make docker-down
```

---

## 📋 Step-by-Step Setup

### Step 1: Prepare Environment

```bash
# Copy environment template
cp env.example .env

# Edit .env with your actual values
nano .env
```

**Minimum required values in `.env`**:
```bash
USUAL_STORE_PORT=4000
API_PORT=4001
INVOICE_PORT=5000
DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
SECRET_FOR_FRONT=your_secret_key_minimum_32_chars
STRIPE_KEY=your_stripe_key
STRIPE_SECRET=your_stripe_secret
FRONT_URL=http://localhost:4000
API_URL=http://localhost:4001
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your_user
SMTP_PASSWORD=your_password
```

### Step 2: Start Services

```bash
# Start all services in detached mode
docker compose up -d

# Wait for services to be healthy (about 10-15 seconds)
docker compose ps
```

**Expected output**:
```
NAME                      STATUS                    PORTS
usualstore-back-end-1     Up                        0.0.0.0:4001->4001/tcp
usualstore-database-1     Up (healthy)              [::1]:5432->5432/tcp, 127.0.0.1:5432->5432/tcp
usualstore-front-end-1    Up                        0.0.0.0:4000->4000/tcp
usualstore-invoice-1      Up                        0.0.0.0:5000->5000/tcp
```

### Step 3: Apply Database Schema

```bash
# Apply schema directly
cat migrations/schema.sql | docker compose exec -T database psql -U postgres -d usualstore

# Or run migrations
docker compose exec database psql -U postgres -d usualstore -f /docker-entrypoint-initdb.d/schema.sql
```

### Step 4: Test Connectivity

```bash
# Comprehensive IPv6 test
./scripts/test-ipv6.sh

# Or use Makefile
make test-ipv6
```

### Step 5: Access Services

**Front-end**:
- IPv4: http://localhost:4000
- IPv6: http://[::1]:4000

**Back-end API**:
- IPv4: http://localhost:4001
- IPv6: http://[::1]:4001

**Invoice Service**:
- IPv4: http://localhost:5000
- IPv6: http://[::1]:5000

**PostgreSQL**:
- IPv4: `postgresql://postgres:password@127.0.0.1:5432/usualstore`
- IPv6: `postgresql://postgres:password@[::1]:5432/usualstore`

---

## 🔍 Verification Checklist

Run these commands to verify everything is working:

```bash
# ✅ 1. Check Docker network has IPv6 enabled
docker network inspect usualstore_usualstore_network | grep EnableIPv6
# Should show: "EnableIPv6": true

# ✅ 2. Check PostgreSQL is listening on IPv6
docker compose exec database ss -tln | grep ::1
# Should show: :::5432 or ::1:5432

# ✅ 3. Test database connection (IPv6)
psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable" -c "SELECT 1;"
# Should show: 1

# ✅ 4. Test database connection (IPv4)
psql "postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable" -c "SELECT 1;"
# Should show: 1

# ✅ 5. Test front-end (IPv6)
curl -I http://[::1]:4000
# Should show: HTTP/1.1 200 OK

# ✅ 6. Test back-end API (IPv6)
curl -I http://[::1]:4001
# Should show: HTTP/1.1 ... (depends on route)

# ✅ 7. Show all container IPs
make show-container-ips
# Should show both IPv4 (172.20.0.x) and IPv6 (fd00:dead:beef::x) addresses
```

---

## 🛠️ Troubleshooting

### Issue: "connection refused"

```bash
# Check if services are running
docker compose ps

# Check logs
docker compose logs database

# Restart services
docker compose restart

# Nuclear option (recreate everything)
docker compose down -v && docker compose up -d
```

### Issue: "IPv6 not enabled"

```bash
# Recreate network with IPv6
docker compose down
docker network prune -f
docker compose up -d

# Verify IPv6 is enabled
docker network inspect usualstore_usualstore_network | grep -A 5 IPv6
```

### Issue: "Cannot connect from container to database"

```bash
# Test DNS resolution from front-end
docker compose exec front-end getent hosts database

# Test ping from front-end to database
docker compose exec front-end ping -c 2 database

# Check database is healthy
docker compose ps database
# Should show: Up (healthy)
```

### Issue: "psql: command not found"

```bash
# Install PostgreSQL client on macOS
brew install postgresql

# Or use docker to connect
docker compose exec database psql -U postgres -d usualstore
```

---

## 📊 Useful Commands Reference

### Docker Compose
```bash
# Start in foreground (see logs in real-time)
docker compose up

# Start in background
docker compose up -d

# Stop services
docker compose down

# Stop and remove volumes (clean slate)
docker compose down -v

# Restart specific service
docker compose restart database

# View logs (all services)
docker compose logs -f

# View logs (specific service)
docker compose logs -f database

# Execute command in container
docker compose exec database bash
```

### Database
```bash
# Connect via psql (from host)
psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable"

# Connect via psql (inside Docker)
docker compose exec database psql -U postgres -d usualstore

# Run SQL file
docker compose exec -T database psql -U postgres -d usualstore < schema.sql

# Backup database
docker compose exec database pg_dump -U postgres usualstore > backup.sql

# Restore database
cat backup.sql | docker compose exec -T database psql -U postgres -d usualstore
```

### Network
```bash
# Inspect network
docker network inspect usualstore_usualstore_network

# List all Docker networks
docker network ls

# Remove unused networks
docker network prune
```

---

## 🎯 Testing Scenarios

### Test 1: Container-to-Container Communication

```bash
# From front-end to database
docker compose exec front-end sh -c 'apk add postgresql-client && psql "$DATABASE_DSN" -c "SELECT version();"'
```

### Test 2: Host-to-Container Communication (IPv6)

```bash
# Test PostgreSQL
psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable" -c "SELECT current_timestamp;"

# Test HTTP
curl http://[::1]:4000
```

### Test 3: Host-to-Container Communication (IPv4)

```bash
# Test PostgreSQL
psql "postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable" -c "SELECT current_timestamp;"

# Test HTTP
curl http://localhost:4000
```

### Test 4: Check IPv6 Addresses

```bash
# Database IPv6 address
docker compose exec database ip -6 addr show eth0 | grep fd00

# Front-end IPv6 address
docker compose exec front-end ip -6 addr show eth0 | grep fd00

# Back-end IPv6 address
docker compose exec back-end ip -6 addr show eth0 | grep fd00
```

---

## 📚 Next Steps

1. **Read the detailed guide**: [IPv6-SETUP.md](./IPv6-SETUP.md)
2. **Configure your application**: Edit `.env` with production values
3. **Set up monitoring**: Consider adding health check endpoints
4. **Deploy**: Follow deployment guide for production setup

---

## 💡 Tips

- Always use service names (`database:5432`) in Docker environment
- Use `[::1]` for IPv6 localhost connections from host
- Use `127.0.0.1` or `localhost` for IPv4 (backward compatible)
- Run `make test-ipv6` after any network changes
- Check logs with `docker compose logs -f` if issues occur
- Use `make help` to see all available commands

---

## ✅ Success Indicators

Your setup is working correctly if:

1. ✅ `docker compose ps` shows all services as `Up` or `Up (healthy)`
2. ✅ `make test-ipv6` completes without errors
3. ✅ You can access http://localhost:4000 in your browser
4. ✅ `make verify-db-ipv6` shows PostgreSQL listening on IPv6
5. ✅ `make show-container-ips` shows both IPv4 and IPv6 addresses

---

**Need help?** Check [IPv6-SETUP.md](./IPv6-SETUP.md) for detailed troubleshooting or open an issue.

