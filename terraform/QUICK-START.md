# Terraform Quick Start Guide

## ⚡ 30-Second Setup

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform

# 1. Initialize
terraform init

# 2. Configure (edit with your settings)
cp terraform.tfvars.example terraform.tfvars

# 3. Deploy!
terraform apply
```

## 🎯 Essential Commands

```bash
# Show what will change
terraform plan

# Deploy infrastructure
terraform apply

# Destroy everything
terraform destroy

# Show current state
terraform show

# Get service URLs
terraform output frontend_urls
```

## 🔧 Makefile Shortcuts

```bash
# Complete setup (install tools + deploy)
make quick-start

# Show all available commands
make help

# Get all service URLs
make urls

# Check service health
make health-check

# View logs
make logs-opa          # OPA server
make logs-enforcer     # Policy enforcer
make logs-database     # Database

# Policy operations
make policy-test              # Test policies
make security-scan            # Security audit
make compliance-report        # Generate report
make network-check SOURCE=react-frontend TARGET=back-end
```

## 📊 Service Ports

| Service | Port | URL |
|---------|------|-----|
| React Frontend | 3000 | http://localhost:3000 |
| TypeScript Frontend | 3001 | http://localhost:3001 |
| Redux Frontend | 3002 | http://localhost:3002 |
| Support Frontend | 3005 | http://localhost:3005 |
| Backend API | 4001 | http://localhost:4001 |
| Support Service | 5001 | http://localhost:5001 |
| Database | 5433 | localhost:5433 |
| Kafka | 9092, 9093 | localhost:9092 |
| Kafka UI | 8090 | http://localhost:8090 |
| Jaeger UI | 16686 | http://localhost:16686 |
| **OPA Server** | 8181 | http://localhost:8181 |
| **Policy Enforcer** | 8080 | http://localhost:8080 |

## 🔍 Policy Examples

### Check Network Access
```bash
curl -X POST http://localhost:8181/v1/data/usualstore/network/allow \
  -H "Content-Type: application/json" \
  -d '{"input": {"source_service": "react-frontend", "target_service": "back-end"}}'
```

### Check User Permissions
```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_api_access \
  -H "Content-Type: application/json" \
  -d '{"input": {"user": {"role": "admin"}, "endpoint": "/api/users", "method": "POST"}}'
```

### Security Scan
```bash
curl http://localhost:8080/enforce
```

## 📁 Key Files

- `main.tf` - Infrastructure definition
- `variables.tf` - Configuration options
- `outputs.tf` - Service URLs and info
- `terraform.tfvars` - Your settings (create from .example)
- `modules/policies/policies/*.rego` - Policy rules

## 🆘 Troubleshooting

### Reset Everything
```bash
make destroy
rm -rf .terraform .terraform.lock.hcl terraform.tfstate*
make init
make apply
```

### View Logs
```bash
docker logs usualstore-opa-server
docker logs usualstore-policy-enforcer
docker logs usualstore-database
```

### Check Container Status
```bash
docker ps | grep usualstore
docker network inspect usualstore_network
```

## 📚 Documentation

- `README.md` - Full documentation
- `POLICY-EXAMPLES.md` - 20+ policy examples
- `../docs/TERRAFORM-INFRASTRUCTURE.md` - Architecture guide

## 🎉 You're Done!

Infrastructure deployed with:
- ✅ Docker containers managed by Terraform
- ✅ Network, resource, and security policies
- ✅ Real-time policy enforcement
- ✅ Role-based access control
- ✅ Monitoring and compliance reporting

**Next**: Try `make security-scan` or `make compliance-report`!

