variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "volume_id" {
  description = "Docker volume name for PostgreSQL data"
  type        = string
}

variable "postgres_user" {
  description = "PostgreSQL user"
  type        = string
  default     = "postgres"
}

variable "postgres_password" {
  description = "PostgreSQL password"
  type        = string
  sensitive   = true
}

variable "postgres_db" {
  description = "PostgreSQL database name"
  type        = string
  default     = "usualstore"
}

