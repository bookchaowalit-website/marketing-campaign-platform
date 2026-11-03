# Core 22 SDLC diagram pack — marketing-campaign-platform (BookMarketing)

ชุดภาพ Core 22 ตาม SDLC Diagram Catalog v2 วาดจากโค้ดและ config ใน repo นี้ (ณ `origin/main` `a4061e5`, 2026-11-03)
repo นี้คือ Tier A interview flagship: polyglot monorepo (`web/` Next.js 16, `api-go/` Go hexagonal, `api-python/` FastAPI, PostgreSQL, gRPC) **ยังไม่มี deployment และไม่มี public URL ที่ลงทะเบียนไว้**

ทุกไฟล์มี `Sources:` ใต้ภาพ และ "Open questions" สำหรับสิ่งที่ยืนยันจากโค้ดไม่ได้ ข้อสังเกตที่มาจากการอ่านโค้ดโดยไม่ได้รันจะเขียนกำกับไว้

| # | Item | Status | File | Sources read |
|---|---|---|---|---|
| 01 | Business process As-Is / To-Be | done | [01-business-process.md](01-business-process.md) | `router.go`, handlers, services, `seo.py`, `analytics_service.py`, `PRODUCT.md` |
| 02 | Use cases | done | [02-use-cases.md](02-use-cases.md) | `router.go`, `middleware.go`, `auth_service.go`, `001_init.up.sql`, `api-python/app/main.py`, `core/auth.py` |
| 03 | User flows | done | [03-user-flows.md](03-user-flows.md) | `web/src/app/**/page.tsx`, `Sidebar.tsx`, `useAuth.tsx`, `web/src/lib/*.ts` |
| 04 | Traceability (RTM) | done | [04-traceability.md](04-traceability.md) | `specs/001-next16-ci-sentry/spec.md`, `PRODUCT.md`, test files, `git log`, `gh run list` |
| 05 | C4 Context | done | [05-c4-context.md](05-c4-context.md) | env gates in `router.go`, provider clients, `seo_data_sources.py`, `sentry.ts`, workflows |
| 06 | C4 Container | done | [06-c4-container.md](06-c4-container.md) | `docker-compose.yml`, Dockerfiles, `next.config.ts`, `main.go`, `config.py`, `grpc_client.py` |
| 07 | C4 Component | done | [07-c4-component.md](07-c4-component.md) | `api-go/internal/**`, `api-python/app/**`, `web/src/**` |
| 08 | ADRs | done (6) | [08-adr/](08-adr/) | `git log`, `README.md`, `PRODUCT.md`, SPEC-001, code and workflows per ADR |
| 09 | ERD | done | [09-erd.md](09-erd.md) | `api-go/migrations/001..016`, repositories, Python models |
| 10 | Sequences | done (4) | [10-sequences.md](10-sequences.md) | auth, SEO, email campaign, social publish code paths |
| 11 | State machines | done | [11-state-machines.md](11-state-machines.md) | status writes in services/repositories + CHECK constraints |
| 12 | API & integrations | done | [12-api-integration.md](12-api-integration.md) | `router.go`, `grpc/server.go`, `marketing.proto`, `api-python/app/api/*.py`, provider clients |
| 13 | Threat model | done | [13-threat-model.md](13-threat-model.md) | middleware, handlers, CORS, migrations, Dockerfiles, compose, `sentry.ts`, `ci.yml` |
| 14 | Backup & DR | done (gap: no backup) | [14-backup-dr.md](14-backup-dr.md) | compose files, `Makefile`, migrations, `RUNBOOK.md`, `.github/workflows/` |
| 15 | CI/CD pipeline | done | [15-cicd-pipeline.md](15-cicd-pipeline.md) | `ci.yml`, `deploy-smoke.yml`, `deploy-smoke.sh`, platform-ci-workflows `87dc5e1c` |
| 16 | Deployment | done (no deployment today) | [16-deployment.md](16-deployment.md) | compose, Dockerfiles, `deploy-smoke.yml`, SPEC-001 §5, GitHub deployments API, DNS/TLS probe |
| 17 | Test coverage matrix | done | [17-test-coverage-matrix.md](17-test-coverage-matrix.md) | `auth_campaign_test.go`, `test_auth.py`, `web/tests/*.test.ts`, `ci.yml` |
| 18 | Git & release workflow | done | [18-git-release-workflow.md](18-git-release-workflow.md) | `git log`, `gh pr list`, branch protection API, workflows |
| 19 | Service dependency map | done | [19-service-dependency-map.md](19-service-dependency-map.md) | `go.mod`, `requirements.in`, `web/package.json`, startup code |
| 20 | Observability | reused + extended | [20-observability.md](20-observability.md) → [`docs/measurement/KPI-PLAN.md`](../measurement/KPI-PLAN.md) | `instrumentation*.ts`, `sentry.ts`, health routes, `gin.Default()` |
| 21 | Incident response | done | [21-incident-response.md](21-incident-response.md) | workflows, `deploy-smoke.sh`, env reads, `Makefile`, `RUNBOOK.md` |
| 22 | DORA | done (0 deployments: 5/5 not computable, CI proxies shown) | [22-dora.md](22-dora.md) | `git log`, `gh run list`, deployments API |

