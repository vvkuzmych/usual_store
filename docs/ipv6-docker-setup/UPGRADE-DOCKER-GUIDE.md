# 🚀 Docker Desktop Upgrade Guide for macOS

## Current Situation
- **Your version**: Docker Desktop 4.38.0
- **Target version**: Docker Desktop 4.42+ (for IPv6 support)
- **Your Mac**: Apple Silicon

---

## ✅ **Method 1: Automatic Update (Easiest)**

### **Step 1: Check for Updates**
1. Click the **Docker icon** in your menu bar (top-right corner)
2. Look for **"Check for Updates"** option
3. Click it

### **Step 2: Download Update**
- If update is available, you'll see:
  ```
  New version available: 4.42.x
  [ Download Update ]
  ```
- Click **"Download Update"**
- Wait for download (may take 2-5 minutes)

### **Step 3: Install Update**
- When download completes, click **"Install Update"**
- Docker Desktop will close
- Installation happens automatically
- Docker Desktop will restart

### **Step 4: Verify**
```bash
docker --version
# Should show: Docker version 27.x.x or higher

# Check Docker Desktop version
system_profiler SPApplicationsDataType | grep -A 3 "Docker:"
# Should show: Version: 4.42 or higher
```

---

## 📥 **Method 2: Manual Download (If auto-update fails)**

### **Step 1: Download Latest Version**

**For Apple Silicon (your Mac):**
```
https://desktop.docker.com/mac/main/arm64/Docker.dmg
```

**Or visit:**
```
https://www.docker.com/products/docker-desktop/
```

### **Step 2: Quit Docker Desktop**
```bash
# Quit Docker Desktop completely
# Menu bar → Docker icon → Quit Docker Desktop

# Or use command:
osascript -e 'quit app "Docker"'
```

### **Step 3: Install New Version**
1. Open the downloaded **Docker.dmg** file
2. Drag **Docker.app** to **Applications** folder
3. When prompted "Docker.app already exists", click **Replace**
4. Authenticate with your Mac password if asked
5. Wait for copy to complete
6. Eject the Docker.dmg

### **Step 4: Launch Docker Desktop**
```bash
# Open Docker Desktop
open -a Docker

# Or double-click Docker.app in Applications folder
```

### **Step 5: Verify Installation**
```bash
# Check Docker version
docker --version

# Check Docker Desktop version
docker version --format '{{.Server.Version}}'

# Full version info
docker info | grep -i version
```

---

## 🍺 **Method 3: Using Homebrew**

If you originally installed Docker Desktop with Homebrew:

### **Step 1: Update Homebrew**
```bash
brew update
```

### **Step 2: Upgrade Docker Desktop**
```bash
brew upgrade --cask docker
```

### **Step 3: Launch Docker**
```bash
open -a Docker
```

### **Step 4: Verify**
```bash
docker --version
```

---

## ⚙️ **After Upgrading: Enable IPv6**

Once you have version 4.42+:

### **Step 1: Open Docker Desktop Settings**
1. Click Docker icon in menu bar
2. Select **"Settings"** (or **"Preferences"**)

### **Step 2: Navigate to Network**
1. Click **"Network"** in the left sidebar
2. You should see new IPv6 options (this is new in 4.42!)

### **Step 3: Enable IPv6**
1. Under **"Default networking mode"**, select:
   - **Dual IPv4/IPv6** ✅ (Recommended)
   - Or **IPv6 only** (Advanced)

2. Under **"DNS resolution behavior"**, select:
   - **Auto** ✅ (Recommended)

### **Step 4: Apply Changes**
1. Click **"Apply & Restart"**
2. Docker Desktop will restart
3. Wait for Docker to fully start

### **Step 5: Verify IPv6 is Enabled**
```bash
# Check Docker daemon info
docker info | grep -i ipv6

# Should show something about IPv6 being enabled
```

---

## 🧪 **Test Your Upgrade**

### **1. Verify Docker Version**
```bash
docker --version
# Expected: Docker version 27.x.x, build xxxxx

docker version
# Check both Client and Server versions
```

### **2. Verify Docker Desktop Version**
```bash
system_profiler SPApplicationsDataType | grep -A 5 "Docker:"
# Should show: Version: 4.42 or higher
```

### **3. Check Docker is Running**
```bash
docker ps
# Should list running containers (or show empty table)
```

### **4. Test IPv6 (After enabling in settings)**
```bash
# Create a test network with IPv6
docker network create --ipv6 --subnet 2001:db8::/64 test-ipv6

# Run a test container
docker run --rm --network test-ipv6 alpine ping -c 2 2001:db8::1

# Clean up
docker network rm test-ipv6
```

---

## ⚠️ **Troubleshooting**

