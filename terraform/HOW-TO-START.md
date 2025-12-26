# How to Start All Docker Containers with Terraform

## 🎯 Quick Start (3 Steps)

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
terraform init      # First time only
terraform apply     # Starts everything
```

## 📋 Detailed Step-by-Step Guide

### Method 1: Interactive Script (Recommended)

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
./START-CONTAINERS.sh
```

This script guides you through each step with confirmations.

### Method 2: Makefile (Fastest)

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
make quick-start    # Does everything automatically
```

### Method 3: Manual Terraform Commands

#### Step 1: Navigate to Terraform Directory

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform
```

#### Step 2: Configure (First Time Only)

```bash
# Copy example configuration
cp terraform.tfvars.example terraform.tfvars

# Edit configuration
vi terraform.tfvars

# Update these values:
# postgres_password = "your_secure_password"
# stripe_key = "pk_test_..." (optional)
# stripe_secret = "sk_test_..." (optional)
```

#### Step 3: Initialize (First Time Only)

```bash
terraform init
```

**Expected output:**
```
Initializing modules...
Initializing the backend...
Initializing provider plugins...
- Installing kreuzwerker/docker v3.6.2...
- Installing hashicorp/local v2.6.1...
- Installing hashicorp/null v3.2.4...

Terraform has been successfully initialized!
```

#### Step 4: Validate Configuration

```bash
terraform validate
```

**Expected output:**
```
Success! The configuration is valid.
```

#### Step 5: Preview Changes

```bash
terraform plan
```

**Expected output:**
```
Plan: 18 to add, 0 to change, 0 to destroy.

Resources to create:
- docker_network.usualstore_network
- docker_volume.db_data
- docker_volume.kafka_data
- docker_volume.zookeeper_data
- module.database.docker_image.postgres
- module.database.docker_container.database
- module.kafka_stack.docker_image.zookeeper
- module.kafka_stack.docker_container.zookeeper
- module.kafka_stack.docker_image.kafka
- module.kafka_stack.docker_container.kafka
- module.kafka_stack.docker_image.kafka_ui
- module.kafka_stack.docker_container.kafka_ui
- module.observability.docker_image.jaeger
- module.observability.docker_container.jaeger
- module.policies.docker_image.opa
- module.policies.docker_container.opa_server
- module.policies.docker_image.policy_enforcer
- module.policies.docker_container.policy_enforcer
```

#### Step 6: Start Containers

```bash
terraform apply
```

Type `yes` when prompted.

**Expected output:**
```
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

docker_network.usualstore_network: Creating...
docker_volume.db_data: Creating...
docker_volume.kafka_data: Creating...
docker_volume.zookeeper_data: Creating...
[... more creation logs ...]

Apply complete! Resources: 18 added, 0 changed, 0 destroyed.

Outputs:

backend_api_url = "http://localhost:4001"
database_connection_string = <sensitive>
frontend_urls = {
  "react" = "http://localhost:3000"
  "redux" = "http://localhost:3002"
  "support_ui" = "http://localhost:3005"
  "typescript" = "http://localhost:3001"
}
jaeger_ui_url = "http://localhost:16686"
kafka_broker = "kafka:9092"
kafka_ui_url = "http://localhost:8090"
policy_server_url = "http://localhost:8181"
```

#### Step 7: Verify Deployment

```bash
# Check running containers
docker ps | grep usualstore

# View service URLs
terraform output

# Check specific service
docker logs usualstore-database

# Check health
docker ps --format "{{.Names}}: {{.Status}}" | grep usualstore
```

## 🐳 What Gets Started

| Container | Port(s) | Purpose |
|-----------|---------|---------|
| usualstore-database | 5433 | PostgreSQL database |
| usualstore-zookeeper | 2181 | Kafka coordination |
| usualstore-kafka | 9092, 9093 | Message broker |
| usualstore-kafka-ui | 8090 | Kafka management UI |
| usualstore-jaeger | 16686, 4318 | Distributed tracing |
| usualstore-opa-server | 8181 | Policy evaluation engine |
| usualstore-policy-enforcer | 8080 | Policy monitoring service |

**Plus:**
- Docker network: `usualstore_network` (IPv4 + IPv6)
- Volumes: `usualstore_db_data`, `usualstore_kafka_data`, `usualstore_zookeeper_data`

## 🔍 Verification Commands

```bash
# List all Terraform-managed containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep usualstore

# Check network
docker network inspect usualstore_network

# Check volumes
docker volume ls | grep usualstore

# View all outputs
terraform output

# Check specific output
terraform output kafka_ui_url

# Access services
curl http://localhost:8181/health  # OPA health
curl http://localhost:8090         # Kafka UI
open http://localhost:16686        # Jaeger UI
```

## 🛑 Stopping Containers

### Option 1: Stop and Keep Data

```bash
terraform destroy
```

This removes containers but keeps volumes (data persists).

### Option 2: Stop with Makefile

```bash
make destroy
```

### Option 3: Remove Everything (Including Data)

```bash
terraform destroy
docker volume rm usualstore_db_data usualstore_kafka_data usualstore_zookeeper_data
```

## 🔄 Restart Containers

```bash
# Stop
terraform destroy

# Start again
terraform apply
```

Or use Makefile:

```bash
make destroy
make apply
```

## 📊 Monitoring

```bash
# View logs
docker logs usualstore-database -f
docker logs usualstore-kafka -f
docker logs usualstore-opa-server -f
docker logs usualstore-policy-enforcer -f

# Check policy enforcer
curl http://localhost:8080/health

# View OPA policies
curl http://localhost:8181/v1/policies

# Check Kafka
open http://localhost:8090
```

## ⚠️ Troubleshooting

### Port Already in Use

```bash
# Find what's using the port
lsof -i :5433

# Stop existing containers
docker ps -a | grep usualstore
docker stop <container_id>
```

### Terraform State Lock

```bash
# If terraform is stuck
rm -rf .terraform.lock.info
```

### Containers Not Starting

```bash
# Check logs
docker logs usualstore-database

# Rebuild
terraform destroy
terraform apply
```

### Configuration Errors

```bash
# Validate configuration
terraform validate

# Format files
terraform fmt -recursive
```

## 🎓 Advanced Usage

### Selective Deployment

```bash
# Deploy only database
terraform apply -target=module.database

# Deploy specific resources
terraform apply -target=docker_network.usualstore_network
```

### Using Different Configurations

```bash
# Use custom variable file
terraform apply -var-file=production.tfvars
```

### Check What Will Change

```bash
# Detailed plan
terraform plan -out=tfplan

# Apply the plan
terraform apply tfplan
```

## 📚 Related Documentation

- [Main Terraform README](README.md) - Complete guide
- [Policy Examples](POLICY-EXAMPLES.md) - OPA policy usage
- [Quick Start](QUICK-START.md) - Quick reference
- [Cheat Sheet](CHEAT-SHEET.md) - Command reference

## 💡 Tips

1. **Always run `terraform plan`** before `apply` to see changes
2. **Use Makefile commands** for convenience: `make apply`, `make destroy`
3. **Check logs** if something fails: `docker logs <container_name>`
4. **Keep backups** of terraform.tfvars (it contains passwords)
5. **Use `terraform output`** to get service URLs anytime

## 🚀 Quick Commands Summary

```bash
# Start everything
cd terraform && terraform apply

# Stop everything
cd terraform && terraform destroy

# Check status
docker ps | grep usualstore

# View logs
docker logs usualstore-database -f

# Get URLs
terraform output frontend_urls

# Health check
make health-check
```

---

**Ready to start?** Run: `cd terraform && ./START-CONTAINERS.sh`
