#!/bin/bash

# Start TypeScript Frontend + Backend Services
# This script starts the TypeScript-based Vite frontend along with necessary backend services

echo "🚀 Starting Usual Store - TypeScript Frontend Edition"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Services to start:"
echo "  ✓ Database (PostgreSQL)"
echo "  ✓ Backend API (Go)"
echo "  ✓ AI Assistant (Go)"
echo "  ✓ Invoice Service (Go)"
echo "  ✓ TypeScript Frontend (Vite + React)"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Navigate to project root
cd "$(dirname "$0")/.." || exit

# Stop any running services
echo "🛑 Stopping existing services..."
docker compose --profile typescript-frontend down 2>/dev/null

# Start services with TypeScript frontend profile
echo ""
echo "▶️  Starting services..."
docker compose --profile typescript-frontend up -d

# Wait for services to be ready
echo ""
echo "⏳ Waiting for services to start..."
sleep 10

# Check service status
echo ""
echo "📊 Service Status:"
docker compose ps

# Display access information
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Services Started Successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🌐 Access URLs:"
echo "  TypeScript Frontend: http://localhost:3001"
echo "  Backend API:         http://localhost:4001"
echo "  AI Assistant:        http://localhost:8080"
echo "  Database:            localhost:5433"
echo ""
echo "🔐 Demo Credentials:"
echo "  Email:    admin@example.com"
echo "  Password: qwerty"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 To stop services:"
echo "  docker compose --profile typescript-frontend down"
echo ""
echo "📋 To view logs:"
echo "  docker compose logs -f typescript-frontend"
echo ""

