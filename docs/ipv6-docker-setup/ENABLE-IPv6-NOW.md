# 🎉 Enable IPv6 on Docker Desktop 4.55.0

**Great news!** Your Docker Desktop 4.55.0 fully supports IPv6!

---

## ✅ **Step 1: Enable IPv6 in Docker Desktop Settings**

### **Open Docker Desktop Settings**

1. Click the **Docker icon** in your menu bar (top-right)
2. Select **"Settings"** or **"Preferences"**

### **Navigate to Network Settings**

1. In the left sidebar, click **"Network"**
2. You should see IPv6 configuration options

### **Enable Dual-Stack Networking**

1. Under **"Default networking mode"**, select:
   - ✅ **"Dual IPv4/IPv6"** (Recommended)
   
   This allows both IPv4 and IPv6 to work together.

2. Under **"DNS resolution behavior"**, select:
   - ✅ **"Auto"** (Recommended)
   
   This automatically handles DNS for both protocols.

### **Apply Changes**

1. Click **"Apply & Restart"**
2. Docker Desktop will restart (takes about 30 seconds)
3. Wait for Docker to fully start

---

## ✅ **Step 2: Restart Your Services**

```bash
# Navigate to your project
cd /Users/vkuzm/Projects/usual_store

# Stop current services
docker compose down

# Start with new IPv6 configuration
docker compose up -d

# Wait for services to be healthy
sleep 10

# Check status
docker compose ps
```

**Expected output:**
```
NAME                      STATUS
usual_store-database-1    Up (healthy)
usual_store-back-end-1    Up
usual_store-front-end-1   Up
usual_store-invoice-1     Up
```

---

## ✅ **Step 3: Verify IPv6 is Working**

### **Test 1: Check Network Configuration**

```bash
# Verify IPv6 is enabled on the network
docker network inspect usual_store_usualstore_network | grep -A 5 IPv6
```

**Expected output:**
```json
"EnableIPv6": true,
"IPAM": {
    ...
    "Config": [
        ...
        {
            "Subnet": "2001:db8:1::/64",
            "Gateway": "2001:db8:1::1"
        }
    ]
}
```

### **Test 2: Check Container IPv6 Addresses**

```bash
# Check database container
docker inspect usual_store-database-1 | grep -E "IPAddress|GlobalIPv6Address"
```

**Expected output:**
```json
"IPAddress": "172.22.0.3",
"GlobalIPv6Address": "2001:db8:1::3",
```

✅ If you see an IPv6 address, it's working!

### **Test 3: Test IPv6 Connection from Your Mac**

```bash
# Test database connection via IPv6
psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable" -c "SELECT '✅ IPv6 works!' as status, version();"
```

**Expected output:**
```
     status      | version
-----------------+---------
 ✅ IPv6 works!  | PostgreSQL 15.x ...
```

✅ If this works, IPv6 is fully functional!

### **Test 4: Test IPv4 Still Works (Backward Compatibility)**

```bash
# Test database connection via IPv4
psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable" -c "SELECT '✅ IPv4 works!' as status;"
```

**Expected output:**
```
     status      
-----------------
 ✅ IPv4 works!
```

✅ Both protocols should work!

### **Test 5: Test Web Application**

```bash
# Test front-end via IPv6
curl -I http://[::1]:4000

# Test front-end via IPv4
curl -I http://localhost:4000

# Both should return HTTP 200 or similar
```

---

## 🎯 **What You Now Have**

```
✅ Docker Desktop 4.55.0 with native IPv6 support
✅ Dual-stack networking (IPv4 + IPv6)
✅ IPv6 addresses assigned to all containers
✅ IPv6 localhost binding working ([::1]:5433)
✅ Container-to-container IPv6 communication
✅ Backward compatible with IPv4
```

---

## 🌐 **Connection Strings**

### **From Your Mac (Host) - Both Work!**

**IPv6:**
```bash
# Database
postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable

# Web app
http://[::1]:4000

# API
http://[::1]:4001
```

**IPv4 (Still works):**
```bash
# Database
postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable

# Web app
http://localhost:4000

# API
http://localhost:4001
```

### **Inside Docker (Container-to-Container)**

```bash
# Use service name - automatically uses best protocol
DATABASE_DSN=postgres://postgres:password@database:5432/usualstore?sslmode=disable
```

---

## 📊 **Network Configuration Summary**

