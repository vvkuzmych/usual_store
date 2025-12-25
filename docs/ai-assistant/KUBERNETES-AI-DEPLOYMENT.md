# ☸️ AI Assistant - Kubernetes Deployment Guide

Deploy the AI assistant to Kubernetes for production-grade scalability and high availability.

---

## 📦 What's Included

**4 New Kubernetes Manifests:**
- `13-ai-assistant-deployment.yaml` - AI service deployment (2 replicas)
- `14-ai-assistant-service.yaml` - Internal service for load balancing
- `15-ai-assistant-secrets.yaml` - OpenAI API key storage
- `16-ai-assistant-ingress.yaml` - External access configuration

---

## 🚀 Quick Start

### **1. Set Up Secrets**

```bash
cd /Users/vkuzm/Projects/usual_store

# Create OpenAI API key secret
kubectl create secret generic ai-assistant-secrets \
  --from-literal=OPENAI_API_KEY=sk-your-actual-key-here \
  -n usualstore

# Or apply the secrets file (after editing)
# Edit k8s/15-ai-assistant-secrets.yaml first!
kubectl apply -f k8s/15-ai-assistant-secrets.yaml
```

### **2. Build and Push Image**

```bash
# Build image
docker build -f Dockerfile.ai-assistant \
  -t YOUR_USERNAME/usualstore-ai-assistant:latest .

# Push to Docker Hub
docker push YOUR_USERNAME/usualstore-ai-assistant:latest

# Update image name in k8s/13-ai-assistant-deployment.yaml
```

### **3. Deploy to Kubernetes**

```bash
# Deploy AI assistant
kubectl apply -f k8s/13-ai-assistant-deployment.yaml
kubectl apply -f k8s/14-ai-assistant-service.yaml
kubectl apply -f k8s/16-ai-assistant-ingress.yaml

# Or deploy all at once
kubectl apply -f k8s/

# Watch deployment
kubectl get pods -n usualstore -w
```

### **4. Verify Deployment**

```bash
# Check pods
kubectl get pods -n usualstore -l app=ai-assistant

# Expected output:
# NAME                            READY   STATUS    RESTARTS
# ai-assistant-xxxxxxxxx-xxxxx    1/1     Running   0
# ai-assistant-xxxxxxxxx-xxxxx    1/1     Running   0

# Check service
kubectl get service ai-assistant -n usualstore

# Check logs
kubectl logs -f deployment/ai-assistant -n usualstore
```

---

## 🏗️ Architecture on Kubernetes

```
┌─────────────────────────────────────────────────────────────┐
│                  Kubernetes Cluster                         │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Ingress (External Access)                           │  │
│  │  → Routes /api/ai/* to ai-assistant service          │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │                                     │
│  ┌────────────────────▼──────────────────┐                 │
│  │  AI Assistant Service (ClusterIP)     │                 │
│  │  → Load balances to ai-assistant pods │                 │
│  └────────────────────┬──────────────────┘                 │
│                       │                                     │
│  ┌────────────────────▼──────────────────┐                 │
│  │  AI Assistant Deployment               │                 │
│  │  ├── Pod 1 (usualstore-ai:8080)       │                 │
│  │  └── Pod 2 (usualstore-ai:8080)       │                 │
│  └────────────────────┬──────────────────┘                 │
│                       │                                     │
│                       ├─ Connects to Database Service      │
│                       └─ Calls OpenAI API (external)       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 📋 Kubernetes Manifests

### **1. Deployment (13-ai-assistant-deployment.yaml)**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-assistant
  namespace: usualstore
spec:
  replicas: 2  # High availability
  selector:
    matchLabels:
      app: ai-assistant
  template:
    spec:
      containers:
      - name: ai-assistant
        image: YOUR_USERNAME/usualstore-ai-assistant:latest
        ports:
        - containerPort: 8080
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-assistant-secrets
              key: OPENAI_API_KEY
        - name: DATABASE_DSN
          valueFrom:
            secretKeyRef:
              name: usualstore-secrets
              key: DATABASE_DSN
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
```

### **2. Service (14-ai-assistant-service.yaml)**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: ai-assistant
  namespace: usualstore
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
  selector:
    app: ai-assistant
```

### **3. Ingress (16-ai-assistant-ingress.yaml)**

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ai-assistant-ingress
  namespace: usualstore
  annotations:
    nginx.ingress.kubernetes.io/enable-cors: "true"
spec:
  rules:
  - host: store.example.com
    http:
      paths:
      - path: /api/ai
        pathType: Prefix
        backend:
          service:
            name: ai-assistant
            port:
              number: 8080
```

