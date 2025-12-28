variable "tenant_key" {
  description = "Unique key for this tenant"
  type        = string
}

variable "database_name" {
  description = "Name of the tenant database (customer chooses this)"
  type        = string
  
  validation {
    condition     = can(regex("^[a-z][a-z0-9_]*$", var.database_name))
    error_message = "Database name must start with a letter and contain only lowercase letters, numbers, and underscores."
  }
}

variable "tenant_name" {
  description = "Display name of the tenant"
  type        = string
}

variable "plan" {
  description = "Subscription plan"
  type        = string
  
  validation {
    condition     = contains(["free", "starter", "professional", "enterprise"], var.plan)
    error_message = "Plan must be one of: free, starter, professional, enterprise."
  }
}

variable "admins" {
  description = "List of admin users for this tenant"
  type = list(object({
    username = string
    password = string
    email    = string
  }))
  default = []
}

variable "developers" {
  description = "List of developer users for this tenant"
  type = list(object({
    username = string
    password = string
    email    = string
  }))
  default = []
}

variable "customers" {
  description = "List of customer users for this tenant"
  type = list(object({
    username = string
    password = string
    email    = string
  }))
  default = []
}

variable "postgres_host" {
  description = "PostgreSQL host"
  type        = string
}

variable "postgres_port" {
  description = "PostgreSQL port"
  type        = number
}

variable "postgres_user" {
  description = "PostgreSQL superuser"
  type        = string
}

variable "postgres_password" {
  description = "PostgreSQL superuser password"
  type        = string
  sensitive   = true
}

