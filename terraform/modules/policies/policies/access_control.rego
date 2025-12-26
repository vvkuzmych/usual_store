# Access Control Policy for Usual Store
# Manages user and service access permissions

package usualstore.access

import future.keywords.if
import future.keywords.in

# User roles
roles := {
    "super_admin": {
        "level": 100,
        "permissions": ["all"]
    },
    "admin": {
        "level": 80,
        "permissions": [
            "manage_users",
            "manage_containers",
            "view_logs",
            "manage_support",
            "view_policies"
        ]
    },
    "supporter": {
        "level": 50,
        "permissions": [
            "view_support_tickets",
            "respond_to_tickets",
            "view_user_info"
        ]
    },
    "user": {
        "level": 10,
        "permissions": [
            "view_products",
            "make_purchases",
            "contact_support"
        ]
    }
}

# Service accounts
service_accounts := {
    "back-end": {
        "permissions": [
            "database_read",
            "database_write",
            "kafka_publish",
            "user_authentication"
        ]
    },
    "support-service": {
        "permissions": [
            "database_read",
            "database_write",
            "websocket_management"
        ]
    },
    "messaging-service": {
        "permissions": [
            "kafka_consume",
            "smtp_send"
        ]
    },
    "policy-enforcer": {
        "permissions": [
            "docker_inspect",
            "docker_events",
            "opa_query"
        ]
    }
}

# Check if user has permission
user_has_permission(user, permission) if {
    user.role
    role_data := roles[user.role]
    "all" in role_data.permissions
}

user_has_permission(user, permission) if {
    user.role
    role_data := roles[user.role]
    permission in role_data.permissions
}

# Check if service has permission
service_has_permission(service, permission) if {
    service_data := service_accounts[service]
    permission in service_data.permissions
}

# API endpoint access control
api_endpoints := {
    "/api/users": {
        "GET": ["admin", "super_admin"],
        "POST": ["admin", "super_admin"],
        "DELETE": ["super_admin"]
    },
    "/api/products": {
        "GET": ["user", "supporter", "admin", "super_admin"],
        "POST": ["admin", "super_admin"],
        "PUT": ["admin", "super_admin"],
        "DELETE": ["super_admin"]
    },
    "/api/support/tickets": {
        "GET": ["supporter", "admin", "super_admin"],
        "POST": ["user", "supporter", "admin", "super_admin"],
        "PUT": ["supporter", "admin", "super_admin"]
    },
    "/api/policies": {
        "GET": ["admin", "super_admin"],
        "POST": ["super_admin"],
        "PUT": ["super_admin"],
        "DELETE": ["super_admin"]
    }
}

# Allow API access
allow_api_access if {
    input.user
    input.endpoint
    input.method
    
    endpoint_config := api_endpoints[input.endpoint]
    allowed_roles := endpoint_config[input.method]
    
    input.user.role in allowed_roles
}

# Container management permissions
allow_container_action if {
    input.user
    input.action
    input.container
    
    # Check user has manage_containers permission
    user_has_permission(input.user, "manage_containers")
    
    # Validate action
    valid_container_action(input.action)
}

valid_container_action("start")
valid_container_action("stop")
valid_container_action("restart")
valid_container_action("inspect")
valid_container_action("logs")

# Deny dangerous actions
deny_container_action if {
    input.action in ["remove", "kill"]
    input.user.role != "super_admin"
}

# Database access control
allow_database_access if {
    input.service
    input.operation
    
    # Check service has appropriate permission
    operation_permission := database_operation_permission(input.operation)
    service_has_permission(input.service, operation_permission)
}

database_operation_permission("SELECT") := "database_read"
database_operation_permission("INSERT") := "database_write"
database_operation_permission("UPDATE") := "database_write"
database_operation_permission("DELETE") := "database_write"

# Kafka topic access control
kafka_topics := {
    "email-notifications": {
        "publishers": ["back-end", "support-service"],
        "consumers": ["messaging-service"]
    },
    "user-events": {
        "publishers": ["back-end"],
        "consumers": ["analytics-service", "messaging-service"]
    },
    "support-messages": {
        "publishers": ["support-service"],
        "consumers": ["back-end"]
    }
}

allow_kafka_publish if {
    input.service
    input.topic
    
    topic_config := kafka_topics[input.topic]
    input.service in topic_config.publishers
}

allow_kafka_consume if {
    input.service
    input.topic
    
    topic_config := kafka_topics[input.topic]
    input.service in topic_config.consumers
}

# Policy management access
allow_policy_management if {
    input.user
    input.policy_action
    
    # Only super_admin can manage policies
    input.user.role == "super_admin"
    
    valid_policy_action(input.policy_action)
}

valid_policy_action("create")
valid_policy_action("update")
valid_policy_action("delete")
valid_policy_action("evaluate")

# Support dashboard access
allow_support_dashboard if {
    input.user
    
    # Support, admin, and super_admin can access
    input.user.role in ["supporter", "admin", "super_admin"]
}

# User management access
allow_user_management if {
    input.user
    input.target_user
    input.action
    
    # Check permission
    user_has_permission(input.user, "manage_users")
    
    # Super admin can manage all users
    input.user.role == "super_admin"
}

allow_user_management if {
    input.user
    input.target_user
    input.action
    
    # Admin can manage users below their level
    input.user.role == "admin"
    user_has_permission(input.user, "manage_users")
    
    target_level := roles[input.target_user.role].level
    user_level := roles[input.user.role].level
    
    target_level < user_level
}

# Prevent deletion of last super_admin
deny_user_deletion if {
    input.action == "delete"
    input.target_user.role == "super_admin"
    count_super_admins <= 1
}

count_super_admins := count([user | user := input.all_users[_]; user.role == "super_admin"])

# Audit log entry
audit_entry := {
    "timestamp": time.now_ns(),
    "user": input.user.email,
    "action": input.action,
    "resource": input.resource,
    "allowed": allow_action,
    "reason": decision_reason
}

allow_action if {
    allow_api_access
}

allow_action if {
    allow_container_action
    not deny_container_action
}

allow_action if {
    allow_database_access
}

allow_action if {
    allow_kafka_publish
}

allow_action if {
    allow_kafka_consume
}

allow_action if {
    allow_policy_management
}

allow_action if {
    allow_support_dashboard
}

allow_action if {
    allow_user_management
    not deny_user_deletion
}

decision_reason := "Access granted by role permissions" if allow_action
decision_reason := "Access denied - insufficient permissions" if not allow_action
decision_reason := "Access denied - last super_admin cannot be deleted" if deny_user_deletion
decision_reason := "Access denied - dangerous container action" if deny_container_action

# Rate limiting (requests per minute)
rate_limits := {
    "user": 60,
    "supporter": 120,
    "admin": 300,
    "super_admin": 1000
}

within_rate_limit if {
    input.user
    input.request_count
    input.time_window_minutes == 1
    
    limit := rate_limits[input.user.role]
    input.request_count <= limit
}

# Session management
valid_session if {
    input.session
    input.current_time
    
    # Session must not be expired
    input.session.expires_at > input.current_time
    
    # Session must belong to the user
    input.session.user_id == input.user.id
    
    # Session must not be revoked
    not input.session.revoked
}

# Multi-factor authentication requirement
requires_mfa if {
    input.user
    input.action
    
    # Super admin always requires MFA
    input.user.role == "super_admin"
}

requires_mfa if {
    input.action in ["delete_user", "modify_policies", "access_database"]
}

allow_with_mfa if {
    requires_mfa
    input.mfa_verified == true
}

allow_with_mfa if {
    not requires_mfa
}

