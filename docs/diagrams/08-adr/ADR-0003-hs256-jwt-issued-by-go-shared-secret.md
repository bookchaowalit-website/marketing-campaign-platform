# ADR-0003 — Auth: Go ออก HS256 JWT + refresh token, Python verify ด้วย secret เดียวกัน

- **Status:** accepted (retro-documented)
- **วันที่ตัดสินใจ:** auth Go มีตั้งแต่ `2c50c4e`; Python เปลี่ยนจาก python-jose เป็น PyJWT ใน `76f10ed` (2026-10-22, SPEC-001 AC-06)

## Context

Dashboard ต้องเรียกสอง API ด้วย session เดียว ไม่มี identity provider ภายนอก

## Decision

- `AuthService` (Go) hash password ด้วย bcrypt, ออก access token HS256 อายุ 24 ชั่วโมง claims `user_id`, `email`, `role` และ refresh token (UUID) อายุ 7 วัน เก็บเป็น sha256 ใน `refresh_tokens`; `/auth/refresh` revoke ตัวเก่าแล้วออกใหม่ (rotation); `/auth/logout` ใส่ hash ของ access token ใน `token_blacklist`
- `api-python/app/core/auth.py` verify ด้วย `jwt.decode(..., settings.JWT_SECRET, algorithms=["HS256"])` ทุก router (`Depends(get_current_user_id)`)
- Browser เก็บ `token`, `refresh_token`, `user` ใน `localStorage` และแนบ `Authorization: Bearer` ผ่าน axios interceptor
- `/api/v1/auth/*` สาธารณะมี rate limit 10 ครั้ง/นาที/IP (in-memory)

## Consequences

- Go และ Python ต้องตั้ง `JWT_SECRET` ค่าเดียวกัน; ทั้งคู่มี default hard-coded (`your-secret-key-change-in-production`) ถ้าไม่ตั้ง env
- Go `ValidateToken` ไม่ได้จำกัด signing method ด้วย `jwt.WithValidMethods` (คืน HMAC secret ให้ทุก alg) ต่างจากฝั่ง Python ที่ล็อก HS256
- JWT ไม่มี `organization_id`: tenancy ไม่ถูกบังคับ (`PRODUCT.md` §Trust boundaries)
- Token ใน `localStorage` อ่านได้จาก script ใดก็ได้ในหน้า (XSS exposure)
- Rate limiter และ blacklist check อยู่ใน process เดียว: ไม่แชร์ข้าม instance

## Evidence

`api-go/internal/application/service/auth_service.go`, `api-go/internal/delivery/http/middleware.go`, `rate_limiter.go`, `router.go`, `api-go/migrations/009_refresh_tokens.up.sql`, `api-python/app/core/auth.py`, `api-python/app/core/config.py`, `api-python/tests/test_auth.py`, `web/src/lib/api-client.ts`, `web/src/hooks/useAuth.tsx`, `specs/001-next16-ci-sentry/spec.md` AC-06
