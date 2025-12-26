#!/bin/bash
# Quick Start script for deploying Usual Store to LOCAL Kubernetes

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Usual Store Kubernetes Deployment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Step 1: Check prerequisites
echo "📋 Step 1: Checking prerequisites..."
echo ""

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed"
    echo "   Install: brew install kubectl"
    exit 1
fi
echo "✅ kubectl found"

# Check if terraform is installed
if ! command -v terraform &> /dev/null; then
    echo "❌ Terraform is not installed"
    echo "   Install: brew install terraform"
    exit 1
fi
echo "✅ Terraform found"

# Check if Kubernetes cluster is running
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Kubernetes cluster is not running"
    echo "   Options:"
    echo "   - Docker Desktop: Enable Kubernetes in Settings"
    echo "   - Minikube: Run 'minikube start'"
    echo "   - Kind: Run 'kind create cluster --name usualstore'"
    exit 1
fi
echo "✅ Kubernetes cluster is running"

# Detect which K8s cluster
CONTEXT=$(kubectl config current-context)
echo "   Using context: $CONTEXT"
echo ""

# Step 2: Build Docker images if needed
echo "📦 Step 2: Checking Docker images..."
echo ""

if ! docker images | grep -q "usual_store-back-end"; then
    echo "⚠️  Docker images not found. Building them..."
    cd .. && docker-compose build
    echo "✅ Docker images built"
else
    echo "✅ Docker images already exist"
fi
echo ""

# Step 3: Initialize Terraform
echo "🔧 Step 3: Initializing Terraform..."
echo ""
cd local
terraform init
echo ""

# Step 4: Deploy
echo "🚀 Step 4: Deploying to Kubernetes..."
echo ""

# Use the detected context
if [[ "$CONTEXT" == "minikube" ]]; then
    terraform apply -var="kube_context=minikube" -auto-approve
elif [[ "$CONTEXT" == kind-* ]]; then
    terraform apply -var="kube_context=$CONTEXT" -auto-approve
else
    terraform apply -auto-approve
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Deployment Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Step 5: Wait for pods to be ready
echo "⏳ Waiting for pods to be ready (this may take 1-2 minutes)..."
echo ""
kubectl wait --for=condition=ready pod --all -n usualstore --timeout=300s 2>/dev/null || true

# Show status
echo ""
echo "📊 Current Status:"
echo ""
kubectl get pods -n usualstore
echo ""

# Step 6: Show access instructions
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 How to Access Your Application"
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

if [[ "$CONTEXT" == "minikube" ]]; then
    echo "Option 2: Minikube service (for Minikube)"
    echo ""
    echo "    minikube service frontend -n usualstore"
    echo ""
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 Useful Commands"
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
echo "  terraform destroy -auto-approve"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Happy deploying!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

