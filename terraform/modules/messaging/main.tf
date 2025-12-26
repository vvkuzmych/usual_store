# Messaging Service Module

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Messaging Service
resource "docker_image" "messaging_service" {
  name         = "usual_store-messaging-service:latest"
  keep_locally = false

  build {
    context    = "${path.root}/.."
    dockerfile = "cmd/messaging-service/Dockerfile"
    tag        = ["usual_store-messaging-service:latest"]
  }
}

resource "docker_container" "messaging_service" {
  name    = "usualstore-messaging-service"
  image   = docker_image.messaging_service.image_id
  restart = "unless-stopped"

  env = [
    "KAFKA_BROKERS=${var.kafka_broker}",
    "SMTP_HOST=sandbox.smtp.mailtrap.io",
    "SMTP_PORT=587",
    "SMTP_USER=d43e9bfd6010ba",
    "SMTP_USERNAME=d43e9bfd6010ba",
    "SMTP_PASSWORD=8d80f34d4bbe3d",
    "SMTP_FROM=noreply@usualstore.com"
  ]

  ports {
    internal = 6001
    external = 6001
  }

  networks_advanced {
    name    = var.network_id
    aliases = ["messaging-service"]
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:6001/health"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "messaging-service"
  }

  labels {
    label = "com.usualstore.tier"
    value = "backend"
  }

  depends_on = [var.kafka_dependency]
}
