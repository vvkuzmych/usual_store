output "gateway_url" {
  value       = "http://localhost:${var.gateway_port}"
  description = "API Gateway proxy URL"
}

output "admin_url" {
  value       = "http://localhost:${var.admin_port}"
  description = "Kong Admin API URL for management"
}

output "gateway_https_url" {
  value       = "https://localhost:${var.gateway_ssl_port}"
  description = "API Gateway HTTPS URL"
}

output "container_name" {
  value       = "api-gateway"
  description = "API Gateway container name"
}

