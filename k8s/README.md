# Kubernetes Configuration Files

This directory contains all Kubernetes manifests for deploying Usual Store.

## 📁 Files

| File | Description |
|------|-------------|
| `01-namespace.yaml` | Creates isolated `usualstore` namespace |
| `02-configmap.yaml` | Non-sensitive configuration |
| `03-secrets.yaml` | Passwords and connection strings |
| `04-database-pvc.yaml` | Persistent storage for PostgreSQL |
| `05-database-deployment.yaml` | PostgreSQL StatefulSet |
| `06-database-service.yaml` | Database network service |
| `07-backend-deployment.yaml` | API server deployment (2 replicas) |
| `08-backend-service.yaml` | Backend network service |
| `09-frontend-deployment.yaml` | Web app deployment (3 replicas) |
| `10-frontend-service.yaml` | Frontend LoadBalancer service |
| `11-invoice-deployment.yaml` | Invoice microservice deployment |
| `12-invoice-service.yaml` | Invoice network service |

## 🚀 Quick Start

### Deploy all at once
```bash
kubectl apply -f k8s/
```

### Deploy step-by-step
```bash
# 1. Namespace and config
kubectl apply -f 01-namespace.yaml
kubectl apply -f 02-configmap.yaml
kubectl apply -f 03-secrets.yaml

# 2. Database
kubectl apply -f 04-database-pvc.yaml
kubectl apply -f 05-database-deployment.yaml
kubectl apply -f 06-database-service.yaml

# Wait for database
kubectl wait --for=condition=ready pod -l app=database -n usualstore --timeout=120s

# 3. Backend
kubectl apply -f 07-backend-deployment.yaml
kubectl apply -f 08-backend-service.yaml

# 4. Frontend
kubectl apply -f 09-frontend-deployment.yaml
kubectl apply -f 10-frontend-service.yaml

# 5. Invoice
kubectl apply -f 11-invoice-deployment.yaml
kubectl apply -f 12-invoice-service.yaml
```

## ✅ Verify

```bash
# Check all pods
kubectl get pods -n usualstore

# Check services
kubectl get services -n usualstore

# Check logs
kubectl logs -f deployment/frontend -n usualstore
```

## 🌐 Access

### Minikube
```bash
minikube service frontend -n usualstore
```

### Cloud providers
```bash
kubectl get service frontend -n usualstore
# Wait for EXTERNAL-IP, then visit http://EXTERNAL-IP
```

## 🔧 Customize

Before deploying:

1. **Replace image names** in deployment files:
   ```yaml
   # Change from:
   image: usual_store-front-end:latest
   
   # To your registry:
   image: YOUR_USERNAME/usualstore-frontend:latest
   ```

2. **Update secrets** (production):
   ```bash
   kubectl create secret generic usualstore-secrets \
     --from-literal=DB_PASSWORD=SECURE_PASSWORD \
     -n usualstore
   ```

3. **Adjust resources** based on your needs:
   ```yaml
   resources:
     requests:
       memory: "128Mi"
       cpu: "100m"
     limits:
       memory: "512Mi"
       cpu: "500m"
   ```

## 📖 More Info

See `../docs/kubernetes/` for detailed guides:
- `KUBERNETES-OVERVIEW.md` - Kubernetes concepts
- `KUBERNETES-DEPLOYMENT.md` - Deployment guide

