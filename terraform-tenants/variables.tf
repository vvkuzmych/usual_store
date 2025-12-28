variable "postgres_host" {
  description = "PostgreSQL host"
  type        = string
  default     = "localhost"
}

variable "postgres_port" {
  description = "PostgreSQL port"
  type        = number
  default     = 5432
}

variable "postgres_user" {
  description = "PostgreSQL superuser"
  type        = string
  default     = "postgres"
}

variable "postgres_password" {
  description = "PostgreSQL superuser password"
  type        = string
  sensitive   = true
}

variable "postgres_sslmode" {
  description = "PostgreSQL SSL mode"
  type        = string
  default     = "disable"
}

variable "master_database" {
  description = "Master database name for metadata"
  type        = string
  default     = "usualstore"
}

variable "tenants" {
  description = "Map of tenant configurations"
  type = map(object({
    database_name = string
    tenant_name   = string
    plan          = string
    
    admins = list(object({
      username = string
      password = string
      email    = string
    }))
    
    developers = list(object({
      username = string
      password = string
      email    = string
    }))
    
    customers = list(object({
      username = string
      password = string
      email    = string
    }))
  }))
  
  default = {}
  
  validation {
    condition = alltrue([
      for tenant_key, tenant in var.tenants : 
      can(regex("^[a-z][a-z0-9_]*$", tenant.database_name))
    ])
    error_message = "Database names must start with a letter and contain only lowercase letters, numbers, and underscores."
  }
  
  validation {
    condition = alltrue([
      for tenant_key, tenant in var.tenants : 
      contains(["free", "starter", "professional", "enterprise"], tenant.plan)
    ])
    error_message = "Plan must be one of: free, starter, professional, enterprise."
  }
}

