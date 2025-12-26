# Frontend Module Variables

variable "network_id" {
  description = "Docker network ID"
  type        = string
}

variable "react_port" {
  description = "React frontend external port"
  type        = number
  default     = 3000
}

variable "typescript_port" {
  description = "TypeScript frontend external port"
  type        = number
  default     = 3001
}

variable "redux_port" {
  description = "Redux frontend external port"
  type        = number
  default     = 3002
}

variable "support_ui_port" {
  description = "Support UI frontend external port"
  type        = number
  default     = 3005
}

variable "api_url" {
  description = "Backend API URL for frontends"
  type        = string
}

variable "support_url" {
  description = "Support service URL for frontends"
  type        = string
}

