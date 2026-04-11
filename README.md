# Marketing SaaS Platform

A comprehensive Marketing SaaS platform built with a polyglot microservices architecture.

## 🚀 Tech Stack

### Frontend
- **Next.js 14** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Material-UI** - Component library

### Backend Services
- **Go (Golang)** - Hexagonal architecture for core marketing services
- **Python FastAPI** - ML/AI-powered marketing analytics and research
- **gRPC** - High-performance inter-service communication
- **PostgreSQL** - Primary database

## 📁 Project Structure

```
marketing-saas-platform/
├── apps/
│   ├── web/              # Next.js frontend
│   ├── api-go/           # Go backend (hexagonal)
│   └── api-python/       # Python FastAPI backend
├── packages/
│   ├── grpc-protos/      # Shared gRPC protobuf definitions
│   └── shared/           # Shared utilities and types
├── docker/               # Docker configurations
└── scripts/              # Development scripts
```

## 🛠️ Features

### Marketing Tools
- **Marketing Research** - Market analysis and insights
- **Competitor Analysis** - Track and analyze competitors
- **SEO Tools** - Search engine optimization
- **SEM Tools** - Search engine marketing
- **Ads Campaign Management** - Multi-platform ad campaigns
- **A/B Testing** - Experimentation and optimization
- **Analytics Dashboard** - Real-time marketing metrics

## 🏃 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- Python 3.11+
- Docker & Docker Compose
- PostgreSQL 15+

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd marketing-saas-platform
```

2. **Install dependencies**
```bash
npm install
```

3. **Start Docker services (PostgreSQL)**
```bash
npm run docker:up
```

4. **Generate gRPC protos**
```bash
npm run proto:generate
```

5. **Run development servers**
```bash
npm run dev
```

Or run individual services:
- Frontend: `npm run dev:web`
- Go API: `npm run dev:go`
- Python API: `npm run dev:python`

## 📖 Documentation

- [Architecture Overview](./docs/architecture.md)
- [API Documentation](./docs/api.md)
- [Development Guide](./docs/development.md)

## 🏗️ Architecture

### Hexagonal Architecture (Go)
```
┌─────────────────┐
│   HTTP/gRPC     │
│     Adapter     │
├─────────────────┤
│    Application  │
│     Service     │
├─────────────────┤
│     Domain      │
│    (Entities)   │
├─────────────────┤
│  Infrastructure │
│  (Repository)   │
└─────────────────┘
```

### Service Communication
```
┌─────────────┐      gRPC      ┌──────────────┐
│   Next.js   │ ◄──────────►   │   Go API     │
│   Frontend  │      gRPC      │  (Core Svcs) │
└─────────────┘                └──────────────┘
       │                              │
       │         gRPC                 │
       ▼                              ▼
┌─────────────┐                ┌──────────────┐
│   Python    │                │  PostgreSQL  │
│   FastAPI   │                │   Database   │
│  (ML/Analytics)              └──────────────┘
└─────────────┘
```

## 📄 License

MIT
