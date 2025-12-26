# Frontend Module Outputs

output "react_container_id" {
  description = "React frontend container ID"
  value       = docker_container.react_frontend.id
}

output "typescript_container_id" {
  description = "TypeScript frontend container ID"
  value       = docker_container.typescript_frontend.id
}

output "redux_container_id" {
  description = "Redux frontend container ID"
  value       = docker_container.redux_frontend.id
}

output "support_container_id" {
  description = "Support UI frontend container ID"
  value       = docker_container.support_frontend.id
}

output "status" {
  description = "Frontend services status"
  value       = "running"
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

