# Frontends module - placeholder
variable "network_id" { type = string }
variable "react_port" { type = number }
variable "typescript_port" { type = number }
variable "redux_port" { type = number }
variable "support_ui_port" { type = number }
variable "api_url" { type = string }
variable "support_url" { type = string }

output "status" { value = "configured" }

