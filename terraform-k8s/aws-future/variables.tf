# Variables for AWS EKS deployment
# 🚨 NOT ACTIVE YET - For future use

variable "aws_region" {
  description = "AWS region for EKS cluster"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
  default     = "usualstore"
}

variable "node_instance_type" {
  description = "EC2 instance type for worker nodes"
  type        = string
  default     = "t3.medium"
  
  # Cost-effective options:
  # - t3.small  (2 vCPU, 2 GB RAM) - $0.0208/hr
  # - t3.medium (2 vCPU, 4 GB RAM) - $0.0416/hr  ← Recommended
  # - t3.large  (2 vCPU, 8 GB RAM) - $0.0832/hr
}

variable "min_nodes" {
  description = "Minimum number of worker nodes"
  type        = number
  default     = 2
}

variable "max_nodes" {
  description = "Maximum number of worker nodes"
  type        = number
  default     = 10
}

variable "desired_nodes" {
  description = "Desired number of worker nodes"
  type        = number
  default     = 3
}

