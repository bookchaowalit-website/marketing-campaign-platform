# Marketing SaaS Platform (BookMarketing)

A comprehensive Marketing SaaS platform built with a polyglot microservices architecture.

**Interview case study:** [`PRODUCT.md`](./PRODUCT.md)  
**Dirty-tree commit clusters:** [`DIRTY-TRIAGE.md`](./DIRTY-TRIAGE.md)

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
├── web/                  # Next.js frontend
├── api-go/               # Go backend (hexagonal)
├── api-python/           # Python FastAPI backend
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

## Notes from a recent audit pass

This became the deepest pass of a portfolio-wide interview-readiness audit,
because **the entire Go backend had apparently never successfully
compiled, and its full migration set had apparently never successfully
applied, before this pass.** What follows is grouped by service.

### Structural

- A dead `apps/` directory (`apps/web`, `apps/api-go`, `apps/api-python`)
  coexisted with the real, working top-level `web/`, `api-go/`,
  `api-python/` — `apps/api-go` and `apps/api-python` were completely
  empty, `apps/web` had one orphaned stale copy of a real component.
  `packages/grpc-protos/Makefile`'s output paths still pointed at the dead
  `apps/` structure. Removed `apps/`, fixed the Makefile, and generated
  the Go gRPC bindings it was supposed to produce (`api-go/pkg/grpc/` —
  `internal/delivery/grpc/server.go` imports this as a local package path;
  since the directory never existed, Go tried to fetch it as an external
  module and failed. Never generated before this pass). This README's own
  "Project Structure" section still described the dead `apps/` layout —
  fixed. It also referenced three docs (`docs/architecture.md`, `api.md`,
  `development.md`) that don't exist — removed the dead links.

### `api-go` (Go backend) — never compiled

Every one of the following was found by actually running `go build ./...`
and fixing errors until it passed, then `go vet`, then applying every
migration to a real (throwaway, local) PostgreSQL database and running the
live server against it:

- A malformed dependency pseudo-version in `go.mod` blocked `go mod tidy`
  outright.
- The generated gRPC code (see above) targeted a newer gRPC-Go API than
  the pinned `google.golang.org/grpc` version supported — bumped it.
- **Two independent duplicate-type-declaration clusters**, both from
  genuinely separate features that happened to reuse the same type names:
  `SocialAccount`/`SocialPost` were declared once for a real OAuth
  publishing feature (`/social-pub`, with real TikTok/Twitter/LinkedIn
  client registrations) and again for a simpler post-scheduling feature
  (`/social`) — renamed the OAuth-publishing ones to
  `SocialPublishAccount`/`SocialPublishPost` throughout. Separately,
  `unified_campaign.go` redeclared six types (`BudgetPeriod`,
  `CalendarEvent`, `ApprovalWorkflow`, etc.) that already had their own
  dedicated, actually-used domain files — those six were genuinely dead
  copies, confirmed by checking nothing in the unified-campaign
  service/handler/repository ever referenced them, and removed.