```
Network: usual_store_usualstore_network
├── Type: Bridge with dual-stack
├── IPv4: 172.22.0.0/16
│   ├── Gateway: 172.22.0.1
│   └── Example IPs: 172.22.0.2, 172.22.0.3, ...
└── IPv6: 2001:db8:1::/64
    ├── Gateway: 2001:db8:1::1
    └── Example IPs: 2001:db8:1::2, 2001:db8:1::3, ...

Port Bindings:
├── Database: 127.0.0.1:5433 AND [::1]:5433 → container:5432
├── Front-end: 0.0.0.0:4000 (all interfaces, both protocols)
├── Back-end: 0.0.0.0:4001 (all interfaces, both protocols)
└── Invoice: 0.0.0.0:5000 (all interfaces, both protocols)
```

---

## ⚠️ **Troubleshooting**

### **Issue: "EnableIPv6": false**

**Problem:** Network doesn't have IPv6 enabled

**Solution:**
```bash
# Recreate the network
docker compose down
docker network prune -f
docker compose up -d
```

### **Issue: No IPv6 address on containers**

**Problem:** Containers not getting IPv6 addresses

**Solution:**
```bash
# Check Docker Desktop settings
# Settings → Network → Ensure "Dual IPv4/IPv6" is selected

# Restart Docker Desktop
osascript -e 'quit app "Docker"'
sleep 5
open -a Docker

# Restart services
cd /Users/vkuzm/Projects/usual_store
docker compose down
docker compose up -d
```

### **Issue: Can't connect via [::1]:5433**

**Problem:** IPv6 localhost not working

**Solution:**
1. Check Docker Desktop version: `docker --version` (should be 27.x+)
2. Check Docker Desktop app version (should be 4.42+)
3. Verify Settings → Network → "Dual IPv4/IPv6" is enabled
4. Test if IPv6 is available on your Mac: `ping6 ::1`
5. Check firewall isn't blocking IPv6

### **Issue: "Connection refused" on IPv6**

**Problem:** Port binding not working

**Solution:**
```bash
# Check what's listening
lsof -nP -iTCP:5433 -sTCP:LISTEN

# Should show Docker binding to both [::1]:5433 and 127.0.0.1:5433

# If not, restart Docker Desktop and services
```

---

## 🧪 **Complete Test Script**

```bash
#!/bin/bash
echo "=== Testing IPv6 Configuration ==="

echo -e "\n1. Docker Desktop Version:"
docker --version
system_profiler SPApplicationsDataType | grep -A 2 "Docker:" | head -3

echo -e "\n2. Network IPv6 Status:"
docker network inspect usual_store_usualstore_network | grep EnableIPv6

echo -e "\n3. Container IPv6 Addresses:"
docker inspect usual_store-database-1 | grep GlobalIPv6Address

echo -e "\n4. Test IPv6 Connection:"
psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable" -c "SELECT 'IPv6 ✅' as status;" 2>&1 | head -5

echo -e "\n5. Test IPv4 Connection (backward compat):"
psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable" -c "SELECT 'IPv4 ✅' as status;" 2>&1 | head -5

echo -e "\n6. Test Web App (IPv6):"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://[::1]:4000

echo -e "\n7. Test Web App (IPv4):"
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:4000

echo -e "\n=== All tests complete! ==="
```

Save this as `test-ipv6.sh` and run:
```bash
chmod +x test-ipv6.sh
./test-ipv6.sh
```

---

## 📚 **What Changed**

### **docker-compose.yml**
```yaml
# Added:
networks:
  usualstore_network:
    enable_ipv6: true  # ← New!
    ipam:
      config:
        - subnet: 172.22.0.0/16
          gateway: 172.22.0.1
        - subnet: 2001:db8:1::/64  # ← New IPv6 subnet!
          gateway: 2001:db8:1::1

# Database ports:
ports:
  - "127.0.0.1:5433:5432"  # IPv4
  - "[::1]:5433:5432"      # ← New IPv6 binding!
```

### **Docker Desktop Settings**
```
Network → Default networking mode → "Dual IPv4/IPv6" ✅
```

---

## 🎉 **You're All Set!**

Your Docker setup now has:
- ✅ Full IPv6 support
- ✅ Dual-stack networking
- ✅ IPv6 addresses on all containers
- ✅ IPv6 localhost binding working
- ✅ Backward compatible with IPv4

**Congratulations! You now have modern dual-stack Docker networking!** 🚀

---

## 📖 **Next Steps**

- Test your application with both IPv4 and IPv6
- Update any documentation that mentions IPv6 limitations
- Consider using IPv6 in production environments

**Need help?** Check the troubleshooting section above or review the other documentation files.

