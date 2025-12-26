# 🚢 Kubernetes Setup for Usual Store

Your application is **Kubernetes-ready**! This file points you to all Kubernetes documentation and configuration.

---

## 📖 Documentation

All Kubernetes documentation is in **`docs/kubernetes/`**:

### **🌟 Start Here**
- **[docs/kubernetes/GETTING-STARTED.md](docs/kubernetes/GETTING-STARTED.md)** - Complete hands-on guide from zero to deployed

### **📚 Full Documentation**
- **[docs/kubernetes/KUBERNETES-OVERVIEW.md](docs/kubernetes/KUBERNETES-OVERVIEW.md)** - What is Kubernetes? Key concepts, architecture
- **[docs/kubernetes/DOCKER-VS-KUBERNETES.md](docs/kubernetes/DOCKER-VS-KUBERNETES.md)** - Comparison, when to use what
- **[docs/kubernetes/KUBERNETES-DEPLOYMENT.md](docs/kubernetes/KUBERNETES-DEPLOYMENT.md)** - Production deployment guide
- **[docs/kubernetes/README.md](docs/kubernetes/README.md)** - Documentation index

---

## 📁 Configuration Files

All Kubernetes YAML files are in **`k8s/`** folder:

```
k8s/
├── README.md                    # Quick reference
├── 01-namespace.yaml            # Isolated environment
├── 02-configmap.yaml            # Configuration
├── 03-secrets.yaml              # Passwords
├── 04-database-pvc.yaml         # Storage
├── 05-database-deployment.yaml  # PostgreSQL
├── 06-database-service.yaml     # DB networking
├── 07-backend-deployment.yaml   # API server (2 replicas)
├── 08-backend-service.yaml      # API networking
├── 09-frontend-deployment.yaml  # Web app (3 replicas)
├── 10-frontend-service.yaml     # Web app LoadBalancer
├── 11-invoice-deployment.yaml   # Invoice service
└── 12-invoice-service.yaml      # Invoice networking
```

---

## 🚀 Quick Start

### **1. Enable Kubernetes Locally**

**Docker Desktop:**
```
Settings → Kubernetes → Enable Kubernetes
```

**Or install Minikube:**
```bash
brew install minikube
minikube start
```

### **2. Deploy Your App**

```bash
# Deploy everything
kubectl apply -f k8s/

# Watch pods start
kubectl get pods -n usualstore -w
```

### **3. Access Your App**

**Docker Desktop:**
```bash
open http://localhost:80
```

**Minikube:**
```bash
minikube service frontend -n usualstore
```

---

## 📊 What You Get

```
Kubernetes Cluster:
┌─────────────────────────────────────────┐
│  • 1x database (PostgreSQL)             │
│  • 2x backend (API, load balanced)      │
│  • 3x frontend (Web, load balanced)     │
│  • 1x invoice (Microservice)            │
│  • LoadBalancer (automatic)             │
│  • Auto-scaling (ready to enable)       │
│  • Self-healing (automatic restarts)    │
└─────────────────────────────────────────┘
```

**Benefits:**
- ✅ High availability (multiple replicas)
- ✅ Auto-scaling (scale up/down automatically)
- ✅ Self-healing (auto-restart failed containers)
- ✅ Zero-downtime updates (rolling deployments)
- ✅ Production-ready configuration

---

## 🎯 Development Workflow

### **Local Development (Docker Compose)**
```bash
# Use Docker Compose for development
make docker-up
# Access: http://localhost:4000

make docker-down
```

### **Test on Local Kubernetes**
```bash
# Deploy to local Kubernetes to test
kubectl apply -f k8s/
# Access: http://localhost:80

# Clean up
kubectl delete -f k8s/
```

### **Production (Cloud Kubernetes)**
```bash
# Push images to registry
docker push YOUR_USERNAME/usualstore-frontend:latest
docker push YOUR_USERNAME/usualstore-backend:latest
docker push YOUR_USERNAME/usualstore-invoice:latest

# Deploy to cloud
kubectl apply -f k8s/
# Access: http://LOAD_BALANCER_IP
```

---

## 💡 Useful Commands

### **Check Status**
```bash
kubectl get pods -n usualstore
kubectl get services -n usualstore
kubectl get all -n usualstore
```

### **View Logs**
```bash
kubectl logs -f deployment/frontend -n usualstore
kubectl logs -f deployment/backend -n usualstore
```

### **Scale Services**
```bash
kubectl scale deployment frontend --replicas=5 -n usualstore
kubectl scale deployment backend --replicas=3 -n usualstore
```

### **Troubleshoot**
```bash
kubectl describe pod POD_NAME -n usualstore
kubectl exec -it POD_NAME -n usualstore -- /bin/sh
```

### **Delete Everything**
```bash
kubectl delete -f k8s/
# or
kubectl delete namespace usualstore
```

---

## 🎓 Learn More

1. **Read the getting started guide:**
   - [docs/kubernetes/GETTING-STARTED.md](docs/kubernetes/GETTING-STARTED.md)

2. **Understand the concepts:**
   - [docs/kubernetes/KUBERNETES-OVERVIEW.md](docs/kubernetes/KUBERNETES-OVERVIEW.md)

3. **Compare with Docker Compose:**
   - [docs/kubernetes/DOCKER-VS-KUBERNETES.md](docs/kubernetes/DOCKER-VS-KUBERNETES.md)

4. **Deploy to production:**
   - [docs/kubernetes/KUBERNETES-DEPLOYMENT.md](docs/kubernetes/KUBERNETES-DEPLOYMENT.md)

---

## 💰 Deployment Options

### **Free (Local Testing)**
- Docker Desktop Kubernetes
- Minikube

### **Cloud (Production)**
- DigitalOcean Kubernetes: ~$30-50/month (easiest)
- Google Kubernetes Engine (GKE): ~$70/month
- Amazon EKS: ~$75/month
- Azure AKS: ~$70/month

---

## ✅ You're Ready!

You have:
- ✅ Complete Kubernetes configuration (k8s/ folder)
- ✅ Comprehensive documentation (docs/kubernetes/)
- ✅ Step-by-step guides (from beginner to production)
- ✅ Working local setup (Docker Compose)
- ✅ Production-ready deployment files

**Next step:** Open [docs/kubernetes/GETTING-STARTED.md](docs/kubernetes/GETTING-STARTED.md) and follow along!

---

**Happy Kubernetes-ing!** 🚀

