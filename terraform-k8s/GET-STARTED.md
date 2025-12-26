# 🚀 Get Started with Kubernetes + Terraform

## 3-Minute Quick Start

### Prerequisites Check

1. **Do you have Docker Desktop?**
   ```bash
   docker --version
   ```
   If not: [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)

2. **Enable Kubernetes in Docker Desktop**:
   - Open Docker Desktop
   - Settings → Kubernetes
   - ✅ Enable Kubernetes
   - Click "Apply & Restart"
   - Wait ~2 minutes for K8s to start

3. **Install Terraform**:
   ```bash
   brew install terraform
   ```

### Deploy Now!

```bash
# Navigate to the terraform-k8s directory
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform-k8s

# Run the automatic deployment script
./QUICK-START.sh
```

That's it! The script will:
1. ✅ Check all prerequisites
2. ✅ Build Docker images (if needed)
3. ✅ Initialize Terraform
4. ✅ Deploy to Kubernetes
5. ✅ Wait for everything to be ready
6. ✅ Show you how to access your app

### Access Your App

After deployment, run:

```bash
# Frontend (in one terminal)
kubectl port-forward svc/frontend 3000:80 -n usualstore

# Backend (in another terminal)
kubectl port-forward svc/backend 4001:4001 -n usualstore
```

Then visit:
- Frontend: http://localhost:3000
- Backend API: http://localhost:4001

---

## Manual Deployment (Step by Step)

If you prefer to do it manually:

### Step 1: Initialize Terraform

```bash
cd terraform-k8s/local
terraform init
```

### Step 2: Preview Changes

```bash
terraform plan
```

This shows what will be created (no changes made yet).

### Step 3: Deploy

```bash
terraform apply
```

Type `yes` when prompted.

### Step 4: Verify

```bash
kubectl get pods -n usualstore
```

You should see:
```
NAME                        READY   STATUS    RESTARTS   AGE
backend-xxxxxxxxx-xxxxx     1/1     Running   0          2m
backend-xxxxxxxxx-xxxxx     1/1     Running   0          2m
database-0                  1/1     Running   0          2m
frontend-xxxxxxxxx-xxxxx    1/1     Running   0          2m
frontend-xxxxxxxxx-xxxxx    1/1     Running   0          2m
frontend-xxxxxxxxx-xxxxx    1/1     Running   0          2m
```

---

## What Just Happened?

Terraform just:

1. ✅ Created a `usualstore` namespace in Kubernetes
2. ✅ Applied all 18 YAML files from `/k8s/`
3. ✅ Started these services:
   - **Frontend**: 3 replicas (for high availability)
   - **Backend**: 2 replicas (load balanced)
   - **Database**: 1 replica (with persistent storage)
   - **AI Assistant**: 1 replica
   - **Invoice Service**: 1 replica
   - **Support Services**: Ready to use

4. ✅ Configured:
   - ConfigMaps (environment settings)
   - Secrets (passwords, API keys)
   - Persistent Volumes (database storage)
   - Services (networking between pods)
   - Health checks (automatic restart if unhealthy)

---

## Common Tasks

### View Logs

```bash
# Backend logs
kubectl logs -f deployment/backend -n usualstore

# Frontend logs
kubectl logs -f deployment/frontend -n usualstore

# Database logs
kubectl logs -f statefulset/database -n usualstore
```

### Check Status

```bash
# All pods
kubectl get pods -n usualstore

# All services
kubectl get svc -n usualstore

# Everything
kubectl get all -n usualstore
```

### Scale Services

```bash
# Scale frontend to 5 replicas
kubectl scale deployment frontend --replicas=5 -n usualstore

# Scale backend to 4 replicas
kubectl scale deployment backend --replicas=4 -n usualstore

# Check the scaling
kubectl get pods -n usualstore -w
```

### Update Application

```bash
# After making code changes and rebuilding images:
kubectl rollout restart deployment/backend -n usualstore
kubectl rollout restart deployment/frontend -n usualstore

# Watch the rollout
kubectl rollout status deployment/backend -n usualstore
```

### Debug a Pod

```bash
# Get shell access to a pod
kubectl exec -it <pod-name> -n usualstore -- /bin/sh

# Example: Access backend pod
kubectl exec -it $(kubectl get pod -n usualstore -l app=backend -o jsonpath='{.items[0].metadata.name}') -n usualstore -- /bin/sh
```

### View Events

```bash
# See what's happening
kubectl get events -n usualstore --sort-by='.lastTimestamp'
```

---

## Cleanup

### Option 1: Use Terraform (Recommended)

```bash
cd terraform-k8s/local
terraform destroy -auto-approve
```

### Option 2: Use kubectl

```bash
kubectl delete namespace usualstore
```

Both options remove everything completely.

---

## Troubleshooting

### Issue: "No Kubernetes context"

**Problem**: Docker Desktop Kubernetes is not enabled.

**Solution**:
1. Open Docker Desktop
2. Settings → Kubernetes
3. ✅ Enable Kubernetes
4. Apply & Restart
5. Wait 2 minutes

Verify:
```bash
kubectl cluster-info
```

### Issue: Pods stuck in "ImagePullBackOff"

**Problem**: Docker images don't exist locally.

**Solution**:
```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store
docker-compose build
```

Then redeploy:
```bash
cd terraform-k8s/local
terraform apply
```

### Issue: Pod stuck in "Pending"

**Problem**: Not enough resources.

**Solution**: Check Docker Desktop resources:
- Docker Desktop → Settings → Resources
- Increase CPUs to 4
- Increase Memory to 8 GB
- Apply & Restart

### Issue: "connection refused" when port forwarding

**Problem**: Service isn't ready yet.

**Solution**: Wait for pods to be ready:
```bash
kubectl get pods -n usualstore
# Wait until all show "Running" and "1/1" or "2/2" ready

# Then try port forwarding again
kubectl port-forward svc/frontend 3000:80 -n usualstore
```

### Issue: Want to use Minikube instead

**Solution**:
```bash
# Start minikube
minikube start

# Deploy with minikube context
cd terraform-k8s/local
terraform apply -var="kube_context=minikube"

# Access frontend
minikube service frontend -n usualstore
```

---

## Next Steps

1. ✅ **Deploy locally** (you just did this!)
2. ⏳ **Explore the app** running on Kubernetes
3. ⏳ **Try scaling** services up and down
4. ⏳ **Read the comparison** in `DOCKER-VS-K8S.md`
5. ⏳ **When ready for AWS**, check `aws-future/` directory

---

## Learn More

- 📖 **Full Documentation**: [README.md](README.md)
- 🐳 **Docker vs K8s**: [DOCKER-VS-K8S.md](DOCKER-VS-K8S.md)
- ☸️ **Kubernetes Basics**: [/k8s/README.md](../k8s/README.md)

---

## Questions?

**"Can I use this for production?"**
- Yes! But deploy to AWS EKS using `aws-future/` directory

**"Does this cost money?"**
- Local: Free (runs on your computer)
- AWS: ~$180/month (only when you deploy to cloud)

**"Do I need to rewrite my code?"**
- No! We use your existing code and Docker images

**"Can I still use Docker Compose?"**
- Yes! Use Docker Compose for development, Kubernetes for testing/production

---

**🎉 Congratulations!** You now have a production-ready Kubernetes deployment!

