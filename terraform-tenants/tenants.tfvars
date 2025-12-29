# Terraform Tenants Configuration
# Customize this file to create your tenant databases

# PostgreSQL connection (can also be set via environment variables)
postgres_password = "yourpassword"  # Change this!

# Define your tenants
tenants = {
  # Example Tenant 1: Building Shop
  building_shop = {
    database_name = "building_shop"      # Customer chooses this name!
    tenant_name   = "Building Shop Inc."
    plan          = "professional"        # free, starter, professional, enterprise
    
    # Admin users (full access: read, write, delete)
    admins = [
      {
        username = "building_admin"
        password = "AdminPass123!"        # Change this!
        email    = "admin@buildingshop.com"
      }
    ]
    
    # Developer users (read/write access, no delete)
    developers = [
      {
        username = "building_dev1"
        password = "DevPass123!"          # Change this!
        email    = "dev1@buildingshop.com"
      },
      {
        username = "building_dev2"
        password = "DevPass456!"          # Change this!
        email    = "dev2@buildingshop.com"
      }
    ]
    
    # Customer users (read-only access)
    customers = [
      {
        username = "building_customer1"
        password = "CustPass123!"         # Change this!
        email    = "customer1@buildingshop.com"
      }
    ]
  }

  # Example Tenant 2: Hardware Store
  # Uncomment to create second tenant
  /*
  hardware_store = {
    database_name = "hardware_store"     # Different database name!
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
    
    customers = [
      {
        username = "hardware_customer"
        password = "HardwareCust123!"
        email    = "customer@hardwarestore.com"
      }
    ]
  }
  */
}

