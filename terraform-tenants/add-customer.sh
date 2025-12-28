#!/bin/bash

# 🚀 Add New Customer Script
# This script helps you add a new customer by generating Terraform configuration

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏢 ADD NEW CUSTOMER"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Get customer information
read -p "Customer Company Name (e.g., Building Shop Inc): " company_name
read -p "Database Name (e.g., building_shop): " db_name
read -p "Customer Key (e.g., building_shop): " customer_key
read -p "Plan (starter/professional/enterprise): " plan

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "👤 ADMIN USER"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
read -p "Admin Username: " admin_user
read -sp "Admin Password: " admin_pass
echo ""
read -p "Admin Email: " admin_email

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "👨‍💻 DEVELOPER USER (Optional - press Enter to skip)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
read -p "Developer Username (or Enter to skip): " dev_user

dev_config=""
if [ ! -z "$dev_user" ]; then
  read -sp "Developer Password: " dev_pass
  echo ""
  read -p "Developer Email: " dev_email
  
  dev_config="    developers = [
      {
        username = \"$dev_user\"
        password = \"$dev_pass\"
        email    = \"$dev_email\"
      }
    ]"
else
  dev_config="    developers = []"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧑 CUSTOMER USER (Optional - press Enter to skip)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
read -p "Customer Username (or Enter to skip): " cust_user

cust_config=""
if [ ! -z "$cust_user" ]; then
  read -sp "Customer Password: " cust_pass
  echo ""
  read -p "Customer Email: " cust_email
  
  cust_config="    customers = [
      {
        username = \"$cust_user\"
        password = \"$cust_pass\"
        email    = \"$cust_email\"
      }
    ]"
else
  cust_config="    customers = []"
fi

# Generate Terraform config
config_file="customers/${customer_key}.tfvars"
mkdir -p customers

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 GENERATING CONFIGURATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

cat > "$config_file" <<EOF
# Customer Configuration: $company_name
# Generated: $(date)

# Database connection (update with your actual DB host/credentials)
db_host     = "localhost"
db_port     = "5432"
db_username = "postgres"
db_password = "yourpassword"  # ⚠️ UPDATE THIS!

# Customer tenant definition
tenants = {
  $customer_key = {
    database_name = "$db_name"
    tenant_name   = "$company_name"
    plan          = "$plan"
    
    admins = [
      {
        username = "$admin_user"
        password = "$admin_pass"
        email    = "$admin_email"
      }
    ]
    
$dev_config
    
$cust_config
  }
}
EOF

echo "✅ Configuration saved to: $config_file"
echo ""

# Generate deployment script
deploy_script="customers/${customer_key}-deploy.sh"

cat > "$deploy_script" <<EOF
#!/bin/bash
# Deployment script for $company_name

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 DEPLOYING: $company_name"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Step 1: Create database with Terraform
echo "Step 1: Creating database '$db_name' with Terraform..."
cd "\$(dirname "\$0")/.."
terraform init
terraform apply -var-file="customers/${customer_key}.tfvars" -auto-approve

echo ""
echo "✅ Database '$db_name' created!"
echo ""

# Step 2: Show connection info
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 CONNECTION INFORMATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Database Name: $db_name"
echo "Admin User:    $admin_user"
echo ""
echo "Backend Environment Variable:"
echo "export DATABASE_DSN=\"host=localhost port=5432 user=$admin_user password=$admin_pass dbname=$db_name sslmode=disable\""
echo ""
echo "Frontend Environment Variable:"
echo "export REACT_APP_API_URL=\"http://localhost:4001\""
echo "export REACT_APP_TENANT_NAME=\"$company_name\""
echo ""

# Step 3: Generate docker-compose file
compose_file="customers/${customer_key}-docker-compose.yml"

echo "Step 2: Generating Docker Compose file..."

cat > "\$compose_file" <<COMPOSE
version: '3.8'

# Docker Compose for $company_name
# Database: $db_name

services:
  backend:
    image: usualstore/api:latest
    environment:
      DATABASE_DSN: "host=db port=5432 user=$admin_user password=$admin_pass dbname=$db_name sslmode=disable"
      API_PORT: "4001"
      TENANT_NAME: "$customer_key"
      STRIPE_KEY: "pk_test_${customer_key}"
      STRIPE_SECRET: "sk_test_${customer_key}"
      SECRET_FOR_FRONT: "your-jwt-secret-here"
      FRONT_URL: "http://localhost:3000"
    ports:
      - "4001:4001"
    depends_on:
      - db
    networks:
      - ${customer_key}_network

  frontend:
    image: usualstore/frontend:latest
    environment:
      REACT_APP_API_URL: "http://backend:4001"
      REACT_APP_TENANT_NAME: "$company_name"
      REACT_APP_STRIPE_KEY: "pk_test_${customer_key}"
    ports:
      - "3000:80"
    depends_on:
      - backend
    networks:
      - ${customer_key}_network

  db:
    image: postgres:14
    environment:
      POSTGRES_DB: $db_name
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: yourpassword
    volumes:
      - ${customer_key}_data:/var/lib/postgresql/data
    networks:
      - ${customer_key}_network
    ports:
      - "5432:5432"

networks:
  ${customer_key}_network:
    driver: bridge

volumes:
  ${customer_key}_data:
COMPOSE

echo "✅ Docker Compose saved to: \$compose_file"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🐳 TO START WITH DOCKER:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "docker-compose -f \$compose_file up -d"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ DEPLOYMENT COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
EOF

chmod +x "$deploy_script"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ CUSTOMER ADDED SUCCESSFULLY!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📁 Files created:"
echo "   - Configuration: $config_file"
echo "   - Deployment script: $deploy_script"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 NEXT STEPS:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Review configuration:"
echo "   vim $config_file"
echo ""
echo "2. Deploy customer database:"
echo "   ./$deploy_script"
echo ""
echo "3. Or manually:"
echo "   cd terraform-tenants"
echo "   terraform apply -var-file=\"$config_file\""
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

