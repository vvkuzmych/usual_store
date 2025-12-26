output "status" {
  description = "Messaging service status"
  value       = "running"
}

output "container_id" {
  description = "Messaging service container ID"
  value       = docker_container.messaging_service.id
}

