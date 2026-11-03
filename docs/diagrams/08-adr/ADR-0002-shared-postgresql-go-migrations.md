# ADR-0002 — PostgreSQL ตัวเดียวร่วมกัน, golang-migrate ใน `api-go/migrations` เป็น schema หลัก

- **Status:** accepted (retro-documented)
- **วันที่ตัดสินใจ:** schema 001–015 มีตั้งแต่ `2c50c4e`; `016_campaigns_organization_id` เพิ่มใน `b934a69` (2026-08-11)

## Context

ทั้ง Go และ Python ต้องอ่าน `users`, `campaigns`, `ab_tests`, `competitors` README audit pass (2026-08) พบว่า migration set ไม่เคย apply ได้ จึงแก้แล้วยืนยันว่า 16 migrations apply บน PostgreSQL ใหม่ได้

## Decision

ใช้ PostgreSQL 15 database `marketing_saas` ตัวเดียว (`DATABASE_URL`) โดยมี `api-go/migrations/NNN_*.up.sql` / `.down.sql` (รันด้วย `make migrate-up` = `migrate -path migrations ...`) เป็นชุดที่ README ยืนยันว่า apply ครบ Go ใช้ `database/sql` + `lib/pq` (pool 25 open / 5 idle); Python ใช้ SQLAlchemy (pool 10 + overflow 20)

## Consequences

- มี schema source อีกสองชุดที่ไม่ตรงกัน: `docker/migrations/001_initial_schema.sql` + `002_seed_data.sql` (compose init) และ `api-python/migrations/versions/001..005` (Alembic) ไม่มีเอกสารบอกว่าใช้ชุดไหนใน environment ใด
- `docker/docker-compose.yml` mount `../api-go/migrations` เป็น sub-directory `go-migrations` ของ `/docker-entrypoint-initdb.d`; entrypoint ของ image postgres รันเฉพาะไฟล์ระดับบนสุด ดังนั้น compose จะได้ schema จาก `docker/migrations` ไม่ใช่ Go migrations (อ่านจาก config ยังไม่ได้รันยืนยัน)
- `013_ppc_ads.up.sql` และ `014_sms_marketing.up.sql` ไม่มีไฟล์ `.down.sql`
- ไม่มี migration runner ตอน start service: ต้องรันเอง

## Evidence

`api-go/migrations/*.sql`, `api-go/Makefile` (`migrate-up`), `api-go/internal/infrastructure/database/postgres.go`, `api-python/app/core/database.py`, `api-python/alembic.ini`, `api-python/migrations/versions/`, `docker/docker-compose.yml`, `docker/migrations/`, `README.md` §api-go
