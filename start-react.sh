#!/bin/bash

# Start React Frontend
echo "🚀 Starting React Frontend..."
echo ""
echo "This will start:"
echo "  - React Frontend (port 3000)"
echo "  - Go Backend API (port 4001)"
echo "  - AI Assistant (port 8080)"
echo "  - PostgreSQL Database (port 5433)"
echo ""

docker compose --profile react-frontend up

echo ""
echo "✅ React Frontend is running!"
echo ""
echo "Access the app:"
echo "  React Frontend: http://localhost:3000"
echo "  Go Backend API: http://localhost:4001"
echo "  AI Assistant:   http://localhost:8080"
echo ""

