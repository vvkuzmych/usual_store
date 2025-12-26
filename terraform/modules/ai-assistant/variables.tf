# Variables for AI Assistant Module

variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "ai_port" {
  description = "AI Assistant service port"
  type        = number
  default     = 8080
}

variable "database_dsn" {
  description = "PostgreSQL connection string"
  type        = string
}

variable "openai_api_key" {
  description = "OpenAI API key for AI assistant"
  type        = string
  sensitive   = true
}