---

## 🔧 Configuration

### **Environment Variables**

Managed via Kubernetes Secrets and ConfigMaps:

```bash
# Required secrets
OPENAI_API_KEY=sk-...    # From ai-assistant-secrets
DATABASE_DSN=postgres://... # From usualstore-secrets

# Optional config
PORT=8080                # Hardcoded in deployment
```

### **Update Secrets**

```bash
# Update OpenAI API key
kubectl create secret generic ai-assistant-secrets \
  --from-literal=OPENAI_API_KEY=sk-new-key-here \
  -n usualstore \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart pods to pick up new secret
kubectl rollout restart deployment/ai-assistant -n usualstore
```

---

## 📊 Scaling

### **Manual Scaling**

```bash
# Scale to 5 replicas
kubectl scale deployment ai-assistant --replicas=5 -n usualstore

# Scale back to 2
kubectl scale deployment ai-assistant --replicas=2 -n usualstore
```

### **Auto-Scaling (HPA)**

Create `k8s/17-ai-assistant-hpa.yaml`:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ai-assistant-hpa
  namespace: usualstore
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ai-assistant
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

Apply:
```bash
kubectl apply -f k8s/17-ai-assistant-hpa.yaml

# Watch auto-scaling
kubectl get hpa -n usualstore -w
```

---

## 🔍 Monitoring

### **Check Status**

```bash
# Pods
kubectl get pods -n usualstore -l app=ai-assistant

# Deployment
kubectl get deployment ai-assistant -n usualstore

# Service
kubectl get service ai-assistant -n usualstore

# Ingress
kubectl get ingress -n usualstore
```

### **View Logs**

```bash
# All pods
kubectl logs -f deployment/ai-assistant -n usualstore

# Specific pod
kubectl logs -f POD_NAME -n usualstore

# Last 100 lines
kubectl logs --tail=100 deployment/ai-assistant -n usualstore

# Since 1 hour ago
kubectl logs --since=1h deployment/ai-assistant -n usualstore
```

### **Describe Resources**

```bash
# Deployment details
kubectl describe deployment ai-assistant -n usualstore

# Pod details
kubectl describe pod POD_NAME -n usualstore

# Service details
kubectl describe service ai-assistant -n usualstore
```

### **Get Metrics**

```bash
# Resource usage
kubectl top pods -n usualstore -l app=ai-assistant

# Expected output:
# NAME                            CPU(cores)   MEMORY(bytes)
# ai-assistant-xxx-xxx            50m          256Mi
# ai-assistant-xxx-xxx            45m          240Mi
```

---

## 🌐 Access Methods

### **From Within Cluster**

```bash
# From another pod (e.g., frontend)
http://ai-assistant.usualstore.svc.cluster.local:8080/api/ai/chat

# Short form (same namespace)
http://ai-assistant:8080/api/ai/chat
```

### **Port Forward (Local Testing)**

```bash
# Forward to local machine
kubectl port-forward service/ai-assistant 8080:8080 -n usualstore

# Access at
curl http://localhost:8080/api/ai/chat
```

### **Via Ingress (Production)**

```bash
# After setting up ingress and DNS
curl https://store.example.com/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"session_id": "test", "message": "Hi"}'
```

---

## 🔄 Updates & Rollbacks

### **Update Deployment**

```bash
# Build new image
docker build -f Dockerfile.ai-assistant \
  -t YOUR_USERNAME/usualstore-ai-assistant:v1.1 .
docker push YOUR_USERNAME/usualstore-ai-assistant:v1.1

# Update deployment
kubectl set image deployment/ai-assistant \
  ai-assistant=YOUR_USERNAME/usualstore-ai-assistant:v1.1 \
  -n usualstore

# Watch rollout
kubectl rollout status deployment/ai-assistant -n usualstore
```

### **Rollback**

```bash
# View rollout history
kubectl rollout history deployment/ai-assistant -n usualstore

# Rollback to previous version
kubectl rollout undo deployment/ai-assistant -n usualstore

# Rollback to specific revision
kubectl rollout undo deployment/ai-assistant \
  --to-revision=2 \
  -n usualstore
```

### **Rolling Update Strategy**

Already configured in deployment:
```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1        # Add 1 new pod before removing old
    maxUnavailable: 0   # Keep all pods available during update
```

**Result:** Zero-downtime deployments! 🎉

---

## 🐛 Troubleshooting

### **Pods Not Starting**

