# 16 — Deployment

**สถานะจริง (2026-11-03): ไม่มี deployment.** GitHub deployments ของ repo = 0, `Deploy smoke` ไม่เคยรัน, `deploy-smoke.yml` ระบุ "This repository has no production host yet" และ SPEC-001 §5 "No deployment" config เดียวที่บอก topology ได้คือ `docker/docker-compose.yml`

## Topology ตาม `docker/docker-compose.yml` (local / self-host)

```mermaid
flowchart TB
  subgraph HOST["Docker host, network marketing-saas-network bridge"]
    WEB["marketing-saas-web: build ../web, node server.js, :3000 to 3000, build args NEXT_PUBLIC_API_GO_URL and NEXT_PUBLIC_API_PYTHON_URL = localhost"]
    GO["marketing-saas-go: build ../api-go, :8080 and 9091, JWT_SECRET placeholder, GRPC_PORT 9091"]
    PY["marketing-saas-python: build ../api-python, uvicorn :8081, GRPC_SERVER_URL api-go:9090"]
    DB[("marketing-saas-db: postgres:15-alpine :5432, postgres/postgres, init from docker/migrations")]
  end
  USER["Browser on the host"] -->|"localhost:3000"| WEB
  WEB -->|"rewrite to NEXT_PUBLIC_API_GO_URL, baked at build"| GO
  WEB -->|"rewrite to NEXT_PUBLIC_API_PYTHON_URL"| PY
  GO -->|"DATABASE_URL db:5432"| DB
  PY -->|"DATABASE_URL db:5432"| DB
  PY -.->|"gRPC api-go:9090, Go listens on 9091"| GO
  GO -->|"depends_on healthy"| DB
  PY -->|"depends_on healthy"| DB
```

Sources: `docker/docker-compose.yml`, `docker/docker-compose.dev.yml`, `web/Dockerfile`, `api-go/Dockerfile`, `api-python/Dockerfile`, `web/next.config.ts`, `api-go/cmd/main.go`

ข้อสังเกต: rewrites ของ web อ่าน `NEXT_PUBLIC_API_GO_URL` ซึ่ง compose ส่งเป็น `http://localhost:8080` ตอน build; ภายใน container `localhost` คือ container web เอง (env `API_GO_URL: http://api-go:8080` ใน compose ไม่มีโค้ดอ่าน) ดังนั้น proxy ใน compose น่าจะต่อ API ไม่ได้ (อ่านจาก config ยังไม่ได้รันยืนยัน)

## เป้าหมายที่โค้ดเตรียมไว้ (ยังไม่ได้ตั้ง)

```mermaid
flowchart LR
  GH["GitHub main"] -.->|"Vercel Git integration, not connected"| V["Vercel project for web, unknown"]
  V -.->|"deployment_status Production"| DS["deploy-smoke.yml"]
  DS -.->|"needs vars.PRODUCTION_URL, not set"| URL["PRODUCTION_URL"]
  DNS["bookmarketing.bookchaowalit.com"] -.->|"CNAME cname.vercel-dns.com, TLS handshake fails"| V
  API["api-go and api-python host: not decided"] -.-> V
  NEON["Database host: not decided"] -.-> API
```

Sources: `.github/workflows/deploy-smoke.yml`, `scripts/deploy-smoke.sh`, `web/src/app/api/health/route.ts` (`VERCEL_GIT_COMMIT_SHA`), `specs/001-next16-ci-sentry/spec.md` §5, `gh api repos/.../deployments` (0), `getent hosts bookmarketing.bookchaowalit.com` + `curl https://bookmarketing.bookchaowalit.com/` (2026-11-03: resolves to Vercel, curl error 35 TLS)

SPEC-001 (2026-10-22) บันทึกว่า `bookmarketing.bookchaowalit.com` ไม่ resolve; วันที่ทำ pack นี้ DNS ชี้ Vercel แล้วแต่ไม่มี certificate/project ตอบ และ `marketing.bookchaowalit.com` เป็นของ repo `marketing-bookmarketing-web` (Vercel project `bookchaowalit-bookmarketing-frontend`)

## Open questions

- ใครเพิ่ม DNS `bookmarketing` → Vercel และผูกกับ Vercel project ใด
- host ของ api-go / api-python / PostgreSQL เมื่อ deploy (Railway? Neon?) และ env ชุดไหน
- จะเลิกใช้ `NEXT_PUBLIC_*` สำหรับ rewrites ฝั่ง server แล้วใช้ env runtime แทนหรือไม่
