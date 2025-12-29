# 🏗️ Multi-Tenant Architecture

**NO master database. Complete separation. Simple isolated databases per customer.**

---

## ❌ What You DON'T Have

**NO master "usualstore" database tracking tenants!**

This architecture does NOT use:
- ❌ Master metadata database
- ❌ Tenants table
- ❌ Complex tenant routing
- ❌ Shared infrastructure

## ✅ What You DO Have

**Complete isolation. Customer-specific databases.**

Each customer gets:
- ✅ Own database (name they choose!)
- ✅ Own users with permissions
- ✅ Complete schema automatically
- ✅ Total isolation from others

---

## 🎯 Architecture Overview

```
┌────────────────────────────────────────────────────────────────┐
│                    TERRAFORM (Infrastructure)                   │
│  Creates databases, users, permissions, schema                 │
│  Runs: When adding new customers                               │
└────────────────────────────────────────────────────────────────┘
                              ↓
                    Creates & Configures
                              ↓
┌────────────────────────────────────────────────────────────────┐
│                         POSTGRESQL                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │building_shop │  │hardware_store│  │  [more...]   │        │
│  │              │  │              │  │              │        │
│  │ - customers  │  │ - customers  │  │ - customers  │        │
│  │ - orders     │  │ - orders     │  │ - orders     │        │
│  │ - widgets    │  │ - widgets    │  │ - widgets    │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
└────────────────────────────────────────────────────────────────┘
                              ↓
                    App Connects To
                              ↓
┌────────────────────────────────────────────────────────────────┐
│                   GO APPLICATION (Business Logic)               │
│  Uses existing databases                                        │
│  Runs: Continuously serving requests                           │
└────────────────────────────────────────────────────────────────┘
```

---

## 📋 Responsibilities

### Terraform's Job

✅ **Create tenant databases with custom names**
   - Customer chooses: `building_shop`, `hardware_store`, etc.
   - Terraform creates PostgreSQL database

✅ **Set up users and permissions**
   - Admin users: Full access (read/write/delete)
   - Developer users: Read/write access (no delete)
   - Customer users: Read-only access

✅ **Apply database schema**
   - Creates all tables (customers, orders, widgets, etc.)
   - Sets up indexes for performance
   - Applies triggers for updated_at columns

### Go Application's Job

✅ **Connect to tenant databases**
   - Reads database connection from environment variable
   - Uses customer's chosen database name
   - Example: `DATABASE_DSN="...dbname=building_shop..."`

✅ **Route requests to correct database**
   - Single application instance
   - Connects to different databases based on configuration
   - No complex routing needed

✅ **Handle business logic**
   - Process orders
   - Manage inventory
   - Handle payments
   - Execute queries on tenant data

---

## 🔄 Complete Workflow

### Step 1: New Customer Signs Up

```bash
# DevOps/Admin edits Terraform configuration
cd terraform-tenants
vim tenants.tfvars

# Add new tenant
tenants = {
  building_shop = {
    database_name = "building_shop"      # ← Customer chooses name
    tenant_name   = "Building Shop Inc."
    plan          = "professional"
    
    admins = [
      { username = "building_admin", password = "...", email = "..." }
    ]
    developers = [
      { username = "building_dev", password = "...", email = "..." }
    ]
    customers = [
      { username = "building_customer", password = "...", email = "..." }
    ]
  }
}

# Apply infrastructure
terraform apply -var-file="tenants.tfvars"
```

**Terraform creates:**
- ✅ Database `building_shop`
- ✅ Users: `building_admin`, `building_dev`, `building_customer`
- ✅ Permissions per role
- ✅ Complete schema with all tables

### Step 2: Configure Application

```bash
# Set environment variable to connect to customer's database
export DATABASE_DSN="host=localhost port=5432 user=building_admin password=SecurePass123! dbname=building_shop sslmode=disable"

# Start application
cd cmd/api
go run *.go
```

**Application now:**
- ✅ Connects to `building_shop` database
- ✅ Serves Building Shop customer's requests
- ✅ All data isolated in their database

### Step 3: Customer Uses Application

