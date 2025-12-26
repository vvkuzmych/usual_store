output "status" {
  description = "Backend service status"
  value       = "running"
}

output "api_url" {
  description = "Backend API URL"
  value       = "http://localhost:${var.api_port}"
}

output "backend_container_id" {
  description = "Backend container ID"
  value       = docker_container.backend.id
}

output "ai_assistant_enabled" {
  description = "AI Assistant service status"
  value       = var.enable_ai_assistant
}

