---
title: "SPEC-001 — Next.js 16 upgrade, shared CI, Sentry and dependency security"
status: implemented
tags: [spec, sdd, tdd, ci, sentry, nextjs, security]
created: 2026-10-22
updated: 2026-10-22
description: "Move the dashboard from Next.js 14 to 16, put all three services under CI and the security scan, and fix what those checks found."
---

# SPEC-001 — Next.js 16 upgrade, shared CI, Sentry and dependency security

Standard: ADR-0004 (TDD + SDD). Plan: Big Loop Plan 2026-10-21, Wave 3 items 1, 3 and 5.
Upgrade guide followed: Next.js "Upgrading to version 16" (bundled with next 16.3.8, `dist/docs/01-app/02-guides/upgrading/version-16.md`).

## Document control

| Field | Value |
|---|---|
| Spec ID / version | SPEC-001 / 1.0 |
| Status | implemented |
| Repository | bookchaowalit-website/marketing-campaign-platform |
| Owner (maker) | Claude Code (Wave 3 session) |
| Approver | Owner (bookchaowalit) |

## 1. Problem and intent

- **Problem:** the web dashboard ran Next.js 14 / React 18 with `next lint` (removed in 16) and a custom webpack hook
  (which makes a Next 16 build fail). No CI ran for any of the three services, and the dependency trees carried
  about 80 known vulnerabilities.
- **Outcome:** web on Next.js 16.3 / React 19.2 builds with Turbopack and renders the same as before; web, Go and
  Python run in CI with the shared security scan green; runtime errors reach Sentry without personal data.

## 2. Scope

| In scope | Out of scope |
|---|---|
| `web/` upgrade, ESLint flat config, Sentry, `.github/workflows/ci.yml`, dependency bumps in all three services, bugs found on the way | New features; MUI 6/7 or Tailwind 4 migrations; making `api-python` start from a checkout (its generated gRPC modules are not in Git); deploying |

## 3. Behaviour — scenarios

| AC | Given / When / Then | Test |
|---|---|---|
| AC-01 | `next.config.ts` has no webpack hook and does not suppress type errors; `/api/go/*` and `/api/python/*` are proxied | `web/tests/next-config.test.ts` |
| AC-02 | CSV exports quote commas, quotes and line breaks (RFC 4180) and neutralise cells starting with `= + - @` tab or CR (CSV injection) | `web/tests/csv.test.ts` |
| AC-03 | Sentry: off without `NEXT_PUBLIC_SENTRY_DSN`; errors only; request, user, breadcrumbs, extra and frame locals stripped | `web/tests/sentry.test.ts` |
| AC-04 | Stat cards resolve theme colours (`success.main`) before tinting, so the dashboard home no longer crashes with MUI error #9 and tinted cards get a valid background | `web/tests/palette.test.ts` + UI review of `/`, `/affiliates`, `/integrations`, `/marketing` |
| AC-05 | Auth pages hydrate without React error #418: Emotion styles are emitted in `<head>` through `useServerInsertedHTML` (`web/src/app/emotion-registry.tsx`) | Playwright: 0 of 24 loads of `/login` and `/signup` at 375/1440 (6 of 16 failed before the registry) |
| AC-06 | `api-python` accepts only unexpired HS256 tokens signed with `JWT_SECRET` (PyJWT replaces python-jose, which had an unfixed algorithm-confusion advisory) | `api-python/tests/test_auth.py` (7 tests) |
| AC-07 | CI: `ci-node` (web), `ci-python` (api-python), an inline Go job (build, vet, test) and `security.yml`, pinned to v1.1.1 `87dc5e1` | GitHub Actions run for this change |

## 4. Interfaces and data

| Interface | Change |
|---|---|
| `web` | next 16.3.8, react 19.2.3, eslint 9 flat config (same `next/core-web-vitals` rules), `@sentry/nextjs` ^11, axios ^1.20, Docker base node:22 |
| `api-go` | x/crypto, x/net, x/text, grpc 1.83.2, golang-jwt v5.3.1; `go 1.26.0`; Docker builder golang:1.26, runtime runs as a non-root user |
| `api-python` | `requirements.in` (top level) → `requirements.txt` (fully pinned, `uv pip compile`); fastapi 0.157, starlette 1.x, python-multipart 0.0.32, PyJWT 2.15.1; `poetry.lock` added; missing `engine_from_config` import in `migrations/env.py` fixed |
| Root | stale `package-lock.json` (old `apps/web` workspace with next 14.1.0) regenerated; concurrently ^9 |

Reviewed scanner exceptions, each with a reason and a 2027-01-17 expiry: `web/osv-scanner.toml` (braces, postcss-selector-parser:
dev-only, no fix in the required major) and `api-go/osv-scanner.toml` (aws-sdk-go v1 s3crypto and x/crypto/openpgp: no fix,
not imported, confirmed by govulncheck).

## 5. Known gaps

- **No deployment.** `https://bookmarketing.bookchaowalit.com` (registry `live_url`) does not resolve, and
  `https://marketing.bookchaowalit.com` is served by a different repository (`marketing-bookmarketing-web`, Vercel project
  `bookchaowalit-bookmarketing-frontend`). No deploy workflow is added. Owner decision: deploy this dashboard, or correct the registry.
- `api-python` cannot import `app.main` from a checkout: `app/grpc_client/*_pb2.py` must be generated from `packages/grpc-protos` first.
- Pre-existing, not changed: horizontal overflow at 375 px on dashboard tables, demo credentials printed on `/login`,
  hard-coded default `JWT_SECRET` in `api-python/app/core/config.py`, the unused `DashboardOverview` component keeps the old colour bug.
