# BookMarketing — Product Brief & Interview Case Study

*Tier A Interview Flagship (depth closer to full-stack product than a
typical showcase). Nested repo:
`projects/product/engineering/book-dev/github/bookchaowalit-website/book-products/bookmarketing`.*

## Who this is for

A small marketing team or solopreneur who needs campaign ops, SEO/SEM
views, email/social publishing, affiliate tracking, and light CRM in one
place — without buying five separate SaaS tools.

## The problem

Marketing tooling fragments across ad platforms, email, social schedulers,
SEO spreadsheets, and a CRM. BookMarketing is a **polyglot monorepo** that
models those domains under one product surface: Next.js dashboard + Go
hexagonal core API + Python ML/analytics service + PostgreSQL + gRPC
protos.

## Architecture (Verified layout)

```text
web/                 Next.js 14 App Router dashboard
api-go/              Go hexagonal HTTP + gRPC core
api-python/          FastAPI analytics / SEO research
packages/grpc-protos Shared protobufs
packages/shared      Shared types/utilities
docker/              Compose / service images
```

There is **no** live `apps/` tree — a dead parallel layout was removed in
the 2026-08-05 audit. Canonical packages are the top-level folders above.

### Trust boundaries

| Surface | Auth / tenancy notes |
|---|---|
| Go JWT login/register | Issues tokens; middleware fills `user_id` / `email` / `role` |
| Multi-tenant org | `organizations` table + `organization_id` on campaigns (added in migration 016) |
| **Gap** | JWT **does not** populate `organization_id` in request context — handlers that read it get `""` |

## Main user flows (Observed)

1. Register / login → JWT for dashboard API calls  
2. SEO dashboard pulls metrics from Python service (response shape must be
   unwrapped — fixed in web client)  
3. Campaigns / affiliates / email / social publish via Go services  
4. Social: two historically colliding domains (OAuth publish vs simple
   schedule) — publish types renamed to `SocialPublish*`  

## Hard problems (Verified by running tooling)

Interview story: **the Go backend had never successfully compiled, and
migrations had never cleanly applied**, until the audit pass.

| Finding | Status |
|---|---|
| Dead `apps/` + missing generated gRPC bindings | Fixed (generate into `api-go/pkg/grpc/`) |
| Duplicate domain types / interface mismatches | Fixed |
| Email service variable shadowing (campaigns could not send) | Fixed |
| Migrations (FK order, CHECK typo, invalid `CREATE TRIGGER IF NOT EXISTS`, table name collisions) | Fixed |
| `campaigns.organization_id` missing despite query layer | Fixed (migration 016) |
| Auth middleware never injects `organization_id` | **Documented limitation** — needs product decision |
| Python service pins need Python 3.11 | Environment mismatch, not code defect |
| Web: missing ESLint; unexported `apiGo`; SEO response shape | Fixed |

Evidence from that pass: `go build ./...`, `go vet ./...`, 16 migrations on
throwaway Postgres, live register → JWT → authenticated request; `next
build` + standalone `node .next/standalone/server.js` route smoke.

## Dirty-tree triage (2026-08-11 snapshot)

**Do not bulk-commit.** Cluster before any nested-repo commit.

| Cluster | Paths | Intent | Commit stance |
|---|---|---|---|
| A. Go compile/migration remediation | `api-go/**` (~26 files) + `api-go/migrations/016_*` | Make backend buildable and schema honest | One commit when reviewed |
| B. Web tooling + SEO/API client | `web/**`, `web/.eslintrc.json` | Lint/build + client contract | Separate commit |
| C. Docs | `README.md`, this `PRODUCT.md` | Interview evidence | Can ship with A or alone |
| D. Protos Makefile | `packages/grpc-protos/Makefile` | Point off dead `apps/` | With cluster A |
| E. Dead tree removal | `apps/web/...` deleted | Structural cleanup | With cluster A |

**Claim hygiene**

- **Verified:** Go builds; migration set applies; web builds with ESLint.  
- **Observed:** Feature surface is large; not every handler is E2E-tested.  
- **Planned / incomplete:** Real multi-tenant org membership in JWT; Python
  live suite on 3.11; automated tests for Go/web (none today).

## Tradeoffs

- Polyglot depth is the portfolio strength and the ops cost (three
  runtimes, gRPC codegen, dual social domains).  
- Multi-tenancy was bolted on after campaigns existed — column fix is
  easy; **auth context wiring is the real product decision**.  
- Pricing/MCP-style marketing claims in older docs are not the evidence
  base for this case study — use live build/migrate/run evidence instead.

## What I'm proud of

Turning a monorepo that looked complete but **could not compile or
migrate** into one that can — and leaving the org-auth gap explicit rather
than hardcoding more `defaultOrganizationID` shortcuts.

## Next improvements

1. Product decision: registration → org creation? multi-org? active org
   header/claim? Then wire middleware + handlers.  
2. `go test` / vitest smoke for auth + campaigns + SEO unwrap.  
3. Python 3.11 Docker target for analytics service.  
4. Selective nested commits for clusters A → B → C (owner-approved).

## Status

- **2026-08-05:** Deep remediation + README audit notes.  
- **2026-08-11:** `PRODUCT.md` + dirty triage; no bulk commit; no force
  redesign of UI.
