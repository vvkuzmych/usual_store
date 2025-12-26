output "opa_container_id" {
  description = "OPA server container ID"
  value       = docker_container.opa_server.id
}

output "enforcer_container_id" {
  description = "Policy enforcer container ID"
  value       = docker_container.policy_enforcer.id
}

output "policy_server_url" {
  description = "OPA policy server URL"
  value       = "http://localhost:8181"
}

output "status" {
  description = "Policy system status"
  value       = "running"
}

