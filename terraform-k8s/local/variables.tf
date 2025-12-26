# Variables for LOCAL Kubernetes deployment

variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "kube_context" {
  description = "Kubernetes context to use"
  type        = string
  default     = "docker-desktop"  # Change to "minikube" if using minikube
  
  # Other common values:
  # - "docker-desktop" (Docker Desktop)
  # - "minikube" (Minikube)
  # - "kind-usualstore" (Kind)
}

