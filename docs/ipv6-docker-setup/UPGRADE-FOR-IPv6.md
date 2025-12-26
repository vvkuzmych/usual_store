# 🚀 Upgrade Docker Desktop for IPv6 Support

## ✨ Good News!

**Docker Desktop 4.42** (released June 2025) added **native IPv6 support** for macOS!

**Your current version**: 4.38.0 ❌  
**Required version**: 4.42+ ✅

---

## 📋 Upgrade Steps

### **1. Update Docker Desktop**

```bash
# Option A: Use Docker Desktop UI
# 1. Click Docker icon in menu bar
# 2. Select "Check for Updates"
# 3. Download and install version 4.42+

# Option B: Download manually
# Visit: https://www.docker.com/products/docker-desktop/
# Download latest version for Apple Silicon
```

### **2. Verify New Version**

```bash
# After upgrade, check version
docker --version
# Should show: Docker version 27.x.x or higher

# Check Docker Desktop version
system_profiler SPApplicationsDataType | grep -A 5 "Docker"
# Should show: Version: 4.42 or higher
```

### **3. Enable IPv6 in Docker Desktop**

1. **Open Docker Desktop**
   - Click Docker icon in menu bar
   - Select "Settings" or "Preferences"

2. **Navigate to Network Tab**
   - Click "Network" in left sidebar
   - This is a NEW feature in 4.42+

3. **Configure Networking Mode**
   - Set **Default networking mode** to:
     - ✅ **Dual IPv4/IPv6** (recommended)
     - or **IPv6 only** (advanced)

4. **Configure DNS Resolution** (optional)
   - Set **DNS resolution behavior**:
     - ✅ **Auto** (recommended)
     - or customize based on your needs

5. **Apply Changes**
   - Click "Apply & Restart"
   - Docker Desktop will restart

---

## 🔧 Update Your Project Configuration

### **Update docker-compose.yml**

```yaml
services:
  front-end:
    # ... existing config ...
    networks:
      - usualstore_network
  
  back-end:
    # ... existing config ...
    networks:
      - usualstore_network
  
  database:
    # ... existing config ...
    ports:
      - "[::1]:5433:5432"       # ✅ IPv6 localhost (now works!)
      - "127.0.0.1:5433:5432"   # ✅ IPv4 localhost (still works)
    networks:
      - usualstore_network
  
  invoice:
    # ... existing config ...
    networks:
      - usualstore_network

networks:
  usualstore_network:
    driver: bridge
    enable_ipv6: true  # ✅ Enable IPv6
    ipam:
      driver: default
      config:
        - subnet: 172.22.0.0/16
          gateway: 172.22.0.1
        - subnet: 2001:db8:1::/64  # ✅ Add IPv6 subnet
          gateway: 2001:db8:1::1

volumes:
  db_data:
```

### **Update database.yml**

```yaml
development:
  dialect: postgres
  database: usualstore
  user: postgres
  password: password
  host: "::1"  # ✅ Can now use IPv6!
  port: 5433
  sslmode: disable

development_ipv4:
  dialect: postgres
  database: usualstore
  user: postgres
  password: password
  host: "127.0.0.1"  # IPv4 fallback
  port: 5433
  sslmode: disable

docker:
  dialect: postgres
  database: usualstore
  user: postgres
  password: password
  host: database  # Service name (works with both)
  port: 5432
  sslmode: disable
```

---

## ✅ Test IPv6 After Upgrade

### **1. Restart Your Services**

```bash
cd /Users/vkuzm/Projects/usual_store
docker compose down
docker compose up -d
```

### **2. Verify IPv6 Network**

```bash
# Check network has IPv6 enabled
docker network inspect usual_store_usualstore_network | grep -A 5 IPv6

# Should show:
# "EnableIPv6": true,
# IPv6 subnet configuration
```

### **3. Check Container IPv6 Addresses**

