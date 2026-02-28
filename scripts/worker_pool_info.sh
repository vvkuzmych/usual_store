#!/bin/bash

# Worker Pool Information
# Shows current configuration and implementation details (no tests)

PROJECT_ROOT="/Users/vkuzm/Projects/usual_store"
cd "$PROJECT_ROOT" || exit 1

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Worker Pool Information           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Implementation stats
echo -e "${YELLOW}📦 Implementation:${NC}"
echo "   Package: internal/workerpool/"
echo "   Lines of code:"
wc -l internal/workerpool/*.go | tail -1 | awk '{printf "     Total: %s lines\n", $1}'
echo "   Files:"
ls internal/workerpool/ | grep -v test | while read f; do echo "     - $f"; done
echo ""

# Test stats
echo -e "${YELLOW}🧪 Tests:${NC}"
test_count=$(grep -c "^func Test" internal/workerpool/workerpool_test.go)
subtest_count=$(grep -c "t.Run(" internal/workerpool/workerpool_test.go)
echo "   Test functions: $test_count"
echo "   Subtests: $subtest_count"
wc -l internal/workerpool/workerpool_test.go | awk '{printf "   Test code: %s lines\n", $1}'
echo ""

# Configuration
echo -e "${YELLOW}⚙️  Configuration:${NC}"
echo "   Default Settings:"
grep "getEnvInt.*EMAIL_WORKER" cmd/messaging-service/main.go | while read line; do
    var=$(echo "$line" | grep -o "EMAIL_WORKER_[A-Z]*")
    val=$(echo "$line" | grep -o "[0-9]*)" | tr -d ')')
    echo "     $var = $val"
done
echo ""

# Integration points
echo -e "${YELLOW}🔗 Integration:${NC}"
echo "   Used in:"
grep -l "workerpool" internal/messaging/*.go cmd/messaging-service/*.go 2>/dev/null | while read f; do
    echo "     - $f"
done
echo ""

# Usage info
echo -e "${YELLOW}💡 Usage:${NC}"
echo "   Start service:"
echo "     ./bin/messaging-service"
echo ""
echo "   Custom configuration:"
echo "     ./bin/messaging-service -workers 20 -buffer 200"
echo ""
echo "   Environment variables:"
echo "     EMAIL_WORKER_COUNT=15 ./bin/messaging-service"
echo ""

# Performance estimates
echo -e "${YELLOW}📊 Expected Performance:${NC}"
default_workers=$(grep "EMAIL_WORKER_COUNT" cmd/messaging-service/main.go | grep -o "[0-9]*)" | head -1 | tr -d ')')
echo "   With $default_workers workers:"
echo "     Sequential: ~2-10 emails/second"
echo "     With pool:  ~20-100 emails/second"
echo "     Speedup:    ~${default_workers}x improvement"
echo ""

# Health check reminder
echo -e "${GREEN}✓ For detailed health check, run:${NC}"
echo "  ./scripts/check_worker_pool.sh"
echo ""
