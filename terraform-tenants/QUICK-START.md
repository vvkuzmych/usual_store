# 🚀 Terraform Multi-Tenant Quick Start

Create isolated tenant databases with custom names using Terraform.

---

## Prerequisites

```bash
# Install Terraform
brew install terraform

# Verify installation
terraform version

# Ensure PostgreSQL is running
psql -U postgres -c "SELECT version();"
```

---

## Step 1: Initialize Terraform (1 minute)

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store/terraform-tenants

# Initialize Terraform (downloads providers)
terraform init
```

**Expected Output:**
```
Initializing modules...
Initializing the backend...
Initializing provider plugins...
- Finding cyrilgdn/postgresql versions matching "~> 1.21"...
- Installing cyrilgdn/postgresql v1.21.0...

Terraform has been successfully initialized!
```

---

## Step 2: Configure Your First Tenant (2 minutes)

Edit `tenants.tfvars`:

```hcl
postgres_password = "yourpostgrespassword"

tenants = {
  building_shop = {
    database_name = "building_shop"      # ← Customer chooses this!
    tenant_name   = "Building Shop Inc."
    plan          = "professional"
    
    admins = [
      {
        username = "building_admin"
        password = "AdminPass123!"
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
    
    customers = [
      {
        username = "building_customer"
        password = "CustPass123!"
        email    = "customer@buildingshop.com"
      }
    ]
  }
}
```

---

## Step 3: Preview Changes (30 seconds)

```bash
terraform plan -var-file="tenants.tfvars"
```

**Expected Output:**
```
Terraform will perform the following actions:

  # module.tenants["building_shop"].postgresql_database.tenant will be created
  + resource "postgresql_database" "tenant" {
      + name = "building_shop"
      ...
    }

  # module.tenants["building_shop"].postgresql_role.admins["building_admin"] will be created
  + resource "postgresql_role" "admins" {
      + name  = "building_admin"
      + login = true
      ...
    }

Plan: 15 to add, 0 to change, 0 to destroy.
```

---

## Step 4: Create Infrastructure (1 minute)

```bash
terraform apply -var-file="tenants.tfvars"
```

Type `yes` when prompted.

**What Terraform Creates:**
1. ✅ Database `building_shop`
2. ✅ Admin user with full access
3. ✅ Developer user with read/write (no delete)
4. ✅ Customer user with read-only
5. ✅ All tables, indexes, triggers
6. ✅ Proper permissions per role

---

## Step 5: Verify Creation (2 minutes)

### Check Databases

```bash
psql -U postgres -l | grep building_shop

# Output:
# building_shop | postgres | UTF8
```

### Check Tables

```bash
psql -U postgres -d building_shop -c "\dt"

# Output:
#  customers
#  orders
#  widgets
#  transactions
#  statuses
#  transaction_statuses
```

### Test Admin Access

```bash
# Admin can do everything
psql -U building_admin -d building_shop -c "
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) 
  VALUES ('Hammer', 2999, 50, NOW(), NOW());
  
  SELECT * FROM widgets;
"

# Success! ✅
```

### Test Developer Access

```bash
# Developer can read and write, but NOT delete
psql -U building_dev -d building_shop -c "
  SELECT * FROM widgets;
"
# ✅ Works!

psql -U building_dev -d building_shop -c "
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) 
  VALUES ('Drill', 8999, 30, NOW(), NOW());
"
# ✅ Works!

psql -U building_dev -d building_shop -c "
  DELETE FROM widgets WHERE id = 1;
"
# ❌ ERROR: permission denied
# ✅ Correct! Developer cannot delete
```

### Test Customer Access

```bash
# Customer can ONLY read
psql -U building_customer -d building_shop -c "
  SELECT * FROM widgets;
"
# ✅ Works!

psql -U building_customer -d building_shop -c "
  INSERT INTO widgets (name, price, inventory_level) 
  VALUES ('Wrench', 1999, 20);
"
# ❌ ERROR: permission denied
# ✅ Correct! Customer is read-only
```

---

## Step 6: View Terraform Outputs

```bash
terraform output

# Shows:
# - Created databases
# - Users per tenant
# - Connection strings
# - Summary
```

---

## Step 7: Add Second Tenant (2 minutes)

Edit `tenants.tfvars` and add:

```hcl
tenants = {
  building_shop = {
    # ... existing config ...
  }
  
  # Add new tenant
  hardware_store = {
    database_name = "hardware_store"     # Different database!
    tenant_name   = "Hardware Store Co."
    plan          = "starter"
    
    admins = [
      {
        username = "hardware_admin"
        password = "HardwareAdmin123!"
        email    = "admin@hardwarestore.com"
      }
    ]
    
    developers = [
      {
        username = "hardware_dev"
        password = "HardwareDev123!"
        email    = "dev@hardwarestore.com"
      }
    ]
    
    customers = []  # No customer users
  }
}
```

Apply changes:

```bash
terraform apply -var-file="tenants.tfvars"
```

**Result**: Second database `hardware_store` created with separate users!

---

## Step 8: Verify Isolation (1 minute)

```bash
# Insert into building_shop
psql -U building_admin -d building_shop -c "
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) 
  VALUES ('Building Tool', 5000, 10, NOW(), NOW());
"

# Insert into hardware_store
psql -U hardware_admin -d hardware_store -c "
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) 
  VALUES ('Hardware Tool', 3000, 20, NOW(), NOW());
"

# Check building_shop
psql -U building_admin -d building_shop -c "SELECT name FROM widgets;"
# Output: Building Tool, Hammer, Drill

