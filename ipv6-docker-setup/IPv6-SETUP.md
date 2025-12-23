# IPv6 Configuration for Usual Store

This document explains the IPv6 configuration for the Usual Store project.

## 📋 Overview

The project now supports both **IPv4** and **IPv6** networking through Docker. Services can communicate using either protocol, with IPv6 enabled by default in the Docker network.

---

## 🌐 Network Architecture

### Docker Network Configuration

```yaml
usualstore_network:
  driver: bridge
  enable_ipv6: true
  ipam:
    config:
      - subnet: 172.20.0.0/16        # IPv4 subnet
        gateway: 172.20.0.1
      - subnet: fd00:dead:beef::/48   # IPv6 subnet (ULA)
        gateway: fd00:dead:beef::1
```

### Service IP Addressing

Each Docker service gets assigned:
- **IPv4 address**: `172.20.0.x`
- **IPv6 address**: `fd00:dead:beef::x`

---

## 🔌 Connection Methods

### 1. **Container-to-Container Communication** (Recommended)

Use Docker service names - they resolve to both IPv4 and IPv6:

```bash
DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
```

✅ **Pros**: Protocol-agnostic, works with both IPv4 and IPv6  
✅ **Docker DNS**: Handles resolution automatically

### 2. **Host to Container (IPv6)**

From your Mac to the Docker container via IPv6 localhost:

```bash
DATABASE_DSN=postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable
```

⚠️ **Note**: Use brackets `[::1]` for IPv6 addresses in connection strings

### 3. **Host to Container (IPv4)** - Backward Compatible

Traditional IPv4 localhost connection:

```bash
DATABASE_DSN=postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable
```

---

## 🚀 Getting Started

### Prerequisites

1. **Docker Desktop**: Ensure IPv6 is enabled
   - macOS: Docker Desktop → Preferences → Docker Engine
   - Add to daemon config:
     ```json
     {
       "ipv6": true,
       "fixed-cidr-v6": "fd00::/80"
     }
     ```

2. **System IPv6**: Verify your Mac has IPv6 enabled
   ```bash
   ping6 ::1
   ```

### Starting the Application

1. **Copy environment file**:
   ```bash
   cp env.example .env
   # Edit .env with your actual values
   ```

2. **Start services**:
   ```bash
   docker compose up -d
   ```

3. **Verify IPv6 connectivity**:
   ```bash
   # Check if PostgreSQL listens on IPv6
   docker compose exec database ss -tln | grep ::1
   
   # Test IPv6 connection from host
   psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable" -c "\conninfo"
   ```

---

## 🔍 Troubleshooting

### Issue: "connection refused" on IPv6

**Symptoms**:
```
dial tcp [::1]:5432: connect: connection refused
```

**Solutions**:

1. **Check Docker IPv6 support**:
   ```bash
   docker network inspect usualstore_usualstore_network | grep EnableIPv6
   ```
   Should show: `"EnableIPv6": true`

2. **Verify PostgreSQL listens on IPv6**:
   ```bash
   docker compose logs database | grep "listening on"
   ```
   Should show: `listening on IPv6 address "::"`

3. **Check port binding**:
   ```bash
   docker compose ps
   ```
   Should show: `[::1]:5432->5432/tcp` and `127.0.0.1:5432->5432/tcp`

### Issue: "network not found"

**Symptoms**:
```
ERROR: network usualstore_network declared as external, but could not be found
```

**Solution**:
```bash
# Recreate the network
docker compose down
docker compose up -d
```

### Issue: Services can't communicate

**Symptoms**:
```
dial tcp: lookup database: no such host
```

**Solution**:

1. Ensure all services are on the same network:
   ```bash
   docker compose ps --format json | jq -r '.[].Networks'
   ```

2. Check network connectivity:
   ```bash
   docker compose exec front-end ping -c 2 database
   docker compose exec front-end ping6 -c 2 database
   ```

---

## 📊 Verifying IPv6 Configuration

### Check Container IPv6 Address

```bash
# Get front-end IPv6 address
docker compose exec front-end ip -6 addr show

# Get database IPv6 address
docker compose exec database ip -6 addr show
```