### **Problem: "Docker Desktop is already running"**
```bash
# Quit Docker completely first
osascript -e 'quit app "Docker"'

# Wait 10 seconds
sleep 10

# Then try installation again
```

### **Problem: "Cannot move to Applications"**
```bash
# Make sure Docker is completely quit
pkill -9 Docker

# Try installation again
```

### **Problem: "Docker daemon not responding after upgrade"**
```bash
# Reset Docker Desktop
# Menu bar → Docker icon → Troubleshoot → Reset to factory defaults

# Or restart your Mac
sudo shutdown -r now
```

### **Problem: "Containers won't start after upgrade"**
```bash
# Restart all containers
docker restart $(docker ps -aq)

# Or restart Docker Desktop
osascript -e 'quit app "Docker"'
sleep 5
open -a Docker
```

### **Problem: "Docker version still shows old version"**
```bash
# Clear Docker cache
rm -rf ~/Library/Caches/com.docker.docker

# Quit and restart Docker
osascript -e 'quit app "Docker"'
sleep 5
open -a Docker
```

---

## 📋 **Pre-Upgrade Checklist**

Before upgrading, make sure to:

- [ ] **Stop your running containers** (if you want to preserve state)
  ```bash
  docker compose -f /Users/vkuzm/Projects/usual_store/docker-compose.yml down
  ```

- [ ] **Backup your Docker data** (optional, but safe)
  ```bash
  # Backup volumes
  docker volume ls
  # Note important volumes
  ```

- [ ] **Save your current configuration**
  ```bash
  # Backup your project
  cp /Users/vkuzm/Projects/usual_store/docker-compose.yml \
     /Users/vkuzm/Projects/usual_store/docker-compose.yml.backup
  ```

- [ ] **Check disk space**
  ```bash
  df -h
  # Need at least 2-3 GB free
  ```

---

## 📋 **Post-Upgrade Checklist**

After upgrading:

- [ ] **Verify Docker version**
  ```bash
  docker --version  # Should be 27.x.x+
  ```

- [ ] **Check Docker Desktop version**
  ```bash
  system_profiler SPApplicationsDataType | grep -A 3 "Docker:"
  # Should show 4.42+
  ```

- [ ] **Enable IPv6** (Settings → Network → Dual IPv4/IPv6)

- [ ] **Restart your project**
  ```bash
  cd /Users/vkuzm/Projects/usual_store
  docker compose up -d
  ```

- [ ] **Test services**
  ```bash
  docker compose ps
  curl http://localhost:4000
  psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable" -c "SELECT 1;"
  ```

- [ ] **Test IPv6** (if enabled)
  ```bash
  # Update docker-compose.yml first (see UPGRADE-FOR-IPv6.md)
  docker compose down
  docker compose up -d
  
  # Test IPv6 connection
  psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable" -c "SELECT 1;"
  ```

---

## 🎯 **Quick Command Summary**

```bash
# Method 1: Auto-update (from Docker Desktop UI)
# Docker menu bar icon → Check for Updates → Install

# Method 2: Manual download
# Download from: https://desktop.docker.com/mac/main/arm64/Docker.dmg
# Install by dragging to Applications

# Method 3: Homebrew
brew upgrade --cask docker
open -a Docker

# After upgrade - verify
docker --version
system_profiler SPApplicationsDataType | grep -A 3 "Docker:"

# Enable IPv6
# Docker Desktop → Settings → Network → Dual IPv4/IPv6 → Apply & Restart

# Restart your project
cd /Users/vkuzm/Projects/usual_store
docker compose down
docker compose up -d

# Test everything
docker compose ps
curl http://localhost:4000
```

---

## 📚 **Additional Resources**

- [Docker Desktop Release Notes](https://docs.docker.com/desktop/release-notes/)
- [Docker Desktop for Mac](https://docs.docker.com/desktop/install/mac-install/)
- [Docker IPv6 Documentation](https://docs.docker.com/engine/daemon/ipv6/)
- [IPv6 Setup Guide](./UPGRADE-FOR-IPv6.md) (in this repo)

---

## 💡 **Tips**

- **Don't worry**: Upgrading Docker Desktop is safe and won't delete your containers/volumes
- **Backup first**: Always a good idea to backup important data
- **Test after**: Verify everything works after upgrade
- **IPv6 is optional**: You can upgrade without enabling IPv6 immediately
- **Rollback**: You can always download and install an older version if needed

---

## 🚀 **Ready to Upgrade?**

1. ✅ Click Docker icon → Check for Updates
2. ✅ Or download from: https://desktop.docker.com/mac/main/arm64/Docker.dmg
3. ✅ Install, restart, verify
4. ✅ Enable IPv6 in Settings (optional)
5. ✅ Test your project

**Good luck with the upgrade!** 🎉

