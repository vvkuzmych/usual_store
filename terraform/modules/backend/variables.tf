variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "api_port" {
  description = "Backend API port"
  type        = number
  default     = 4001
}

variable "database_dsn" {
  description = "Database connection string"
  type        = string
  sensitive   = true
}

variable "stripe_key" {
  description = "Stripe publishable key"
  type        = string
  default     = ""
}

variable "stripe_secret" {
  description = "Stripe secret key"
  type        = string
  sensitive   = true
  default     = ""
}

variable "openai_api_key" {
  description = "OpenAI API key for AI assistant"
  type        = string
  sensitive   = true
  default     = ""
}

variable "enable_ai_assistant" {
  description = "Enable AI assistant service"
  type        = bool
  default     = false
}

variable "database_dependency" {
  description = "Database container dependency"
  type        = any
}

