# Support Service Module

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Support Service
resource "docker_image" "support_service" {
  name         = "usual_store-support-service:latest"
  keep_locally = false

  build {
    context    = "${path.root}/.."
    dockerfile = "cmd/support-service/Dockerfile"
    tag        = ["usual_store-support-service:latest"]
  }
}

resource "docker_container" "support_service" {
  name    = "usualstore-support-service"
  image   = docker_image.support_service.image_id
  restart = "unless-stopped"

  env = [
    "SUPPORT_SERVICE_PORT=${var.support_port}",
    "DATABASE_DSN=${var.database_dsn}"
  ]

  ports {
    internal = var.support_port
    external = var.support_port
    ip       = "127.0.0.1"
  }

  ports {
    internal = var.support_port
    external = var.support_port
    ip       = "::1"
  }

  networks_advanced {
    name    = var.network_id
    aliases = ["support-service"]
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:${var.support_port}/api/support/health"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "support-service"
  }

  labels {
    label = "com.usualstore.tier"
    value = "backend"
  }

  depends_on = [var.database_dependency]
}
