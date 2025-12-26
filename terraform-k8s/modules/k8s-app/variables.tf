# Module variables

variable "manifests_dir" {
  description = "Directory containing Kubernetes YAML manifests"
  type        = string
}

variable "environment" {
  description = "Environment name (local, dev, staging, prod)"
  type        = string
  default     = "local"
}

variable "use_loadbalancer" {
  description = "Whether to use LoadBalancer service type (for cloud deployments)"
  type        = bool
  default     = false
}

variable "enable_ingress" {
  description = "Whether to enable ingress (for cloud deployments)"
  type        = bool
  default     = false
}