```bash
# Check pod status
kubectl get pods -n usualstore -l app=ai-assistant

# Describe pod for events
kubectl describe pod POD_NAME -n usualstore

# Common issues:
# 1. Image pull error
# 2. Missing secrets
# 3. Resource limits
```

### **Health Check Failing**

```bash
# Check logs
kubectl logs POD_NAME -n usualstore

# Test health endpoint from pod
kubectl exec POD_NAME -n usualstore -- \
  wget -qO- http://localhost:8080/health

# Check liveness/readiness probes
kubectl describe pod POD_NAME -n usualstore | grep -A 10 "Liveness\|Readiness"
```

### **Cannot Connect to Database**

```bash
# Test database connectivity from pod
kubectl exec -it POD_NAME -n usualstore -- sh

# Inside pod:
# apt-get update && apt-get install -y postgresql-client
# psql "postgres://postgres:password@database:5432/usualstore"

# Check database service
kubectl get service database -n usualstore

# Check database pod
kubectl get pods -n usualstore -l app=database
```

### **High CPU/Memory Usage**

```bash
# Check current usage
kubectl top pods -n usualstore -l app=ai-assistant

# Adjust resource limits in deployment
resources:
  limits:
    memory: "1Gi"  # Increase from 512Mi
    cpu: "1000m"   # Increase from 500m
```

---

## 🔒 Production Best Practices

### **1. Security**

```bash
# Use sealed secrets (better than plain secrets)
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.18.0/controller.yaml

# Create sealed secret
kubeseal --format=yaml < k8s/15-ai-assistant-secrets.yaml > k8s/15-ai-assistant-sealed-secrets.yaml

# Use network policies
kubectl apply -f k8s/18-network-policy.yaml
```

### **2. Resource Quotas**

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: usualstore-quota
  namespace: usualstore
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
    persistentvolumeclaims: "10"
```

### **3. Pod Disruption Budget**

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: ai-assistant-pdb
  namespace: usualstore
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: ai-assistant
```

### **4. Monitoring & Alerts**

```bash
# Install Prometheus
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml

# Install Grafana
kubectl apply -f https://raw.githubusercontent.com/grafana/grafana/main/deploy/kubernetes/deployment.yaml

# Configure alerts for:
# - High error rate
# - API cost exceeding budget
# - Pod crashes
# - High latency
```

---

## 📊 Cost Optimization

### **1. Right-Size Resources**

```bash
# Monitor actual usage
kubectl top pods -n usualstore -l app=ai-assistant

# Adjust based on metrics
resources:
  requests:
    memory: "128Mi"  # Start small
    cpu: "100m"
  limits:
    memory: "512Mi"  # Prevent OOM
    cpu: "500m"
```

### **2. Use Node Affinity**

```yaml
affinity:
  nodeAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 1
      preference:
        matchExpressions:
        - key: node.kubernetes.io/instance-type
          operator: In
          values:
          - t3.small  # Use cheaper instances
```

### **3. Monitor OpenAI Costs**

```bash
# Get cost stats
kubectl exec -it POD_NAME -n usualstore -- \
  wget -qO- http://localhost:8080/api/ai/stats?days=7

# Set up cost alerts in your monitoring
```

---

## 🎯 Complete Deployment Checklist

### **Pre-Deployment**
- [ ] OpenAI API key obtained
- [ ] Image built and pushed to registry
- [ ] Secrets created in Kubernetes
- [ ] Database migrations applied
- [ ] Resource limits configured

### **Deployment**
- [ ] Apply all YAML manifests
- [ ] Verify pods are running
- [ ] Check health endpoints
- [ ] Test internal service access
- [ ] Configure ingress/LoadBalancer

### **Post-Deployment**
- [ ] Monitor logs for errors
- [ ] Test chat functionality
- [ ] Check API costs
- [ ] Set up auto-scaling (HPA)
- [ ] Configure alerts
- [ ] Document access URLs

---

## 📚 Next Steps

1. ✅ Kubernetes deployment complete
2. Set up monitoring (Prometheus/Grafana)
3. Configure CI/CD pipeline
4. Implement cost alerts
5. Scale based on traffic

**See also:**
- [DOCKER-DEPLOYMENT.md](DOCKER-DEPLOYMENT.md) - Docker Compose setup
- [QUICK-START.md](QUICK-START.md) - Local development
- [AI-ASSISTANT-OVERVIEW.md](AI-ASSISTANT-OVERVIEW.md) - Architecture

---

**Your AI assistant is now running on Kubernetes!** ☸️🎉

