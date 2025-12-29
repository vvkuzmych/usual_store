# 🏢 Multi-Tenant Infrastructure - Terraform

**Create isolated databases for each customer with one command.**

---

## 🚀 Deployment Options

| Environment | Guide | Best For |
|-------------|-------|----------|
| **Local/Docker** | This README | Development, testing |
| **Kubernetes** | `../terraform-k8s/README.md` | On-premise, self-hosted |
| **AWS Lambda** | `AWS-LAMBDA-GUIDE.md` ⭐ | Cloud, serverless, auto-scaling |

---

## 🎯 Architecture

**NO master database. Complete separation between customers.**

```
Customer 1         Customer 2         Customer 3
    ↓                  ↓                  ↓
building_shop DB   hardware_store DB   bookstore DB
    ↑                  ↑                  ↑
   App 1              App 2              App 3
(Same Code!)       (Same Code!)       (Same Code!)
```

See [SIMPLIFIED-ARCHITECTURE.md](./SIMPLIFIED-ARCHITECTURE.md) for details.

---

## 🚀 Quick Start

### Add Your First Customer

```bash
cd terraform-tenants

# Interactive script
./add-customer.sh
```

**You'll be prompted for:**
- Customer name (e.g., "Building Shop Inc")
- Database name (e.g., "building_shop")
- Admin user credentials
- Developer user credentials (optional)
- Customer user credentials (optional)

**Script generates:**
1. Terraform configuration: `customers/{name}.tfvars`
2. Deployment script: `customers/{name}-deploy.sh`
3. Docker Compose file: `customers/{name}-docker-compose.yml`

### Deploy Customer Database

```bash
# Option 1: Run generated deployment script (easiest)
./customers/building_shop-deploy.sh

# Option 2: Run Terraform manually
terraform init
terraform apply -var-file="customers/building_shop.tfvars"
```

**Terraform creates:**
- ✅ PostgreSQL database with custom name
- ✅ Admin, developer, and customer users
- ✅ Complete schema (tables, indexes, triggers)
- ✅ Proper permissions per role

### Start Application for Customer

```bash
# Use generated Docker Compose file
docker-compose -f customers/building_shop-docker-compose.yml up -d

# Or configure manually
export DATABASE_DSN="host=localhost port=5432 user=building_admin password=*** dbname=building_shop"
cd ../cmd/api
go run *.go
```

---

## 📋 Complete Example

### Scenario: Onboard 3 Customers

#### Customer 1: Building Shop

```bash
./add-customer.sh
# Enter: Building Shop Inc, building_shop, admin credentials

./customers/building_shop-deploy.sh
docker-compose -f customers/building_shop-docker-compose.yml up -d

# ✅ Live at http://localhost:3000
```

#### Customer 2: Hardware Store

```bash
./add-customer.sh
# Enter: Hardware Store LLC, hardware_store, admin credentials

./customers/hardware_store-deploy.sh
docker-compose -f customers/hardware_store-docker-compose.yml up -d

# ✅ Live at http://localhost:3001
```

#### Customer 3: BookStore

```bash
./add-customer.sh
# Enter: BookStore Co, bookstore_2025, admin credentials

./customers/bookstore_2025-deploy.sh
docker-compose -f customers/bookstore_2025-docker-compose.yml up -d

# ✅ Live at http://localhost:3002
```

**Result:** 3 isolated databases, 3 application instances, same code!

---

## 📂 File Structure

```
terraform-tenants/
├── add-customer.sh              ← Interactive script to add customers
├── main.tf                      ← Terraform main configuration
├── variables.tf                 ← Input variables
├── outputs.tf                   ← Output values
│
├── modules/
│   └── tenant/                  ← Reusable module per customer
│       ├── main.tf              ← Creates DB, users, permissions
│       ├── variables.tf
│       ├── outputs.tf
│       └── database-schema.sql  ← Schema applied to each database
│
├── customers/                   ← Generated per customer (gitignored)
│   ├── building_shop.tfvars
│   ├── building_shop-deploy.sh
│   ├── building_shop-docker-compose.yml
│   │
│   ├── hardware_store.tfvars
│   ├── hardware_store-deploy.sh
│   └── hardware_store-docker-compose.yml
│
├── examples/
│   ├── single-tenant.tfvars     ← Example: One customer
│   └── multiple-tenants.tfvars  ← Example: Multiple customers
│
├── SIMPLIFIED-ARCHITECTURE.md   ← Architecture explanation
├── QUICK-START.md               ← Manual setup guide
└── README.md                    ← This file
```

---

## 🔧 What Gets Created

### Per Customer Database

When you run `./add-customer.sh` and deploy, Terraform creates:

#### 1. Database
```sql
CREATE DATABASE building_shop;
```

#### 2. Admin Role & User
```sql
CREATE ROLE building_shop_admin;
GRANT ALL PRIVILEGES ON DATABASE building_shop TO building_shop_admin;

CREATE USER building_admin WITH PASSWORD '***' IN ROLE building_shop_admin;
```

#### 3. Developer Role & User (Optional)
```sql
CREATE ROLE building_shop_developer;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES TO building_shop_developer;

CREATE USER building_dev WITH PASSWORD '***' IN ROLE building_shop_developer;
```

#### 4. Customer Role & User (Optional)
```sql
CREATE ROLE building_shop_customer;
GRANT SELECT ON ALL TABLES TO building_shop_customer;

CREATE USER building_customer WITH PASSWORD '***' IN ROLE building_shop_customer;
```

