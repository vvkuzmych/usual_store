# OPA Policy Examples

This document provides practical examples of using OPA policies to manage your Usual Store infrastructure.

## 📚 Table of Contents

- [Network Policy Examples](#network-policy-examples)
- [Resource Management Examples](#resource-management-examples)
- [Security Policy Examples](#security-policy-examples)
- [Access Control Examples](#access-control-examples)
- [Custom Policy Development](#custom-policy-development)

## 🌐 Network Policy Examples

### Example 1: Check if Frontend Can Access Backend

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/network/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "source_service": "react-frontend",
      "target_service": "back-end",
      "source_network": "usualstore_network",
      "target_network": "usualstore_network"
    }
  }'
```

**Expected Result**: `{"result": true}` (Allowed: Presentation → Application tier)

### Example 2: Check if Frontend Can Directly Access Database

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/network/deny \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "source_service": "react-frontend",
      "target_service": "database"
    }
  }'
```

**Expected Result**: `{"result": true}` (Denied: Direct Presentation → Data tier access)

### Example 3: Validate Port Configuration

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/network/valid_port_config \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "service": "back-end",
      "port": 4001
    }
  }'
```

**Expected Result**: `{"result": true}` (Valid port for backend service)

### Example 4: Check Network Isolation Requirements

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/network/requires_isolation \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "service": "database"
    }
  }'
```

**Expected Result**: `{"result": true}` (Database requires isolation)

## 💾 Resource Management Examples

### Example 5: Validate Container Resources

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/resources/valid_resources \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "container_name": "back-end",
      "resources": {
        "cpu": "0.5",
        "memory": "256m"
      }
    }
  }'
```

**Expected Result**: `{"result": true}` (Within application tier limits)

### Example 6: Get Resource Utilization Recommendations

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/resources/recommendations \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "container_name": "database",
      "resources": {
        "cpu": "0.5",
        "memory": "500m"
      }
    }
  }'
```

**Expected Result**: Recommendations for resource optimization

### Example 7: Validate Health Check Configuration

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/resources/valid_health_check \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "health_check": {
        "interval": "10s",
        "timeout": "5s",
        "retries": 3
      }
    }
  }'
```

**Expected Result**: `{"result": true}` (Valid health check config)

### Example 8: Check Volume Size Limits

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/resources/valid_volume_size \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "volume_name": "db_data",
      "size_gb": 50
    }
  }'
```

**Expected Result**: `{"result": true}` (Within database volume limits)

## 🔒 Security Policy Examples

### Example 9: Check Container Security Compliance

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/security/secure_container \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "container_name": "back-end",
      "privileged": false,
      "user": "node",
      "read_only_rootfs": true,
      "security_opts": [
        "no-new-privileges:true",
        "seccomp:default"
      ]
    }
  }'
```

**Expected Result**: `{"result": true}` (Secure configuration)

### Example 10: Validate Image Source

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/security/valid_image \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "image": "postgres:15",
      "environment": "development"
    }
  }'
```

**Expected Result**: `{"result": true}` (Trusted image)

### Example 11: Check for Insecure Environment Variables

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/security/insecure_env_vars \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "env_vars": [
        {"name": "DATABASE_PASSWORD", "value": "password123"},
        {"name": "API_KEY", "value": "secret_key"}
      ]
    }
  }'
```

**Expected Result**: List of insecure environment variables

### Example 12: Calculate Security Score

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/security/security_score \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "container_name": "database",
      "privileged": false,
      "user": "postgres",
      "security_opts": ["no-new-privileges:true", "seccomp:default"],
      "capabilities": ["CHOWN", "FOWNER"],
      "image": "postgres:15",
      "env_vars": [],
      "network_config": {"name": "usualstore_network", "internal": true},
      "exposed_ports": [5433],
      "volumes": [{"host_path": "/data/postgres"}],
      "health_check": {"test": ["CMD", "pg_isready"], "user": "postgres"}
    }
  }'
```

**Expected Result**: Security score (0-100)

## 👤 Access Control Examples

### Example 13: Check User Permission

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/user_has_permission \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": {"role": "admin", "email": "admin@example.com"},
      "permission": "manage_users"
    }
  }'
```

**Expected Result**: `{"result": true}` (Admin has manage_users permission)

### Example 14: Check API Endpoint Access

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_api_access \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": {"role": "supporter", "email": "support@example.com"},
      "endpoint": "/api/support/tickets",
      "method": "GET"
    }
  }'
```

**Expected Result**: `{"result": true}` (Supporter can view support tickets)

