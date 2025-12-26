# Network Policy for Usual Store
# Defines allowed network connections between services

package usualstore.network

import future.keywords.if
import future.keywords.in

# Default deny all connections
default allow := false

# Service tier definitions
service_tiers := {
    "database": "data",
    "kafka": "messaging",
    "zookeeper": "messaging",
    "back-end": "application",
    "support-service": "application",
    "messaging-service": "application",
    "react-frontend": "presentation",
    "typescript-frontend": "presentation",
    "redux-frontend": "presentation",
    "support-frontend": "presentation",
    "jaeger": "observability",
    "opa-server": "governance"
}

# Network communication rules
# Allow if source and target are in valid communication pattern
allow if {
    input.source_service
    input.target_service
    
    # Get service tiers
    source_tier := service_tiers[input.source_service]
    target_tier := service_tiers[input.target_service]
    
    # Check if communication is allowed
    valid_communication(source_tier, target_tier)
}

# Valid tier-to-tier communication patterns
valid_communication("presentation", "application")
valid_communication("application", "data")
valid_communication("application", "messaging")
valid_communication("application", "observability")
valid_communication("messaging", "data")
valid_communication("governance", _)  # Policy service can access all

# Allow internal cluster communication
allow if {
    input.source_network == "usualstore_network"
    input.target_network == "usualstore_network"
}

# Deny direct frontend to database access
deny if {
    service_tiers[input.source_service] == "presentation"
    service_tiers[input.target_service] == "data"
}

# Required ports per service
required_ports := {
    "database": [5432, 5433],
    "kafka": [9092, 9093],
    "zookeeper": [2181],
    "back-end": [4001],
    "support-service": [5001],
    "react-frontend": [3000],
    "typescript-frontend": [3001],
    "redux-frontend": [3002],
    "support-frontend": [3005],
    "jaeger": [16686, 4318],
    "opa-server": [8181]
}

# Validate port configuration
valid_port_config if {
    input.service
    input.port
    input.port in required_ports[input.service]
}

# Network isolation rules
isolated_services := ["database", "kafka", "zookeeper"]

# Services that must be isolated from internet
requires_isolation(service) if {
    service in isolated_services
}

# Validate container network configuration
valid_network_config if {
    input.container_name
    input.networks
    
    # Must be in usualstore network
    "usualstore_network" in input.networks
    
    # Check isolation requirements
    not requires_isolation(input.container_name)
}

# Validate IPv6 configuration
valid_ipv6_config if {
    input.network_config
    input.network_config.ipv6_enabled == true
    startswith(input.network_config.ipv6_subnet, "2001:db8:")
}

# Audit log for policy decisions
audit_log := {
    "allowed": allow,
    "source": input.source_service,
    "target": input.target_service,
    "reason": reason
}

reason := "Valid tier-to-tier communication" if allow
reason := "Direct frontend to database access denied" if deny
reason := "Unknown communication pattern" if not allow

