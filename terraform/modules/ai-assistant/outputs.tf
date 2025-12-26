# Outputs for AI Assistant Module

output "ai_url" {
  description = "AI Assistant service URL"
  value       = "http://ai-assistant:8080"
}

output "container_id" {
  description = "AI Assistant container ID"
  value       = docker_container.ai_assistant.id
}

