terraform {
  required_providers {
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "~> 1.21"
    }
  }
}

# Create tenant database
resource "postgresql_database" "tenant" {
  name              = var.database_name
  owner             = var.postgres_user
  template          = "template0"
  lc_collate        = "en_US.UTF-8"
  lc_ctype          = "en_US.UTF-8"
  connection_limit  = local.connection_limit
  allow_connections = true
  
  lifecycle {
    prevent_destroy = false
  }
}

# Create database schema
resource "null_resource" "schema" {
  depends_on = [postgresql_database.tenant]
  
  provisioner "local-exec" {
    command = <<-EOT
      PGPASSWORD=${var.postgres_password} psql \
        -h ${var.postgres_host} \
        -p ${var.postgres_port} \
        -U ${var.postgres_user} \
        -d ${var.database_name} \
        -f ${path.module}/database-schema.sql
    EOT
  }
  
  triggers = {
    database_name = var.database_name
    schema_hash   = filemd5("${path.module}/database-schema.sql")
  }
}

# Create admin users (full access)
resource "postgresql_role" "admins" {
  for_each = { for idx, admin in var.admins : admin.username => admin }
  
  name     = each.value.username
  login    = true
  password = each.value.password
  
  # Admin privileges
  superuser         = false
  create_database   = false
  create_role       = false
  inherit           = true
  replication       = false
  bypass_row_level_security = false
  connection_limit  = -1
  
  depends_on = [postgresql_database.tenant]
}

# Grant admin access to tenant database
resource "postgresql_grant" "admin_database" {
  for_each = { for idx, admin in var.admins : admin.username => admin }
  
  database    = var.database_name
  role        = postgresql_role.admins[each.key].name
  object_type = "database"
  privileges  = ["CREATE", "CONNECT", "TEMPORARY"]
  
  depends_on = [
    postgresql_database.tenant,
    postgresql_role.admins,
    null_resource.schema
  ]
}

# Grant admin access to all tables
resource "postgresql_grant" "admin_tables" {
  for_each = { for idx, admin in var.admins : admin.username => admin }
  
  database    = var.database_name
  role        = postgresql_role.admins[each.key].name
  schema      = "public"
  object_type = "table"
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE", "REFERENCES", "TRIGGER"]
  
  depends_on = [
    postgresql_grant.admin_database,
    null_resource.schema
  ]
}

# Grant admin access to sequences
resource "postgresql_grant" "admin_sequences" {
  for_each = { for idx, admin in var.admins : admin.username => admin }
  
  database    = var.database_name
  role        = postgresql_role.admins[each.key].name
  schema      = "public"
  object_type = "sequence"
  privileges  = ["SELECT", "UPDATE", "USAGE"]
  
  depends_on = [
    postgresql_grant.admin_database,
    null_resource.schema
  ]
}

# Create developer users (read/write access, no delete)
resource "postgresql_role" "developers" {
  for_each = { for idx, dev in var.developers : dev.username => dev }
  
  name     = each.value.username
  login    = true
  password = each.value.password
  
  # Developer privileges
  superuser         = false
  create_database   = false
  create_role       = false
  inherit           = true
  replication       = false
  bypass_row_level_security = false
  connection_limit  = -1
  
  depends_on = [postgresql_database.tenant]
}

# Grant developer database access
resource "postgresql_grant" "developer_database" {
  for_each = { for idx, dev in var.developers : dev.username => dev }
  
  database    = var.database_name
  role        = postgresql_role.developers[each.key].name
  object_type = "database"
  privileges  = ["CONNECT", "TEMPORARY"]
  
  depends_on = [
    postgresql_database.tenant,
    postgresql_role.developers,
    null_resource.schema
  ]
}

# Grant developer table access (read/write, no delete)
resource "postgresql_grant" "developer_tables" {
  for_each = { for idx, dev in var.developers : dev.username => dev }
  
  database    = var.database_name
  role        = postgresql_role.developers[each.key].name
  schema      = "public"
  object_type = "table"
  privileges  = ["SELECT", "INSERT", "UPDATE"]  # No DELETE
  
  depends_on = [
    postgresql_grant.developer_database,
    null_resource.schema
  ]
}

# Grant developer sequence access
resource "postgresql_grant" "developer_sequences" {
  for_each = { for idx, dev in var.developers : dev.username => dev }
  
  database    = var.database_name
  role        = postgresql_role.developers[each.key].name
  schema      = "public"
  object_type = "sequence"
  privileges  = ["SELECT", "USAGE"]
  
  depends_on = [
    postgresql_grant.developer_database,
    null_resource.schema
  ]
}

# Create customer users (read-only access)
resource "postgresql_role" "customers" {
  for_each = { for idx, cust in var.customers : cust.username => cust }
  
  name     = each.value.username
  login    = true
  password = each.value.password
  
  # Customer privileges (minimal)
  superuser         = false
  create_database   = false
  create_role       = false
  inherit           = true
  replication       = false
  bypass_row_level_security = false
  connection_limit  = -1
  
  depends_on = [postgresql_database.tenant]
}

# Grant customer database access
resource "postgresql_grant" "customer_database" {
  for_each = { for idx, cust in var.customers : cust.username => cust }
  
  database    = var.database_name
  role        = postgresql_role.customers[each.key].name
  object_type = "database"
  privileges  = ["CONNECT"]  # Only connect
  
  depends_on = [
    postgresql_database.tenant,
    postgresql_role.customers,
    null_resource.schema
  ]
}

# Grant customer table access (read-only)
resource "postgresql_grant" "customer_tables" {
  for_each = { for idx, cust in var.customers : cust.username => cust }
  
  database    = var.database_name
  role        = postgresql_role.customers[each.key].name
  schema      = "public"
  object_type = "table"
  privileges  = ["SELECT"]  # Only SELECT
  
  depends_on = [
    postgresql_grant.customer_database,
    null_resource.schema
  ]
}

# Local variables
locals {
  connection_limit = {
    "free"         = 5
    "starter"      = 10
    "professional" = 25
    "enterprise"   = 50
  }[var.plan]
}

