# Frontend Applications Module

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# React Frontend
resource "docker_image" "react_frontend" {
  name         = "usual_store-react-frontend:latest"
  keep_locally = false

  build {
    context    = "${path.root}/../react-frontend"
    dockerfile = "Dockerfile"
    tag        = ["usual_store-react-frontend:latest"]
  }
}

resource "docker_container" "react_frontend" {
  name  = "usualstore-react-frontend"
  image = docker_image.react_frontend.image_id

  restart = "unless-stopped"

  env = [
    "REACT_APP_API_URL=${var.api_url}",
    "NODE_ENV=production"
  ]

  ports {
    internal = 3000
    external = var.react_port
    ip       = "127.0.0.1"
  }

  ports {
    internal = 3000
    external = var.react_port
    ip       = "::1"
  }

  networks_advanced {
    name = var.network_id
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3000"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "react-frontend"
  }

  labels {
    label = "com.usualstore.tier"
    value = "frontend"
  }
}

# TypeScript Frontend
resource "docker_image" "typescript_frontend" {
  name         = "usual_store-typescript-frontend:latest"
  keep_locally = false

  build {
    context    = "${path.root}/../typescript-frontend"
    dockerfile = "Dockerfile"
    tag        = ["usual_store-typescript-frontend:latest"]
  }
}

resource "docker_container" "typescript_frontend" {
  name  = "usualstore-typescript-frontend"
  image = docker_image.typescript_frontend.image_id

  restart = "unless-stopped"

  env = [
    "VITE_API_URL=${var.api_url}",
    "NODE_ENV=production"
  ]

  ports {
    internal = 3001
    external = var.typescript_port
    ip       = "0.0.0.0"
  }

  networks_advanced {
    name = var.network_id
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3001"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "typescript-frontend"
  }

  labels {
    label = "com.usualstore.tier"
    value = "frontend"
  }
}

# Redux Frontend
resource "docker_image" "redux_frontend" {
  name         = "usual_store-redux-frontend:latest"
  keep_locally = false

  build {
    context    = "${path.root}/../react-redux-frontend"
    dockerfile = "Dockerfile"
    tag        = ["usual_store-redux-frontend:latest"]
  }
}

resource "docker_container" "redux_frontend" {
  name  = "usualstore-redux-frontend"
  image = docker_image.redux_frontend.image_id

  restart = "unless-stopped"

  env = [
    "REACT_APP_API_URL=${var.api_url}",
    "NODE_ENV=production"
  ]

  ports {
    internal = 3002
    external = var.redux_port
    ip       = "127.0.0.1"
  }

  ports {
    internal = 3002
    external = var.redux_port
    ip       = "::1"
  }

  networks_advanced {
    name = var.network_id
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3002/"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "redux-frontend"
  }

  labels {
    label = "com.usualstore.tier"
    value = "frontend"
  }
}

# Support UI Frontend
resource "docker_image" "support_frontend" {
  name         = "usual_store-support-frontend:latest"
  keep_locally = false

  build {
    context    = "${path.root}/../support-frontend"
    dockerfile = "Dockerfile"
    tag        = ["usual_store-support-frontend:latest"]
  }
}

resource "docker_container" "support_frontend" {
  name  = "usualstore-support-frontend"
  image = docker_image.support_frontend.image_id

  restart = "unless-stopped"

  env = [
    "REACT_APP_API_URL=${var.api_url}",
    "REACT_APP_SUPPORT_URL=${var.support_url}",
    "NODE_ENV=production"
  ]

  ports {
    internal = 3005
    external = var.support_ui_port
    ip       = "127.0.0.1"
  }

  ports {
    internal = 3005
    external = var.support_ui_port
    ip       = "::1"
  }

  networks_advanced {
    name = var.network_id
  }

  healthcheck {
    test     = ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3005"]
    interval = "30s"
    timeout  = "3s"
    retries  = 3
  }

  labels {
    label = "com.usualstore.service"
    value = "support-frontend"
  }

  labels {
    label = "com.usualstore.tier"
    value = "frontend"
  }
}
