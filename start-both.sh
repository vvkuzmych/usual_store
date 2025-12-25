#!/bin/bash

# Start Both Frontends
echo "🚀 Starting BOTH Frontends (Go + React)..."
echo ""
echo "This will start:"
echo "  - Go Frontend (port 4000)"
echo "  - React Frontend (port 3000)"
echo "  - Go Backend API (port 4001)"
echo "  - AI Assistant (port 8080)"
echo "  - PostgreSQL Database (port 5433)"
echo ""

docker compose --profile go-frontend --profile react-frontend up

echo ""
echo "✅ Both frontends are running!"
echo ""
echo "Access the apps:"
echo "  Go Frontend:    http://localhost:4000"
echo "  React Frontend: http://localhost:3000"
echo "  Go Backend API: http://localhost:4001"
echo "  AI Assistant:   http://localhost:8080"
echo ""
echo "Compare them side-by-side! 🎨"
echo ""

