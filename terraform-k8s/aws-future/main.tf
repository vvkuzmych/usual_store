# Terraform configuration for AWS EKS deployment
# 🚨 NOT ACTIVE YET - For future AWS deployment
# 
# This file is prepared but NOT deployed yet.
# When you're ready to deploy to AWS:
# 1. Create an AWS account
# 2. Install AWS CLI and configure credentials
# 3. Uncomment the resources below
# 4. Run: terraform init && terraform apply

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = "~> 1.14"
    }
  }
}

# AWS Provider (uncomment when ready)
# provider "aws" {
#   region = var.aws_region
# }

# EKS Cluster (uncomment when ready)
# module "eks" {
#   source  = "terraform-aws-modules/eks/aws"
#   version = "~> 19.0"
#   
#   cluster_name    = "usualstore-${var.environment}"
#   cluster_version = "1.28"
#   
#   vpc_id     = module.vpc.vpc_id
#   subnet_ids = module.vpc.private_subnets
#   
#   eks_managed_node_groups = {
#     main = {
#       min_size     = 2
#       max_size     = 10
#       desired_size = 3
#       
#       instance_types = ["t3.medium"]
#       capacity_type  = "ON_DEMAND"
#     }
#   }
# }

# VPC for EKS (uncomment when ready)
# module "vpc" {
#   source  = "terraform-aws-modules/vpc/aws"
#   version = "~> 5.0"
#   
#   name = "usualstore-vpc"
#   cidr = "10.0.0.0/16"
#   
#   azs             = ["${var.aws_region}a", "${var.aws_region}b", "${var.aws_region}c"]
#   private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
#   public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
#   
#   enable_nat_gateway = true
#   single_nat_gateway = false
#   
#   public_subnet_tags = {
#     "kubernetes.io/role/elb" = "1"
#   }
#   
#   private_subnet_tags = {
#     "kubernetes.io/role/internal-elb" = "1"
#   }
# }

# Kubernetes Provider for EKS (uncomment when ready)
# provider "kubernetes" {
#   host                   = module.eks.cluster_endpoint
#   cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
#   
#   exec {
#     api_version = "client.authentication.k8s.io/v1beta1"
#     command     = "aws"
#     args = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
#   }
# }

# provider "kubectl" {
#   host                   = module.eks.cluster_endpoint
#   cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
#   load_config_file       = false
#   
#   exec {
#     api_version = "client.authentication.k8s.io/v1beta1"
#     command     = "aws"
#     args = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
#   }
# }

# Deploy application to EKS (uncomment when ready)
# module "usualstore_app" {
#   source = "../modules/k8s-app"
#   
#   manifests_dir = "${path.root}/../../k8s"
#   environment   = var.environment
#   
#   # AWS-specific settings
#   use_loadbalancer = true   # Use AWS LoadBalancer
#   enable_ingress   = true   # Use ALB Ingress
#   
#   depends_on = [module.eks]
# }

# Outputs (uncomment when ready)
# output "cluster_endpoint" {
#   description = "EKS cluster endpoint"
#   value       = module.eks.cluster_endpoint
# }
# 
# output "cluster_name" {
#   description = "EKS cluster name"
#   value       = module.eks.cluster_name
# }
# 
# output "configure_kubectl" {
#   description = "Command to configure kubectl"
#   value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
# }

