# ADR-0004 — Browser เรียก API ผ่าน Next.js rewrites แบบ same-origin

- **Status:** accepted (retro-documented)
- **วันที่ตัดสินใจ:** `api-client.ts` ("Shared API client using Next.js API proxy (avoids CORS)") มีก่อน audit; ยืนยันเป็น AC-01 ใน `c7e232c` (2026-10-22)

## Context

Dashboard อยู่ที่ :3000 ส่วน API อยู่ที่ :8080 และ :8081 การเรียกข้าม origin ต้องพึ่ง CORS ของทั้งสอง API

## Decision

`web/next.config.ts` rewrite `/api/go/:path*` → `${NEXT_PUBLIC_API_GO_URL}/api/:path*` และ `/api/python/:path*` → `${NEXT_PUBLIC_API_PYTHON_URL}/api/:path*` (default localhost) client ใช้ `baseURL: '/api/go'` / `'/api/python'` และ path `/v1/...`

## Consequences

- Next.js server ต้องเข้าถึง API ได้ตอน runtime; ค่าที่ใช้คือ `NEXT_PUBLIC_*` ซึ่งใน `web/Dockerfile` ถูกฝังเป็น build arg (ค่า `API_GO_URL` / `API_PYTHON_URL` ใน compose ไม่มีโค้ดใดอ่าน)
- ไฟล์ client รุ่นเก่า (`affiliate-api.ts`, `external-affiliate-api.ts`, `marketing-features-api.ts`) ยังต่อ `/api/v1/...` บน base `/api/go` ได้ path ซ้อน `/api/api/v1` (ดู 03)
- API ทั้งสองยังเปิด CORS กว้าง (Go reflect `Origin` + `Allow-Credentials: true`; Python `allow_origins=["*"]`) แม้ proxy จะทำให้ไม่จำเป็น

## Evidence

`web/next.config.ts`, `web/tests/next-config.test.ts`, `web/src/lib/api-client.ts`, `web/src/lib/api.ts`, `web/Dockerfile`, `docker/docker-compose.yml`, `api-go/internal/delivery/http/router.go` (CORS), `api-python/app/main.py`
