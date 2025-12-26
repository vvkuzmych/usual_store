# Terraform Infrastructure for Usual Store

## 🎯 Overview

This document describes the Terraform-based infrastructure management system for Usual Store, including Docker container orchestration and Open Policy Agent (OPA) policy enforcement.

## 📁 Project Structure

```
terraform/
├── main.tf                          # Main Terraform configuration
├── variables.tf                     # Input variables
├── outputs.tf                       # Output values
├── terraform.tfvars.example         # Example variables file
├── Makefile                         # Convenience commands
├── .gitignore                       # Terraform-specific ignores
├── README.md                        # Main documentation
├── POLICY-EXAMPLES.md               # Policy usage examples
│
├── modules/                         # Terraform modules
│   ├── database/                    # PostgreSQL module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   │
│   ├── kafka/                       # Kafka stack module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   │
│   ├── policies/                    # ⭐ Policy management module
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   │
│   │   ├── policies/                # OPA policy files
│   │   │   ├── network.rego         # Network policies
│   │   │   ├── resource_limits.rego # Resource management
│   │   │   ├── security.rego        # Security policies
│   │   │   └── access_control.rego  # Access control
│   │   │
│   │   └── policy-enforcer/         # Policy enforcement service
│   │       ├── main.go              # Go application
│   │       ├── Dockerfile
│   │       └── go.mod
│   │
│   ├── backend/                     # Backend API (placeholder)
│   ├── support/                     # Support service (placeholder)
│   ├── frontends/                   # Frontend services (placeholder)
│   ├── messaging/                   # Messaging service (placeholder)
│   └── observability/               # Jaeger tracing
│
└── backups/                         # State backups (created by Makefile)
```

## 🚀 Quick Start

### 1. Install Prerequisites

```bash
# Install Terraform, OPA, and other tools
cd terraform
make install-tools
```

### 2. Configure Variables

```bash
# Copy example config
cp terraform.tfvars.example terraform.tfvars

# Edit with your values
vi terraform.tfvars
```

### 3. Deploy Infrastructure

```bash
# Option 1: Step by step
make init
make plan
make apply

# Option 2: Quick start (all in one)
make quick-start
```

### 4. Verify Deployment

```bash
# Show all service URLs
make urls

# Check service health
make health-check

# View service status
make status
```

## 🏗️ Infrastructure Components

### Core Infrastructure

1. **Docker Network**: `usualstore_network`
   - IPv4: 172.22.0.0/16
   - IPv6: 2001:db8:1::/64
   - Bridge driver with full isolation

2. **Docker Volumes**:
   - `usualstore_db_data`: PostgreSQL data
   - `usualstore_kafka_data`: Kafka logs and data
   - `usualstore_zookeeper_data`: Zookeeper data

### Services

| Service | Port(s) | Purpose |
|---------|---------|---------|
| PostgreSQL | 5433 | Database |
| Kafka | 9092, 9093 | Message broker |
| Zookeeper | 2181 | Kafka coordination |
| Kafka UI | 8090 | Kafka management UI |
| Backend API | 4001 | REST API |
| Support Service | 5001 | Support/chat service |
| React Frontend | 3000 | Main web UI |
| TypeScript Frontend | 3001 | Alternative UI |
| Redux Frontend | 3002 | Redux-based UI |
| Support Frontend | 3005 | Support dashboard |
| Jaeger | 16686, 4318 | Distributed tracing |
| **OPA Server** | 8181 | **Policy evaluation** |
| **Policy Enforcer** | 8080 | **Policy monitoring** |

## 🔒 Policy Management System

### Architecture

```
┌─────────────────────────────────────────────────┐
│           Docker Container Events               │
│  (create, start, stop, network, volume)         │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│         Policy Enforcer (Go Service)            │
│  • Monitors Docker events in real-time          │
│  • Inspects container configurations            │
│  • Queries OPA for policy decisions             │
│  • Logs violations and recommendations          │
│  • Runs periodic compliance checks              │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│         OPA Server (Port 8181)                  │
│  • Evaluates policies written in Rego           │
│  • REST API for policy queries                  │
│  • Maintains policy bundles                     │
│  • Provides decision logs                       │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│         Policy Files (.rego)                    │
│  ├─ network.rego         (Layer 3 controls)     │
│  ├─ resource_limits.rego (CPU, memory, disk)    │
│  ├─ security.rego        (Security best practices) │
│  └─ access_control.rego  (RBAC, permissions)    │
└─────────────────────────────────────────────────┘
```

