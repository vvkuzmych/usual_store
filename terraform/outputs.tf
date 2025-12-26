# Outputs for Usual Store infrastructure

output "network_id" {
  description = "Docker network ID"
  value       = docker_network.usualstore_network.id
}

output "network_name" {
  description = "Docker network name"
  value       = docker_network.usualstore_network.name
}

output "database_connection_string" {
  description = "PostgreSQL connection string"
  value       = module.database.connection_string
  sensitive   = true
}

output "database_container_id" {
  description = "PostgreSQL container ID"
  value       = module.database.container_id
}

output "kafka_broker" {
  description = "Kafka broker address"
  value       = module.kafka_stack.kafka_broker
}

output "kafka_ui_url" {
  description = "Kafka UI URL"
  value       = "http://localhost:8090"
}

output "backend_api_url" {
  description = "Backend API URL"
  value       = "http://localhost:${var.api_port}"
}

output "support_service_url" {
  description = "Support service URL"
  value       = "http://localhost:${var.support_port}"
}

output "frontend_urls" {
  description = "Frontend application URLs"
  value = {
    react      = "http://localhost:${var.react_port}"
    typescript = "http://localhost:${var.typescript_port}"
    redux      = "http://localhost:${var.redux_port}"
    support_ui = "http://localhost:${var.support_ui_port}"
  }
}

output "jaeger_ui_url" {
  description = "Jaeger UI URL for tracing"
  value       = "http://localhost:16686"
}

output "policy_server_url" {
  description = "OPA policy server URL"
  value       = module.policies.policy_server_url
}

output "service_status" {
  description = "Status of all services"
  value = {
    database    = module.database.status
    kafka       = module.kafka_stack.status
    backend_api = module.backend_api.status
    support     = module.support_service.status
    messaging   = module.messaging_service.status
  }
}

output "volumes" {
  description = "Created volumes"
  value = {
    database  = docker_volume.db_data.name
    kafka     = docker_volume.kafka_data.name
    zookeeper = docker_volume.zookeeper_data.name
  }
}

