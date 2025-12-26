# Security Policy for Usual Store
# Enforces security best practices for containers and services

package usualstore.security

import future.keywords.if
import future.keywords.in

# Default deny privileged containers
default allow_privileged := false

# Container security requirements
secure_container if {
    input.container_name
    input.security_opts
    
    # Must not run as privileged
    not input.privileged
    
    # Must not run as root (except database)
    not runs_as_root
    
    # Must have read-only root filesystem (where applicable)
    appropriate_filesystem_mode
}

# Check if container runs as root
runs_as_root if {
    input.user == "root"
}

runs_as_root if {
    input.user == "0"
}

runs_as_root if {
    not input.user
    input.container_name != "database"
}

# Containers that can have writable filesystem
writable_filesystem_allowed := {
    "database",
    "kafka",
    "zookeeper",
    "messaging-service"
}

appropriate_filesystem_mode if {
    input.container_name in writable_filesystem_allowed
}

appropriate_filesystem_mode if {
    input.read_only_rootfs == true
    not input.container_name in writable_filesystem_allowed
}

# Required security options
required_security_opts := {
    "no-new-privileges:true",
    "seccomp:default"
}

valid_security_opts if {
    input.security_opts
    count({opt | opt := input.security_opts[_]; opt in required_security_opts}) == count(required_security_opts)
}

# Capability restrictions
# Only allow necessary capabilities
allowed_capabilities := {
    "database": ["CHOWN", "FOWNER", "DAC_OVERRIDE"],
    "kafka": ["NET_BIND_SERVICE"],
    "zookeeper": ["NET_BIND_SERVICE"]
}

default_denied_capabilities := [
    "NET_ADMIN",
    "SYS_ADMIN",
    "SYS_MODULE",
    "SYS_RAWIO",
    "SYS_PTRACE",
    "SYS_BOOT",
    "MAC_ADMIN",
    "MAC_OVERRIDE"
]

valid_capabilities if {
    input.container_name
    input.capabilities
    
    # Get allowed caps for this service
    allowed := allowed_capabilities[input.container_name]
    
    # Check all requested capabilities are allowed
    count({cap | cap := input.capabilities[_]; not cap in allowed}) == 0
    
    # Check no denied capabilities
    count({cap | cap := input.capabilities[_]; cap in default_denied_capabilities}) == 0
}

# Image security
trusted_registries := [
    "docker.io",
    "gcr.io",
    "ghcr.io",
    "registry.hub.docker.com"
]

valid_image if {
    input.image
    
    # Must be from trusted registry or locally built
    registry_trusted(input.image)
    
    # Must have specific tag (not :latest in production)
    not uses_latest_tag
}

registry_trusted(image) if {
    some registry in trusted_registries
    startswith(image, registry)
}

registry_trusted(image) if {
    # Local images (built with docker build)
    not contains(image, "/")
}

uses_latest_tag if {
    endswith(input.image, ":latest")
    input.environment == "production"
}

# Environment variable security
# Sensitive data should not be in plain text
sensitive_env_vars := [
    "PASSWORD",
    "SECRET",
    "TOKEN",
    "KEY",
    "PRIVATE"
]

insecure_env_vars contains var if {
    some env in input.env_vars
    some sensitive in sensitive_env_vars
    contains(upper(env.name), sensitive)
    not uses_secure_source(env)
    var := env
}

uses_secure_source(env) if {
    # Should use secrets or files, not plain values
    startswith(env.value, "/run/secrets/")
}

uses_secure_source(env) if {
    startswith(env.value, "file://")
}

# Network security
secure_network_config if {
    input.network_config
    
    # Must use custom network, not default bridge
    input.network_config.name != "bridge"
    
    # Must have network isolation
    input.network_config.internal == true
}

# Port exposure security
# Only expose necessary ports to host
exposed_ports := {
    "react-frontend": [3000],
    "typescript-frontend": [3001],
    "redux-frontend": [3002],
    "support-frontend": [3005],
    "back-end": [4001],
    "support-service": [5001],
    "database": [5433],  # External port
    "kafka": [9092, 9093],
    "jaeger": [16686, 4318],
    "kafka-ui": [8090],
    "opa-server": [8181]
}

valid_port_exposure if {
    input.container_name
    input.exposed_ports
    
    allowed := exposed_ports[input.container_name]
    
    # All exposed ports must be in allowed list
    count({port | port := input.exposed_ports[_]; not port in allowed}) == 0
}

# Volume mount security
# Sensitive host paths should not be mounted
dangerous_mounts := [
    "/",
    "/etc",
    "/usr",
    "/bin",
    "/sbin",
    "/boot",
    "/sys"
]

safe_volume_mounts if {
    input.volumes
    
    # Check no dangerous mounts
    count({vol | vol := input.volumes[_]; vol.host_path in dangerous_mounts}) == 0
}

# Only policy-enforcer can mount Docker socket
docker_socket_mount if {
    input.volumes
    some vol in input.volumes
    vol.host_path == "/var/run/docker.sock"
}

allow_docker_socket_mount if {
    docker_socket_mount
    input.container_name == "policy-enforcer"
}

deny_docker_socket_mount if {
    docker_socket_mount
    input.container_name != "policy-enforcer"
}

allow_docker_socket_mount if {
    not deny_docker_socket_mount
}

allow_docker_socket_mount if {
    not docker_socket_mount
}

# Health check security
valid_health_check_security if {
    input.health_check
    
    # Health check should not run as root
    not input.health_check.user == "root"
    
    # Should use minimal commands
    valid_health_check_command
}

valid_health_check_command if {
    input.health_check.test
    allowed_health_commands := ["CMD", "CMD-SHELL", "NONE"]
    input.health_check.test[0] in allowed_health_commands
}

# Security scan results
security_violations contains violation if {
    not secure_container
    violation := {
        "type": "insecure_container",
        "severity": "high",
        "message": "Container does not meet security requirements"
    }
}

security_violations contains violation if {
    count(insecure_env_vars) > 0
    violation := {
        "type": "insecure_env_vars",
        "severity": "critical",
        "message": sprintf("Found %d insecure environment variables", [count(insecure_env_vars)]),
        "details": insecure_env_vars
    }
}

security_violations contains violation if {
    deny_docker_socket_mount
    violation := {
        "type": "unauthorized_docker_socket",
        "severity": "critical",
        "message": "Unauthorized Docker socket mount detected"
    }
}

security_violations contains violation if {
    not safe_volume_mounts
    violation := {
        "type": "dangerous_mount",
        "severity": "critical",
        "message": "Dangerous host path mounted"
    }
}

# Security score calculation
security_score := score if {
    total_checks := 10
    passed_checks := count([
        secure_container,
        valid_security_opts,
        valid_capabilities,
        valid_image,
        count(insecure_env_vars) == 0,
        secure_network_config,
        valid_port_exposure,
        safe_volume_mounts,
        allow_docker_socket_mount,
        valid_health_check_security
    ])
    score := (passed_checks / total_checks) * 100
}

# Compliance check
compliant if {
    security_score >= 80
    count(security_violations) == 0
}