```bash
# Customer makes request
curl -X POST http://localhost:4001/api/orders \
  -H "Authorization: Bearer TOKEN" \
  -d '{"product_id": 1, "quantity": 2}'

# Application:
# 1. Receives request
# 2. Connects to building_shop database (from DATABASE_DSN)
# 3. Executes: INSERT INTO orders (...)
# 4. Returns response

# Data is stored in building_shop database only!
```

---

## 🔐 Access Control

### Terraform Creates These Users Per Tenant

| User Type | Username Example | Permissions | Use Case |
|-----------|------------------|-------------|----------|
| **Admin** | `building_admin` | Full: SELECT, INSERT, UPDATE, DELETE | Tenant owner, full access |
| **Developer** | `building_dev` | Limited: SELECT, INSERT, UPDATE | Dev team, no delete |
| **Customer** | `building_customer` | Read-only: SELECT | Reports, analytics |

### Application Connects As

Choose which user based on context:

```bash
# For admin operations
DATABASE_DSN="...user=building_admin...dbname=building_shop..."

# For developer operations
DATABASE_DSN="...user=building_dev...dbname=building_shop..."

# For customer read-only
DATABASE_DSN="...user=building_customer...dbname=building_shop..."
```

---

## 📂 File Structure

```
project/
├── terraform-tenants/              # ← Infrastructure management
│   ├── main.tf                     # Terraform config
│   ├── tenants.tfvars             # Customer definitions
│   └── modules/tenant/            # Reusable module
│       ├── main.tf                # Creates database + users
│       └── database-schema.sql    # Schema to apply
│
├── cmd/api/                        # ← Application code
│   ├── api.go                      # Uses DATABASE_DSN
│   └── handlers-api.go             # Business logic
│
├── internal/
│   ├── driver/driver.go           # Database connection
│   └── models/models.go           # Data models
│
└── docker-compose.yml             # Development environment
```

---

## 🚀 Development Workflow

### Local Development

```bash
# 1. Create test tenant database
cd terraform-tenants
terraform apply -var-file="tenants.tfvars"

# 2. Set environment to use test database
export DATABASE_DSN="host=localhost port=5432 user=test_admin password=test dbname=test_tenant sslmode=disable"

# 3. Run application
cd ../cmd/api
go run *.go

# 4. Test
curl http://localhost:4001/api/orders
```

### Production Deployment

```bash
# 1. Add production tenant via Terraform
cd terraform-tenants
vim tenants.tfvars  # Add production tenant
terraform apply -var-file="tenants.tfvars"

# 2. Configure application with production database
# In Docker, Kubernetes, etc:
env:
  - name: DATABASE_DSN
    value: "host=prod-db port=5432 user=prod_admin password=xxx dbname=customer_db"

# 3. Deploy application
docker-compose up -d
# or
kubectl apply -f deployment.yaml
```

---

## 💡 Key Benefits

### Separation of Concerns

✅ **Infrastructure (Terraform)**
   - Declarative configuration
   - Version controlled
   - Preview changes before applying
   - Managed by DevOps/Platform team

✅ **Application (Go)**
   - Business logic only
   - No database creation code
   - Simpler codebase
   - Managed by Development team

### Security

✅ **Least Privilege**
   - Application doesn't need permissions to create databases
   - Different users for different access levels
   - Credentials managed separately

✅ **Isolation**
   - Each customer has separate physical database
   - No shared tables
   - No risk of data leakage

### Scalability

✅ **Unlimited Tenants**
   - Add new tenant: Edit Terraform config + apply
   - No application code changes needed
   - Each tenant scales independently

✅ **Easy Management**
   - Terraform state tracks all infrastructure
   - Easy to audit who has what access
   - Simple to backup/restore per tenant

---

## 📝 Common Operations

### Add New Tenant

```bash
# 1. Edit tenants.tfvars
vim terraform-tenants/tenants.tfvars

# 2. Apply
cd terraform-tenants
terraform apply -var-file="tenants.tfvars"

# Done! Database created with all permissions
```

### Connect Application to Different Tenant

```bash
# Just change environment variable
export DATABASE_DSN="host=localhost port=5432 user=another_admin password=xxx dbname=another_tenant"

# Restart application
cd cmd/api
go run *.go
```

