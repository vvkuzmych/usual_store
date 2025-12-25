#!/bin/bash

# Start Go Frontend
echo "🚀 Starting Go Frontend..."
echo ""
echo "This will start:"
echo "  - Go Frontend (port 4000)"
echo "  - Go Backend API (port 4001)"
echo "  - AI Assistant (port 8080)"
echo "  - PostgreSQL Database (port 5433)"
echo ""

docker compose --profile go-frontend up

echo ""
echo "✅ Go Frontend is running!"
echo ""
echo "Access the app:"
echo "  Go Frontend:    http://localhost:4000"
echo "  Go Backend API: http://localhost:4001"
echo "  AI Assistant:   http://localhost:8080"
echo ""

