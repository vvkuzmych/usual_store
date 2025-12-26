# Variables for Usual Store Terraform configuration

variable "docker_host" {
  description = "Docker daemon host"
  type        = string
  default     = "unix:///var/run/docker.sock"
}

# Network Configuration
variable "ipv4_subnet" {
  description = "IPv4 subnet for the Docker network"
  type        = string
  default     = "172.22.0.0/16"
}

variable "ipv4_gateway" {
  description = "IPv4 gateway for the Docker network"
  type        = string
  default     = "172.22.0.1"
}

variable "ipv6_subnet" {
  description = "IPv6 subnet for the Docker network"
  type        = string
  default     = "2001:db8:1::/64"
}

variable "ipv6_gateway" {
  description = "IPv6 gateway for the Docker network"
  type        = string
  default     = "2001:db8:1::1"
}

# Database Configuration
variable "postgres_user" {
  description = "PostgreSQL user"
  type        = string
  default     = "postgres"
}

variable "postgres_password" {
  description = "PostgreSQL password"
  type        = string
  default     = "password"
  sensitive   = true
}

variable "postgres_db" {
  description = "PostgreSQL database name"
  type        = string
  default     = "usualstore"
}

# Service Ports
variable "api_port" {
  description = "Backend API port"
  type        = number
  default     = 4001
}

variable "support_port" {
  description = "Support service port"
  type        = number
  default     = 5001
}

variable "ai_port" {
  description = "AI Assistant service port"
  type        = number
  default     = 8080
}

variable "react_port" {
  description = "React frontend port"
  type        = number
  default     = 3000
}

variable "typescript_port" {
  description = "TypeScript frontend port"
  type        = number
  default     = 3001
}

variable "redux_port" {
  description = "Redux frontend port"
  type        = number
  default     = 3002
}

variable "support_ui_port" {
  description = "Support UI port"
  type        = number
  default     = 3005
}

# Stripe Configuration
variable "stripe_key" {
  description = "Stripe publishable key"
  type        = string
  default     = ""
  sensitive   = true
}

variable "stripe_secret" {
  description = "Stripe secret key"
  type        = string
  default     = ""
  sensitive   = true
}

# OpenAI Configuration
variable "openai_api_key" {
  description = "OpenAI API key for AI assistant"
  type        = string
  default     = ""
  sensitive   = true
}

# Environment
variable "environment" {
  description = "Environment name (development, staging, production)"
  type        = string
  default     = "development"
}

# Feature Flags
variable "enable_kafka" {
  description = "Enable Kafka messaging infrastructure"
  type        = bool
  default     = true
}

variable "enable_observability" {
  description = "Enable observability stack (Jaeger, Prometheus)"
  type        = bool
  default     = true
}

variable "enable_ai_assistant" {
  description = "Enable AI assistant service"
  type        = bool
  default     = true
}

variable "enable_support" {
  description = "Enable support service and UI"
  type        = bool
  default     = true
}