#### 5. Complete Schema
```sql
CREATE TABLE customers (...);
CREATE TABLE orders (...);
CREATE TABLE widgets (...);
CREATE TABLE transactions (...);
CREATE TABLE statuses (...);
-- + indexes, triggers, constraints
```

---

## 🎯 Access Levels

Each customer gets 3 types of users:

| Role | Permissions | Use Case |
|------|-------------|----------|
| **Admin** | Full: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP | Tenant owner, full control |
| **Developer** | Limited: SELECT, INSERT, UPDATE | Dev team, no delete or schema changes |
| **Customer** | Read-only: SELECT | Analytics, reporting, read-only access |

---

## 🐳 Docker Deployment

### Generated Docker Compose

The `add-customer.sh` script generates a complete Docker Compose file:

```yaml
# customers/building_shop-docker-compose.yml
services:
  backend:
    image: usualstore/api:latest
    environment:
      DATABASE_DSN: "host=db port=5432 user=building_admin password=*** dbname=building_shop"
      API_PORT: "4001"
      TENANT_NAME: "building_shop"
    ports:
      - "4001:4001"

  frontend:
    image: usualstore/frontend:latest
    environment:
      REACT_APP_API_URL: "http://backend:4001"
      REACT_APP_TENANT_NAME: "Building Shop Inc"
    ports:
      - "3000:80"

  db:
    image: postgres:14
    environment:
      POSTGRES_DB: building_shop
    volumes:
      - building_shop_data:/var/lib/postgresql/data
```

**Start:**
```bash
docker-compose -f customers/building_shop-docker-compose.yml up -d
```

---

## ☸️ Kubernetes Deployment

For Kubernetes, modify the generated Docker Compose to K8s manifests:

```yaml
# k8s/building-shop/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: building-shop-backend
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: backend
        image: usualstore/api:latest
        env:
        - name: DATABASE_DSN
          value: "host=postgres port=5432 user=building_admin password=*** dbname=building_shop"
```

---

## 🔐 Security Best Practices

### 1. Sensitive Data

**Don't commit passwords!**

```bash
# Add to .gitignore
echo "customers/" >> .gitignore
```

The `customers/` directory contains passwords and should NOT be committed.

### 2. Use Secrets Management

For production:

```bash
# Use Terraform variables from environment
export TF_VAR_db_password="$(cat /secure/db-password)"
terraform apply -var-file="customers/building_shop.tfvars"
```

### 3. Separate Postgres Instances

For production, consider separate PostgreSQL instances per customer:

```hcl
db_host = "building-shop.postgres.rds.amazonaws.com"  # Customer 1
db_host = "hardware-store.postgres.rds.amazonaws.com"  # Customer 2
```

---

## 📊 Management Operations

### List All Customer Databases

```bash
psql -U postgres -c "\l" | grep -E "building_shop|hardware_store|bookstore"
```

### Connect to Customer Database

```bash
psql -h localhost -p 5432 -U building_admin -d building_shop
```

### Check Database Size

```sql
SELECT 
  datname as database,
  pg_size_pretty(pg_database_size(datname)) as size
FROM pg_database
WHERE datname IN ('building_shop', 'hardware_store', 'bookstore_2025');
```

### Backup Customer Database

```bash
pg_dump -U building_admin -d building_shop > building_shop_backup.sql
```

### Remove Customer

```bash
# 1. Stop application
docker-compose -f customers/building_shop-docker-compose.yml down -v

# 2. Remove from Terraform
rm customers/building_shop*

# 3. Drop database manually (Terraform won't do this automatically)
psql -U postgres -c "DROP DATABASE building_shop;"
```

---

## 🚀 Scaling

### Same Server, Multiple Databases

**Current setup** - Good for:
- Development
- Small to medium deployments
- Cost optimization

```
PostgreSQL Server (one instance)
  ├── building_shop DB
  ├── hardware_store DB
  └── bookstore_2025 DB
```

### Separate Servers Per Customer

**Production recommendation** - Best for:
- Large customers
- Compliance requirements
- Performance isolation

```
Customer 1 → RDS Instance 1 → building_shop DB
Customer 2 → RDS Instance 2 → hardware_store DB
Customer 3 → RDS Instance 3 → bookstore_2025 DB
```

Update `customers/{name}.tfvars`:
```hcl
db_host = "customer1.us-east-1.rds.amazonaws.com"
```

---

## 📖 Documentation

- **[SIMPLIFIED-ARCHITECTURE.md](./SIMPLIFIED-ARCHITECTURE.md)** - Complete architecture explanation
- **[QUICK-START.md](./QUICK-START.md)** - Manual setup guide
- **[examples/](./examples/)** - Example configurations

---

## 🎯 Summary

### To Add a New Customer:

```bash
./add-customer.sh                                    # Interactive prompts
./customers/{name}-deploy.sh                         # Deploy database
docker-compose -f customers/{name}-docker-compose.yml up -d  # Start app
```

### What You Get:

✅ Isolated database per customer  
✅ Custom database name (customer chooses!)  
✅ Admin, developer, customer users  
✅ Complete schema automatically applied  
✅ Docker Compose file ready to use  
✅ Same application code for all customers  
✅ NO master database tracking metadata  

---

**Simple. Isolated. Scalable.** 🚀
