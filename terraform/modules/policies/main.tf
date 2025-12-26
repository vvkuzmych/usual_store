# OPA (Open Policy Agent) Policy Management Module

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Pull OPA image
resource "docker_image" "opa" {
  name         = "openpolicyagent/opa:latest"
  keep_locally = false
}

# OPA Server container
resource "docker_container" "opa_server" {
  name  = "usualstore-opa-server"
  image = docker_image.opa.image_id

  restart = "unless-stopped"

  command = [
    "run",
    "--server",
    "--addr", "0.0.0.0:8181",
    "--log-level", "debug",
    "/policies"
  ]

  ports {
    internal = 8181
    external = 8181
  }

  networks_advanced {
    name = var.network_id
  }

  upload {
    content = file("${path.module}/policies/network.rego")
    file    = "/policies/network.rego"
  }

  upload {
    content = file("${path.module}/policies/resource_limits.rego")
    file    = "/policies/resource_limits.rego"
  }

  upload {
    content = file("${path.module}/policies/security.rego")
    file    = "/policies/security.rego"
  }

  upload {
    content = file("${path.module}/policies/access_control.rego")
    file    = "/policies/access_control.rego"
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8181/health"]
    interval = "10s"
    timeout  = "5s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "opa-policy-server"
  }

  labels {
    label = "com.usualstore.tier"
    value = "governance"
  }
}

# Policy Enforcer Service (sidecar pattern)
resource "docker_image" "policy_enforcer" {
  name         = "usualstore/policy-enforcer:latest"
  keep_locally = false

  # Build from Dockerfile
  build {
    context = "${path.module}/policy-enforcer"
    tag     = ["usualstore/policy-enforcer:latest"]
  }
}

resource "docker_container" "policy_enforcer" {
  name  = "usualstore-policy-enforcer"
  image = docker_image.policy_enforcer.image_id

  restart = "unless-stopped"

  env = [
    "OPA_SERVER=http://opa-server:8181",
    "DOCKER_HOST=unix:///var/run/docker.sock"
  ]

  networks_advanced {
    name = var.network_id
  }

  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }

  depends_on = [docker_container.opa_server]

  labels {
    label = "com.usualstore.service"
    value = "policy-enforcer"
  }

  labels {
    label = "com.usualstore.tier"
    value = "governance"
  }
}

