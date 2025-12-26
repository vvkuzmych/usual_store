# Backend API Module

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Backend API Service
resource "docker_image" "backend" {
  name         = "usual_store-back-end:latest"
  keep_locally = false

  build {
    context    = "${path.root}/.."
    dockerfile = "Dockerfile"
    tag        = ["usual_store-back-end:latest"]
  }
}

resource "docker_container" "backend" {
  name  = "usualstore-back-end"
  image = docker_image.backend.image_id

  restart = "unless-stopped"

  command = ["./usualstore_api", "-port=${var.api_port}"]

  env = [
    "STRIPE_SECRET=${var.stripe_secret}",
    "STRIPE_KEY=${var.stripe_key}",
    "API_PORT=${var.api_port}",
    "DATABASE_DSN=${var.database_dsn}",
    "SECRET_FOR_FRONT=your-secret-key-change-in-production",
    "FRONT_URL=http://localhost:3000",
    "SMTP_HOST=sandbox.smtp.mailtrap.io",
    "SMTP_PORT=587",
    "SMTP_USER=d43e9bfd6010ba",
    "SMTP_USERNAME=d43e9bfd6010ba",
    "SMTP_PASSWORD=8d80f34d4bbe3d",
    "SMTP_FROM=noreply@usualstore.com",
    "OTEL_ENABLED=true",
    "OTEL_SERVICE_NAME=usual-store-api",
    "OTEL_SERVICE_VERSION=1.0.0",
    "OTEL_ENVIRONMENT=development",
    "OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4318"
  ]

  ports {
    internal = var.api_port
    external = var.api_port
    ip       = "127.0.0.1"
  }

  ports {
    internal = var.api_port
    external = var.api_port
    ip       = "::1"
  }

  networks_advanced {
    name    = var.network_id
    aliases = ["back-end"]
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:${var.api_port}/api/health"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "backend-api"
  }

  labels {
    label = "com.usualstore.tier"
    value = "backend"
  }

  depends_on = [var.database_dependency]
}

# AI Assistant Service (optional)
resource "docker_image" "ai_assistant" {
  count        = var.enable_ai_assistant ? 1 : 0
  name         = "usual_store-ai-assistant:latest"
  keep_locally = false

  build {
    context    = "${path.root}/.."
    dockerfile = "Dockerfile.ai-assistant"
    tag        = ["usual_store-ai-assistant:latest"]
  }
}

resource "docker_container" "ai_assistant" {
  count   = var.enable_ai_assistant ? 1 : 0
  name    = "usualstore-ai-assistant"
  image   = docker_image.ai_assistant[0].image_id
  restart = "unless-stopped"

  env = [
    "PORT=8080",
    "DATABASE_DSN=${var.database_dsn}",
    "OPENAI_API_KEY=${var.openai_api_key}"
  ]

  ports {
    internal = 8080
    external = 8080
  }

  networks_advanced {
    name    = var.network_id
    aliases = ["ai-assistant"]
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "ai-assistant"
  }

  labels {
    label = "com.usualstore.tier"
    value = "backend"
  }

  depends_on = [var.database_dependency]
}
