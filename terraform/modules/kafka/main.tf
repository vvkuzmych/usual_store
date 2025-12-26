# Kafka Stack module - Zookeeper, Kafka, and Kafka UI

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

# Pull images
resource "docker_image" "zookeeper" {
  name         = "confluentinc/cp-zookeeper:7.5.0"
  keep_locally = false
}

resource "docker_image" "kafka" {
  name         = "confluentinc/cp-kafka:7.5.0"
  keep_locally = false
}

resource "docker_image" "kafka_ui" {
  name         = "provectuslabs/kafka-ui:latest"
  keep_locally = false
}

# Zookeeper container
resource "docker_container" "zookeeper" {
  name  = "usualstore-zookeeper"
  image = docker_image.zookeeper.image_id

  restart = "unless-stopped"

  env = [
    "ZOOKEEPER_CLIENT_PORT=2181",
    "ZOOKEEPER_TICK_TIME=2000"
  ]

  ports {
    internal = 2181
    external = 2181
  }

  networks_advanced {
    name = var.network_id
  }

  volumes {
    volume_name    = var.zookeeper_volume_id
    container_path = "/var/lib/zookeeper/data"
  }

  healthcheck {
    test     = ["CMD", "nc", "-z", "localhost", "2181"]
    interval = "10s"
    timeout  = "5s"
    retries  = 5
  }

  labels {
    label = "com.usualstore.service"
    value = "zookeeper"
  }

  labels {
    label = "com.usualstore.tier"
    value = "messaging"
  }
}

# Kafka container
resource "docker_container" "kafka" {
  name  = "usualstore-kafka"
  image = docker_image.kafka.image_id

  restart = "unless-stopped"

  env = [
    "KAFKA_BROKER_ID=1",
    "KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181",
    "KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT",
    "KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:9093",
    "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
    "KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1",
    "KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1",
    "KAFKA_AUTO_CREATE_TOPICS_ENABLE=true"
  ]

  ports {
    internal = 9092
    external = 9092
  }

  ports {
    internal = 9093
    external = 9093
  }

  networks_advanced {
    name = var.network_id
  }

  volumes {
    volume_name    = var.kafka_volume_id
    container_path = "/var/lib/kafka/data"
  }

  healthcheck {
    test         = ["CMD", "kafka-broker-api-versions", "--bootstrap-server", "localhost:9092"]
    interval     = "10s"
    timeout      = "10s"
    retries      = 5
    start_period = "20s"
  }

  depends_on = [docker_container.zookeeper]

  labels {
    label = "com.usualstore.service"
    value = "kafka"
  }

  labels {
    label = "com.usualstore.tier"
    value = "messaging"
  }
}

# Kafka UI container
resource "docker_container" "kafka_ui" {
  name  = "usualstore-kafka-ui"
  image = docker_image.kafka_ui.image_id

  restart = "unless-stopped"

  env = [
    "KAFKA_CLUSTERS_0_NAME=usual-store-cluster",
    "KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092",
    "KAFKA_CLUSTERS_0_ZOOKEEPER=zookeeper:2181"
  ]

  ports {
    internal = 8080
    external = 8090
  }

  networks_advanced {
    name = var.network_id
  }

  depends_on = [docker_container.kafka]

  labels {
    label = "com.usualstore.service"
    value = "kafka-ui"
  }

  labels {
    label = "com.usualstore.tier"
    value = "messaging"
  }
}

