# Main Terraform configuration for Usual Store local infrastructure
terraform {
  required_version = ">= 1.0"

  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2"
    }
  }
}

# Docker provider configuration
provider "docker" {
  host = var.docker_host
}

# Create custom network with IPv4 and IPv6 support
resource "docker_network" "usualstore_network" {
  name   = "usualstore_network"
  driver = "bridge"

  ipam_config {
    subnet  = var.ipv4_subnet
    gateway = var.ipv4_gateway
  }

  ipam_config {
    subnet  = var.ipv6_subnet
    gateway = var.ipv6_gateway
  }

  ipv6     = true
  internal = false

  labels {
    label = "com.usualstore.network"
    value = "main"
  }
}

# Create volumes
resource "docker_volume" "db_data" {
  name = "usualstore_db_data"

  labels {
    label = "com.usualstore.volume"
    value = "database"
  }
}

resource "docker_volume" "kafka_data" {
  name = "usualstore_kafka_data"

  labels {
    label = "com.usualstore.volume"
    value = "kafka"
  }
}

resource "docker_volume" "zookeeper_data" {
  name = "usualstore_zookeeper_data"

  labels {
    label = "com.usualstore.volume"
    value = "zookeeper"
  }
}

# Import modules
module "database" {
  source = "./modules/database"

  network_id = docker_network.usualstore_network.id
  volume_id  = docker_volume.db_data.name

  postgres_user     = var.postgres_user
  postgres_password = var.postgres_password
  postgres_db       = var.postgres_db
}

module "kafka_stack" {
  source = "./modules/kafka"

  network_id          = docker_network.usualstore_network.id
  kafka_volume_id     = docker_volume.kafka_data.name
  zookeeper_volume_id = docker_volume.zookeeper_data.name

  depends_on = [docker_network.usualstore_network]
}

module "backend_api" {
  source = "./modules/backend"

  network_id = docker_network.usualstore_network.id

  api_port             = var.api_port
  database_dsn         = module.database.connection_string
  stripe_key           = var.stripe_key
  stripe_secret        = var.stripe_secret
  openai_api_key       = var.openai_api_key
  enable_ai_assistant  = var.enable_ai_assistant
  database_dependency  = module.database.container_id

  depends_on = [module.database]
}

module "support_service" {
  source = "./modules/support"

  network_id = docker_network.usualstore_network.id

  support_port        = var.support_port
  database_dsn        = module.database.connection_string
  database_dependency = module.database.container_id

  depends_on = [module.database]
}

module "ai_assistant" {
  source = "./modules/ai-assistant"

  network_id = docker_network.usualstore_network.id

  ai_port        = var.ai_port
  database_dsn   = module.database.connection_string
  openai_api_key = var.openai_api_key

  depends_on = [module.database]
}

module "frontends" {
  source = "./modules/frontends"

  network_id = docker_network.usualstore_network.id

  react_port      = var.react_port
  typescript_port = var.typescript_port
  redux_port      = var.redux_port
  support_ui_port = var.support_ui_port

  api_url     = "http://back-end:${var.api_port}"
  support_url = "http://support-service:${var.support_port}"

  depends_on = [module.backend_api, module.support_service, module.ai_assistant]
}

module "messaging_service" {
  source = "./modules/messaging"

  network_id       = docker_network.usualstore_network.id
  kafka_broker     = "kafka:9092"
  kafka_dependency = module.kafka_stack.kafka_container_id

  depends_on = [module.kafka_stack]
}

module "observability" {
  source = "./modules/observability"

  network_id = docker_network.usualstore_network.id
}

# Policy enforcement with OPA
module "policies" {
  source = "./modules/policies"

  network_id = docker_network.usualstore_network.id
}

# API Gateway (Kong)
module "api_gateway" {
  source = "./modules/api-gateway"

  network_name  = docker_network.usualstore_network.name
  gateway_port  = 8000
  admin_port    = 8001
  backend_port  = var.api_port

  depends_on = [module.backend_api, module.ai_assistant, module.support_service]
}

# Outputs
output "api_gateway_url" {
  value       = module.api_gateway.gateway_url
  description = "API Gateway URL - use this for all API requests"
}

output "api_gateway_admin_url" {
  value       = module.api_gateway.admin_url
  description = "Kong Admin API URL for configuration"
}
