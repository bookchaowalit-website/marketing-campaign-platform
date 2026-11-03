# 06 — C4 Container

Container ตาม `docker/docker-compose.yml` และ Dockerfile ของแต่ละ service (ยังไม่มี hosting จริง ดู 16)

```mermaid
C4Container
  title BookMarketing containers
  Person(user, "Marketer or admin")
  System_Boundary(bm, "BookMarketing monorepo") {
    Container(web, "web", "Next.js 16, React 19, MUI 5, standalone node server.js :3000", "Dashboard; rewrites /api/go and /api/python; /api/health")
    Container(go, "api-go", "Go 1.26, Gin, golang-jwt, :8080 HTTP and gRPC GRPC_PORT", "Auth, campaigns, email, SMS, PPC, affiliates, social, budget, calendar, approvals")
    Container(py, "api-python", "Python 3.11, FastAPI, SQLAlchemy, :8081", "SEO, competitors, analytics, gRPC client")
    ContainerDb(db, "PostgreSQL 15", "postgres:15-alpine", "marketing_saas database shared by both APIs")
  }
  System_Ext(ext, "Email, SMS, social, marketplace, SEO providers")
  System_Ext(sentry, "Sentry")
  Rel(user, web, "HTTPS")
  Rel(web, go, "rewrite /api/go/:path to /api/:path", "HTTP, Bearer JWT")
  Rel(web, py, "rewrite /api/python/:path to /api/:path", "HTTP, Bearer JWT")
  Rel(go, db, "database/sql lib/pq", "DATABASE_URL")
  Rel(py, db, "SQLAlchemy", "DATABASE_URL")
  Rel(py, go, "gRPC insecure channel", "GRPC_SERVER_URL")
  Rel(go, ext, "HTTPS REST")
  Rel(py, ext, "httpx")
  Rel(web, sentry, "errors only")
```

Sources: `docker/docker-compose.yml`, `web/Dockerfile`, `api-go/Dockerfile`, `api-python/Dockerfile`, `web/next.config.ts`, `web/src/lib/api-client.ts`, `api-go/cmd/main.go`, `api-go/internal/infrastructure/database/postgres.go`, `api-python/app/core/database.py`, `api-python/app/core/config.py`, `api-python/app/api/grpc_client.py`, `web/package.json`, `api-go/go.mod`

ข้อสังเกต (จาก config ไม่ได้รัน)

- Port gRPC ไม่ตรงกัน: compose ตั้ง `api-go` `GRPC_PORT: 9091` แต่ `api-python` `GRPC_SERVER_URL: api-go:9090`; default ในโค้ด Go คือ `9090` ส่วน Python คือ `localhost:9091`
- `api-python/app/api/grpc_client.py` import `app.grpc_client.campaign_pb2` ที่ไม่อยู่ใน Git (`.gitignore` ตัด `*_pb2.py`) และ `packages/grpc-protos/Makefile` generate ไปที่ `app/grpc/marketing_pb2.py` ซึ่งเป็นคนละชื่อ: `app.main` import ไม่ผ่านจาก checkout (SPEC-001 §5)
- `packages/shared` (`@bookmarketing/shared-types`) เป็น TypeScript types ที่ web ใช้ผ่าน `tsconfig.json` paths ไม่ใช่ container

## Open questions

- gRPC port ที่ถูกต้องคือ 9090 หรือ 9091
- Python service จะ generate gRPC stubs ตอน build Docker หรือ commit ไว้
