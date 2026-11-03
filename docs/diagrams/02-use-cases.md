# 02 — Use cases

Actors มาจาก auth level ที่โค้ดบังคับจริง: route สาธารณะ, `AuthMiddleware` (JWT ใดก็ได้), `AdminOnlyMiddleware` (`role == "admin"`) ใน `api-go` และ `get_current_user_id` (JWT ใดก็ได้) ใน `api-python`
`register` สร้าง `role: "user"` เสมอ (`auth_service.go`); admin มีเฉพาะ seed `admin@marketing.com` ใน `api-go/migrations/001_init.up.sql`
ตาราง `organization_members.role` (`owner/admin/member/viewer`) มีใน schema แต่ไม่มีโค้ดใดอ่าน

```mermaid
flowchart LR
  V["Anonymous visitor"]
  U["Signed-in user, role user"]
  A["Admin, role admin"]
  R["Email recipient"]
  P["Referral traffic, partner site"]
  O["Ops: CI and deploy-smoke"]

  subgraph PUB["Public routes"]
    UC1(["Register, login, refresh"])
    UC2(["Forgot or reset password, verify email"])
    UC3(["Track email open pixel and click redirect"])
    UC4(["Track affiliate click, signup, conversion"])
    UC5(["Health: Go /health, Python /health, web /api/health"])
  end

  subgraph AUTH["Any valid JWT"]
    UC10(["View dashboard overview"])
    UC11(["Manage classic campaigns: CRUD, pause, activate"])
    UC12(["Manage unified campaigns: CRUD, submit"])
    UC13(["Approve or reject unified campaign"])
    UC14(["A/B tests: CRUD, start, stop, own tests only"])
    UC15(["Budget periods, allocations, transactions, alerts"])
    UC16(["Calendar events and views"])
    UC17(["Approval workflows and responses"])
    UC18(["Social posts, landing pages, automation, influencers, reporting"])
    UC19(["Social publish via OAuth, only if FACEBOOK_APP_ID set"])
    UC20(["Join affiliate program, links, commissions"])
    UC21(["SEO metrics, keywords, audit, backlinks"])
    UC22(["Competitors CRUD and analysis"])
    UC23(["Analytics summary, channels, funnel, predict"])
    UC24(["Logout, resend verification"])
  end

  subgraph ADM["Admin only"]
    UC30(["Approve or suspend affiliates, affiliate settings"])
    UC31(["Connect Shopee or TikTok accounts, sync, OAuth"])
    UC32(["Email lists, subscribers, send campaign, email A/B"])
    UC33(["PPC platforms and ad campaigns"])
    UC34(["SMS lists, subscribers, send campaign"])
  end

  V --> UC1
  V --> UC2
  V --> UC5
  R --> UC3
  P --> UC4
  O --> UC5
  U --> UC10
  U --> UC11
  U --> UC12
  U --> UC13
  U --> UC14
  U --> UC15
  U --> UC16
  U --> UC17
  U --> UC18
  U --> UC19
  U --> UC20
  U --> UC21
  U --> UC22
  U --> UC23
  U --> UC24
  A --> UC30
  A --> UC31
  A --> UC32
  A --> UC33
  A --> UC34
  A -.->|"inherits"| U
```

Sources: `api-go/internal/delivery/http/router.go`, `api-go/internal/delivery/http/middleware.go`, `api-go/internal/application/service/auth_service.go` (`Role: "user"`), `api-go/migrations/001_init.up.sql` (admin seed), `api-go/migrations/005_organizations.up.sql` (`organization_members`), `api-python/app/main.py`, `api-python/app/core/auth.py`, `web/src/app/api/health/route.ts`, `.github/workflows/deploy-smoke.yml`

ข้อสังเกต

- `require_admin` ใน `api-python/app/core/auth.py` ถูกประกาศแต่ไม่มี route ใดใช้
- หน้า web `/email-campaigns`, `/sms`, `/ppc`, `/integrations` เรียก endpoint admin-only: user ปกติได้ 403
- `UC19` ถูก register เฉพาะเมื่อมี `FACEBOOK_APP_ID` และ `FACEBOOK_APP_SECRET` (`router.go`)

## Open questions

- Role model ที่ต้องการ: ใช้ `users.role` (`user/admin`) ต่อไป หรือย้ายไป `organization_members.role`
- การ approve unified campaign ควรจำกัดเฉพาะ admin หรือ approver ใน `approval_workflows` หรือไม่
