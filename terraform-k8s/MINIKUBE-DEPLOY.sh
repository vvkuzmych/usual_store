#!/bin/bash
# Minikube-specific deployment script
# This handles building images in Minikube's Docker environment

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Usual Store - Minikube Deployment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Step 1: Check if Minikube is running
echo "📋 Step 1: Checking Minikube status..."
if ! minikube status &> /dev/null; then
    echo "⚠️  Minikube is not running. Starting..."
    minikube start
    echo "✅ Minikube started"
else
    echo "✅ Minikube is running"
fi
echo ""

# Step 2: Configure Docker to use Minikube's daemon
echo "📋 Step 2: Configuring Docker environment..."
eval $(minikube docker-env)
echo "✅ Docker configured to use Minikube"
echo ""

# Step 3: Build images
echo "📋 Step 3: Building Docker images..."
echo "⏰ This will take 5-10 minutes..."
cd ..
docker-compose build
echo "✅ Images built successfully"
echo ""

# Step 4: Deploy with Terraform
echo "📋 Step 4: Deploying with Terraform..."
cd terraform-k8s/local
terraform init -upgrade
terraform apply -var="kube_context=minikube" -auto-approve
echo ""

#Step 5: Wait for pods
echo "📋 Step 5: Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod --all -n usualstore --timeout=300s 2>/dev/null || true
echo ""

# Step 6: Show status
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Deployment Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 Pod Status:"
kubectl get pods -n usualstore
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 How to Access Your Application:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Option 1: Port forwarding (recommended)"
echo ""
echo "  Frontend:"
echo "    kubectl port-forward svc/frontend 3000:80 -n usualstore"
echo "    Then visit: http://localhost:3000"
echo ""
echo "  Backend API:"
echo "    kubectl port-forward svc/backend 4001:4001 -n usualstore"
echo "    Then visit: http://localhost:4001"
echo ""
echo "Option 2: Minikube service"
echo ""
echo "    minikube service frontend -n usualstore"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 Useful Commands:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "View logs:"
echo "  kubectl logs -f deployment/backend -n usualstore"
echo "  kubectl logs -f deployment/frontend -n usualstore"
echo ""
echo "Get all resources:"
echo "  kubectl get all -n usualstore"
echo ""
echo "Destroy everything:"
echo "  cd terraform-k8s/local && terraform destroy -auto-approve"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Happy deploying!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

