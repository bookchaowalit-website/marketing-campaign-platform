# 21 — Incident response

Solo operator ไม่มี on-call หลายคน (`RUNBOOK.md` §Owner) ไม่มี production จึงยัง **ไม่มี alert channel ที่ live**: สัญญาณที่มีจริงคือ GitHub Actions run ที่ fail (`ci.yml`, `deploy-smoke.yml`) และ Sentry ของ web เมื่อมีการตั้ง `NEXT_PUBLIC_SENTRY_DSN` ไม่พบ template postmortem หรือ incident log ใน repo
ขั้นตอนด้านล่างอ้างชื่อ job, step, env และคำสั่งตามที่อยู่ในโค้ด และแยกให้ชัดว่าการเปลี่ยน env ต้อง restart/redeploy แล้วตรวจแบบ functional ไม่ใช่แค่ GET `/health`

```mermaid
flowchart TD
  A1["GitHub Actions: CI run failed, job web, api-python, api-go or security"] --> T
  A2["GitHub Actions: Deploy smoke failed, error lines start with FAIL"] --> T
  A3["Sentry issue from web, only when NEXT_PUBLIC_SENTRY_DSN is set at build"] --> T
  A4["User report: login or data pages fail"] --> T
  T{"Triage: which signal?"}
  T -->|"CI"| C1["gh run view RUN_ID --log-failed"]
  C1 --> C2{"security / osv-scanner?"}
  C2 -->|"yes"| C3["Bump dependency, or add reviewed exception with ignoreUntil in web or api-go osv-scanner.toml"]
  C2 -->|"no"| C4["Reproduce locally: npm test and npm run build in web, go build, vet, test in api-go, pytest in api-python"]
  C3 --> F["Fix forward: commit to main, wait for CI green"]
  C4 --> F
  T -->|"Deploy smoke"| D1{"Which FAIL line?"}
  D1 -->|"No production URL"| D2["Set repo variable PRODUCTION_URL, then gh workflow run deploy-smoke.yml"]
  D1 -->|"FAIL code on / or /login, or marker MarketPro missing"| D3["Check host logs and the web build"]
  D1 -->|"production does not run the expected commit"| D4["Check Vercel alias moved; script retries SMOKE_ATTEMPTS 20 x SMOKE_DELAY 15 s"]
  T -->|"Auth or data broken"| U1{"Which layer?"}
  U1 -->|"all API calls 401"| U2["JWT_SECRET differs between api-go and api-python, or was rotated"]
  U1 -->|"api-go exits at start"| U3["Check DATABASE_URL, db.Ping, EMAIL_PROVIDER sendgrid or ses needs SENDGRID_API_KEY or EMAIL_FROM_ADDRESS"]
  U1 -->|"provider feature misbehaves"| U4["Kill switch: unset FACEBOOK_APP_ID, TWILIO_ACCOUNT_SID or SHOPEE_PARTNER_ID, or set EMAIL_PROVIDER stub"]
  U2 --> E["Env change: restart api-go and api-python; NEXT_PUBLIC values need a web rebuild"]
  U3 --> E
  U4 --> E
  E --> V["Functional check: POST /api/v1/auth/login with a test account, then GET /api/v1/auth/me and GET /api/v1/seo/sources with that token, expect 200; old token expect 401 after rotation"]
  D2 --> V2["Smoke green: ok lines for /, /api/health commit, /login"]
  D3 --> V2
  D4 --> V2
  F --> V3["CI green on the fix commit"]
  V --> R{"Need rollback?"}
  V2 --> R
  R -->|"code"| RB1["Redeploy previous image or commit, RUNBOOK Rollback"]
  R -->|"schema"| RB2["migrate -path migrations -database DATABASE_URL down 1, not make migrate-down"]
  R -->|"no"| PM["Postmortem: no template in repo"]
  V3 --> PM
  RB1 --> PM
  RB2 --> PM
```

Sources: `.github/workflows/ci.yml` (job names), `.github/workflows/deploy-smoke.yml` (step "Require a production URL", `vars.PRODUCTION_URL`), `scripts/deploy-smoke.sh` (`HOME_MARKER`, `HEALTH_PATH`, `MAIN_ROUTE`, `SMOKE_ATTEMPTS`, `SMOKE_DELAY`, `::error::FAIL` messages), `web/osv-scanner.toml`, `api-go/osv-scanner.toml`, `web/src/lib/sentry.ts`, `api-go/cmd/main.go`, `api-go/internal/delivery/http/router.go` (env read once in `NewRouter`), `api-go/internal/application/service/email_factory.go`, `twilio_service.go`, `api-python/app/core/config.py`, `api-python/app/core/auth.py`, `api-python/app/api/seo.py` (`/sources`), `api-go/Makefile` (`migrate-up`, `migrate-down`), `RUNBOOK.md` §Rollback, §Owner

รายละเอียดที่ต้องระวัง

- **Env ของ api-go อ่านครั้งเดียวตอน start** (`NewRouter`, `NewEmailService`, `NewTwilioClient`) และ `api-python` อ่าน `Settings()` ตอน import: เปลี่ยน env แล้วต้อง restart/redeploy ทั้ง process
- `NEXT_PUBLIC_API_GO_URL`, `NEXT_PUBLIC_API_PYTHON_URL` (rewrites) และ `NEXT_PUBLIC_SENTRY_DSN` ถูกใช้ตอน build (`web/Dockerfile` build args): ต้อง rebuild web
- `GET /health` ของ Go และ Python ตอบ `ok` โดยไม่แตะ DB หรือ provider: ไม่ใช่หลักฐานว่าแก้สำเร็จ ให้ใช้ functional check ข้างบน และสำหรับ email ให้เรียก `POST /api/v1/auth/forgot-password` กับบัญชีทดสอบแล้วดูว่า provider ได้รับจริง (stub แค่ log)
- Rotate `JWT_SECRET` ต้องตั้ง **ค่าเดียวกัน** ใน api-go และ api-python; token ทั้งหมดใช้ไม่ได้ทันที; refresh token ใน DB ยังอยู่แต่ access token ใหม่จะออกด้วย secret ใหม่
- `make migrate-down` เรียก `migrate ... down` แบบไม่ระบุจำนวนและ hard-code DSN `localhost`: จะย้อน **ทุก** migration บน DB local; ใน incident ให้ใช้ `down 1` กับ `DATABASE_URL` ที่ถูกต้อง และ `013`, `014` ไม่มีไฟล์ down
- Kill switch ของ social publish คือเอา `FACEBOOK_APP_ID`/`FACEBOOK_APP_SECRET` ออกแล้ว restart: route `/api/v1/social-pub/*` จะไม่ถูก register (404) และ timer ของ scheduled posts ที่ค้างอยู่หายไปพร้อม process

## Open questions

- ช่องทาง alert ที่ต้องการเมื่อ deploy (Sentry email, Telegram ตาม pattern ของ repo อื่น?)
- ต้องการบัญชีทดสอบสำหรับ functional check ใน production และ seed admin จะถูกจัดการอย่างไร
- ที่เก็บ postmortem / incident log
