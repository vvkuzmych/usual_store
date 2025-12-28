# 🏗️ Simplified Architecture - NO Master Database

**Complete separation. No metadata tracking. Just isolated customer databases.**

---

## 🎯 What You Have

```
┌─────────────────────────────────────────────────────────────┐
│                         TERRAFORM                            │
│  Creates separate, isolated databases for each customer     │
└─────────────────────────────────────────────────────────────┘
                            ↓
            Creates Independent Databases
                            ↓
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ building_shop│  │hardware_store│  │ bookstore_db │
├──────────────┤  ├──────────────┤  ├──────────────┤
│ • customers  │  │ • customers  │  │ • customers  │
│ • orders     │  │ • orders     │  │ • orders     │
│ • widgets    │  │ • widgets    │  │ • widgets    │
│ • statuses   │  │ • statuses   │  │ • statuses   │
│              │  │              │  │              │
│ ISOLATED     │  │ ISOLATED     │  │ ISOLATED     │
└──────────────┘  └──────────────┘  └──────────────┘
       ↑                 ↑                 ↑
       │                 │                 │
       │                 │                 │
┌──────┴─────┐  ┌────────┴────┐  ┌────────┴────┐
│   App 1    │  │    App 2    │  │    App 3    │
│ (Building) │  │ (Hardware)  │  │ (BookStore) │
└────────────┘  └─────────────┘  └─────────────┘

SAME APPLICATION CODE - Different DATABASE_DSN!
```

---

## ❌ What You DON'T Have

**NO master database** like this:

```
❌ usualstore (master metadata database)
   ├── tenants table
   ├── tenant_users table
   └── mapping to tenant databases
```

**You DON'T need this!** Each customer is completely separate.

---

## ✅ How It Works

### 1. New Customer Signs Up

You fill in customer information and run a script:

```bash
cd terraform-tenants
./add-customer.sh
```

**Interactive prompts:**
```
Customer Company Name: Building Shop Inc
Database Name: building_shop
Customer Key: building_shop
Plan: professional

Admin Username: building_admin
Admin Password: ••••••••••
Admin Email: admin@buildingshop.com

Developer Username: building_dev
Developer Password: ••••••••••
Developer Email: dev@buildingshop.com
```

### 2. Script Generates Terraform Config

**Creates: `customers/building_shop.tfvars`**

```hcl
# Customer Configuration: Building Shop Inc

db_host     = "localhost"
db_port     = "5432"
db_username = "postgres"
db_password = "yourpassword"

tenants = {
  building_shop = {
    database_name = "building_shop"
    tenant_name   = "Building Shop Inc"
    plan          = "professional"
    
    admins = [
      {
        username = "building_admin"
        password = "SecurePass123!"
        email    = "admin@buildingshop.com"
      }
    ]
    
    developers = [
      {
        username = "building_dev"
        password = "DevPass123!"
        email    = "dev@buildingshop.com"
      }
    ]
    
    customers = []
  }
}
```

### 3. Run Terraform to Create Database

**Automatically or manually:**

```bash
# Option 1: Run generated deployment script
./customers/building_shop-deploy.sh

# Option 2: Run Terraform manually
terraform apply -var-file="customers/building_shop.tfvars"
```

**Terraform creates:**
- ✅ Database: `building_shop`
- ✅ Admin user: `building_admin` (full access)
- ✅ Developer user: `building_dev` (read/write)
- ✅ All tables: `customers`, `orders`, `widgets`, `statuses`, etc.
- ✅ Indexes, triggers, constraints

### 4. Deploy Application

**Generated Docker Compose file: `customers/building_shop-docker-compose.yml`**

```yaml
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
```

**Start:**
```bash
docker-compose -f customers/building_shop-docker-compose.yml up -d
```

---

## 🎯 Complete Example

### Add Building Shop (Customer 1)

```bash
# 1. Run script
./add-customer.sh

# Fill in:
# - Company: Building Shop Inc
# - Database: building_shop
# - Admin: building_admin / SecurePass123!

# 2. Deploy
./customers/building_shop-deploy.sh

# 3. Start application
docker-compose -f customers/building_shop-docker-compose.yml up -d

# ✅ Building Shop is live at http://localhost:3000
```

### Add Hardware Store (Customer 2)

```bash
# 1. Run script AGAIN (same script!)
./add-customer.sh

# Fill in:
# - Company: Hardware Store LLC
# - Database: hardware_store
# - Admin: hardware_admin / HardwarePass123!

# 2. Deploy
./customers/hardware_store-deploy.sh

# 3. Start application (different ports)
docker-compose -f customers/hardware_store-docker-compose.yml up -d

# ✅ Hardware Store is live at http://localhost:3001
```

**Result:**
- Two separate databases
- Two separate application instances
- Same application code!
- Complete isolation

---

## 📂 File Structure After Adding Customers

