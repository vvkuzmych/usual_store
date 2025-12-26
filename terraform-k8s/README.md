# 🚀 Terraform Kubernetes Deployment

Deploy Usual Store to Kubernetes using Terraform - locally or in the cloud.

## 📁 Directory Structure

```
terraform-k8s/
├── local/              ← Deploy to LOCAL Kubernetes (Docker Desktop, Minikube, Kind)
├── aws-future/         ← Deploy to AWS EKS (prepared for future)
└── modules/
    └── k8s-app/        ← Reusable deployment module
```

---

## 🏠 LOCAL DEPLOYMENT (Start Here!)

### Prerequisites

1. **Kubernetes cluster running locally**:
   - **Docker Desktop** (recommended): Enable Kubernetes in Settings
   - **Minikube**: `minikube start`
   - **Kind**: `kind create cluster --name usualstore`

2. **Install tools**:
   ```bash
   # Install Terraform
   brew install terraform
   
   # Verify kubectl is installed
   kubectl version --client
   
   # Check your cluster is running
   kubectl cluster-info
   ```

### Quick Start

```bash
# 1. Navigate to local deployment
cd terraform-k8s/local

# 2. Initialize Terraform
terraform init

# 3. Preview what will be created
terraform plan

# 4. Deploy!
terraform apply -auto-approve

# 5. Access your app
kubectl port-forward svc/frontend 3000:80 -n usualstore
# Visit http://localhost:3000
```

### Configuration

If you're using **Minikube** or **Kind**, update the context:

```bash
# For Minikube:
terraform apply -var="kube_context=minikube"

# For Kind:
terraform apply -var="kube_context=kind-usualstore"

# Or edit local/variables.tf and change the default
```

### Verify Deployment

```bash
# Check all pods are running
kubectl get pods -n usualstore

# Expected output:
# NAME                        READY   STATUS    RESTARTS   AGE
# backend-xxxxxxxxx-xxxxx     1/1     Running   0          2m
# backend-xxxxxxxxx-xxxxx     1/1     Running   0          2m
# database-0                  1/1     Running   0          2m
# frontend-xxxxxxxxx-xxxxx    1/1     Running   0          2m
# frontend-xxxxxxxxx-xxxxx    1/1     Running   0          2m
# frontend-xxxxxxxxx-xxxxx    1/1     Running   0          2m

# Check services
kubectl get svc -n usualstore

# View logs
kubectl logs -f deployment/backend -n usualstore
```

### Access Your Application

**Method 1: Port forwarding** (recommended for local)
```bash
# Frontend
kubectl port-forward svc/frontend 3000:80 -n usualstore

# Backend API
kubectl port-forward svc/backend 4001:4001 -n usualstore

# Database (for debugging)
kubectl port-forward svc/database 5432:5432 -n usualstore
```

**Method 2: Minikube service** (Minikube only)
```bash
minikube service frontend -n usualstore
```

### Cleanup

```bash
# Destroy all resources
terraform destroy -auto-approve

# Or use kubectl
kubectl delete namespace usualstore
```

---

## ☁️ AWS DEPLOYMENT (Future)

> 🚨 **Not active yet!** This is prepared for when you get an AWS account.

### What's Ready

The `aws-future/` directory contains a **complete AWS EKS configuration** that:

- ✅ Creates a production-ready EKS cluster
- ✅ Sets up VPC with private/public subnets
- ✅ Configures auto-scaling node groups
- ✅ Deploys your app with LoadBalancers
- ✅ Uses AWS Application Load Balancer for ingress

### When You're Ready for AWS

1. **Create AWS account** and configure credentials:
   ```bash
   aws configure
   ```

2. **Edit `aws-future/main.tf`**:
   - Uncomment all resources (remove `#` from lines)

3. **Deploy**:
   ```bash
   cd terraform-k8s/aws-future
   terraform init
   terraform apply
   ```

4. **Estimated AWS costs**:
   - EKS cluster: ~$0.10/hour (~$73/month)
   - 3 × t3.medium nodes: ~$0.12/hour (~$90/month)
   - LoadBalancer: ~$0.025/hour (~$18/month)
   - **Total: ~$180/month**

5. **Cost-saving tips**:
   - Use `t3.small` instances instead of `t3.medium`
   - Reduce `desired_nodes` to 2
   - Use spot instances for non-production
   - Stop cluster when not in use

---

