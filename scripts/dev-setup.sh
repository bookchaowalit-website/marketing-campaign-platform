#!/bin/bash

# Quick Development Setup Script
# This sets up the development environment WITHOUT Docker (except database)

set -e

echo "🚀 Marketing SaaS - Local Development Setup"
echo "============================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

cd /home/bookchaowalit/book_dev/personal/saas/bookmarketing

# Step 1: Start Database
echo "📦 Step 1: Starting PostgreSQL database..."
docker compose -f docker/docker-compose.dev.yml up -d
echo "Waiting for database to be ready..."
sleep 5
echo -e "${GREEN}✓ Database started${NC}"
echo ""

# Step 2: Setup Frontend
echo "📦 Step 2: Setting up Next.js frontend..."
cd web
if [ ! -d "node_modules" ]; then
    echo "Installing npm packages (this may take a few minutes)..."
    npm install
    echo -e "${GREEN}✓ Frontend dependencies installed${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend dependencies already exist${NC}"
fi
cd ../..
echo ""

# Step 3: Setup Python
echo "📦 Step 3: Setting up Python backend..."
cd api-python
if [ ! -d "venv" ]; then
    python3 -m venv venv
    echo "Created Python virtual environment"
fi
source venv/bin/activate
pip install -q --upgrade pip
pip install -r requirements.txt
echo -e "${GREEN}✓ Python dependencies installed${NC}"
deactivate
cd ../..
echo ""

# Step 4: Check Go
echo "📦 Step 4: Checking Go installation..."
if command -v go &> /dev/null; then
    cd api-go
    go mod download
    echo -e "${GREEN}✓ Go dependencies downloaded${NC}"
    cd ..
else
    echo -e "${YELLOW}⚠️  Go is not installed. Go API will only run in Docker.${NC}"
    echo "   Install Go from: https://go.dev/dl/"
fi
echo ""

# Summary
echo "✅ Setup Complete!"
echo ""
echo "📖 How to run the application:"
echo ""
echo "Option 1: Run all services together (recommended)"
echo "  npm run dev"
echo ""
echo "Option 2: Run services in separate terminals"
echo "  Terminal 1 (Frontend):"
echo "    cd apps/web && npm run dev"
echo ""
echo "  Terminal 2 (Go API):"
echo "    cd apps/api-go && make run"
echo ""
echo "  Terminal 3 (Python API):"
echo "    cd apps/api-python && source venv/bin/activate && make dev"
echo ""
echo "🌐 Access Points:"
echo "  Frontend:    http://localhost:3000"
echo "  Go API:      http://localhost:8080"
echo "  Python API:  http://localhost:8081/docs"
echo "  PostgreSQL:  localhost:5432"
echo ""
echo "🛑 To stop:"
echo "  - Press Ctrl+C in each terminal"
echo "  - Run: npm run docker:down  (to stop database)"
echo ""
