variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "policy_files" {
  description = "List of policy files to load"
  type        = list(string)
  default     = []
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "development"
}