```bash
# Database IPv6 address
docker inspect usual_store-database-1 | grep GlobalIPv6Address

# Should show something like:
# "GlobalIPv6Address": "2001:db8:1::2"
```

### **4. Test IPv6 Connection from Mac**

```bash
# Test database via IPv6
psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable" -c "SELECT version();"

# Should connect successfully! ✅

# Test web app via IPv6
curl http://[::1]:4000

# Should work! ✅
```

### **5. Test Container-to-Container**

```bash
# From front-end to database
docker compose exec front-end getent hosts database

# Should show both IPv4 and IPv6 addresses
```

---

## ⚠️ Known Issues & Workarounds

### **Issue 1: DNS Resolution Problems**

**Problem**: `localhost` resolves to `::1` when you expect `127.0.0.1`

**Workaround**:
- Use explicit IP addresses: `[::1]` or `127.0.0.1`
- Configure DNS resolution behavior in Docker Desktop settings
- Use "Auto" mode in DNS settings (recommended)

### **Issue 2: Port Binding Issues**

**Problem**: Services fail to bind to `[::]:port`

**Workaround**:
- Explicitly bind to both: `[::1]:port` AND `127.0.0.1:port`
- Or use `0.0.0.0:port` (binds to all interfaces)

### **Issue 3: Performance Issues**

**Problem**: Slower networking with IPv6 enabled

**Workaround**:
- Switch to "Dual IPv4/IPv6" mode (not IPv6-only)
- Monitor resource usage
- Report issues to Docker GitHub

### **Issue 4: Connection Refused**

**Problem**: Still can't connect via IPv6 after upgrade

**Check**:
1. Docker Desktop settings → Network → Dual IPv4/IPv6 enabled?
2. Firewall allowing IPv6 connections?
3. Container has IPv6 address assigned?
4. Network has IPv6 subnet configured?

**Debug**:
```bash
# Check Docker daemon config
docker info | grep -i ipv6

# Check network settings
docker network inspect usual_store_usualstore_network

# Check container networking
docker inspect usual_store-database-1 | grep -A 10 Networks
```

---

## 📚 Resources

- [Docker Desktop 4.42 Release Notes](https://www.infoq.com/news/2025/07/docker-desktop-442/)
- [Docker Desktop Networking Guide](https://docs.docker.com/desktop/features/networking/networking-how-tos/)
- [Docker Engine IPv6 Documentation](https://docs.docker.com/engine/daemon/ipv6/)
- [GitHub Issues for macOS IPv6](https://github.com/docker/for-mac/issues/7269)

---

## 🎯 Quick Checklist

Before upgrading:
- [ ] Backup your current docker-compose.yml
- [ ] Backup your .env file
- [ ] Note your current working configuration

After upgrading:
- [ ] Docker Desktop version 4.42+ installed
- [ ] Network settings configured (Dual IPv4/IPv6)
- [ ] docker-compose.yml updated with IPv6 subnet
- [ ] Services restarted
- [ ] IPv6 connectivity tested
- [ ] Container-to-container communication verified

---

## 💡 Recommendation

1. **Backup first**: Save your current working configuration
2. **Upgrade Docker Desktop** to 4.42 or later
3. **Test thoroughly**: IPv6 support is new and may have issues
4. **Keep IPv4 working**: Don't remove IPv4 configuration yet
5. **Monitor stability**: Watch for networking issues
6. **Report bugs**: If you find issues, report to Docker GitHub

---

## 🚀 Next Steps

```bash
# 1. Upgrade Docker Desktop
# Download from https://www.docker.com/products/docker-desktop/

# 2. Enable IPv6 in Settings
# Docker Desktop → Settings → Network → Dual IPv4/IPv6

# 3. Update your project
cd /Users/vkuzm/Projects/usual_store
# Edit docker-compose.yml (add IPv6 config)

# 4. Restart services
docker compose down
docker compose up -d

# 5. Test IPv6
psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable"
curl http://[::1]:4000
```

**Good luck with IPv6!** 🎉