- A missing `domain.` package prefix, a rune passed where a function
  wanted a string (×3), `http.Response.Close()` called instead of
  `.Body.Close()` (×2, doesn't exist on `Response`), 2 unused variables, a
  fundamentally broken timer-cleanup method (tried to `Reset()` a
  `time.Timer` using a `.Duration` field that doesn't exist — and was
  never even called from anywhere), and several unused imports.
- Two Social-platform clients (`FacebookClient`, `InstagramClient`) each
  returned their own bespoke anonymous struct from `GetAnalytics` instead
  of the shared `domain.SocialPostAnalytics` type their own interface
  (`SocialPublishClient`) required — neither had ever actually satisfied
  the interface they were written to implement. Rewrote both.
- A real `domain.EmailService` (SendGrid/SES/stub, chosen by
  `EMAIL_PROVIDER` env var) was constructed correctly for auth emails, but
  a second, unrelated `*EmailMarketingService` (list/subscriber/campaign
  CRUD only, no send capability) was declared later in the same function
  under the same variable name, shadowing the real mailer — so email
  campaigns were wired to something that could never send email. Renamed
  the shadowing variable and fixed the wiring rather than adding a stub
  `SendEmail` method.
- Four places in `affiliate_handler.go` reached directly into another
  package's unexported struct fields (`h.service.affiliateRepo`, etc.) —
  added the missing public methods on `AffiliateService` instead.
- **Migrations, once the code compiled, didn't apply either**: a table
  had a foreign key to another table defined later in the same file
  (reordered); a `CHECK` constraint referenced the wrong column name
  (typo'd against its own table); `CREATE TRIGGER IF NOT EXISTS` was used
  twice — not valid PostgreSQL syntax (only `DROP TRIGGER IF EXISTS` +
  `CREATE TRIGGER` is) — fixed to that pattern; two migrations created
  tables named `social_accounts`/`social_posts` identical to names
  already claimed by an earlier migration for the *other* social feature
  above — since `CREATE TABLE IF NOT EXISTS` is silently a no-op on a
  name collision, the second feature's schema would never actually have
  applied — renamed to `social_publish_*` throughout (migration + Go SQL
  queries) to match the Go-level rename.
- A demo organization ID (`00000000-...-000001`) was referenced by
  seed data in two separate migrations, but no migration ever created that
  `organizations` row — every FK-constrained seed insert referencing it
  was unrunnable as written. Added the actual seed row where the
  `organizations` table is created, and restored the two seed inserts
  that had to be stripped out mid-pass before that fix existed.
- **Live-testing a real authenticated request surfaced one more**: `GET
  /api/v1/campaigns` 500'd with `column "organization_id" does not
  exist`. The `campaigns` table was created in the very first migration,
  before multi-tenancy (the `organizations` table) existed; the
  repository layer was later rewritten to query by `organization_id`, but
  no migration ever added that column to the pre-existing table. Added it.
- **Found, documented, not fixed** — a real architectural gap, not a
  bug with a clear one-line fix: the auth middleware's JWT claims only
  ever populate `user_id`/`email`/`role` in request context, never
  `organization_id`. Every handler that reads
  `c.GetString("organization_id")` (e.g. `unified_campaign_handler.go`)
  always gets an empty string for every user, which then fails as
  `invalid input syntax for type uuid: ""` against the database. The
  "legacy" campaigns handler avoids this by using a hardcoded
  `defaultOrganizationID` constant instead of real auth context — real
  multi-tenant org-membership resolution was never wired into
  authentication. Fixing this properly requires a product decision (does
  registration create an org? can a user belong to several? how is the
  active one selected?) that's out of scope for an audit pass.
- Verified live end-to-end: real user registration → real bcrypt
  password + JWT issuance → login → an authenticated request correctly
  401s without a token and succeeds with one → all 16 migrations apply
  cleanly to a fresh database → the previously-500ing campaigns endpoint
  now returns real data. `go build ./...` and `go vet ./...` both pass
  clean; no tests exist in this service.

### `api-python` (Python backend) — syntax-verified only

All `.py` files pass `python -m py_compile` cleanly. A full dependency
install and live verification wasn't possible in this environment: the
project pins packages (`grpcio==1.60.0`, `pandas==2.2.0`,
`pydantic==2.5.3`, `psycopg2-binary`) against the declared Python 3.11,
and none of those old versions have prebuilt wheels for the much newer
Python available in this sandbox — several (`pandas`, `pydantic-core`)
failed to even build from source against it. This is an environment
mismatch, not a defect in the code; a real verification pass needs
Python 3.11 (`pyenv`/Docker). No tests exist in this service either.

### `web` (Next.js frontend)

- `next lint` had no ESLint config at all, and `eslint`/`eslint-config-next`
  weren't even declared as dependencies — `next lint` could never have
  run. Added both; the two real errors it then caught (unescaped JSX
  quotes, a missing `Suspense` boundary around `useSearchParams` that
  made `/reset-password` fail to prerender) are fixed.
- `next build` failed outright with real TypeScript errors: an `apiGo`
  axios instance existed in `api-client.ts` but was never `export`ed,
  breaking 3 files/4 dashboard pages that imported it; the SEO dashboard's
  `useQuery` was typed against a local flat shape that didn't match what
  `seoApi.getMetrics()` actually returns (a wrapped, partly-differently-named
  object — ground-truthed against the real Python endpoint's Pydantic
  schema, not guessed) — fixed both files' data-fetching to unwrap and
  correctly map the real response shape, including a `search_volume` vs
  `volume` field-name mismatch in the keywords list that would have
  silently rendered as blank/zero even after the build was fixed.
- `next` was bumped `14.1.0` → `^14.2.35`.
- Verified live against the real production entrypoint
  (`node .next/standalone/server.js`, matching what the repo's own
  Dockerfile actually runs — `next start` alone logs a warning and isn't
  representative): all 12+ routes return 200, auth-gated pages render.

## 📄 License

MIT