### Example 15: Validate Container Management Action

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_container_action \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": {"role": "admin", "email": "admin@example.com"},
      "action": "restart",
      "container": "back-end"
    }
  }'
```

**Expected Result**: `{"result": true}` (Admin can restart containers)

### Example 16: Check Kafka Topic Access

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_kafka_publish \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "service": "back-end",
      "topic": "email-notifications"
    }
  }'
```

**Expected Result**: `{"result": true}` (Backend can publish to email notifications)

### Example 17: Validate User Management Action

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/allow_user_management \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": {"role": "admin", "email": "admin@example.com"},
      "target_user": {"role": "supporter"},
      "action": "create"
    }
  }'
```

**Expected Result**: `{"result": true}` (Admin can create supporter accounts)

### Example 18: Check Super Admin Deletion Protection

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/deny_user_deletion \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "action": "delete",
      "target_user": {"role": "super_admin"},
      "all_users": [
        {"role": "super_admin"},
        {"role": "admin"},
        {"role": "user"}
      ]
    }
  }'
```

**Expected Result**: `{"result": true}` (Cannot delete last super_admin)

### Example 19: Check Rate Limiting

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/within_rate_limit \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "user": {"role": "user"},
      "request_count": 50,
      "time_window_minutes": 1
    }
  }'
```

**Expected Result**: `{"result": true}` (Within rate limit for user role)

### Example 20: Validate Session

```bash
curl -X POST http://localhost:8181/v1/data/usualstore/access/valid_session \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "session": {
        "user_id": 123,
        "expires_at": 1735200000,
        "revoked": false
      },
      "user": {"id": 123},
      "current_time": 1735100000
    }
  }'
```

**Expected Result**: `{"result": true}` (Valid session)

## 🔨 Custom Policy Development

### Creating a New Policy

1. **Create Policy File**:

```bash
cd terraform/modules/policies/policies
touch custom_policy.rego
```

2. **Write Policy**:

```rego
# custom_policy.rego
package usualstore.custom

import future.keywords.if

# Your custom policy logic
allow_custom_action if {
    input.user.role == "admin"
    input.action == "custom_action"
}
```

3. **Test Policy**:

```bash
# Install OPA CLI
brew install opa

# Test policy
opa eval -d custom_policy.rego -i input.json 'data.usualstore.custom.allow_custom_action'
```

4. **Load into OPA Server**:

The policy will be automatically loaded when you run `terraform apply`.

### Testing Policies Locally

Create `input.json`:

```json
{
  "user": {
    "role": "admin",
    "email": "admin@example.com"
  },
  "action": "custom_action"
}
```

Test with OPA CLI:

```bash
opa eval -d terraform/modules/policies/policies/ -i input.json 'data.usualstore.custom.allow_custom_action'
```

### Policy Testing Best Practices

1. **Unit Tests**: Create test files for each policy
2. **Integration Tests**: Test policy combinations
3. **Edge Cases**: Test boundary conditions
4. **Performance**: Test with large input sets

Example test file (`network_test.rego`):

```rego
package usualstore.network_test

import future.keywords.if

test_frontend_to_backend_allowed if {
    allow with input as {
        "source_service": "react-frontend",
        "target_service": "back-end"
    }
}

test_frontend_to_database_denied if {
    deny with input as {
        "source_service": "react-frontend",
        "target_service": "database"
    }
}
```

Run tests:

```bash
opa test terraform/modules/policies/policies/
```

## 📊 Policy Monitoring

### View All Policy Decisions

```bash
# Get all policies
curl http://localhost:8181/v1/policies

# Get specific policy
curl http://localhost:8181/v1/policies/network.rego
```

### Audit Logging

The policy enforcer logs all decisions:

```bash
docker logs usualstore-policy-enforcer | grep "violation"
```

### Policy Performance Metrics

```bash
curl http://localhost:8181/metrics
```

## 🎯 Real-World Scenarios

### Scenario 1: New Service Deployment

When deploying a new service, validate it against all policies:

```bash
./validate-service.sh new-service
```

### Scenario 2: Security Audit

Run a comprehensive security audit:

```bash
./audit-security.sh > security-report.json
```

### Scenario 3: Compliance Check

Check if all containers are compliant:

```bash
curl http://localhost:8080/audit
```

### Scenario 4: Resource Optimization

Get optimization recommendations:

```bash
./optimize-resources.sh
```

## 📚 Additional Resources

- [OPA Playground](https://play.openpolicyagent.org/) - Interactive policy testing
- [Rego Cheat Sheet](https://www.openpolicyagent.org/docs/latest/policy-cheatsheet/)
- [Policy Library](https://github.com/open-policy-agent/library) - Community policies

