# Module outputs

output "manifests_applied" {
  description = "List of Kubernetes manifests applied"
  value       = keys(kubectl_manifest.usualstore)
}

output "environment" {
  description = "Deployment environment"
  value       = var.environment
}

