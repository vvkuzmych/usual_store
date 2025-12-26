# Observability module - Jaeger tracing

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

variable "network_id" { type = string }

resource "docker_image" "jaeger" {
  name         = "jaegertracing/all-in-one:latest"
  keep_locally = false
}

resource "docker_container" "jaeger" {
  name  = "usualstore-jaeger"
  image = docker_image.jaeger.image_id

  restart = "unless-stopped"

  env = ["COLLECTOR_OTLP_ENABLED=true"]

  ports {
    internal = 16686
    external = 16686
    ip       = "127.0.0.1"
  }

  ports {
    internal = 4318
    external = 4318
    ip       = "127.0.0.1"
  }

  networks_advanced {
    name = var.network_id
  }

  healthcheck {
    test         = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:16686"]
    interval     = "30s"
    timeout      = "3s"
    retries      = 3
    start_period = "5s"
  }

  labels {
    label = "com.usualstore.service"
    value = "jaeger"
  }
}

output "status" { value = "running" }
output "jaeger_ui_url" { value = "http://localhost:16686" }

