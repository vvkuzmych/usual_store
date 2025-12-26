# Terraform configuration for LOCAL Kubernetes deployment
# Works with: Docker Desktop K8s, Minikube, or Kind

terraform {
  required_version = ">= 1.0"
  
  required_providers {
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

# Provider for local Kubernetes cluster
provider "kubernetes" {
  config_path    = var.kubeconfig_path
  config_context = var.kube_context
}

provider "kubectl" {
  config_path    = var.kubeconfig_path
  config_context = var.kube_context
}

# Deploy all Kubernetes manifests
module "usualstore_app" {
  source = "../modules/k8s-app"
  
  manifests_dir = "${path.root}/../../k8s"
  environment   = "local"
  
  # Local-specific settings
  use_loadbalancer = false  # Use NodePort for local
  enable_ingress   = false  # No ingress for local
}

# Output useful information
output "cluster_context" {
  description = "Kubernetes context being used"
  value       = var.kube_context
}

output "namespace" {
  description = "Application namespace"
  value       = "usualstore"
}

output "access_instructions" {
  description = "How to access your application"
  value       = <<-EOT
  
  ✅ Deployment complete!
  
  Check status:
    kubectl get pods -n usualstore
    kubectl get services -n usualstore
  
  Access frontend:
    kubectl port-forward svc/frontend 3000:80 -n usualstore
    Then visit: http://localhost:3000
  
  Access backend API:
    kubectl port-forward svc/backend 4001:4001 -n usualstore
    Then visit: http://localhost:4001
  
  View logs:
    kubectl logs -f deployment/backend -n usualstore
    kubectl logs -f deployment/frontend -n usualstore
  
  EOT
}

