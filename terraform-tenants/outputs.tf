output "tenants" {
  description = "Created tenant databases"
  value = {
    for tenant_key, tenant in module.tenants : tenant_key => {
      database_name = tenant.database_name
      tenant_name   = tenant.tenant_name
      plan          = tenant.plan
      admin_users   = tenant.admin_usernames
      dev_users     = tenant.developer_usernames
      customer_users = tenant.customer_usernames
    }
  }
}

output "connection_strings" {
  description = "Connection strings for each tenant database (sensitive)"
  value = {
    for tenant_key, tenant in module.tenants : tenant_key => {
      admin_dsn = "host=${var.postgres_host} port=${var.postgres_port} user=${tenant.admin_usernames[0]} password=*** dbname=${tenant.database_name} sslmode=${var.postgres_sslmode}"
      dev_dsn   = length(tenant.developer_usernames) > 0 ? "host=${var.postgres_host} port=${var.postgres_port} user=${tenant.developer_usernames[0]} password=*** dbname=${tenant.database_name} sslmode=${var.postgres_sslmode}" : null
    }
  }
  sensitive = true
}

output "summary" {
  description = "Summary of created infrastructure"
  value = {
    total_tenants    = length(var.tenants)
    master_database  = var.master_database
    postgres_host    = var.postgres_host
    tenant_databases = [for tenant in module.tenants : tenant.database_name]
  }
}

