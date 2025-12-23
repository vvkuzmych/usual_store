#!/bin/bash

# ==============================================================================
# IPv6 Connectivity Test Script for Usual Store
# ==============================================================================
# This script tests IPv6 connectivity for all services in the Usual Store
# application stack.
#
# Usage: ./scripts/test-ipv6.sh
# ==============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

print_header() {
    echo ""
    echo "=================================================="
    print_status "$BLUE" "$1"
    echo "=================================================="
}

print_success() {
    print_status "$GREEN" "✅ $1"
}

print_error() {
    print_status "$RED" "❌ $1"
}

print_warning() {
    print_status "$YELLOW" "⚠️  $1"
}

print_info() {
    print_status "$BLUE" "ℹ️  $1"
}

# Check if Docker is running
check_docker() {
    print_header "Checking Docker Status"
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker Desktop."
        exit 1
    fi
    print_success "Docker is running"
}

# Check if docker-compose.yml exists
check_compose_file() {
    print_header "Checking Docker Compose Configuration"
    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found. Please run this script from the project root."
        exit 1
    fi
    print_success "docker-compose.yml found"
}

# Check if services are running
check_services() {
    print_header "Checking Running Services"
    
    services=("database" "back-end" "front-end" "invoice")
    all_running=true
    
    for service in "${services[@]}"; do
        if docker compose ps --status running | grep -q "$service"; then
            print_success "$service is running"
        else
            print_warning "$service is not running"
            all_running=false
        fi
    done
    
    if [ "$all_running" = false ]; then
        print_warning "Not all services are running. Start them with: docker compose up -d"
        read -p "Do you want to start services now? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker compose up -d
            sleep 5
        else
            exit 1
        fi
    fi
}

# Check network IPv6 configuration
check_network_ipv6() {
    print_header "Checking Docker Network IPv6 Configuration"
    
    network_name=$(docker compose config | grep -A 5 "networks:" | grep -v "networks:" | head -1 | tr -d ' ' | cut -d: -f1)
    full_network_name="${PWD##*/}_${network_name}"
    
    if docker network inspect "$full_network_name" 2>/dev/null | grep -q '"EnableIPv6": true'; then
        print_success "IPv6 is enabled on network: $full_network_name"
        
        # Get IPv6 subnet
        ipv6_subnet=$(docker network inspect "$full_network_name" | grep -A 1 '"Subnet"' | grep "fd00" | awk -F'"' '{print $4}')
        if [ -n "$ipv6_subnet" ]; then
            print_info "IPv6 Subnet: $ipv6_subnet"
        fi
    else
        print_error "IPv6 is not enabled on network: $full_network_name"
        print_info "Run 'docker compose down && docker compose up -d' to recreate the network"
        return 1
    fi
}

# Check container IPv6 addresses
check_container_ipv6() {
    print_header "Checking Container IPv6 Addresses"
    
    services=("database" "back-end" "front-end" "invoice")
    
    for service in "${services[@]}"; do
        container_name="${PWD##*/}-${service}-1"
        if docker ps --format '{{.Names}}' | grep -q "$container_name"; then
            ipv6_addr=$(docker exec "$container_name" ip -6 addr show eth0 2>/dev/null | grep "fd00" | awk '{print $2}' | cut -d/ -f1)
            if [ -n "$ipv6_addr" ]; then
                print_success "$service: $ipv6_addr"
            else
                print_warning "$service: No IPv6 address found"
            fi
        fi
    done
}

# Test PostgreSQL IPv6 connectivity from host
test_postgres_ipv6_host() {
    print_header "Testing PostgreSQL IPv6 Connectivity (Host)"
    
    # Test IPv6 localhost
    print_info "Testing connection to [::1]:5432..."
    if command -v psql > /dev/null 2>&1; then
        if psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable" -c "SELECT version();" > /dev/null 2>&1; then
            print_success "PostgreSQL IPv6 connection successful from host"
        else
            print_error "Failed to connect to PostgreSQL via IPv6"
            print_info "Make sure PostgreSQL port is bound to [::1]:5432"
        fi
    else
        print_warning "psql not installed. Skipping host connection test"
        print_info "Install: brew install postgresql"
    fi
    
    # Test IPv4 localhost (backward compatibility)
    print_info "Testing connection to 127.0.0.1:5432..."
    if command -v psql > /dev/null 2>&1; then
        if psql "postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable" -c "SELECT version();" > /dev/null 2>&1; then
            print_success "PostgreSQL IPv4 connection successful from host"
        else
            print_warning "Failed to connect to PostgreSQL via IPv4"
        fi
    fi
}

