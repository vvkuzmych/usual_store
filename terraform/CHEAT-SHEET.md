# Terraform Cheat Sheet - Usual Store

## 🚀 Getting Started (First Time)

```bash
cd /path/to/usual_store/terraform

# 1. Install tools (macOS)
make install-tools

# 2. Initialize Terraform
make init

# 3. Create configuration
cp terraform.tfvars.example terraform.tfvars
vi terraform.tfvars  # Edit with your settings

# 4. Preview changes
make plan

# 5. Deploy!
make apply
```

## 📝 Daily Commands

### Infrastructure Management
```bash
make plan           # Preview changes before applying
make apply          # Deploy infrastructure
make destroy        # Remove everything
make status         # Show service status
make urls           # Show all service URLs
make health-check   # Check all services
```

### Viewing Logs
```bash
make logs-database     # Database logs
make logs-kafka        # Kafka logs
make logs-opa          # OPA policy server
make logs-enforcer     # Policy enforcer
```

### Policy Management
```bash
make policy-test       # Test all policies
make policy-check      # Validate policies
make policy-fmt        # Format policy files
make security-scan     # Run security audit
make compliance-report # Generate compliance report
```

### Network Policy Checks
```bash
# Check if service A can access service B
make network-check SOURCE=react-frontend TARGET=back-end
make network-check SOURCE=back-end TARGET=database
make network-check SOURCE=react-frontend TARGET=database  # Should be denied
```

## 🔍 Direct Terraform Commands

```bash
terraform init              # Initialize (download providers)
terraform plan              # Preview changes
terraform apply             # Deploy (asks for confirmation)
terraform apply -auto-approve  # Deploy without asking
terraform destroy           # Remove all infrastructure
terraform show              # Show current state
terraform state list        # List all resources
terraform output            # Show all outputs
terraform output database_connection_string  # Specific output
terraform fmt -recursive    # Format all .tf files
terraform validate          # Check syntax
```

## 📊 OPA Policy Queries

Once deployed (after `make apply`), query policies:

### Network Policies
```bash
# Check if frontend can access backend
curl -X POST http://localhost:8181/v1/data/usualstore/network/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "source_service": "react-frontend",
      "target_service": "back-end"
    }
  }'

# Check if frontend can access database (should be denied)
curl -X POST http://localhost:8181/v1/data/usualstore/network/deny \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "source_service": "react-frontend",
      "target_service": "database"
    }
  }'
```

### Security Policies
```bash
# Get security score for a container
curl -X POST http://localhost:8181/v1/data/usualstore/security/security_score \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "container_name": "database",
      "privileged": false,
      "user": "postgres"
    }
  }'
```

### Access Control
```bash
# Check user API access
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_api_access \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": {"role": "admin", "email": "admin@example.com"},
      "endpoint": "/api/users",
      "method": "POST"
    }
  }'
```

## 🐳 Docker Commands

```bash
# List Terraform-managed containers
docker ps | grep usualstore

# View specific container
docker logs usualstore-database -f
docker inspect usualstore-database

# Network inspection
docker network inspect usualstore_network

# Volume inspection
docker volume ls | grep usualstore
docker volume inspect usualstore_db_data
```

## 🔧 Troubleshooting

### Reset Everything
```bash
make destroy
make clean
rm terraform.tfstate*
make init
make apply
```

### Check Policy Enforcer
```bash
# View real-time policy enforcement
docker logs usualstore-policy-enforcer -f

# Manual policy check
curl -X POST http://localhost:8080/enforce

# View audit logs
curl http://localhost:8080/audit
```

### Terraform State Issues
```bash
# Backup state
make backup-state

# Refresh state from actual infrastructure
terraform refresh

# Import existing resource
terraform import docker_network.usualstore_network usualstore_network

# Restore from backup
make restore-state BACKUP=terraform.tfstate.20250126_120000
```

## 📈 Workflow Examples

### Scenario 1: Deploy Fresh Infrastructure
```bash
make quick-start  # Does everything!
```

### Scenario 2: Update Configuration
```bash
vi terraform.tfvars     # Change settings
make plan               # Preview changes
make apply              # Apply changes
```

### Scenario 3: Add New Service
```bash
vi main.tf                   # Add new module
terraform fmt                # Format
terraform validate           # Check syntax
make plan                    # Preview
make apply                   # Deploy
```

### Scenario 4: Security Audit
```bash
make security-scan           # Run scan
make compliance-report       # Generate report
cat compliance-report.txt    # Review
```

### Scenario 5: Debug Policy Issues
```bash
make logs-enforcer           # Check enforcer logs
make logs-opa               # Check OPA logs
curl http://localhost:8181/v1/policies  # List all policies
```

## 🎯 Common Outputs

### Get Service URLs
```bash
terraform output frontend_urls
# Output:
# {
#   "react" = "http://localhost:3000"
#   "typescript" = "http://localhost:3001"
#   "redux" = "http://localhost:3002"
#   "support_ui" = "http://localhost:3005"
# }
```

### Get Service Status
```bash
terraform output service_status
# Output shows running/stopped status of all services
```

### Get Database Connection
```bash
terraform output database_connection_string
# Output: postgres://postgres:password@database:5432/usualstore?sslmode=disable
```

## 🔐 Environment Variables

```bash
# Set Terraform log level
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log

# Use different variable file
terraform apply -var-file=prod.tfvars
```

## 📚 Quick Reference

| Command | What It Does |
|---------|--------------|
| `make help` | Show all commands |
| `make quick-start` | Complete setup + deploy |
| `make init` | Initialize Terraform |
| `make plan` | Preview changes |
| `make apply` | Deploy infrastructure |
| `make destroy` | Remove everything |
| `make urls` | Show service URLs |
| `make health-check` | Check all services |
| `make security-scan` | Security audit |
| `make policy-test` | Test policies |
| `make logs-*` | View service logs |
| `make backup-state` | Backup state file |

## 🆘 Help

- `make help` - Show all Makefile commands
- `terraform/README.md` - Full documentation
- `terraform/POLICY-EXAMPLES.md` - Policy examples
- `terraform/QUICK-START.md` - Quick reference

## 💡 Pro Tips

1. **Always backup state before major changes**: `make backup-state`
2. **Review plan before applying**: `make plan`
3. **Use workspaces for multiple environments**: `terraform workspace new staging`
4. **Format files before committing**: `make fmt`
5. **Test policies after changes**: `make policy-test`
6. **Monitor enforcer in real-time**: `make logs-enforcer`
7. **Generate reports regularly**: `make compliance-report`

---

**Need more help?** Run `make help` or check the documentation files!