```
terraform-tenants/
├── add-customer.sh                    ← Run this to add new customer
├── main.tf                            ← Terraform configuration
├── variables.tf
├── outputs.tf
├── modules/
│   └── tenant/
│       ├── main.tf
│       ├── variables.tf
│       └── database-schema.sql        ← Schema applied to each DB
│
└── customers/                         ← Generated per customer
    ├── building_shop.tfvars           ← Terraform config
    ├── building_shop-deploy.sh        ← Deployment script
    ├── building_shop-docker-compose.yml  ← Docker Compose
    │
    ├── hardware_store.tfvars
    ├── hardware_store-deploy.sh
    ├── hardware_store-docker-compose.yml
    │
    ├── bookstore.tfvars
    ├── bookstore-deploy.sh
    └── bookstore-docker-compose.yml
```

---

## 🔧 What Each File Does

### `add-customer.sh`
Interactive script that:
1. Prompts for customer information
2. Generates Terraform configuration
3. Generates deployment script
4. Generates Docker Compose file

### `customers/{name}.tfvars`
Terraform configuration for ONE customer:
- Database name (customer chooses!)
- Admin/developer/customer users
- Emails, passwords, plan

### `customers/{name}-deploy.sh`
Automated deployment script:
1. Runs Terraform to create database
2. Shows connection information
3. Generates Docker Compose file
4. Provides next steps

### `customers/{name}-docker-compose.yml`
Docker Compose file to run application:
- Backend with `DATABASE_DSN` pointing to customer DB
- Frontend with customer branding
- PostgreSQL with customer database

---

## 🚀 Typical Workflow

### Day 1: First Customer

```bash
cd terraform-tenants

# Add customer
./add-customer.sh
# -> Enter: Building Shop, building_shop, admin info

# Deploy
./customers/building_shop-deploy.sh

# Start application
docker-compose -f customers/building_shop-docker-compose.yml up -d

# Access: http://localhost:3000
```

**Result:**
- Database `building_shop` created
- Application running
- Customer can start using it

### Day 2: Second Customer

```bash
# Add another customer (SAME SCRIPT!)
./add-customer.sh
# -> Enter: Hardware Store, hardware_store, admin info

# Deploy
./customers/hardware_store-deploy.sh

# Start application (different ports in compose file)
docker-compose -f customers/hardware_store-docker-compose.yml up -d

# Access: http://localhost:3001
```

**Result:**
- Database `hardware_store` created (separate!)
- Another application instance running
- Both customers completely isolated

### Day 30: Tenth Customer

```bash
# Still the SAME SCRIPT!
./add-customer.sh
# -> Enter: BookStore, bookstore_2025, admin info

./customers/bookstore_2025-deploy.sh
docker-compose -f customers/bookstore_2025-docker-compose.yml up -d
```

**Scalability: Unlimited customers!**

---

## 🔍 Example Database State

After adding 3 customers, you have **3 completely separate databases**:

### PostgreSQL Instance

```
postgres=# \l
                                  List of databases
      Name       | Owner    | Encoding |   Collate   |    Ctype    
-----------------+----------+----------+-------------+-------------
 building_shop   | postgres | UTF8     | en_US.UTF-8 | en_US.UTF-8
 hardware_store  | postgres | UTF8     | en_US.UTF-8 | en_US.UTF-8
 bookstore_2025  | postgres | UTF8     | en_US.UTF-8 | en_US.UTF-8
```

### Each database has the SAME schema

```
building_shop=# \dt
              List of relations
 Schema |      Name       | Type  |     Owner      
--------+-----------------+-------+----------------
 public | customers       | table | building_admin
 public | orders          | table | building_admin
 public | widgets         | table | building_admin
 public | transactions    | table | building_admin
 public | statuses        | table | building_admin
```

### But DIFFERENT data

```
building_shop=> SELECT * FROM customers;
 Alice Johnson, alice@building.com
 Bob Smith, bob@building.com

hardware_store=> SELECT * FROM customers;
 Jane Doe, jane@hardware.com
 John Wilson, john@hardware.com

bookstore_2025=> SELECT * FROM customers;
 Mary Brown, mary@bookstore.com
```

**Complete isolation!**

---

## 💡 Key Points

### ✅ What Terraform Does

1. Creates database with customer-chosen name
2. Creates admin/developer/customer users
3. Applies complete schema (tables, indexes, triggers)
4. Sets up proper permissions per role

### ✅ What Application Does

1. Reads `DATABASE_DSN` from environment
2. Connects to customer's database
3. Handles business logic
4. Serves requests

### ✅ What You Get

- **No master database** tracking tenants
- **Complete separation** between customers
- **Same application code** for all
- **Easy onboarding** with `add-customer.sh`
- **Automated deployment** with generated scripts
- **Infinite scalability** - add unlimited customers

---

## 🎯 Summary

```
┌─────────────────────────────────────────────────────────┐
│  1. Run ./add-customer.sh                               │
│     → Fill in customer info (name, database, users)    │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  2. Script generates Terraform config                   │
│     → customers/{name}.tfvars                           │
│     → customers/{name}-deploy.sh                        │
│     → customers/{name}-docker-compose.yml               │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  3. Run ./customers/{name}-deploy.sh                    │
│     → Terraform creates database                        │
│     → Shows connection info                             │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  4. Start application                                   │
│     → docker-compose -f customers/{name}-docker...      │
│     → Application uses customer's database              │
└─────────────────────────────────────────────────────────┘

✅ Customer is live! Repeat for next customer.
```

---

**NO master database. NO metadata tracking. Just isolated customer databases.** 🎯

