# 🏗️ Architecture Summary

**Complete separation. No master database. Simple multi-tenant setup.**

---

## 🎯 Current Architecture

### Multi-Tenant System

**NO master database tracking tenants. Each customer = isolated database.**

```
Customer 1 → building_shop DB → App Instance 1
Customer 2 → hardware_store DB → App Instance 2  
Customer 3 → bookstore_2025 DB → App Instance 3

✅ Same application code
✅ Different databases
✅ Complete isolation
```

---

## 📂 Key Documentation

### Multi-Tenancy (Terraform-Based)

| Document | Description |
|----------|-------------|
| [terraform-tenants/SIMPLIFIED-ARCHITECTURE.md](terraform-tenants/SIMPLIFIED-ARCHITECTURE.md) | **START HERE** - Complete architecture explanation |
| [terraform-tenants/README.md](terraform-tenants/README.md) | Usage guide with `add-customer.sh` script |
| [terraform-tenants/add-customer.sh](terraform-tenants/add-customer.sh) | Interactive script to add new customers |
| [terraform-tenants/MONGODB-GUIDE.md](terraform-tenants/MONGODB-GUIDE.md) | **MongoDB option** - Use MongoDB instead of PostgreSQL 🍃 |
| [MULTI-TENANT-ARCHITECTURE.md](MULTI-TENANT-ARCHITECTURE.md) | Responsibilities: Terraform vs Go App |

### Infrastructure

| Document | Description |
|----------|-------------|
| [terraform/README.md](terraform/README.md) | Docker deployment with Terraform |
| [terraform-k8s/README.md](terraform-k8s/README.md) | Kubernetes deployment with Terraform |
| [terraform-tenants/AWS-LAMBDA-GUIDE.md](terraform-tenants/AWS-LAMBDA-GUIDE.md) | **AWS Lambda + RDS deployment** ⭐ |
| [terraform-tenants/MONGODB-GUIDE.md](terraform-tenants/MONGODB-GUIDE.md) | **MongoDB option** (local & AWS) 🍃 |
| [docs/guides/DEVELOPMENT-WORKFLOW.md](docs/guides/DEVELOPMENT-WORKFLOW.md) | Update code after changes |

---

## 🚀 Quick Start

### Add Your First Customer

```bash
# 1. Run interactive script
cd terraform-tenants
./add-customer.sh

# Fill in prompts:
# - Company: Building Shop Inc
# - Database: building_shop
# - Admin credentials

# 2. Deploy database
./customers/building_shop-deploy.sh

# 3. Start application
docker-compose -f customers/building_shop-docker-compose.yml up -d

# ✅ Customer is live!
```

### Add More Customers

```bash
# Same script, different customer
./add-customer.sh
# Fill in: Hardware Store, hardware_store...

./customers/hardware_store-deploy.sh
docker-compose -f customers/hardware_store-docker-compose.yml up -d

# ✅ Another customer with complete isolation!
```

---

## 🔑 Key Points

### ❌ What You DON'T Have

- NO master "usualstore" database
- NO tenants metadata table
- NO complex tenant routing in Go code

### ✅ What You DO Have

- Isolated database per customer
- Customer chooses database name
- Terraform creates database + users + schema
- Same application code for all
- Environment variables configure which DB
- Interactive script to add customers

---

## 🎯 Responsibilities

### Terraform's Job

```hcl
# Customer fills this in Terraform config
tenants = {
  building_shop = {
    database_name = "building_shop"  # ← Customer chooses!
    admins = [...]
    developers = [...]
  }
}
```

**Terraform creates:**
1. ✅ Database with custom name
2. ✅ Admin/developer/customer users
3. ✅ Complete schema (tables, indexes, triggers)
4. ✅ Proper permissions per role

### Go Application's Job

```bash
# Environment variable tells app which database
export DATABASE_DSN="host=localhost port=5432 user=building_admin password=*** dbname=building_shop"

# Start application
cd cmd/api
go run *.go
```

**Application does:**
1. ✅ Connects to specified database
2. ✅ Handles business logic
3. ✅ Executes queries
4. ✅ Serves customer requests

---

## 📊 Example Setup

### Three Customers Running

```
PostgreSQL Server
├── building_shop DB     → App Instance 1 (port 4001)
├── hardware_store DB    → App Instance 2 (port 4002)
└── bookstore_2025 DB    → App Instance 3 (port 4003)

All using: usualstore/api:latest (SAME IMAGE!)
```

**Each has:**
- Own database with custom name
- Own application instance
- Own frontend
- Complete isolation

---

## 🏢 File Structure

```
usual_store/
├── terraform-tenants/              ← Multi-tenant infrastructure
│   ├── add-customer.sh             ← Interactive script ✨
│   ├── main.tf                     ← Terraform config
│   ├── modules/tenant/             ← Reusable module
│   │   └── database-schema.sql     ← Schema applied to each DB
│   ├── customers/                  ← Generated configs (gitignored)
│   │   ├── building_shop.tfvars
│   │   ├── building_shop-deploy.sh
│   │   └── building_shop-docker-compose.yml
│   ├── SIMPLIFIED-ARCHITECTURE.md  ← Architecture guide
│   └── README.md                   ← Usage guide
│
├── terraform/                      ← Docker deployment
│   └── main.tf                     ← Deploy with Docker
│
├── terraform-k8s/                  ← Kubernetes deployment
│   └── main.tf                     ← Deploy with K8s
│
├── cmd/api/                        ← Go application (reusable!)
├── internal/                       ← Business logic
├── react-frontend/                 ← Frontend
└── typescript-frontend/            ← TypeScript frontend
```

---

## 🔧 Common Operations

### Add Customer

```bash
cd terraform-tenants
./add-customer.sh
```

### Deploy Customer Database

```bash
./customers/{name}-deploy.sh
```

### Start Application for Customer

```bash
docker-compose -f customers/{name}-docker-compose.yml up -d
```

### Check Running Customers

```bash
psql -U postgres -c "\l" | grep -E "building_shop|hardware_store|bookstore"
```

### Backup Customer Database

```bash
pg_dump -U building_admin -d building_shop > backup.sql
```

---

## 🎯 This Is What You Have

1. **Simple Architecture** - No master database complexity
2. **Customer Choice** - They choose database name
3. **Easy Onboarding** - Interactive script handles everything
4. **Complete Isolation** - Each customer separate
5. **Reusable Code** - Same app for all customers
6. **Scalable** - Add unlimited customers

---

**Documentation**: Start with [terraform-tenants/SIMPLIFIED-ARCHITECTURE.md](terraform-tenants/SIMPLIFIED-ARCHITECTURE.md)

**Add Customer**: Run `./terraform-tenants/add-customer.sh`

**Questions**: Check [terraform-tenants/README.md](terraform-tenants/README.md)