## 🔧 How It Works

### Architecture

```
terraform-k8s/
│
├── local/main.tf          ← Configures local K8s provider
│   └── Uses module "k8s-app"
│
├── aws-future/main.tf     ← Creates EKS + configures AWS K8s provider
│   └── Uses module "k8s-app"
│
└── modules/k8s-app/       ← Reusable deployment logic
    ├── Applies all YAML files from /k8s/
    ├── Waits for deployments to be ready
    └── Works with ANY Kubernetes cluster
```

### Key Benefits

1. **Same deployment code** for local and cloud
2. **Your existing YAML files** are reused (no rewriting!)
3. **Infrastructure as Code** - version controlled and reproducible
4. **Easy migration** from local to AWS

### What Terraform Manages

```
Local:
  ✓ Applies all k8s/*.yaml manifests
  ✓ Waits for deployments to be ready
  ✓ Provides access instructions

AWS (when ready):
  ✓ Creates VPC and networking
  ✓ Provisions EKS cluster
  ✓ Manages node groups (auto-scaling)
  ✓ Applies all k8s/*.yaml manifests
  ✓ Configures LoadBalancers
  ✓ Sets up ALB Ingress Controller
```

---

## 📖 Common Commands

### Terraform

```bash
# Initialize (first time only)
terraform init

# See what will change
terraform plan

# Apply changes
terraform apply

# Destroy everything
terraform destroy

# Show current state
terraform show

# List all resources
terraform state list
```

### Kubernetes

```bash
# Get all resources
kubectl get all -n usualstore

# Describe a pod
kubectl describe pod <pod-name> -n usualstore

# View logs
kubectl logs -f <pod-name> -n usualstore

# Execute command in pod
kubectl exec -it <pod-name> -n usualstore -- /bin/sh

# Get events
kubectl get events -n usualstore --sort-by='.lastTimestamp'

# Scale deployment
kubectl scale deployment backend --replicas=3 -n usualstore
```

---

## 🐛 Troubleshooting

### Issue: "cluster-info: connection refused"

**Solution**: Your Kubernetes cluster is not running.

```bash
# Docker Desktop: Enable Kubernetes in Docker Desktop settings
# Minikube: 
minikube start

# Kind:
kind create cluster --name usualstore
```

### Issue: "context 'docker-desktop' not found"

**Solution**: Update the context in `local/variables.tf`:

```terraform
variable "kube_context" {
  default = "minikube"  # or "kind-usualstore"
}
```

Or use command line:
```bash
terraform apply -var="kube_context=minikube"
```

### Issue: Pods stuck in "Pending" or "ImagePullBackOff"

**Solution**: Check your Docker images exist:

```bash
# List images
docker images | grep usual_store

# If missing, build them
cd /Users/vkuzm/Projects/UsualStore/usual_store
docker-compose build
```

### Issue: "No resources found in usualstore namespace"

**Solution**: Namespace might not be created yet:

```bash
kubectl create namespace usualstore
terraform apply
```

---

## 🎯 Next Steps

1. ✅ **Deploy locally** using `terraform-k8s/local/`
2. ⏳ **Test your application** thoroughly
3. ⏳ **When ready for AWS**: Uncomment and deploy `aws-future/`

---

## 📚 Additional Resources

- [Terraform Kubernetes Provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs)
- [AWS EKS Module](https://registry.terraform.io/modules/terraform-aws-modules/eks/aws/latest)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)

---

## 💡 Pro Tips

1. **Use workspaces** for multiple environments:
   ```bash
   terraform workspace new dev
   terraform workspace new staging
   terraform workspace select dev
   ```

2. **Store state in S3** (for team collaboration):
   ```terraform
   terraform {
     backend "s3" {
       bucket = "usualstore-terraform-state"
       key    = "k8s/terraform.tfstate"
       region = "us-east-1"
     }
   }
   ```

3. **Use terraform.tfvars** for custom settings:
   ```bash
   # local/terraform.tfvars
   kube_context = "minikube"
   ```

4. **Automate with Makefile**:
   ```makefile
   .PHONY: deploy
   deploy:
       cd terraform-k8s/local && terraform apply -auto-approve
   
   .PHONY: destroy
   destroy:
       cd terraform-k8s/local && terraform destroy -auto-approve
   ```

---

**Questions?** Check `/k8s/README.md` for Kubernetes-specific documentation.

