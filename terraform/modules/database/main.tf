# Database module - PostgreSQL container

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Pull PostgreSQL image
resource "docker_image" "postgres" {
  name         = "postgres:15"
  keep_locally = false
}

# PostgreSQL container
resource "docker_container" "database" {
  name  = "usualstore-database"
  image = docker_image.postgres.image_id

  restart = "unless-stopped"

  env = [
    "POSTGRES_USER=${var.postgres_user}",
    "POSTGRES_PASSWORD=${var.postgres_password}",
    "POSTGRES_DB=${var.postgres_db}"
  ]

  command = [
    "postgres",
    "-c", "listen_addresses=*",
    "-c", "max_connections=200"
  ]

  ports {
    internal = 5432
    external = 5433
    ip       = "127.0.0.1"
  }

  ports {
    internal = 5432
    external = 5433
    ip       = "::1"
  }

  networks_advanced {
    name = var.network_id
  }

  volumes {
    volume_name    = var.volume_id
    container_path = "/var/lib/postgresql/data"
  }

  healthcheck {
    test         = ["CMD", "pg_isready", "-U", var.postgres_user, "-d", var.postgres_db, "-h", "localhost"]
    interval     = "10s"
    timeout      = "5s"
    retries      = 5
    start_period = "10s"
  }

  labels {
    label = "com.usualstore.service"
    value = "database"
  }

  labels {
    label = "com.usualstore.tier"
    value = "data"
  }
}

