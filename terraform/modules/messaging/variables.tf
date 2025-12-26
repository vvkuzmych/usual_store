variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "kafka_broker" {
  description = "Kafka broker address"
  type        = string
  default     = "kafka:9092"
}

variable "smtp_host" {
  description = "SMTP server host"
  type        = string
  default     = "sandbox.smtp.mailtrap.io"
}

variable "smtp_port" {
  description = "SMTP server port"
  type        = string
  default     = "2525"
}

variable "smtp_user" {
  description = "SMTP username"
  type        = string
  default     = ""
}

variable "smtp_password" {
  description = "SMTP password"
  type        = string
  sensitive   = true
  default     = ""
}

variable "smtp_from" {
  description = "Email from address"
  type        = string
  default     = "noreply@usualstore.com"
}

variable "kafka_dependency" {
  description = "Kafka service dependency"
  type        = any
}

