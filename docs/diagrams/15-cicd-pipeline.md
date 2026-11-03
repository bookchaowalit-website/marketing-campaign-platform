# 15 — CI/CD pipeline

สอง workflow ใน `.github/workflows/`: `ci.yml` (build/test/security, "CI never deploys") และ `deploy-smoke.yml` (ตรวจหลัง deploy ซึ่งยังไม่เคยเกิด)

```mermaid
flowchart TD
  subgraph CI["ci.yml: push main, pull_request, workflow_dispatch; permissions contents read; concurrency ci-ref cancel in progress"]
    T(["trigger"]) --> W["web: ci-node.yml at 87dc5e1c v1.1.1"]
    T --> PY["api-python: ci-python.yml at 87dc5e1c v1.1.1"]
    T --> GO["api-go: inline job ubuntu-24.04"]
    T --> SEC["security: security.yml at 87dc5e1c v1.1.1"]
    W --> W1["Checkout, detect package manager, setup Node 22, npm install"]
    W1 --> W2["Lint: npm run lint, eslint"]
    W2 --> W3["Typecheck: npm run typecheck"]
    W3 --> W4["Test: npm test, node --test tests, require-tests true"]
    W4 --> W5["Build: npm run build"]
    PY --> P1["Checkout, setup-uv, Python 3.11"]
    P1 --> P2["Install: uv pip install -r requirements.txt"]
    P2 --> P3["Lint: ruff check --select E9,F63,F7,F82"]
    P3 --> P4["Test: pytest -q tests and compileall app migrations"]
    GO --> G1["actions/checkout v4, actions/setup-go v5 from go.mod"]
    G1 --> G2["go build ./..."]
    G2 --> G3["go vet ./..."]
    G3 --> G4["go test ./..."]
    SEC --> S1["gitleaks"]
    SEC --> S2["semgrep"]
    SEC --> S3["osv-scanner, exceptions in web and api-go osv-scanner.toml"]
    SEC --> S4["trivy-fs"]
  end
  subgraph DS["deploy-smoke.yml: deployment_status, workflow_dispatch; concurrency deploy-smoke"]
    D0(["Production deployment success, or manual"]) --> D1["Checkout v7.0.1 SHA-pinned, no credentials"]
    D1 --> D2{"vars.PRODUCTION_URL set?"}
    D2 -->|"no"| D3["error, exit 1"]
    D2 -->|"yes"| D4["scripts/deploy-smoke.sh URL EXPECTED_SHA"]
    D4 --> D5["GET / shows MarketPro"]
    D5 --> D6["GET /api/health commit matches, 20 tries x 15 s"]
    D6 --> D7["GET /login shows MarketPro"]
  end
  CI -.->|"no needs or workflow_run link"| DS
```

Sources: `.github/workflows/ci.yml`, `.github/workflows/deploy-smoke.yml`, `scripts/deploy-smoke.sh`, `web/package.json` (scripts), `api-go/go.mod`, `api-python/requirements.txt`, `web/osv-scanner.toml`, `api-go/osv-scanner.toml`, platform-ci-workflows `87dc5e1c` `.github/workflows/ci-node.yml` (step names, `lint/typecheck/test-command: auto`), `ci-python.yml`, `security.yml` (job names)

ข้อเท็จจริงจาก `gh run list` (2026-11-03): CI รัน 8 ครั้ง ทั้งหมดเป็น `push` บน `main` (ไม่มี PR run); ล้ม 1 ครั้งที่ `1095ae3` (jobs `api-go`, `web`, `security / osv-scanner` fail) แก้ใน `c796b31` ภายใน ~14 นาที; `Deploy smoke` ไม่เคยรัน; GitHub deployments ของ repo = 0

## Open questions

- จะ deploy ไปที่ไหน (deploy-smoke เขียนรอ Vercel `deployment_status`) และ api-go / api-python จะไป host ใด
- ต้องการ pin `actions/checkout@v4` / `actions/setup-go@v5` ใน job Go เป็น SHA ให้ตรงกับ policy ของ shared workflows หรือไม่
- เปิด coverage gate (`coverage-threshold`) ของ shared workflows หรือไม่ (ตอนนี้ไม่ได้ส่ง input)
