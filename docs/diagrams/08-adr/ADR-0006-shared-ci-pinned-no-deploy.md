# ADR-0006 — Shared CI pinned ที่ platform-ci-workflows v1.1.1, CI ไม่ deploy

- **Status:** accepted (retro-documented)
- **วันที่ตัดสินใจ:** `1095ae3` (2026-10-22 local), green ครั้งแรก `c796b31`; deploy-smoke เพิ่มใน `8e1d0f7` / `60379ea` (2026-10-24)

## Context

ก่อน 2026-10-22 ไม่มี CI เลย และ dependency มีช่องโหว่ราว 80 รายการ (SPEC-001 §1) ยังไม่มี host ของ dashboard นี้

## Decision

`ci.yml` เรียก reusable workflows ที่ pin SHA `87dc5e1c7e0416886a6371d5a101bb2ccba38edd` (v1.1.1): `ci-node.yml` (web, `require-tests: true`, `npm run build`), `ci-python.yml` (api-python, Python 3.11, Ruff เฉพาะ `E9,F63,F7,F82`, pytest + compileall), `security.yml` (gitleaks, semgrep, osv-scanner, trivy-fs) และ Go job inline (`go build`, `go vet`, `go test`) CI "never deploys" `deploy-smoke.yml` รอ `deployment_status` จาก Vercel และต้องมี repo variable `PRODUCTION_URL`

## Consequences

- ไม่มี branch protection บน `main` (`gh api .../branches/main/protection` = 404): CI เป็นสัญญาณหลัง push ไม่ใช่ gate
- Scanner exceptions มีวันหมดอายุ 2027-01-17 (`web/osv-scanner.toml`, `api-go/osv-scanner.toml`)
- deploy-smoke ไม่เคยรัน (ไม่มี deployment บน GitHub); manual run จะ fail จนกว่าจะตั้ง `PRODUCTION_URL`

## Evidence

`.github/workflows/ci.yml`, `.github/workflows/deploy-smoke.yml`, `scripts/deploy-smoke.sh`, `web/osv-scanner.toml`, `api-go/osv-scanner.toml`, `specs/001-next16-ci-sentry/spec.md`, `gh run list` (8 CI runs), platform-ci-workflows `87dc5e1c` `.github/workflows/ci-node.yml`, `ci-python.yml`, `security.yml`