# Test PostgreSQL connectivity from containers
test_postgres_containers() {
    print_header "Testing PostgreSQL Connectivity (Containers)"
    
    services=("back-end" "front-end")
    
    for service in "${services[@]}"; do
        container_name="${PWD##*/}-${service}-1"
        if docker ps --format '{{.Names}}' | grep -q "$container_name"; then
            print_info "Testing from $service..."
            
            # Test DNS resolution
            if docker exec "$container_name" getent hosts database > /dev/null 2>&1; then
                print_success "$service can resolve 'database' hostname"
                
                # Show resolved addresses
                addresses=$(docker exec "$container_name" getent hosts database | awk '{print $1}')
                print_info "Resolved addresses: $addresses"
            else
                print_error "$service cannot resolve 'database' hostname"
            fi
            
            # Test ping (IPv4)
            if docker exec "$container_name" ping -c 1 -W 2 database > /dev/null 2>&1; then
                print_success "$service can ping database (IPv4)"
            else
                print_warning "$service cannot ping database (IPv4)"
            fi
            
            # Test ping6 (IPv6) - if available
            if docker exec "$container_name" sh -c "command -v ping6" > /dev/null 2>&1; then
                if docker exec "$container_name" ping6 -c 1 -W 2 database > /dev/null 2>&1; then
                    print_success "$service can ping database (IPv6)"
                else
                    print_warning "$service cannot ping database (IPv6)"
                fi
            fi
        fi
    done
}

# Check PostgreSQL configuration
check_postgres_config() {
    print_header "Checking PostgreSQL IPv6 Configuration"
    
    container_name="${PWD##*/}-database-1"
    
    if docker ps --format '{{.Names}}' | grep -q "$container_name"; then
        # Check listen_addresses
        listen_addr=$(docker exec "$container_name" psql -U postgres -d usualstore -t -c "SHOW listen_addresses;" 2>/dev/null | tr -d ' ')
        print_info "listen_addresses: $listen_addr"
        
        if [ "$listen_addr" == "*" ] || [ "$listen_addr" == "0.0.0.0,::" ]; then
            print_success "PostgreSQL is listening on all addresses (IPv4 and IPv6)"
        else
            print_warning "PostgreSQL listen_addresses: $listen_addr"
        fi
        
        # Check active connections
        print_info "Active connections:"
        docker exec "$container_name" ss -tln 2>/dev/null | grep ":5432" || print_warning "Could not query listening ports"
    fi
}

# Test HTTP services
test_http_services() {
    print_header "Testing HTTP Services"
    
    # Get ports from docker-compose
    print_info "Testing front-end service..."
    if curl -s "http://localhost:4000" > /dev/null 2>&1; then
        print_success "Front-end is accessible on IPv4 localhost"
    else
        print_warning "Front-end not accessible on IPv4 localhost"
    fi
    
    if curl -s "http://[::1]:4000" > /dev/null 2>&1; then
        print_success "Front-end is accessible on IPv6 localhost"
    else
        print_warning "Front-end not accessible on IPv6 localhost"
    fi
    
    print_info "Testing back-end API..."
    if curl -s "http://localhost:4001" > /dev/null 2>&1; then
        print_success "Back-end API is accessible on IPv4 localhost"
    else
        print_warning "Back-end API not accessible on IPv4 localhost"
    fi
    
    if curl -s "http://[::1]:4001" > /dev/null 2>&1; then
        print_success "Back-end API is accessible on IPv6 localhost"
    else
        print_warning "Back-end API not accessible on IPv6 localhost"
    fi
}

# System IPv6 check
check_system_ipv6() {
    print_header "Checking System IPv6 Configuration"
    
    # Test IPv6 loopback
    if ping6 -c 1 ::1 > /dev/null 2>&1; then
        print_success "System IPv6 loopback is working"
    else
        print_error "System IPv6 loopback is not working"
        print_info "IPv6 may not be enabled on your system"
    fi
    
    # Check Docker daemon IPv6 config
    print_info "Checking Docker daemon IPv6 configuration..."
    if docker info 2>/dev/null | grep -q "IPv6"; then
        print_success "Docker daemon has IPv6 configuration"
    else
        print_warning "Docker daemon may not have IPv6 enabled globally"
        print_info "This is OK if you're using per-network IPv6 configuration"
    fi
}

# Main execution
main() {
    clear
    print_header "🔍 Usual Store IPv6 Connectivity Test"
    print_info "Starting comprehensive IPv6 connectivity tests..."
    
    check_docker
    check_compose_file
    check_system_ipv6
    check_services
    check_network_ipv6
    check_container_ipv6
    check_postgres_config
    test_postgres_ipv6_host
    test_postgres_containers
    test_http_services
    
    print_header "✅ IPv6 Connectivity Test Complete"
    print_info "Review the results above for any issues"
    print_info ""
    print_info "For more information, see: IPv6-SETUP.md"
}

# Run main function
main

