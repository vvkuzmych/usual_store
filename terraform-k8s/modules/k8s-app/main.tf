# Reusable module for deploying Usual Store to Kubernetes
# Works with both local and cloud Kubernetes clusters

terraform {
  required_providers {
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = "~> 1.14"
    }
  }
}

# Get all YAML files from k8s directory
locals {
  k8s_files = fileset(var.manifests_dir, "*.yaml")
  
  # Sort files to ensure correct deployment order
  sorted_files = sort(local.k8s_files)
}

# Apply each Kubernetes manifest
resource "kubectl_manifest" "usualstore" {
  for_each = { for idx, file in local.sorted_files : file => file }
  
  yaml_body = file("${var.manifests_dir}/${each.value}")
  
  # Don't fail if resource already exists
  force_conflicts = true
  
  # Wait for resources to be ready
  wait = true
  
  # Server-side apply for better handling
  server_side_apply = true
  
  # Override existing resources if needed
  override_namespace = "usualstore"
}

# Wait for all deployments to be ready
resource "null_resource" "wait_for_deployments" {
  depends_on = [kubectl_manifest.usualstore]
  
  provisioner "local-exec" {
    command = <<-EOT
      echo "Waiting for deployments to be ready..."
      kubectl wait --for=condition=available --timeout=300s \
        deployment --all -n usualstore 2>/dev/null || true
    EOT
  }
}