### Test Database Connection

```bash
# From host machine (IPv6)
psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable" -c "SELECT version();"

# From host machine (IPv4)
psql "postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable" -c "SELECT version();"

# From within front-end container
docker compose exec front-end sh -c 'apk add postgresql-client && psql "$DATABASE_DSN" -c "SELECT version();"'
```

### Monitor Network Traffic

```bash
# IPv6 connections to PostgreSQL
docker compose exec database ss -tn6 | grep 5432

# IPv4 connections to PostgreSQL
docker compose exec database ss -tn4 | grep 5432
```

---

## 🛠️ Advanced Configuration

### Custom IPv6 Subnet

Edit `docker-compose.yml`:

```yaml
networks:
  usualstore_network:
    enable_ipv6: true
    ipam:
      config:
        - subnet: fd12:3456:789a::/48  # Your custom ULA
          gateway: fd12:3456:789a::1
```

### Disable IPv6 (Fallback to IPv4 only)

1. Remove IPv6 subnet from `docker-compose.yml`:
   ```yaml
   networks:
     usualstore_network:
       enable_ipv6: false
       ipam:
         config:
           - subnet: 172.20.0.0/16
             gateway: 172.20.0.1
   ```

2. Update port bindings in `database` service:
   ```yaml
   ports:
     - "5432:5432"  # Binds to 0.0.0.0 (IPv4 only)
   ```

### Force IPv6-only Connections

Update `DATABASE_DSN` to use the container's actual IPv6 address:

```bash
# Get database IPv6 address
DB_IPV6=$(docker compose exec database ip -6 addr show eth0 | grep 'inet6 fd00' | awk '{print $2}' | cut -d/ -f1)

# Update connection string
DATABASE_DSN="postgres://postgres:password@[$DB_IPV6]:5432/usualstore?sslmode=disable"
```

---

## 📝 Configuration Files Reference

### `docker-compose.yml`
- Defines IPv6-enabled network
- Configures PostgreSQL with IPv6 listening
- Binds ports to both IPv4 and IPv6 localhost

### `database.yml`
- Development: Uses `::1` (IPv6 localhost)
- Docker: Uses `database` service name

### `postgres-init/01-enable-ipv6.sql`
- Initializes PostgreSQL with proper permissions
- Logs IPv6 configuration status

### `env.example`
- Template for environment variables
- Includes IPv6 DSN examples

---

## 🎯 Best Practices

1. **Use Service Names**: Prefer `database:5432` over IP addresses
2. **Bracket Notation**: Always use `[::1]` for IPv6 addresses in URIs
3. **Health Checks**: Include IPv6 in health check configurations
4. **Testing**: Test both IPv4 and IPv6 connectivity
5. **Documentation**: Keep connection string formats documented

---

## 🔗 Useful Commands

```bash
# View all network details
docker compose exec database ip addr show

# Test DNS resolution (both IPv4 and IPv6)
docker compose exec front-end getent hosts database

# Check listening ports
docker compose exec database netstat -tulpn | grep postgres

# Monitor PostgreSQL logs
docker compose logs -f database

# Restart with clean state
docker compose down -v && docker compose up -d

# Check IPv6 routing
docker compose exec front-end ip -6 route
```

---

## 📚 Additional Resources

- [Docker IPv6 Documentation](https://docs.docker.com/config/daemon/ipv6/)
- [PostgreSQL Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)
- [IPv6 Addressing Guide](https://www.rfc-editor.org/rfc/rfc4291)
- [Unique Local Addresses (ULA)](https://www.rfc-editor.org/rfc/rfc4193)

---

## ✅ Summary

Your Usual Store application now supports:
- ✅ Dual-stack networking (IPv4 + IPv6)
- ✅ Container-to-container communication via service names
- ✅ Host access via both `[::1]` (IPv6) and `127.0.0.1` (IPv4)
- ✅ PostgreSQL listening on both protocols
- ✅ Backward compatibility with IPv4-only setups

For questions or issues, check the troubleshooting section or create an issue in the repository.

