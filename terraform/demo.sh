#!/bin/bash

# Terraform Practical Demo Script
# This script demonstrates how to use Terraform for Usual Store

set -e

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  Terraform Infrastructure Demo - Usual Store              ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# Step 1: Show current directory
echo "📁 Step 1: Current Location"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
pwd
echo ""

# Step 2: Show available Makefile commands
echo "🔧 Step 2: Available Commands"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Quick commands:"
echo "  make init           - Initialize Terraform"
echo "  make plan           - Preview changes"
echo "  make apply          - Deploy infrastructure"
echo "  make urls           - Show service URLs"
echo "  make health-check   - Check all services"
echo "  make security-scan  - Run security audit"
echo "  make policy-test    - Test policies"
echo ""

# Step 3: Show configuration
echo "⚙️  Step 3: Configuration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f terraform.tfvars ]; then
    echo "✅ terraform.tfvars exists"
    echo "Configuration loaded:"
    grep -E "^[a-z_]+ = " terraform.tfvars | head -10
else
    echo "❌ terraform.tfvars not found"
    echo "Run: cp terraform.tfvars.example terraform.tfvars"
fi
echo ""

# Step 4: Show infrastructure state
echo "🏗️  Step 4: Infrastructure State"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f terraform.tfstate ]; then
    echo "✅ Infrastructure is deployed"
    echo "Resources: $(terraform state list 2>/dev/null | wc -l)"
else
    echo "⚪ Infrastructure not deployed yet"
    echo "Run 'make apply' to deploy"
fi
echo ""

# Step 5: Show modules
echo "📦 Step 5: Terraform Modules"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
ls -1 modules/
echo ""

# Step 6: Show policies
echo "🔒 Step 6: Policy Files"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Policy files (Rego):"
ls -lh modules/policies/policies/*.rego | awk '{print "  " $9 " (" $5 ")"}'
echo ""
echo "Total policy lines: $(wc -l modules/policies/policies/*.rego | tail -1 | awk '{print $1}')"
echo ""

# Step 7: Test policies (if OPA is available)
echo "🧪 Step 7: Policy Testing"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if command -v opa &> /dev/null; then
    echo "✅ OPA is installed"
    echo "Testing policy syntax..."
    cd modules/policies/policies
    opa check *.rego && echo "✅ All policies are valid!"
    cd ../../..
else
    echo "⚪ OPA not installed"
    echo "Install with: brew install opa"
fi
echo ""

# Step 8: Show what would be created
echo "📋 Step 8: Terraform Plan Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if terraform plan -detailed-exitcode > /dev/null 2>&1; then
    echo "No changes needed"
else
    echo "Resources to be created:"
    terraform plan 2>&1 | grep "will be created" | wc -l | xargs echo "  •"
fi
echo ""

# Step 9: Show example policy queries
echo "🔍 Step 9: Example Policy Queries"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat << 'EXAMPLES'
Once deployed, you can query policies:

# Check network access:
curl -X POST http://localhost:8181/v1/data/usualstore/network/allow \
  -d '{"input": {"source_service": "react-frontend", "target_service": "back-end"}}'

# Check user permissions:
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_api_access \
  -d '{"input": {"user": {"role": "admin"}, "endpoint": "/api/users", "method": "POST"}}'

# Get security score:
curl -X POST http://localhost:8181/v1/data/usualstore/security/security_score \
  -d '{"input": {"container_name": "database"}}'
EXAMPLES
echo ""

# Step 10: Next steps
echo "🚀 Step 10: Next Steps"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat << 'NEXT'
To deploy the infrastructure:

1. Configure your settings:
   vi terraform.tfvars

2. Initialize Terraform:
   make init

3. Preview changes:
   make plan

4. Deploy everything:
   make apply

5. Verify deployment:
   make urls
   make health-check

6. Monitor policies:
   make logs-enforcer

7. Run security scan:
   make security-scan

8. Generate compliance report:
   make compliance-report
NEXT
echo ""

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║  Demo Complete!                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""
echo "📚 Documentation:"
echo "  • terraform/README.md              - Full guide"
echo "  • terraform/POLICY-EXAMPLES.md     - Policy examples"
echo "  • terraform/QUICK-START.md         - Quick reference"
echo ""
echo "💡 Quick command: make help"

