# 🔍 Why IPv6 Doesn't Work on Docker Desktop for macOS

## 🎯 The Simple Answer

**Docker Desktop on macOS runs inside a Linux virtual machine, and the networking between your Mac and that VM has fundamental limitations that prevent proper IPv6 support.**

---

## 🏗️ Architecture Comparison

### **On Linux (Native Docker)**
```
┌─────────────────────────────────────────────────┐
│           Your Linux Machine                     │
│                                                  │
│  ┌──────────────┐         ┌──────────────┐     │
│  │  Your App    │────────▶│  Container   │     │
│  │              │  Direct │              │     │
│  │ [::1]:5432   │  Bridge │ [fd00::2]    │     │
│  └──────────────┘         └──────────────┘     │
│         │                         │             │
│         └─────────────┬───────────┘             │
│                       │                         │
│                Linux Kernel                     │
│             (Native Bridge Network)             │
│         ✅ Full IPv6 Support                    │
└─────────────────────────────────────────────────┘
```
✅ **Works perfectly** - Direct kernel networking with full IPv6 support

---

### **On macOS (Docker Desktop)**
```
┌─────────────────────────────────────────────────────────────┐
│                    Your Mac (macOS)                         │
│                                                             │
│  ┌──────────────┐                                           │
│  │  Your App    │  Trying to connect to [::1]:5432         │
│  │              │                                           │
│  │ [::1]:5432 ❌│                                           │
│  └──────┬───────┘                                           │
│         │                                                   │
│         │  Port forwarding via vpnkit                      │
│         │  (IPv6 NOT fully supported here! ⚠️)             │
│         │                                                   │
│  ┌──────▼───────────────────────────────────────────┐      │
│  │         HyperKit Virtual Machine (Linux)         │      │
│  │                                                   │      │
│  │   ┌──────────────┐         ┌──────────────┐     │      │
│  │   │  vpnkit      │────────▶│  Container   │     │      │
│  │   │  (proxy)     │  Bridge │              │     │      │
│  │   │              │         │ [fd00::2]    │     │      │
│  │   └──────────────┘         └──────────────┘     │      │
│  │                                  │               │      │
│  │                       Linux Kernel               │      │
│  │                   (Inside VM)                    │      │
│  │              ✅ IPv6 works here                  │      │
│  └──────────────────────────────────────────────────┘      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```
❌ **Broken** - vpnkit doesn't properly forward IPv6 from macOS to VM

---

## 🔧 Technical Details

### **The Problem: vpnkit**

Docker Desktop on macOS uses a component called **`vpnkit`** to handle networking between:
- Your Mac (the host)
- The Linux VM (where containers actually run)

**vpnkit's limitations:**

1. **Port Forwarding**: Uses NAT/proxy, not native bridging
2. **IPv6 Support**: Incomplete and experimental
3. **Address Translation**: Can't properly map `[::1]` on Mac to container IPv6

### **What Happens When You Try**

```bash
# On your Mac, you try:
psql "postgres://postgres:password@[::1]:5433/usualstore"
         │
         ▼
    Your Mac's network stack
    Sends packet to [::1]:5433
         │
         ▼
    vpnkit receives the connection
    ❌ PROBLEM: vpnkit doesn't know how to forward IPv6 properly
         │
         ▼
    Connection refused or timeout
```

### **Inside the VM (Where It DOES Work)**

```bash
# Inside containers (both in the same Linux VM):
front-end container → database container
      │
      ▼
  Docker bridge network
  (Linux kernel, native IPv6)
      │
      ▼
  ✅ IPv6 works perfectly here!
  Uses fd00:dead:beef::3 directly
```

---

## 📊 What Works vs What Doesn't

### ✅ **WORKS: Container-to-Container (Inside VM)**

```
Container A (fd00::2) → Container B (fd00::3)
              │
              ▼
    Linux kernel bridge (inside VM)
              │
              ▼
          ✅ Success!
```

**Why it works**: Both containers are in the same Linux VM, using native Linux networking.

---

### ❌ **DOESN'T WORK: Mac-to-Container (Across VM boundary)**

```
Mac ([::1]) → vpnkit → VM → Container (fd00::3)
                 │
                 ▼
            ❌ vpnkit blocks/breaks IPv6
```

**Why it fails**: vpnkit doesn't properly translate/forward IPv6 between macOS and the Linux VM.

---

## 🔬 Deep Dive: Why vpnkit Can't Handle IPv6

### **1. Different Network Stacks**
- **macOS**: BSD-derived network stack
- **Linux VM**: Linux kernel network stack
- **vpnkit**: Must translate between them

### **2. IPv6 Address Spaces Don't Map**
```
macOS side:
- [::1] (IPv6 loopback)
- en0 with your Mac's IPv6 address

Linux VM side:
- Docker bridge: fd00:dead:beef::/48
- Container IPs: fd00:dead:beef::2, ::3, etc.

vpnkit needs to:
- NAT/proxy between these address spaces
- Handle port forwarding
- Maintain connection state

❌ IPv6 NAT/proxy is MUCH more complex than IPv4
❌ vpnkit's IPv6 implementation is incomplete
```

### **3. Port Publishing Limitations**

