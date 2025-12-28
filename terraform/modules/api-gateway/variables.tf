variable "network_name" {
  description = "Docker network name"
  type        = string
}

variable "gateway_port" {
  description = "API Gateway external port (HTTP)"
  type        = number
  default     = 8000
}

variable "gateway_ssl_port" {
  description = "API Gateway external port (HTTPS)"
  type        = number
  default     = 8443
}

variable "admin_port" {
  description = "Kong Admin API port"
  type        = number
  default     = 8001
}

variable "backend_port" {
  description = "Backend API port"
  type        = number
  default     = 4001
}