### Add User to Existing Tenant

```bash
# 1. Edit tenant config
vim terraform-tenants/tenants.tfvars

# Add to developers list
developers = [
  { username = "dev1", password = "...", email = "..." },
  { username = "dev2", password = "...", email = "..." }  # ← New
]

# 2. Apply
terraform apply -var-file="tenants.tfvars"

# Done! New user created with correct permissions
```

### Remove Tenant

```bash
# 1. Remove from tenants.tfvars
vim terraform-tenants/tenants.tfvars
# Delete tenant block

# 2. Apply (will destroy database!)
terraform apply -var-file="tenants.tfvars"

# ⚠️ WARNING: This permanently deletes the database!
```

---

## 🔍 Example Scenarios

### Scenario 1: Single Application, Multiple Customers

```
One Go Application Deployment
├── Configure with: DATABASE_DSN=...dbname=building_shop
├── Serves: Building Shop customers
└── All data in building_shop database

Another Go Application Deployment
├── Configure with: DATABASE_DSN=...dbname=hardware_store
├── Serves: Hardware Store customers
└── All data in hardware_store database
```

### Scenario 2: Multi-Tenant SaaS

```
Load Balancer
    ↓
Application (looks at subdomain/header)
    ├── building.usualstore.com → DATABASE_DSN=...dbname=building_shop
    ├── hardware.usualstore.com → DATABASE_DSN=...dbname=hardware_store
    └── [more subdomains] → [more databases]
```

---

## 🎯 Summary

| Concern | Handled By | How |
|---------|------------|-----|
| Create database | **Terraform** | `terraform apply` |
| Create users | **Terraform** | Defined in `tenants.tfvars` |
| Set permissions | **Terraform** | Automatic per role |
| Apply schema | **Terraform** | `database-schema.sql` |
| Connect to DB | **Go App** | `DATABASE_DSN` env var |
| Execute queries | **Go App** | Standard database operations |
| Handle business logic | **Go App** | Application code |

---

## 🚀 Reusable Application Code

### Your "usual_store" is NOT Hardcoded!

The application code is **completely reusable** for any customer. The database name is configured via environment variables:

```
Same Application Code:
  ├── Building Shop    → DATABASE_DSN="...dbname=building_shop..."
  ├── Hardware Store   → DATABASE_DSN="...dbname=hardware_store..."
  └── BookStore        → DATABASE_DSN="...dbname=bookstore_2025..."

✅ ONE codebase → INFINITE customers
```

### Deployment Pattern

```bash
# Customer 1: Building Shop
docker run \
  -e DATABASE_DSN="...dbname=building_shop..." \
  -p 4001:4001 \
  usualstore/api:latest

# Customer 2: Hardware Store (SAME IMAGE!)
docker run \
  -e DATABASE_DSN="...dbname=hardware_store..." \
  -p 4002:4002 \
  usualstore/api:latest
```

### Frontend Configuration

Frontend also uses environment variables:

```bash
# Building Shop frontend
docker run \
  -e REACT_APP_API_URL="https://api.buildingshop.com" \
  -e REACT_APP_TENANT_NAME="Building Shop" \
  -p 3000:80 \
  usualstore/frontend:latest

# Hardware Store frontend (SAME IMAGE!)
docker run \
  -e REACT_APP_API_URL="https://api.hardwarestore.com" \
  -e REACT_APP_TENANT_NAME="Hardware Store" \
  -p 3001:80 \
  usualstore/frontend:latest
```

---

## 📚 Documentation

- **Simplified Architecture** ⭐: `terraform-tenants/SIMPLIFIED-ARCHITECTURE.md` - **START HERE**
- **Terraform Setup**: `terraform-tenants/README.md`
- **Interactive Script**: `terraform-tenants/add-customer.sh`
- **Quick Start**: `terraform-tenants/QUICK-START.md`
- **Examples**: `terraform-tenants/examples/`

---

## 🎯 Key Principle

**NO master database. Each customer = isolated database with custom name.**

```
Customer chooses DB name → Terraform creates DB → App connects via env var
```

**Clean Architecture** = Infrastructure as Code (Terraform) + Business Logic (Go App) 🎯

