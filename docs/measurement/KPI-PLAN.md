---
title: "BookMarketing — KPI and Instrumentation Plan"
status: draft
created: 2026-10-22
updated: 2026-10-22
description: "What BookMarketing would measure once it runs somewhere users can reach it."
template: architecture-docs templates/sdlc/kpi-and-instrumentation-plan.md
---

# BookMarketing — KPI and Instrumentation Plan

## Document control

| Field | Value |
|---|---|
| Product | bookchaowalit-website/marketing-campaign-platform |
| Version / status | 0.1 / **draft** — the product is not deployed (see `specs/001-next16-ci-sentry/spec.md` section 5), so nothing can be measured yet |
| Metric owner | Owner (bookchaowalit) |
| AI maker | Claude Code (Wave 3 session, 2026-10-22) |
| Analytics source of truth | the Go API's PostgreSQL (campaigns, users, organisations) for server-side counts; Sentry for errors |
| Next review | when a deployment target is decided |

## KPI tree

```text
OBJ-01 A small team runs its campaigns from one dashboard instead of several tools
  └─ KPI-O-01 Weekly active organisations (≥ 1 campaign created or updated in the week)
       ├─ KPI-L-01 Sign-up → first campaign created within 7 days
       ├─ KPI-L-02 Organisations using ≥ 2 channel modules (email, social, SEO, PPC…) per week
       └─ KPI-G-01 Web error events per 1,000 sessions (Sentry)
```

## KPI definitions

| KPI ID | Type | Formula | Source | Baseline | Target | Status |
|---|---|---|---|---|---|---|
| KPI-O-01 | outcome | distinct `organization_id` with a campaign insert/update in the ISO week | PostgreSQL `campaigns.updated_at` (server side) | not measurable: no deployment | set after 4 weeks of data | blocked by deployment |
| KPI-L-01 | leading | organisations with a campaign within 7 days of sign-up / sign-ups | PostgreSQL `users.created_at`, `campaigns.created_at` | not measurable | — | blocked by deployment |
| KPI-L-02 | leading | organisations with writes in ≥ 2 module tables in the week / active organisations | PostgreSQL | not measurable | — | blocked by deployment |
| KPI-G-01 | guardrail | Sentry error events / sessions × 1000 | Sentry (code merged; DSN not set) | not measurable | must not rise after a release | planned |

## Event and data specification

All outcome and leading KPIs are computed from existing tables on the server, so no client analytics events are
needed. No personal data leaves the database; reports aggregate by organisation.

## Instrumentation verification (before G10)

Not started: there is no environment with users. Checks 1–7 of the template run after the first deployment.

## Review (2026-10-24): changes requested

Checker: independent Claude reviewer subagent (`solo:solo-reviewer`, Sonnet, read-only, no maker context). Status unchanged until fixed.

- Gap: KPI-G-01 (Sentry errors per 1,000 sessions) has no session source while the plan says no client analytics; KPI-L-02 relies on unspecified 'module tables'.
