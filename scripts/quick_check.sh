#!/bin/bash

# Quick Worker Pool Check
# Fast verification that worker pool is functional

PROJECT_ROOT="/Users/vkuzm/Projects/usual_store"
cd "$PROJECT_ROOT" || exit 1

echo "🔍 Quick Worker Pool Check..."
echo ""

# Quick test
echo "▶ Running tests..."
if go test ./internal/workerpool/... -run TestWorkerPool_JobProcessing/BasicFunctionality > /dev/null 2>&1; then
    echo "✅ Tests passing"
else
    echo "❌ Tests failing"
    exit 1
fi

# Quick build
echo "▶ Building service..."
if go build -o /tmp/msg-test cmd/messaging-service/*.go 2>/dev/null; then
    echo "✅ Build successful"
    rm /tmp/msg-test
else
    echo "❌ Build failed"
    exit 1
fi

# Quick coverage
echo "▶ Checking coverage..."
coverage=$(go test ./internal/workerpool/... -cover 2>&1 | grep "coverage:" | awk '{print $5}')
echo "   Coverage: $coverage"

echo ""
echo "✅ Worker pool is healthy!"
