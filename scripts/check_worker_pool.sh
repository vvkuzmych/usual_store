#!/bin/bash

# Worker Pool Health Check Script
# This script verifies the worker pool implementation is working correctly

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root
PROJECT_ROOT="/Users/vkuzm/Projects/usual_store"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Worker Pool Health Check${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print success
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to print error
error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Function to print info
info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Function to print section header
section() {
    echo ""
    echo -e "${BLUE}▶ $1${NC}"
    echo "-------------------------------------------"
}

# Change to project directory
cd "$PROJECT_ROOT" || error "Failed to change to project directory"

# 1. Check if worker pool package exists
section "1. Checking Worker Pool Package"
if [ -d "internal/workerpool" ]; then
    success "Worker pool package exists"
    ls -la internal/workerpool/
else
    error "Worker pool package not found at internal/workerpool"
fi

# 2. Check required files
section "2. Verifying Required Files"
required_files=(
    "internal/workerpool/workerpool.go"
    "internal/workerpool/job.go"
    "internal/workerpool/workerpool_test.go"
    "internal/messaging/consumer_pool.go"
    "cmd/messaging-service/main.go"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        success "$file exists"
    else
        error "$file not found"
    fi
done

# 3. Run worker pool tests
section "3. Running Worker Pool Tests"
info "Running tests..."
if go test ./internal/workerpool/... -v; then
    success "All tests passed"
else
    error "Tests failed"
fi

# 4. Check test coverage
section "4. Checking Test Coverage"
info "Generating coverage report..."
go test ./internal/workerpool/... -coverprofile=coverage.out > /dev/null 2>&1
coverage=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}')
coverage_num=$(echo "$coverage" | sed 's/%//')

echo "Coverage: $coverage"
# Use awk for float comparison
if awk "BEGIN {exit !($coverage_num >= 90)}"; then
    success "Excellent coverage: $coverage"
elif awk "BEGIN {exit !($coverage_num >= 70)}"; then
    success "Good coverage: $coverage"
else
    error "Coverage too low: $coverage (minimum 70% recommended)"
fi

# 5. Verify build
section "5. Building Messaging Service"
info "Building binary..."
if go build -o bin/messaging-service-test cmd/messaging-service/*.go; then
    success "Binary built successfully"
    size=$(ls -lh bin/messaging-service-test | awk '{print $5}')
    echo "Binary size: $size"
    rm bin/messaging-service-test
else
    error "Build failed"
fi

# 6. Check configuration flags
section "6. Verifying Configuration Flags"
info "Building and checking flags..."
go build -o bin/messaging-service-test cmd/messaging-service/*.go
if ./bin/messaging-service-test -h 2>&1 | grep -q "workers"; then
    success "Worker pool flags available:"
    ./bin/messaging-service-test -h 2>&1 | grep -E "(workers|buffer)" | sed 's/^/  /'
else
    error "Worker pool flags not found"
fi
rm bin/messaging-service-test

# 7. Check for proper imports
section "7. Checking Code Quality"
info "Verifying imports..."
if grep -q "usual_store/internal/workerpool" cmd/messaging-service/main.go; then
    success "Worker pool properly imported in main.go"
else
    error "Worker pool not imported in main.go"
fi

if grep -q "workerPool.*workerpool.Pool" internal/messaging/consumer_pool.go; then
    success "Worker pool properly used in consumer"
else
    error "Worker pool not properly integrated in consumer"
fi

# 8. Summary of implementation
section "8. Implementation Summary"
echo ""
echo "📦 Worker Pool Components:"
echo "  - Package: internal/workerpool/"
echo "  - Implementation: $(wc -l < internal/workerpool/workerpool.go) lines"
echo "  - Tests: $(wc -l < internal/workerpool/workerpool_test.go) lines"
echo "  - Integration: internal/messaging/consumer_pool.go"
echo ""

# 9. Configuration info
section "9. Configuration Options"
echo ""
echo "Environment Variables:"
echo "  EMAIL_WORKER_COUNT   - Number of concurrent workers (default: 10)"
echo "  EMAIL_WORKER_BUFFER  - Job queue buffer size (default: 100)"
echo ""
echo "Command-line Flags:"
echo "  -workers <N>  - Set worker count"
echo "  -buffer <N>   - Set buffer size"
echo ""

# 10. Quick usage example
section "10. Usage Example"
cat << 'EOF'
# Start with default settings (10 workers, 100 buffer)
./bin/messaging-service

# Start with custom settings
./bin/messaging-service -workers 20 -buffer 200

# Use environment variables
EMAIL_WORKER_COUNT=15 EMAIL_WORKER_BUFFER=150 ./bin/messaging-service
EOF

# Final summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ Worker Pool Health Check Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}All checks passed successfully!${NC}"
echo -e "Worker pool is ready for production use."
echo ""

# Cleanup
rm -f coverage.out 2>/dev/null

exit 0