### Policy Categories

#### 1. Network Policies (`network.rego`)

**Purpose**: Control service-to-service communication using a tier-based model.

**Tiers**:
- Presentation (Frontends) → Application (APIs)
- Application → Data (Database)
- Application ↔ Messaging (Kafka)
- All → Observability (Jaeger)
- Governance (OPA) → All

**Rules**:
- ✅ Frontend can access Backend
- ❌ Frontend cannot directly access Database
- ✅ Backend can access Database
- ✅ Backend can publish to Kafka
- ✅ All services can report to Jaeger

**Example**:
```bash
# Check if react-frontend can access back-end
make network-check SOURCE=react-frontend TARGET=back-end
# Result: Allowed (Presentation → Application)

# Check if react-frontend can access database
make network-check SOURCE=react-frontend TARGET=database
# Result: Denied (Presentation ↛ Data)
```

#### 2. Resource Limit Policies (`resource_limits.rego`)

**Purpose**: Enforce CPU, memory, and disk limits per service tier.

**Default Limits**:
- Data tier (PostgreSQL): 2 CPU, 2GB RAM, 20GB disk
- Messaging tier (Kafka): 1 CPU, 1GB RAM, 10GB disk
- Application tier (APIs): 1 CPU, 512MB RAM, 5GB disk
- Presentation tier (Frontends): 0.5 CPU, 256MB RAM, 2GB disk

**Features**:
- Resource utilization monitoring
- Automatic recommendations for optimization
- Health check validation
- Volume size constraints

#### 3. Security Policies (`security.rego`)

**Purpose**: Enforce container security best practices.

**Requirements**:
- Containers must not run as privileged
- Must not run as root (except database)
- Read-only root filesystem (where applicable)
- Security options: `no-new-privileges`, `seccomp:default`
- Images must be from trusted registries
- No plain-text secrets in environment variables
- Limited capabilities (no NET_ADMIN, SYS_ADMIN, etc.)
- Health checks must not run as root

**Security Score**: 0-100 based on compliance

#### 4. Access Control Policies (`access_control.rego`)

**Purpose**: Implement Role-Based Access Control (RBAC).

**Roles**:
- `super_admin`: Full access (level 100)
- `admin`: User and container management (level 80)
- `supporter`: Support ticket access (level 50)
- `user`: Basic app access (level 10)

**Controls**:
- API endpoint permissions
- Container management actions
- Database access control
- Kafka topic ACLs
- User management rules
- Rate limiting
- Session validation
- MFA requirements for sensitive actions

## 💻 Common Operations

### Infrastructure Management

```bash
# Deploy everything
make apply

# Destroy everything
make destroy

# Update specific module
terraform apply -target=module.database

# Show infrastructure state
make show

# Get all service URLs
make urls
```

### Policy Operations

```bash
# Test all policies
make policy-test

# Format policy files
make policy-fmt

# Check policies for errors
make policy-check

# Query specific policy
make policy-query POLICY=usualstore/network/allow INPUT='{"source_service":"react-frontend","target_service":"back-end"}'

# Check network communication
make network-check SOURCE=react-frontend TARGET=back-end

# Run security scan
make security-scan
```

### Service Monitoring

```bash
# Health check all services
make health-check

# View logs
make logs-database
make logs-kafka
make logs-opa
make logs-enforcer

# Generate compliance report
make compliance-report
```

### State Management

```bash
# Backup state
make backup-state

# Restore state
make restore-state BACKUP=terraform.tfstate.20250101_120000

# Refresh state from actual infrastructure
terraform refresh
```

## 📊 Policy Enforcement Examples

### Example 1: Network Policy Violation

