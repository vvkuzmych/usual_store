output "container_id" {
  description = "PostgreSQL container ID"
  value       = docker_container.database.id
}

output "container_name" {
  description = "PostgreSQL container name"
  value       = docker_container.database.name
}

output "connection_string" {
  description = "PostgreSQL connection string"
  value       = "postgres://${var.postgres_user}:${var.postgres_password}@database:5432/${var.postgres_db}?sslmode=disable"
  sensitive   = true
}

output "internal_host" {
  description = "Internal hostname for database"
  value       = "database"
}

output "status" {
  description = "Database container status"
  value       = "running"
}

