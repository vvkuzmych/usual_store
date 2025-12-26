# 🎨 Kubernetes Dashboard Guide

Complete guide to viewing and managing your Kubernetes pods and services using various dashboard options.

---

## 📋 Table of Contents

1. [Minikube Dashboard (Easiest)](#minikube-dashboard)
2. [k9s Terminal Dashboard](#k9s-terminal-dashboard)
3. [Lens IDE (Most Powerful)](#lens-ide)
4. [Command Line Monitoring](#command-line-monitoring)
5. [Dashboard Features](#dashboard-features)
6. [Troubleshooting](#troubleshooting)

---

## 🚀 Minikube Dashboard

### What is it?

The official Kubernetes web UI that runs in your browser. Perfect for beginners and visual learners.

### Quick Start

```bash
# Open dashboard (auto-opens browser)
minikube dashboard

# Run in background (keeps terminal free)
minikube dashboard &

# Get dashboard URL only (don't open browser)
minikube dashboard --url
```

### First Time Setup

```bash
# Enable dashboard addon (usually already enabled)
minikube addons enable dashboard

# Enable metrics for CPU/Memory graphs (optional but recommended)
minikube addons enable metrics-server
```

---

## 📊 How to Use the Dashboard

### Step 1: Access the Dashboard

Run:
```bash
minikube dashboard
```

Your browser will automatically open to something like:
```
http://127.0.0.1:xxxxx/api/v1/namespaces/kubernetes-dashboard/services/...
```

### Step 2: Select Your Namespace

1. Look at the **top of the page**
2. Find the **namespace dropdown** (usually shows "default")
3. Click it and select **"usualstore"**
4. Now you'll see only your application's resources

### Step 3: Navigate the Dashboard

**Left Sidebar Sections:**

```
📦 Workloads
   ├── Deployments    (your app components)
   ├── Pods           (running containers) ← START HERE
   ├── ReplicaSets    (pod management)
   └── StatefulSets   (database)

🌐 Service & Discovery
   ├── Services       (networking between pods)
   └── Ingresses      (external access)

⚙️  Config & Storage
   ├── ConfigMaps     (configuration)
   ├── Secrets        (passwords, keys)
   └── PersistentVolumeClaims (storage)

🔧 Cluster
   ├── Namespaces     (isolated environments)
   └── Nodes          (servers)
```

---

## 🔍 Viewing Your Application

### See All Pods

1. Click **"Pods"** in the left sidebar
2. You should see:

```
NAME                                  READY   STATUS
backend-xxx-xxx                       1/1     Running ✅
backend-xxx-xxx                       1/1     Running ✅
frontend-xxx-xxx                      1/1     Running ✅
frontend-xxx-xxx                      1/1     Running ✅
frontend-xxx-xxx                      1/1     Running ✅
database-0                            1/1     Running ✅
typescript-frontend-xxx-xxx           1/1     Running ✅
redux-frontend-xxx-xxx                1/1     Running ✅
support-frontend-xxx-xxx              1/1     Running ✅
ai-assistant-xxx-xxx                  1/1     Running ✅
invoice-xxx-xxx                       1/1     Running ✅
```

**Status Colors:**
- 🟢 **Green** = Running (healthy)
- 🟡 **Yellow** = Pending/Starting
- 🔴 **Red** = Error/Failed/CrashLoopBackOff

### View Pod Details

1. **Click on any pod name**
2. You'll see:
   - **Details tab**: Configuration, status, IP address
   - **Logs tab**: Real-time application logs
   - **Events tab**: What happened to this pod
   - **Exec tab**: Open a terminal inside the pod

### View Pod Logs

**Method 1: In Dashboard**
1. Click on a pod
2. Click the **"Logs"** icon (📄) at the top right
3. Real-time logs appear
4. Use the search box to filter

**Method 2: Download logs**
1. Click the download icon
2. Save logs as a file

### Execute Commands in a Pod

1. Click on a pod
2. Click the **"Exec"** icon (>_) at the top right
3. A terminal opens
4. Run commands:
   ```bash
   # Example commands
   ls -la
   ps aux
   env
   cat /app/config.txt
   ```

---

## 🌐 Viewing Services

### See All Services

1. Click **"Services"** in the left sidebar
2. You'll see:

```
NAME                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)
frontend              ClusterIP   10.96.xxx.xxx   <none>        80/TCP
backend               ClusterIP   10.96.xxx.xxx   <none>        4001/TCP
database              ClusterIP   10.96.xxx.xxx   <none>        5432/TCP
typescript-frontend   ClusterIP   10.96.xxx.xxx   <none>        80/TCP
support-frontend      ClusterIP   10.96.xxx.xxx   <none>        80/TCP
```

### Service Details

Click on a service to see:
- **Endpoints**: Which pods it connects to
- **Selectors**: How it finds pods
- **Ports**: What ports are exposed

---

## 📈 Enable Resource Metrics

### Install Metrics Server

```bash
minikube addons enable metrics-server
```

Wait 1-2 minutes, then refresh the dashboard.

### View Metrics

You'll now see:
- **CPU usage** graphs
- **Memory usage** graphs
- **Network I/O** statistics
- **Resource limits** vs actual usage

This appears on:
- Dashboard overview page
- Individual pod pages
- Deployment pages

---

## 🎯 Common Tasks

### Scale a Deployment

**In Dashboard:**
1. Go to **Deployments**
2. Click the **⋮** (three dots) next to a deployment
3. Click **"Scale"**
4. Change replica count
5. Click **"Scale"**

**With kubectl:**
```bash
# Scale frontend to 5 replicas
kubectl scale deployment frontend --replicas=5 -n usualstore

# Check status
kubectl get deployment frontend -n usualstore
```

### Restart a Pod

**In Dashboard:**
1. Go to **Pods**
2. Click the **⋮** (three dots) next to a pod
3. Click **"Delete"**
4. Kubernetes automatically creates a new one

**With kubectl:**
```bash
# Restart all pods in a deployment
kubectl rollout restart deployment/backend -n usualstore
```

### View Events

1. Go to **Events** in the left sidebar
2. See all cluster events in real-time
3. Filter by resource type
4. Useful for debugging issues

---

## 🖥️ k9s Terminal Dashboard

### What is k9s?

A powerful, terminal-based Kubernetes dashboard. Very popular with developers!

### Installation

```bash
# Install k9s
brew install k9s

# Run k9s
k9s
```

### Quick Navigation

```bash
# When k9s is open:

:pods              # View all pods
:svc               # View services  
:deploy            # View deployments
:ns                # View namespaces
:logs              # View logs

# Press '?' for help
# Press '/' to filter
# Press 'n' to switch namespace
# Press 'l' on a pod to see logs
# Press 'Enter' to see details
# Press 'd' to describe resource
# Press 'Ctrl+D' to delete (careful!)
```

### k9s Features

- ✅ **Real-time updates** - Everything auto-refreshes
- ✅ **Keyboard shortcuts** - Super fast navigation
- ✅ **Log streaming** - Live logs with color coding
- ✅ **Resource usage** - Built-in CPU/memory monitoring
- ✅ **Context switching** - Easy cluster switching
- ✅ **Search/filter** - Find anything quickly

### Example k9s Workflow

```bash
# Start k9s
k9s

# Switch to usualstore namespace
# Press 'n' then type 'usualstore'

# View pods
:pods

# Find a specific pod
# Press '/' then type 'backend'

# View logs of highlighted pod
# Press 'l'

# Go back
# Press 'Esc'

# Describe pod
# Press 'd'

# Quit
# Press ':quit' or 'Ctrl+C'
```

---

## 🎨 Lens IDE (Professional Tool)

### What is Lens?

The most powerful Kubernetes IDE. Beautiful UI, built-in terminal, metrics, and more.

### Installation

**Method 1: Download**
- Visit: https://k8slens.dev/
- Download for macOS
- Install like any app

**Method 2: Homebrew**
```bash
brew install --cask lens
```

### Setup

1. Open Lens
2. It **auto-detects** your Minikube cluster
3. Click on the cluster
4. You're in!

### Lens Features

- ✅ **Multi-cluster** - Manage multiple Kubernetes clusters
- ✅ **Built-in terminal** - Terminal inside the app
- ✅ **Resource charts** - Beautiful CPU/memory graphs
- ✅ **Log viewer** - Powerful log searching
- ✅ **Helm charts** - Install apps with one click
- ✅ **Extensions** - Add custom functionality
- ✅ **Metrics** - Prometheus integration
- ✅ **Port forwarding** - Easy port forwarding UI

### Lens Workspace

```
┌─────────────────────────────────────────────────────┐
│ Lens IDE                                     [⚙️]   │
├──────────────┬──────────────────────────────────────┤
│ Clusters     │   usualstore namespace              │
│ ├─ minikube  │                                     │
│ └─ docker... │   📦 Workloads                      │
│              │   ├─ Pods (13)          [All OK]    │
│ Namespaces   │   ├─ Deployments (8)   [Healthy]   │
│ ├─ default   │   └─ StatefulSets (1)  [Ready]     │
│ ├─ usualst...│                                     │
│              │   🌐 Network                        │
│ Workloads    │   ├─ Services (8)                   │
│ ├─ Pods      │   └─ Ingresses (0)                 │
│ ├─ Deploy... │                                     │
│              │   ⚙️  Config                        │
│ Network      │   ├─ ConfigMaps (1)                │
│ ├─ Services  │   └─ Secrets (2)                   │
│              │                                     │
│ [+] Terminal │   📊 CPU: ████░░ 45%               │
│ [+] Metrics  │   💾 Memory: ██░░░░ 23%            │
└──────────────┴──────────────────────────────────────┘
```

---

## 💻 Command Line Monitoring

### Watch Pods (Real-time)

```bash
# Watch all pods in usualstore namespace
kubectl get pods -n usualstore --watch

# More details (shows node, IP, etc.)
kubectl get pods -n usualstore -o wide --watch

# Watch with auto-refresh every 2 seconds
watch -n 2 kubectl get pods -n usualstore
```

### One-time Status Check

```bash
# All pods
kubectl get pods -n usualstore

# All resources
kubectl get all -n usualstore

# Specific deployment
kubectl get deployment backend -n usualstore

# Services
kubectl get svc -n usualstore
```

### Pod Logs

```bash
# View logs of a pod
kubectl logs <pod-name> -n usualstore

# Follow logs (live stream)
kubectl logs -f <pod-name> -n usualstore

# Logs from specific container in pod
kubectl logs <pod-name> -c <container-name> -n usualstore

# Last 50 lines
kubectl logs <pod-name> -n usualstore --tail=50

# Logs from all pods in deployment
kubectl logs deployment/backend -n usualstore

# Previous pod logs (if pod crashed)
kubectl logs <pod-name> -n usualstore --previous
```

### Describe Resources

```bash
# Detailed pod information
kubectl describe pod <pod-name> -n usualstore

# Deployment details
kubectl describe deployment backend -n usualstore

# Service details
kubectl describe service frontend -n usualstore

# Events at the bottom show what happened
```

### Execute Commands in Pod

```bash
# Open shell in pod
kubectl exec -it <pod-name> -n usualstore -- /bin/sh

# Or bash if available
kubectl exec -it <pod-name> -n usualstore -- /bin/bash

# Run single command
kubectl exec <pod-name> -n usualstore -- ls -la /app

# Check environment variables
kubectl exec <pod-name> -n usualstore -- env
```

---

## 🔍 Dashboard Features Comparison

| Feature | Minikube Dashboard | k9s | Lens | kubectl |
|---------|-------------------|-----|------|---------|
| **Ease of Use** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Speed** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Visual** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐ |
| **Power** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Learning Curve** | Easy | Medium | Easy | Hard |
| **Installation** | Built-in | `brew install` | Download app | Built-in |
| **Best For** | Beginners | Pros | Everyone | Automation |

---

## 🛠️ Troubleshooting

### Dashboard Won't Open

**Problem:** `minikube dashboard` doesn't open browser

**Solutions:**
```bash
# 1. Get the URL manually
minikube dashboard --url

# Copy the URL and paste in browser

# 2. Check if dashboard is enabled
minikube addons list | grep dashboard

# 3. Enable dashboard if disabled
minikube addons enable dashboard

# 4. Restart Minikube
minikube stop
minikube start
```

### Dashboard Shows Empty

**Problem:** No pods/services visible

**Solution:**
1. Check you're in the correct namespace
2. Click the namespace dropdown at top
3. Select **"usualstore"** not "default"

### Metrics Not Showing

**Problem:** No CPU/Memory graphs

**Solution:**
```bash
# Enable metrics-server
minikube addons enable metrics-server

# Wait 1-2 minutes
# Refresh dashboard
```

### Can't Access Pod Logs

**Problem:** Logs tab is empty or shows error

**Solution:**
```bash
# Check if pod is running
kubectl get pods -n usualstore

# View logs in terminal
kubectl logs <pod-name> -n usualstore

# If pod crashed recently, view previous logs
kubectl logs <pod-name> -n usualstore --previous
```

### Dashboard Crashes or Freezes

**Solution:**
```bash
# Restart dashboard
# Press Ctrl+C to stop
# Then run again:
minikube dashboard
```

---

## 📚 Additional Resources

### Minikube Commands

```bash
# Dashboard
minikube dashboard              # Open dashboard
minikube dashboard --url        # Get URL only
minikube addons enable dashboard # Enable addon

# Metrics
minikube addons enable metrics-server

# Cluster
minikube status                 # Check status
minikube stop                   # Stop cluster
minikube start                  # Start cluster
minikube delete                 # Delete cluster
```

### Useful kubectl Commands

```bash
# Get everything
kubectl get all -n usualstore

# Watch resources
kubectl get pods -n usualstore --watch

# Port forwarding
kubectl port-forward svc/frontend 3000:80 -n usualstore

# Logs
kubectl logs -f deployment/backend -n usualstore

# Describe
kubectl describe pod <pod-name> -n usualstore

# Events
kubectl get events -n usualstore --sort-by='.lastTimestamp'

# Execute command
kubectl exec -it <pod-name> -n usualstore -- /bin/sh
```

---

## 🎯 Quick Reference

### Access Your Application

**After pods are running:**

```bash
# Option 1: Minikube service (easiest)
minikube service frontend -n usualstore

# Option 2: Port forwarding
kubectl port-forward svc/frontend 3000:80 -n usualstore
# Then visit: http://localhost:3000

# Option 3: Get service list
minikube service list -n usualstore
```

### Monitor Deployment Progress

```bash
# Watch all pods
kubectl get pods -n usualstore --watch

# Or use k9s
k9s
# Then type: :pods
```

### Debug Issues

```bash
# 1. Check pod status
kubectl get pods -n usualstore

# 2. View logs
kubectl logs <pod-name> -n usualstore

# 3. Describe pod (see events)
kubectl describe pod <pod-name> -n usualstore

# 4. Check events
kubectl get events -n usualstore --sort-by='.lastTimestamp'
```

---

## 🎓 Pro Tips

1. **Use k9s for daily work** - Much faster than web UI once you learn it
2. **Keep dashboard open in background** - Great for visual overview
3. **Enable metrics-server** - Essential for production monitoring
4. **Use Lens for complex debugging** - Best tool for deep investigation
5. **Alias common commands** - Add to your `~/.zshrc`:
   ```bash
   alias kgp='kubectl get pods -n usualstore'
   alias kgpw='kubectl get pods -n usualstore --watch'
   alias kl='kubectl logs -f'
   alias ke='kubectl exec -it'
   ```

---

## 📖 Next Steps

1. **Try each dashboard** - See which you prefer
2. **Explore metrics** - Enable metrics-server and view graphs
3. **Practice debugging** - Use logs and describe commands
4. **Learn k9s shortcuts** - Become a Kubernetes ninja
5. **Install Lens** - Best tool for learning Kubernetes

---

## 🔗 Official Documentation

- **Kubernetes Dashboard**: https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/
- **Minikube Dashboard**: https://minikube.sigs.k8s.io/docs/handbook/dashboard/
- **k9s**: https://k9scli.io/
- **Lens**: https://docs.k8slens.dev/

---

**Happy Kubernetes exploring!** 🚀