```bash
# Attempt: Frontend tries to directly connect to database
# Input: {source: "react-frontend", target: "database"}
# Result: DENIED (violation of tier model)
# Action: Policy enforcer logs violation
```

### Example 2: Security Policy Check

```bash
# Scan container configuration
# Check: Privileged mode, root user, capabilities
# Result: Security score 85/100
# Recommendation: Remove unused capabilities
```

### Example 3: Resource Optimization

```bash
# Monitor: Backend using 0.3 CPU (allocated 1 CPU)
# Analysis: 30% utilization
# Recommendation: Reduce to 0.5 CPU to save resources
```

### Example 4: Access Control Enforcement

```bash
# Request: Admin attempts to delete super_admin
# Check: Is this the last super_admin?
# Result: DENIED (protection rule)
# Reason: Cannot delete last super_admin
```

## 🔍 Troubleshooting

### Issue: Terraform state is corrupted

```bash
# Solution 1: Restore from backup
make restore-state BACKUP=terraform.tfstate.20250126_120000

# Solution 2: Import existing resources
terraform import docker_network.usualstore_network usualstore_network
```

### Issue: OPA policies not loading

```bash
# Check OPA logs
make logs-opa

# Validate policy syntax
make policy-check

# Rebuild OPA container
terraform destroy -target=module.policies
terraform apply -target=module.policies
```

### Issue: Policy enforcer not detecting violations

```bash
# Check enforcer status
docker ps | grep policy-enforcer

# View enforcer logs
make logs-enforcer

# Manually trigger enforcement
curl -X POST http://localhost:8080/enforce
```

### Issue: Container fails security check

```bash
# Get detailed security report
curl -X POST http://localhost:8181/v1/data/usualstore/security/security_violations \
  -d '{"input": {"container_name": "my-service"}}'

# View specific violations
make security-scan | jq '.violations'
```

## 🎓 Advanced Topics

### Custom Policy Development

See [POLICY-EXAMPLES.md](POLICY-EXAMPLES.md) for detailed examples.

### Extending Modules

To add a new service:

1. Create module directory: `modules/myservice/`
2. Add main.tf, variables.tf, outputs.tf
3. Reference in main.tf
4. Add policies if needed

### Integration with CI/CD

```yaml
# Example GitHub Actions
- name: Terraform Plan
  run: |
    cd terraform
    terraform init
    terraform plan

- name: Validate Policies
  run: |
    cd terraform
    make policy-test
```

### Multi-Environment Setup

Use workspace or separate directories:

```bash
# Workspaces
terraform workspace new production
terraform workspace select production
terraform apply -var-file=prod.tfvars

# Or separate directories
terraform/
  ├── dev/
  ├── staging/
  └── production/
```

## 📚 Resources

- [Main README](../terraform/README.md) - Detailed documentation
- [Policy Examples](../terraform/POLICY-EXAMPLES.md) - 20+ policy examples
- [OPA Documentation](https://www.openpolicyagent.org/docs/latest/)
- [Terraform Docker Provider](https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs)

## 🤝 Contributing

To contribute to the infrastructure:

1. Test changes locally
2. Validate policies: `make policy-test`
3. Format code: `make fmt`
4. Run linter: `make lint`
5. Submit PR with description

## 📝 Best Practices

1. **Always backup state** before major changes
2. **Test policies** before deploying
3. **Use version control** for policy files
4. **Monitor policy enforcer logs** regularly
5. **Review compliance reports** weekly
6. **Keep policies up to date** with architecture changes
7. **Document custom policies** clearly

## 🎉 Summary

The Terraform infrastructure provides:

- ✅ **Declarative infrastructure**: Version-controlled, repeatable deployments
- ✅ **Policy enforcement**: Automated compliance and security checks
- ✅ **Resource optimization**: Intelligent resource allocation
- ✅ **Access control**: Fine-grained RBAC
- ✅ **Audit logging**: Complete visibility into policy decisions
- ✅ **Self-service**: Developers can manage infrastructure safely

This system ensures your Usual Store infrastructure is secure, compliant, and efficiently managed!

