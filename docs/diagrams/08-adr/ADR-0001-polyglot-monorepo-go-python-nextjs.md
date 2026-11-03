# ADR-0001 — Polyglot monorepo: Next.js dashboard + Go hexagonal core + Python analytics

- **Status:** accepted (retro-documented)
- **วันที่ตัดสินใจ:** ก่อน `2c50c4e` (init 2026-04-12); layout ปัจจุบันยืนยันใน `b934a69` (2026-08-11, ลบ `apps/` ที่ตายแล้ว)

## Context

ผลิตภัณฑ์รวมหลาย domain (campaign ops, email, SMS, social, affiliates, PPC, SEO, analytics) ในที่เดียว (`PRODUCT.md` §The problem) และเป็น portfolio flagship ที่ต้องแสดงความลึกหลาย runtime

## Decision

แยกเป็น 3 service ใน repo เดียว: `web/` (Next.js App Router), `api-go/` (Gin, hexagonal `delivery → application → domain ← infrastructure`, HTTP + gRPC), `api-python/` (FastAPI สำหรับ SEO/competitor/analytics) โดยมี `packages/grpc-protos/marketing.proto` เป็นสัญญาระหว่าง Go กับ Python และ `packages/shared/types.ts` เป็น types ของ web root `package.json` ใช้ `concurrently` รันทั้งสาม

## Consequences

- ต้นทุน ops 3 runtime + gRPC codegen (`PRODUCT.md` §Tradeoffs); CI จึงมี 3 job (`ci.yml`)
- Go gRPC bindings ถูก commit (`api-go/pkg/grpc/*.pb.go`, ยกเว้นใน `.gitignore` ตั้งแต่ `c796b31`) แต่ Python stubs ไม่ถูก commit: `api-python` start จาก checkout ไม่ได้ (SPEC-001 §5)
- ทั้งสอง API ต่อ DB เดียวกัน: ownership ของ schema ไม่ชัด (ดู ADR-0002)

## Evidence

`README.md` §Tech Stack, §Project Structure, `PRODUCT.md` §Architecture, `package.json`, `api-go/cmd/main.go`, `api-go/internal/`, `api-python/app/main.py`, `packages/grpc-protos/Makefile`, `.gitignore`, `git show b934a69`
