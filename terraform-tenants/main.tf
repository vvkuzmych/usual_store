terraform {
  required_version = ">= 1.0"
  
  required_providers {
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "~> 1.21"
    }
  }
}

# Configure PostgreSQL provider
provider "postgresql" {
  host            = var.postgres_host
  port            = var.postgres_port
  username        = var.postgres_user
  password        = var.postgres_password
  database        = var.master_database
  sslmode         = var.postgres_sslmode
  connect_timeout = 15
  superuser       = false
}

# Create tenant databases and access control
module "tenants" {
  source = "./modules/tenant"
  
  for_each = var.tenants
  
  # Tenant configuration
  tenant_key    = each.key
  database_name = each.value.database_name
  tenant_name   = each.value.tenant_name
  plan          = each.value.plan
  
  # Users configuration
  admins     = each.value.admins
  developers = each.value.developers
  customers  = each.value.customers
  
  # Database connection
  postgres_host = var.postgres_host
  postgres_port = var.postgres_port
  postgres_user = var.postgres_user
  postgres_password = var.postgres_password
  
  depends_on = [
    # Ensure master database exists
  ]
}

# Store tenant metadata in master database
resource "postgresql_database" "master" {
  name              = var.master_database
  owner             = var.postgres_user
  template          = "template0"
  lc_collate        = "en_US.UTF-8"
  lc_ctype          = "en_US.UTF-8"
  connection_limit  = -1
  allow_connections = true
}

