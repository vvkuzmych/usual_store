# AI Assistant Service Module

terraform {
  required_providers {
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

resource "docker_image" "ai_assistant" {
  name         = "usual_store-ai-assistant:latest"
  keep_locally = false

  build {
    context    = "${path.root}/../"
    dockerfile = "Dockerfile.ai-assistant"
    tag        = ["usual_store-ai-assistant:latest"]
  }
}

resource "docker_container" "ai_assistant" {
  name  = "usualstore-ai-assistant"
  image = docker_image.ai_assistant.image_id

  restart = "unless-stopped"

  env = [
    "PORT=8080",
    "DATABASE_DSN=${var.database_dsn}",
    "OPENAI_API_KEY=${var.openai_api_key}"
  ]

  ports {
    internal = 8080
    external = var.ai_port
    ip       = "127.0.0.1"
  }

  ports {
    internal = 8080
    external = var.ai_port
    ip       = "::1"
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
}

