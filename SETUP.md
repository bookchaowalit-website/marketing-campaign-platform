# Marketing SaaS Platform - Complete Setup Guide

## Overview

This is a comprehensive Marketing SaaS platform built with:
- **Frontend**: Next.js 14 + TypeScript + Material-UI
- **Backend 1**: Go (Golang) with Hexagonal Architecture
- **Backend 2**: Python FastAPI for ML/Analytics
- **Communication**: gRPC + REST
- **Database**: PostgreSQL 15

## Features

### Marketing Tools
- ✅ Campaign Management (SEO, SEM, Social, Email, Display)
- ✅ A/B Testing Dashboard
- ✅ Competitor Analysis
- ✅ SEO & SEM Analytics
- ✅ Marketing Research
- ✅ Performance Dashboards
- ✅ Predictive Analytics (ML-powered)

## Quick Start (5 minutes)

### Option 1: Automated Setup

```bash
# Run the setup script
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### Option 2: Manual Setup

#### 1. Start Database

```bash
docker-compose -f docker/docker-compose.yml up -d db
```

#### 2. Install Dependencies

```bash
# Root
npm install

# Go
cd apps/api-go && go mod download

# Python
cd apps/api-python
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

#### 3. Generate gRPC Code

```bash
npm run proto:generate
```

#### 4. Start All Services

```bash
npm run dev
```

## Access Points

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | Next.js Web App |
| Go API | http://localhost:8080 | Core Marketing Services |
| Python API | http://localhost:8081 | Analytics & ML |
| Python API Docs | http://localhost:8081/docs | Swagger UI |
| PostgreSQL | localhost:5432 | Database |

## Project Structure

```
marketing-saas-platform/
├── apps/
│   ├── web/                    # Next.js Frontend
│   │   ├── src/
│   │   │   ├── app/           # Next.js App Router
│   │   │   ├── components/    # React Components
│   │   │   ├── lib/           # Utilities
│   │   │   └── types/         # TypeScript Types
│   │   └── package.json
│   │
│   ├── api-go/                # Go Backend (Hexagonal)
│   │   ├── cmd/
│   │   │   └── main.go       # Entry point
│   │   ├── internal/
│   │   │   ├── domain/       # Business entities
│   │   │   ├── application/  # Business logic
│   │   │   ├── infrastructure/ # External adapters
│   │   │   └── delivery/     # HTTP/gRPC handlers
│   │   └── go.mod
│   │
│   └── api-python/           # Python FastAPI Backend
│       ├── app/
│       │   ├── main.py      # Entry point
│       │   ├── api/         # Routes
│       │   ├── models/      # SQLAlchemy models
│       │   ├── schemas/     # Pydantic schemas
│       │   └── services/    # Business logic
│       └── requirements.txt
│
├── packages/
│   └── grpc-protos/         # Shared gRPC definitions
│       └── marketing.proto
│
├── docker/
│   ├── docker-compose.yml
│   └── migrations/
│       └── 001_initial_schema.sql
│
└── docs/
    ├── architecture.md
    ├── development.md
    └── api.md
```

## Development Commands

### All Services
```bash
npm run dev          # Start all services
npm run build        # Build all services
```

### Frontend
```bash
cd apps/web
npm run dev         # Development server
npm run build       # Production build
npm run lint        # Run ESLint
```

### Go Backend
```bash
cd apps/api-go
make run           # Run application
make build         # Build binary
make test          # Run tests
```

### Python Backend
```bash
cd apps/api-python
make run          # Run application
make dev          # Development mode (auto-reload)
make test         # Run tests
make migrate      # Run database migrations
```

### Docker
```bash
npm run docker:up    # Start all containers
npm run docker:down  # Stop all containers
```

## API Examples

### Create Campaign (Go API)

```bash
curl -X POST http://localhost:8080/api/v1/campaigns \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Summer Sale",
    "type": "SEM",
    "budget": 5000,
    "start_date": "2024-07-01"
  }'
```

### Get SEO Metrics (Python API)

```bash
curl "http://localhost:8081/api/v1/seo/metrics?domain=example.com"
```

### Analyze Competitor (Python API)

```bash
curl "http://localhost:8081/api/v1/competitors/{id}/analyze"
```

## Database Schema

Key tables:
- `users` - User accounts
- `organizations` - Multi-tenant orgs
- `campaigns` - Marketing campaigns
- `ab_tests` - A/B tests
- `competitors` - Competitor tracking
- `seo_keywords` - SEO tracking
- `analytics_daily` - Analytics aggregates

## Environment Variables

### Frontend (.env)
```env
NEXT_PUBLIC_API_GO_URL=http://localhost:8080
NEXT_PUBLIC_API_PYTHON_URL=http://localhost:8081
```

### Go (.env)
```env
PORT=8080
GRPC_PORT=9090
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/marketing_saas
JWT_SECRET=your-secret-key
```

### Python (.env)
```env
PORT=8081
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/marketing_saas
SECRET_KEY=your-secret-key
```

## Testing

```bash
# Frontend
cd apps/web && npm test

# Go
cd apps/api-go && go test ./...

# Python
cd apps/api-python && pytest
```

## Production Deployment

### Build Docker Images

```bash
docker-compose -f docker/docker-compose.yml build
```

### Deploy

```bash
docker-compose -f docker/docker-compose.yml up -d
```

## Troubleshooting

### Port Already in Use
```bash
lsof -i :3000  # Find process
kill -9 <PID>  # Kill process
```

### Database Connection Failed
```bash
docker-compose -f docker/docker-compose.yml restart db
```

### gRPC Errors
```bash
# Regenerate protobuf code
npm run proto:generate
```

## Documentation

- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)
- [API Documentation](docs/api.md)

## Support

For issues or questions:
1. Check documentation in `/docs`
2. Review existing issues
3. Create a new issue with details

## License

MIT

---

**Happy Marketing! 🚀**
