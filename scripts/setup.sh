#!/bin/bash

# Marketing SaaS Platform - Quick Start Script
# This script helps you set up and run the development environment

set -e

echo "🚀 Marketing SaaS Platform - Quick Start"
echo "========================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js is not installed. Please install Node.js 18+${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Node.js:${NC} $(node --version)"

# Check npm
if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ npm:${NC} $(npm --version)"

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}⚠️  Go is not installed (optional for frontend-only development)${NC}"
else
    echo -e "${GREEN}✓ Go:${NC} $(go version)"
fi

# Check Python
if ! command -v python3 &> /dev/null; then
    echo -e "${YELLOW}⚠️  Python 3 is not installed (optional for frontend-only development)${NC}"
else
    echo -e "${GREEN}✓ Python:${NC} $(python3 --version)"
fi

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker is not installed (recommended for database)${NC}"
else
    echo -e "${GREEN}✓ Docker:${NC} $(docker --version)"
fi

echo ""
echo "📦 Installing dependencies..."

# Install root dependencies
echo "Installing root dependencies..."
npm install

# Install Go dependencies
if command -v go &> /dev/null; then
    echo "Installing Go dependencies..."
    cd api-go
    go mod download
    cd ..
fi

# Install Python dependencies
if command -v python3 &> /dev/null; then
    echo "Installing Python dependencies..."
    cd api-python
    if [ ! -d "venv" ]; then
        python3 -m venv venv
    fi
    source venv/bin/activate
    pip install -r requirements.txt
    deactivate
    cd ..
fi

echo ""
echo "🏗️  Generating gRPC code..."
if [ -f "packages/grpc-protos/Makefile" ]; then
    make -C packages/grpc-protos generate || echo -e "${YELLOW}⚠️  gRPC generation skipped (protoc may not be installed)${NC}"
fi

echo ""
echo "🐳 Starting Docker services..."
if command -v docker &> /dev/null; then
    docker-compose -f docker/docker-compose.yml up -d db
    echo "Waiting for database to be ready..."
    sleep 5
else
    echo -e "${YELLOW}⚠️  Docker not available. Please start PostgreSQL manually.${NC}"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "📖 Next steps:"
echo "   1. Start development servers:"
echo "      npm run dev"
echo ""
echo "   2. Or start services individually:"
echo "      npm run dev:web      # Frontend (http://localhost:3000)"
echo "      npm run dev:go       # Go API (http://localhost:8080)"
echo "      npm run dev:python   # Python API (http://localhost:8081)"
echo ""
echo "   3. View API documentation:"
echo "      Python API: http://localhost:8081/docs"
echo ""
echo "   4. Check the documentation:"
echo "      docs/architecture.md"
echo "      docs/development.md"
echo "      docs/api.md"
echo ""
echo "🎉 Happy coding!"
