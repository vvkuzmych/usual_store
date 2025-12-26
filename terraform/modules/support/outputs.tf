output "status" {
  description = "Support service status"
  value       = "running"
}

output "support_url" {
  description = "Support service URL"
  value       = "http://localhost:${var.support_port}"
}

output "container_id" {
  description = "Support service container ID"
  value       = docker_container.support_service.id
}

