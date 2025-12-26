variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "support_port" {
  description = "Support service port"
  type        = number
  default     = 5001
}

variable "database_dsn" {
  description = "Database connection string"
  type        = string
  sensitive   = true
}

variable "database_dependency" {
  description = "Database container dependency"
  type        = any
}

