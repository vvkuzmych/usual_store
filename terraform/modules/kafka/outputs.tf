output "zookeeper_container_id" {
  description = "Zookeeper container ID"
  value       = docker_container.zookeeper.id
}

output "kafka_container_id" {
  description = "Kafka container ID"
  value       = docker_container.kafka.id
}

output "kafka_ui_container_id" {
  description = "Kafka UI container ID"
  value       = docker_container.kafka_ui.id
}

output "kafka_broker" {
  description = "Kafka broker address"
  value       = "kafka:9092"
}

output "status" {
  description = "Kafka stack status"
  value       = "running"
}

