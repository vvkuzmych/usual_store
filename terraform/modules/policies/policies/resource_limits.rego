# Resource Limits Policy for Usual Store
# Enforces CPU, memory, and disk limits for containers

package usualstore.resources

import future.keywords.if
import future.keywords.in

# Default resource limits by tier
default_limits := {
    "data": {
        "cpu": "2",
        "memory": "2g",
        "disk": "20g"
    },
    "messaging": {
        "cpu": "1",
        "memory": "1g",
        "disk": "10g"
    },
    "application": {
        "cpu": "1",
        "memory": "512m",
        "disk": "5g"
    },
    "presentation": {
        "cpu": "0.5",
        "memory": "256m",
        "disk": "2g"
    },
    "observability": {
        "cpu": "1",
        "memory": "1g",
        "disk": "10g"
    },
    "governance": {
        "cpu": "0.5",
        "memory": "256m",
        "disk": "1g"
    }
}

# Service to tier mapping (from network policy)
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

# Validate container resources
valid_resources if {
    input.container_name
    input.resources
    
    # Get expected tier
    tier := service_tiers[input.container_name]
    limits := default_limits[tier]
    
    # Check CPU limit
    parse_cpu(input.resources.cpu) <= parse_cpu(limits.cpu)
    
    # Check memory limit
    parse_memory(input.resources.memory) <= parse_memory(limits.memory)
}

# Helper to parse CPU values (cores)
parse_cpu(cpu) := to_number(trim_suffix(cpu, "m")) / 1000 if {
    endswith(cpu, "m")
}

parse_cpu(cpu) := to_number(cpu) if {
    not endswith(cpu, "m")
}

# Helper to parse memory values (bytes)
parse_memory(mem) := to_number(trim_suffix(mem, "g")) * 1024 * 1024 * 1024 if {
    endswith(mem, "g")
}

parse_memory(mem) := to_number(trim_suffix(mem, "m")) * 1024 * 1024 if {
    endswith(mem, "m")
}

parse_memory(mem) := to_number(trim_suffix(mem, "k")) * 1024 if {
    endswith(mem, "k")
}

parse_memory(mem) := to_number(mem) if {
    not endswith(mem, "g")
    not endswith(mem, "m")
    not endswith(mem, "k")
}

# Require resource limits for all containers
requires_limits if {
    input.container_name
    not input.resources
}

# Warning: container exceeds recommended limits
warning_exceeds_limits := msg if {
    input.container_name
    tier := service_tiers[input.container_name]
    limits := default_limits[tier]
    
    parse_cpu(input.resources.cpu) > parse_cpu(limits.cpu)
    
    msg := sprintf("Container %s exceeds CPU limit: %s > %s", 
        [input.container_name, input.resources.cpu, limits.cpu])
}

# Validate volume sizes
valid_volume_size if {
    input.volume_name
    input.size_gb
    
    # Database volumes can be larger
    input.volume_name == "db_data"
    input.size_gb <= 100
}

valid_volume_size if {
    input.volume_name
    input.size_gb
    
    # Other volumes
    input.volume_name != "db_data"
    input.size_gb <= 20
}

# Enforce health check timeouts
valid_health_check if {
    input.health_check
    
    # Interval should be reasonable
    parse_duration(input.health_check.interval) >= 5
    parse_duration(input.health_check.interval) <= 60
    
    # Timeout should be less than interval
    parse_duration(input.health_check.timeout) < parse_duration(input.health_check.interval)
    
    # Retries should be reasonable
    input.health_check.retries >= 1
    input.health_check.retries <= 10
}

# Helper to parse duration strings (seconds)
parse_duration(duration) := to_number(trim_suffix(duration, "s"))

# Calculate resource utilization
utilization := {
    "cpu_percent": (parse_cpu(input.resources.cpu) / parse_cpu(default_limits[service_tiers[input.container_name]].cpu)) * 100,
    "memory_percent": (parse_memory(input.resources.memory) / parse_memory(default_limits[service_tiers[input.container_name]].memory)) * 100
}

# Recommendations for resource optimization
recommendations contains msg if {
    utilization.cpu_percent < 50
    msg := sprintf("Container %s CPU is underutilized: %.1f%%. Consider reducing allocation.", 
        [input.container_name, utilization.cpu_percent])
}

recommendations contains msg if {
    utilization.memory_percent < 50
    msg := sprintf("Container %s memory is underutilized: %.1f%%. Consider reducing allocation.", 
        [input.container_name, utilization.memory_percent])
}

recommendations contains msg if {
    utilization.cpu_percent > 90
    msg := sprintf("Container %s CPU is highly utilized: %.1f%%. Consider increasing allocation.", 
        [input.container_name, utilization.cpu_percent])
}

recommendations contains msg if {
    utilization.memory_percent > 90
    msg := sprintf("Container %s memory is highly utilized: %.1f%%. Consider increasing allocation.", 
        [input.container_name, utilization.memory_percent])
}

