#!/bin/bash

# Terraform Docker Startup Script
# This script demonstrates how to start all Docker containers using Terraform

set -e

echo "╔═══════════════════════════════════════════════════════════════════════╗"
echo "║  🚀 STARTING ALL DOCKER CONTAINERS WITH TERRAFORM                     ║"
echo "╚═══════════════════════════════════════════════════════════════════════╝"
echo ""

# Navigate to terraform directory
cd "$(dirname "$0")"
echo "📍 Working directory: $(pwd)"
echo ""

# ============================================================================
# STEP 1: Configuration Check
# ============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "STEP 1: Configuration Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f terraform.tfvars ]; then
    echo "⚠️  terraform.tfvars not found!"
    echo "Creating from example..."
    cp terraform.tfvars.example terraform.tfvars
    echo "✅ Created terraform.tfvars"
    echo ""
    echo "⚠️  IMPORTANT: Edit terraform.tfvars and update:"
    echo "   • postgres_password"
    echo "   • stripe_key (optional)"
    echo "   • stripe_secret (optional)"
    echo ""
    read -p "Press Enter to continue with default values, or Ctrl+C to edit first..."
else
    echo "✅ Configuration file exists"
fi
echo ""

# ============================================================================
# STEP 2: Initialize Terraform
# ============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "STEP 2: Initialize Terraform"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -d .terraform ]; then
    echo "Initializing Terraform..."
    terraform init
    echo "✅ Terraform initialized"
else
    echo "✅ Terraform already initialized"
fi
echo ""

# ============================================================================
# STEP 3: Validate Configuration
# ============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "STEP 3: Validate Configuration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if terraform validate; then
    echo "✅ Configuration is valid"
else
    echo "❌ Configuration validation failed!"
    exit 1
fi
echo ""

# ============================================================================
# STEP 4: Show Plan
# ============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "STEP 4: Preview Infrastructure Changes"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "Running terraform plan..."
echo ""
terraform plan -out=tfplan 2>&1 | grep -E "(will be created|Plan:)" || true
echo ""

# ============================================================================
# STEP 5: Apply (Start Containers)
# ============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "STEP 5: Start Docker Containers"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "This will create:"
echo "  • Docker network (usualstore_network)"
echo "  • Docker volumes (db_data, kafka_data, zookeeper_data)"
echo "  • PostgreSQL container (port 5433)"
echo "  • Kafka + Zookeeper containers (ports 9092, 9093, 2181)"
echo "  • Kafka UI (port 8090)"
echo "  • Jaeger tracing (ports 16686, 4318)"
echo "  • OPA policy server (port 8181)"
echo "  • Policy Enforcer (port 8080)"
echo ""

read -p "Do you want to proceed? (yes/no): " -r
echo ""

if [[ $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
    echo "🚀 Starting containers..."
    echo ""
    terraform apply tfplan
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✅ Containers started successfully!"
    else
        echo ""
        echo "❌ Failed to start containers"
        exit 1
    fi
else
    echo "⏸️  Deployment cancelled"
    rm -f tfplan
    exit 0
fi

# ============================================================================
# STEP 6: Verify Deployment
# ============================================================================
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "STEP 6: Verify Deployment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "🐳 Running containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep usualstore || echo "No containers found"
echo ""

echo "🌐 Service URLs:"
terraform output -json frontend_urls 2>/dev/null | grep -v "^null$" || echo "N/A"
echo ""

echo "📊 Service Status:"
terraform output service_status 2>/dev/null || echo "N/A"
echo ""

# ============================================================================
# Summary
# ============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ DEPLOYMENT COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Quick Commands:"
echo "   • Check containers:    docker ps | grep usualstore"
echo "   • View logs:           docker logs usualstore-database"
echo "   • Stop containers:     terraform destroy"
echo "   • Service URLs:        terraform output urls"
echo "   • Health check:        make health-check"
echo ""
echo "📚 Documentation:"
echo "   • terraform/README.md"
echo "   • terraform/QUICK-START.md"
echo "   • terraform/CHEAT-SHEET.md"
echo ""
echo "🎉 Your infrastructure is running!"
echo ""

