# Backend API module - placeholder for custom builds
# You can extend this to build and run your backend API container

variable "network_id" { type = string }
variable "api_port" { type = number }
variable "database_dsn" { type = string }
variable "stripe_key" { type = string }
variable "stripe_secret" { type = string }

output "status" { value = "configured" }
output "api_url" { value = "http://localhost:${var.api_port}" }

# TODO: Add docker_image and docker_container resources for your backend
# Example:
# resource "docker_image" "backend" {
#   name = "usualstore/backend:latest"
#   build { context = "../../" }
# }

