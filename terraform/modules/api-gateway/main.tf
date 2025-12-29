# API Gateway Module for UsualStore
# Routes requests to backend services with authentication and rate limiting

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Kong API Gateway (Open Source)
resource "docker_image" "kong" {
  name = "kong:3.4"
}

# PostgreSQL for Kong (Kong requires a database)
resource "docker_container" "kong_db" {
  name  = "kong-database"
  image = "postgres:14-alpine"
  
  env = [
    "POSTGRES_USER=kong",
    "POSTGRES_DB=kong",
    "POSTGRES_PASSWORD=kong_password"
  ]
  
  networks_advanced {
    name = var.network_name
  }
  
  volumes {
    volume_name    = "kong_data"
    container_path = "/var/lib/postgresql/data"
  }
}

# Kong Database Migration (run once)
resource "docker_container" "kong_migration" {
  name  = "kong-migration"
  image = docker_image.kong.image_id
  
  command = ["kong", "migrations", "bootstrap"]
  
  env = [
    "KONG_DATABASE=postgres",
    "KONG_PG_HOST=kong-database",
    "KONG_PG_USER=kong",
    "KONG_PG_PASSWORD=kong_password"
  ]
  
  networks_advanced {
    name = var.network_name
  }
  
  depends_on = [docker_container.kong_db]
  
  # This container should run once and exit
  must_run = false
}

# Kong API Gateway Container
resource "docker_container" "kong" {
  name  = "api-gateway"
  image = docker_image.kong.image_id
  
  env = [
    "KONG_DATABASE=postgres",
    "KONG_PG_HOST=kong-database",
    "KONG_PG_USER=kong",
    "KONG_PG_PASSWORD=kong_password",
    "KONG_PROXY_ACCESS_LOG=/dev/stdout",
    "KONG_ADMIN_ACCESS_LOG=/dev/stdout",
    "KONG_PROXY_ERROR_LOG=/dev/stderr",
    "KONG_ADMIN_ERROR_LOG=/dev/stderr",
    "KONG_ADMIN_LISTEN=0.0.0.0:8001",
    "KONG_PROXY_LISTEN=0.0.0.0:8000, 0.0.0.0:8443 ssl"
  ]
  
  ports {
    internal = 8000  # Proxy port (HTTP)
    external = var.gateway_port
  }
  
  ports {
    internal = 8443  # Proxy port (HTTPS)
    external = var.gateway_ssl_port
  }
  
  ports {
    internal = 8001  # Admin API
    external = var.admin_port
  }
  
  networks_advanced {
    name = var.network_name
  }
  
  depends_on = [
    docker_container.kong_db,
    docker_container.kong_migration
  ]
}

# Configure Kong using local-exec (alternative to using Kong Admin API directly)
resource "null_resource" "kong_config" {
  depends_on = [docker_container.kong]
  
  # Wait for Kong to be ready
  provisioner "local-exec" {
    command = <<-EOT
      echo "Waiting for Kong to be ready..."
      for i in {1..30}; do
        if curl -s http://localhost:${var.admin_port} > /dev/null; then
          echo "Kong is ready!"
          break
        fi
        echo "Waiting... ($i/30)"
        sleep 2
      done
    EOT
  }
  
  # Create service for backend API
  provisioner "local-exec" {
    command = <<-EOT
      curl -i -X POST http://localhost:${var.admin_port}/services/ \
        --data name=backend-api \
        --data url='http://back-end:${var.backend_port}'
    EOT
  }
  
  # Create route for backend API
  provisioner "local-exec" {
    command = <<-EOT
      curl -i -X POST http://localhost:${var.admin_port}/services/backend-api/routes \
        --data 'paths[]=/api' \
        --data 'strip_path=false'
    EOT
  }
  
  # Enable rate limiting
  provisioner "local-exec" {
    command = <<-EOT
      curl -i -X POST http://localhost:${var.admin_port}/services/backend-api/plugins \
        --data name=rate-limiting \
        --data config.minute=100 \
        --data config.hour=1000
    EOT
  }
  
  # Enable CORS
  provisioner "local-exec" {
    command = <<-EOT
      curl -i -X POST http://localhost:${var.admin_port}/services/backend-api/plugins \
        --data name=cors \
        --data 'config.origins=*' \
        --data 'config.methods=GET,POST,PUT,DELETE,OPTIONS' \
        --data 'config.headers=Accept,Authorization,Content-Type' \
        --data 'config.exposed_headers=X-Auth-Token' \
        --data 'config.credentials=true' \
        --data 'config.max_age=3600'
    EOT
  }
  
  # Enable request/response logging
  provisioner "local-exec" {
    command = <<-EOT
      curl -i -X POST http://localhost:${var.admin_port}/services/backend-api/plugins \
        --data name=file-log \
        --data 'config.path=/tmp/kong-access.log'
    EOT
  }
}