## ADR list

- [ADR-0001](08-adr/ADR-0001-polyglot-monorepo-go-python-nextjs.md) Polyglot monorepo: Next.js + Go hexagonal + Python analytics
- [ADR-0002](08-adr/ADR-0002-shared-postgresql-go-migrations.md) PostgreSQL ร่วมกัน, `api-go/migrations` เป็น schema หลัก
- [ADR-0003](08-adr/ADR-0003-hs256-jwt-issued-by-go-shared-secret.md) HS256 JWT จาก Go, Python verify ด้วย secret เดียวกัน
- [ADR-0004](08-adr/ADR-0004-same-origin-proxy-via-next-rewrites.md) Same-origin proxy ผ่าน Next.js rewrites
- [ADR-0005](08-adr/ADR-0005-providers-opt-in-by-env-with-stubs.md) Providers เปิดด้วย env, ไม่มี = stub
- [ADR-0006](08-adr/ADR-0006-shared-ci-pinned-no-deploy.md) Shared CI pinned v1.1.1, CI ไม่ deploy

## Findings worth the owner's attention

- ไม่มี deployment; `bookmarketing.bookchaowalit.com` ตอนนี้ CNAME ไป Vercel แต่ TLS handshake fail (ดู 16)
- Tenancy: JWT ไม่มี `organization_id` → เฉพาะ unified campaigns ได้ `""` (insert UUID fail); SMS / social / ads fallback เป็น demo org `00000000-0000-0000-0000-000000000001` และ approval / affiliate / classic campaigns ใช้ `defaultOrganizationID` ตรง ๆ — ผู้ใช้ทุกคนจึงเห็นและแก้ข้อมูล demo org ชุดเดียวกัน; classic campaigns ไม่เช็ก owner บน get/update/delete (ดู 09, 11, 13)
- `web/src/lib/affiliate-api.ts`, `external-affiliate-api.ts`, `marketing-features-api.ts` ต่อ path เป็น `/api/go/api/v1/...` → Go `/api/api/v1/...` (น่าจะ 404 ทั้งหน้า `/affiliates`, `/integrations`, `/marketing`) (ดู 03)
- Public `/api/v1/email/track/click` เป็น open redirect; gRPC เปิดแบบ insecure + reflection; seed admin `admin@marketing.com` ยังอยู่ใน `001_init.up.sql`; stub mailer log token reset/verify (ดู 13)
- `email_campaigns` เขียน status `failed` ที่ไม่อยู่ใน CHECK; SMS campaign ค้าง `sending` เมื่อไม่มี list (ดู 11)
- Compose: volume path `/var/libgresql/data` สะกดผิด, Go migrations ไม่ถูกรันโดย initdb (อยู่ใน sub-directory), gRPC port 9090/9091 ไม่ตรงกัน, rewrites ชี้ `localhost` ใน container (ดู 06, 14, 16)
- `api-python` start จาก checkout ไม่ได้ (import `app.grpc_client.*_pb2` ที่ไม่อยู่ใน Git และ Makefile generate ไปคนละ path) (ดู 06)
- ไม่มี backup, RTO/RPO, branch protection หรือ alert channel ที่ live (ดู 14, 18, 21)

Validate: `node ~/work/core22/tools/validate.mjs docs/diagrams/*.md docs/diagrams/08-adr/*.md`
