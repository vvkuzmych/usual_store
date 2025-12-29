output "database_name" {
  description = "Name of the created database"
  value       = postgresql_database.tenant.name
}

output "tenant_name" {
  description = "Display name of the tenant"
  value       = var.tenant_name
}

output "plan" {
  description = "Subscription plan"
  value       = var.plan
}

output "connection_limit" {
  description = "Maximum number of connections"
  value       = local.connection_limit
}

output "admin_usernames" {
  description = "List of admin usernames"
  value       = [for admin in var.admins : admin.username]
}

output "developer_usernames" {
  description = "List of developer usernames"
  value       = [for dev in var.developers : dev.username]
}

output "customer_usernames" {
  description = "List of customer usernames"
  value       = [for cust in var.customers : cust.username]
}