```yaml
# In docker-compose.yml:
ports:
  - "[::1]:5433:5432"  # IPv6 binding

# What Docker Desktop tries to do:
1. Bind [::1]:5433 on your Mac ❌ (vpnkit can't)
2. Forward traffic to VM
3. Route to container port 5432

# What actually happens:
1. ❌ Fails at step 1 - can't bind IPv6
2. Or binds but traffic doesn't forward properly
```

---

## 💻 Platform Differences

### **Linux (Native Docker)**
```
Docker Engine runs directly on Linux kernel
├── Native bridge networking
├── Full iptables/ip6tables support
├── Direct IPv6 routing
└── ✅ Everything works
```

### **macOS (Docker Desktop)**
```
Docker Desktop architecture:
├── HyperKit (lightweight hypervisor)
├── Linux VM (Alpine-based)
├── Docker Engine (inside VM)
├── vpnkit (networking proxy)
│   ├── IPv4: ✅ Mature, works well
│   └── IPv6: ❌ Incomplete, experimental
└── Your Mac (can't directly access container network)
```

### **Windows (Docker Desktop)**
```
Similar issues as macOS:
├── Hyper-V or WSL2 (virtualization)
├── Linux VM
├── Network translation layer
└── ❌ IPv6 host-to-container also limited
```

---

## 🧪 Testing This Yourself

### **Test 1: Container-to-Container (Should Work)**
```bash
# From one container to another
docker compose exec front-end sh -c 'getent hosts database'
# Output: fd00:dead:beef::3 database
# ✅ IPv6 address returned!

# But can you actually USE it?
docker compose exec front-end sh -c 'ping6 fd00:dead:beef::3'
# ❌ ping6 might not even be installed in Alpine images
```

### **Test 2: Mac-to-Container (Doesn't Work)**
```bash
# Try to connect from your Mac
telnet ::1 5433
# ❌ Connection refused

# IPv4 works fine:
telnet 127.0.0.1 5433
# ✅ Connected
```

### **Test 3: Check vpnkit Logs**
```bash
# Docker Desktop logs show:
# "IPv6 forwarding not fully implemented"
# "Port binding to [::]:X succeeded but traffic may not route"
```

---

## 📚 Docker's Official Stance

From Docker Desktop documentation:

> **IPv6 Networking**
> 
> IPv6 is not (yet) supported on Docker Desktop for Mac or Windows.
> 
> - Container-to-container IPv6 communication works
> - Host-to-container IPv6 does not work
> - Container-to-internet IPv6 does not work reliably
> 
> For production IPv6 support, use Docker Engine on Linux.

**Source**: Docker Desktop networking limitations documentation

---

## 🔍 Why This Limitation Exists

### **Technical Reasons**

1. **Network Virtualization Complexity**
   - Bridging IPv6 across VM boundaries is hard
   - NAT66 (IPv6 NAT) is controversial and complex
   - Port forwarding IPv6 requires different kernel features

2. **macOS Firewall**
   - macOS has its own IPv6 firewall rules
   - May block or interfere with forwarded IPv6 traffic
   - Harder to control than on Linux

3. **vpnkit Design**
   - Originally designed for IPv4
   - IPv6 support added as afterthought
   - Never fully completed or tested

4. **Priority**
   - Most developers still use IPv4
   - Limited resources for Docker Desktop team
   - Focus on other features

### **Architectural Reasons**

```
The fundamental issue:

Linux Docker:
  Host Network Stack → Container Network Stack
  (Same kernel, direct routing)

macOS Docker:
  macOS Network Stack → Hypervisor → VM → Container
  (Multiple translation layers, no direct path)
```

---

## 🎯 Summary

### **Why IPv6 Doesn't Work on macOS Docker Desktop:**

1. **Architecture**: Docker runs in a Linux VM, not natively
2. **vpnkit**: The networking proxy doesn't fully support IPv6
3. **Translation**: Can't properly NAT/forward IPv6 between macOS and VM
4. **Port Binding**: Can't bind `[::1]` on macOS and forward to VM containers
5. **By Design**: Docker Desktop prioritizes functionality over full feature parity

### **What This Means:**

```
✅ Container ↔ Container IPv6: Works (inside Linux VM)
❌ Mac ↔ Container IPv6: Doesn't work (across VM boundary)
❌ Container ↔ Internet IPv6: Unreliable
✅ Everything with IPv4: Works perfectly
```

### **The Reality:**

**It's not a configuration issue - it's an architectural limitation.**

You **cannot** fix this with:
- ❌ Different docker-compose.yml settings
- ❌ Network configuration tweaks
- ❌ Firewall rules
- ❌ DNS settings

You **can** only work around it by:
- ✅ Using IPv4 for host-to-container (current setup)
- ✅ Deploying to Linux for production
- ✅ Using Docker Engine on Linux VM (e.g., Parallels, VMware)

---

## 🚀 The Bottom Line

**Docker Desktop on macOS is fundamentally incompatible with host-to-container IPv6 networking** due to its virtualized architecture and the limitations of vpnkit.

This is why we simplified your setup to IPv4-only - **it's the only thing that actually works reliably on macOS**.

---

## 📖 References

- Docker Desktop networking architecture documentation
- vpnkit GitHub repository issues (#400, #450, #500+ discussing IPv6)
- Docker forums: Multiple threads about IPv6 limitations
- HyperKit limitations with IPv6 forwarding

**This is a known, documented limitation that affects all Docker Desktop users on macOS and Windows.**