# Check hardware_store
psql -U hardware_admin -d hardware_store -c "SELECT name FROM widgets;"
# Output: Hardware Tool

# ✅ COMPLETE ISOLATION! Each tenant sees only their data!
```

---

## Quick Reference Commands

### Infrastructure Management

```bash
# Preview changes
terraform plan -var-file="tenants.tfvars"

# Apply changes
terraform apply -var-file="tenants.tfvars"

# Destroy everything (DANGER!)
terraform destroy -var-file="tenants.tfvars"

# Show current state
terraform show

# List resources
terraform state list
```

### Database Management

```bash
# List all databases
psql -U postgres -l

# Connect as admin
psql -U building_admin -d building_shop

# Connect as developer
psql -U building_dev -d building_shop

# Connect as customer
psql -U building_customer -d building_shop
```

### Access Testing

```bash
# Test admin (full access)
psql -U building_admin -d building_shop -c "
  SELECT 'READ' AS test;
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) VALUES ('Test', 100, 1, NOW(), NOW());
  DELETE FROM widgets WHERE name = 'Test';
"
# ✅ All succeed

# Test developer (read/write, no delete)
psql -U building_dev -d building_shop -c "
  SELECT 'READ' AS test;
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) VALUES ('Test', 100, 1, NOW(), NOW());
  DELETE FROM widgets WHERE name = 'Test';
"
# ✅ SELECT works
# ✅ INSERT works
# ❌ DELETE fails (permission denied)

# Test customer (read-only)
psql -U building_customer -d building_shop -c "
  SELECT 'READ' AS test;
  INSERT INTO widgets (name, price, inventory_level, created_at, updated_at) VALUES ('Test', 100, 1, NOW(), NOW());
"
# ✅ SELECT works
# ❌ INSERT fails (permission denied)
```

---

## Connection Strings for Your Application

### Admin Connection
```
host=localhost port=5432 user=building_admin password=AdminPass123! dbname=building_shop sslmode=disable
```

### Developer Connection
```
host=localhost port=5432 user=building_dev password=DevPass123! dbname=building_shop sslmode=disable
```

### Customer Connection
```
host=localhost port=5432 user=building_customer password=CustPass123! dbname=building_shop sslmode=disable
```

---

## Updating Tenants

### Add New User

Edit `tenants.tfvars`:
```hcl
developers = [
  {
    username = "building_dev"
    password = "DevPass123!"
    email    = "dev@buildingshop.com"
  },
  {
    username = "building_dev2"    # ← New developer
    password = "DevPass456!"
    email    = "dev2@buildingshop.com"
  }
]
```

Apply:
```bash
terraform apply -var-file="tenants.tfvars"
```

### Change User Password

Edit password in `tenants.tfvars`, then:
```bash
terraform apply -var-file="tenants.tfvars"
```

### Remove User

Remove user block from `tenants.tfvars`, then:
```bash
terraform apply -var-file="tenants.tfvars"
```

### Delete Tenant

Remove entire tenant block from `tenants.tfvars`, then:
```bash
terraform apply -var-file="tenants.tfvars"
```

**⚠️ WARNING**: This will DELETE the database and all data!

---

## Best Practices

### 1. Use Strong Passwords
```hcl
# ❌ Bad
password = "123456"

# ✅ Good
password = "SecureP@ssw0rd!2025"
```

### 2. Version Control
```bash
# Add to .gitignore
echo "tenants.tfvars" >> .gitignore
echo "terraform.tfstate*" >> .gitignore
echo ".terraform/" >> .gitignore
```

### 3. Use Environment Variables
```bash
export TF_VAR_postgres_password="your_password"
terraform apply -var-file="tenants.tfvars"
```

### 4. Backup Before Changes
```bash
# Backup all tenant databases
for db in $(psql -U postgres -t -c "SELECT datname FROM pg_database WHERE datname LIKE '%shop%' OR datname LIKE '%store%'"); do
  pg_dump -U postgres $db > backup_${db}_$(date +%Y%m%d).sql
done
```

---

## Troubleshooting

### Issue: "Error: could not connect to server"

**Solution**: Check PostgreSQL is running
```bash
brew services start postgresql@14
# or
pg_ctl -D /usr/local/var/postgres start
```

### Issue: "Error: role already exists"

**Solution**: Terraform state is out of sync
```bash
# Import existing role
terraform import 'module.tenants["building_shop"].postgresql_role.admins["building_admin"]' building_admin

# Or destroy and recreate
terraform destroy -target='module.tenants["building_shop"].postgresql_role.admins["building_admin"]'
terraform apply -var-file="tenants.tfvars"
```

### Issue: "Error: database is being accessed by other users"

**Solution**: Close all connections
```bash
psql -U postgres -c "
  SELECT pg_terminate_backend(pg_stat_activity.pid)
  FROM pg_stat_activity
  WHERE pg_stat_activity.datname = 'building_shop'
  AND pid <> pg_backend_pid();
"
```

---

## Summary

**You've successfully:**
- ✅ Created isolated tenant databases with custom names
- ✅ Set up role-based access control
- ✅ Configured admin (full), developer (read/write), and customer (read-only) users
- ✅ Verified complete data isolation
- ✅ Used infrastructure as code with Terraform

**Next Steps:**
1. Add more tenants by editing `tenants.tfvars`
2. Connect your application to tenant databases
3. Implement tenant selection in your app
4. Set up automated backups

---

**Total Time**: ~10 minutes to full multi-tenant infrastructure! 🎉

