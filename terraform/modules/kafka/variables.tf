variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "kafka_volume_id" {
  description = "Docker volume name for Kafka data"
  type        = string
}

variable "zookeeper_volume_id" {
  description = "Docker volume name for Zookeeper data"
  type        = string
}

