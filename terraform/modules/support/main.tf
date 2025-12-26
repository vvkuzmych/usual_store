# Support Service module - placeholder
variable "network_id" { type = string }
variable "support_port" { type = number }
variable "database_dsn" { type = string }

output "status" { value = "configured" }
output "support_url" { value = "http://localhost:${var.support_port}" }

