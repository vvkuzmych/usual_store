# Terraform Infrastructure for Usual Store

This directory contains Terraform configuration for managing Usual Store's local Docker infrastructure with policy enforcement using Open Policy Agent (OPA).

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Policy Management](#policy-management)
- [Usage](#usage)
- [Modules](#modules)
- [Policy Enforcement](#policy-enforcement)
- [Troubleshooting](#troubleshooting)

## 🎯 Overview

This Terraform setup provides:

- **Infrastructure as Code**: Manage Docker containers, networks, and volumes declaratively
- **Policy as Code**: Enforce security, network, resource, and access control policies using OPA
- **Service Orchestration**: Automated deployment of all Usual Store services
- **Compliance Monitoring**: Continuous policy enforcement and audit logging
- **Resource Management**: Automated resource allocation and limits

## 📦 Prerequisites

1. **Terraform** >= 1.0
   ```bash
   brew install terraform  # macOS
   ```

2. **Docker** and **Docker Desktop** >= 4.42 (for IPv6 support)
   
3. **OPA** (Open Policy Agent) - pulled automatically as Docker image

4. **Go** >= 1.23 (for policy enforcer)

## 🚀 Quick Start

### 1. Initialize Terraform

```bash
cd terraform
terraform init
```

### 2. Configure Variables

Copy the example variables file:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your configuration:

```hcl
postgres_password = "your_secure_password"
stripe_key        = "pk_test_..."
stripe_secret     = "sk_test_..."
openai_api_key    = "sk-..."  # Optional

# Feature flags
enable_kafka          = true
enable_observability  = true
enable_ai_assistant   = false
enable_support        = true
```

### 3. Plan and Apply

```bash
# Preview changes
terraform plan

# Apply configuration
terraform apply
```

### 4. Verify Deployment

```bash
# Check service status
terraform output service_status

# Get service URLs
terraform output frontend_urls
terraform output backend_api_url

# Access OPA policy server
curl http://localhost:8181/v1/policies
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Terraform Main Config                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Network: usualstore_network (IPv4 + IPv6)          │   │
│  │  Volumes: db_data, kafka_data, zookeeper_data       │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
           │
           ├──> Database Module (PostgreSQL)
           │
           ├──> Kafka Stack Module (Zookeeper, Kafka, Kafka UI)
           │
           ├──> Backend API Module
           │
           ├──> Support Service Module
           │
           ├──> Frontends Module (React, TypeScript, Redux, Support UI)
           │
           ├──> Messaging Service Module
           │
           ├──> Observability Module (Jaeger)
           │
           └──> Policies Module ⭐ (OPA + Policy Enforcer)
                 ├─ OPA Server (port 8181)
                 ├─ Policy Enforcer (monitors Docker events)
                 └─ Policy Files:
                    ├─ network.rego         (Network policies)
                    ├─ resource_limits.rego (Resource management)
                    ├─ security.rego        (Security policies)
                    └─ access_control.rego  (Access policies)
```

## 🔒 Policy Management

### Policy Types

1. **Network Policies** (`network.rego`)
   - Service-to-service communication rules
   - Tier-based access control (presentation → application → data)
   - Port validation
   - Network isolation

2. **Resource Limits** (`resource_limits.rego`)
   - CPU and memory limits by service tier
   - Health check configuration
   - Volume size constraints
   - Resource utilization recommendations

3. **Security Policies** (`security.rego`)
   - Container security best practices
   - Image registry validation
   - Capability restrictions
   - Volume mount security
   - Secrets management

4. **Access Control** (`access_control.rego`)
   - Role-based access control (RBAC)
   - API endpoint permissions
   - Database access control
   - Kafka topic ACLs
   - User management rules

### Policy Evaluation

Query policies using OPA REST API:

```bash
# Check if a network connection is allowed
curl -X POST http://localhost:8181/v1/data/usualstore/network/allow \
  -d '{
    "input": {
      "source_service": "react-frontend",
      "target_service": "back-end"
    }
  }'

# Check container security compliance
curl -X POST http://localhost:8181/v1/data/usualstore/security/security_score \
  -d '{
    "input": {
      "container_name": "database",
      "privileged": false,
      "user": "postgres"
    }
  }'

# Check API access
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_api_access \
  -d '{
    "input": {
      "user": {"role": "admin", "email": "admin@example.com"},
      "endpoint": "/api/users",
      "method": "POST"
    }
  }'
```

## 💻 Usage

### Managing Infrastructure

```bash
# Show current state
terraform show

# List all resources
terraform state list

# Get specific output
terraform output database_connection_string

# Update specific resource
terraform apply -target=module.database

# Destroy specific module
terraform destroy -target=module.kafka_stack

# Destroy everything
terraform destroy
```

### Policy Enforcement

The Policy Enforcer automatically:
- Monitors Docker events in real-time
- Validates container configurations against policies
- Logs violations and recommendations
- Runs periodic compliance checks (every 5 minutes)

View policy enforcer logs:

```bash
docker logs usualstore-policy-enforcer -f
```

### Manual Policy Check

```bash
# Check a specific container
curl -X POST http://localhost:8080/enforce \
  -d '{"container_id": "abc123"}'

# View audit logs
curl http://localhost:8080/audit
```

## 📚 Modules

### Database Module
- **Purpose**: PostgreSQL database container
- **Resources**: Container, health checks
- **Outputs**: Connection string, container ID

### Kafka Stack Module
- **Purpose**: Kafka messaging infrastructure
- **Resources**: Zookeeper, Kafka broker, Kafka UI
- **Outputs**: Kafka broker address, UI URL

### Policies Module ⭐
- **Purpose**: Policy enforcement and compliance
- **Resources**:
  - OPA Server: Policy evaluation engine
  - Policy Enforcer: Monitors and enforces policies
  - Policy Files: Rego policy definitions
- **Outputs**: OPA server URL, policy status

### Backend/Frontend/Support Modules
- **Purpose**: Application services (placeholders)
- **Note**: Extend these with your custom container definitions

### Observability Module
- **Purpose**: Distributed tracing
- **Resources**: Jaeger all-in-one
- **Outputs**: Jaeger UI URL

## 🛡️ Policy Enforcement

### Service Tier Model

The system uses a tier-based architecture for policy enforcement:

```
┌─────────────────────┐
│  Presentation Tier  │  (Frontends)
│  ↓ allowed          │
├─────────────────────┤
│  Application Tier   │  (Backend APIs, Services)
│  ↓ allowed          │
├─────────────────────┤
│  Data Tier          │  (Database)
│  ✗ no direct access │
│    from Presentation│
├─────────────────────┤
│  Messaging Tier     │  (Kafka, Zookeeper)
│  ↓ app tier only    │
├─────────────────────┤
│  Observability Tier │  (Jaeger, Monitoring)
│  ← all can report   │
├─────────────────────┤
│  Governance Tier    │  (OPA, Policy Enforcer)
│  → can access all   │
└─────────────────────┘
```

### Default Resource Limits

| Tier          | CPU | Memory | Disk |
|---------------|-----|--------|------|
| Data          | 2   | 2g     | 20g  |
| Messaging     | 1   | 1g     | 10g  |
| Application   | 1   | 512m   | 5g   |
| Presentation  | 0.5 | 256m   | 2g   |
| Observability | 1   | 1g     | 10g  |
| Governance    | 0.5 | 256m   | 1g   |

### Security Requirements

All containers must:
- ✅ Run as non-root (except database)
- ✅ Use read-only root filesystem (where applicable)
- ✅ Have `no-new-privileges` security option
- ✅ Use seccomp default profile
- ✅ Not run in privileged mode
- ✅ Use images from trusted registries
- ✅ Have health checks configured
- ✅ Store secrets securely (not in env vars)

## 🔧 Troubleshooting

### Terraform Issues

**Problem**: Provider initialization fails
```bash
# Solution: Clear cache and re-initialize
rm -rf .terraform .terraform.lock.hcl
terraform init
```

**Problem**: State is out of sync
```bash
# Solution: Refresh state
terraform refresh

# Or import existing resource
terraform import docker_network.usualstore_network usualstore_network
```

### Policy Issues

**Problem**: OPA server not responding
```bash
# Check OPA container
docker logs usualstore-opa-server

# Restart OPA
docker restart usualstore-opa-server
```

**Problem**: Policy violations not detected
```bash
# Check policy enforcer logs
docker logs usualstore-policy-enforcer -f

# Manually trigger compliance check
curl -X POST http://localhost:8080/enforce
```

### Docker Issues

**Problem**: Containers can't communicate
```bash
# Verify network
docker network inspect usualstore_network

# Check if containers are on the network
docker network inspect usualstore_network | jq '.[0].Containers'
```

**Problem**: Port already in use
```bash
# Find process using port
lsof -i :5432

# Change port in terraform.tfvars
api_port = 4002  # Use different port
```

## 📖 Further Reading

- [OPA Documentation](https://www.openpolicyagent.org/docs/latest/)
- [Rego Language Guide](https://www.openpolicyagent.org/docs/latest/policy-language/)
- [Docker Provider for Terraform](https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs)
- [Terraform Best Practices](https://www.terraform.io/docs/cloud/guides/recommended-practices/index.html)

## 🤝 Contributing

To extend this infrastructure:

1. **Add new service**: Create a new module in `modules/`
2. **Add new policy**: Create `.rego` file in `modules/policies/policies/`
3. **Modify policy enforcer**: Edit `modules/policies/policy-enforcer/main.go`

## 📝 License

Same as Usual Store project.

