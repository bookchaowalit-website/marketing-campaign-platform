# 03 — User flows

Flow มาจาก App Router pages ใน `web/src/app/`, nav ใน `web/src/components/Sidebar.tsx`, auth state ใน `web/src/hooks/useAuth.tsx` และ API client ใน `web/src/lib/*.ts`
ทุกหน้าใต้ `(dashboard)` ถูกครอบด้วย `web/src/app/(dashboard)/layout.tsx`: ถ้า `useAuth()` ไม่มี `user` จะ `router.push('/login')` (guard ฝั่ง client เท่านั้น ไม่มี `middleware.ts`)

## Flow A — Visitor: sign up, login, password reset

```mermaid
flowchart TD
  S(["Open any dashboard URL"]) --> G{"token and user in localStorage?"}
  G -->|"no"| L["/login"]
  G -->|"yes"| RF["POST /v1/auth/refresh with stored refresh_token"]
  RF -->|"ok"| H["Dashboard home /"]
  RF -->|"fail"| CL["clearAuth, stay logged out"]
  CL --> L
  L -->|"link"| SU["/signup"]
  SU --> SU2["POST /v1/auth/register"]
  SU2 --> H
  L --> LI["POST /v1/auth/login"]
  LI -->|"200: store token, refresh_token, user"| H
  LI -->|"401"| L
  FP["/forgot-password"] --> FP2["POST /v1/auth/forgot-password, always 200"]
  FP2 --> MAIL["Email link FRONTEND_URL/reset-password?token=..."]
  MAIL --> RP["/reset-password, inside Suspense"]
  RP --> RP2["POST /v1/auth/reset-password"]
  RP2 --> L
  SU2 -.-> VE["Verification email link FRONTEND_URL/verify-email?token=..."]
  VE -.-> NF["No /verify-email page in web/src/app: not-found"]
```

Sources: `web/src/app/(dashboard)/layout.tsx`, `web/src/hooks/useAuth.tsx`, `web/src/app/login/page.tsx`, `web/src/app/signup/page.tsx`, `web/src/app/forgot-password/page.tsx`, `web/src/app/reset-password/page.tsx`, `web/src/app/not-found.tsx`, `web/src/lib/api-client.ts` (`authApi`, 401 interceptor), `api-go/internal/application/service/auth_service.go` (`resetURL`, `verificationURL`)

## Flow B — Signed-in user: navigate the dashboard

```mermaid
flowchart LR
  H["/ Dashboard: dashboardApi.getOverview + campaignApi"] --> N{"Sidebar"}
  N --> C1["/unified-campaigns: CRUD, submit, approve, reject"]
  N --> C2["/campaigns: classic CRUD, pause, activate"]
  N --> C3["/ab-testing: list A/B tests"]
  N --> M1["/marketing: email lists, social posts, landing pages, automation, influencers"]
  N --> M2["/ppc: adsApi, admin only"]
  N --> M3["/sms: smsApi, admin only"]
  N --> M4["/email-campaigns: send, analytics, A/B, admin only"]
  N --> M5["/video: unifiedCampaignApi"]
  N --> M6["/seo: seoApi via /api/python"]
  N --> M7["/social: no page.tsx"]
  N --> R1["/competitors: competitorApi via /api/python"]
  N --> R2["/affiliates: affiliate-api"]
  N --> P1["/budget"]
  N --> P2["/calendar"]
  N --> P3["/approvals: list and respond"]
  N --> A1["/analytics: analyticsApi via /api/python"]
  N --> A2["/integrations: external-affiliate-api, admin only"]
  N --> S1["/settings: profile and password calls commented out"]
  N --> D1["/docs: no page.tsx"]
```

Sources: `web/src/components/Sidebar.tsx`, `web/src/app/(dashboard)/*/page.tsx`, `web/src/lib/api-client.ts`, `web/src/lib/affiliate-api.ts`, `web/src/lib/external-affiliate-api.ts`, `web/src/lib/marketing-features-api.ts`, `api-go/internal/delivery/http/router.go` (admin groups)

## Flow C — Admin: send an email campaign

```mermaid
flowchart TD
  A["/email-campaigns"] --> B["Create campaign: marketingApi.createEmailCampaign"]
  B --> C["Send: emailCampaignApi.sendCampaign = POST /v1/email/campaigns/:id/send"]
  C --> D{"Campaign has linked lists and active subscribers?"}
  D -->|"no"| E["error, status set to failed"]
  D -->|"yes"| F["Status sending, one email_sends row per subscriber"]
  F --> G["SendEmail via SendGrid, SES or stub"]
  G --> H["Status sent, sent_count updated"]
  H --> I["Analytics: GET /v1/email/campaigns/:id/analytics"]
```

Sources: `web/src/app/(dashboard)/email-campaigns/page.tsx`, `web/src/lib/api-client.ts` (`emailCampaignApi`), `api-go/internal/application/service/email_campaign_service.go`

## ข้อสังเกตที่กระทบ flow (จากการอ่านโค้ด ยังไม่ได้รันยืนยัน)

- `web/src/lib/affiliate-api.ts`, `external-affiliate-api.ts`, `marketing-features-api.ts` import `apiGo` (baseURL `/api/go`) แล้วต่อ path `/api/v1/...` ได้ `/api/go/api/v1/...` ซึ่ง rewrite ใน `web/next.config.ts` เป็น `<go>/api/api/v1/...`: ไม่มี route นี้ใน `router.go` ส่วน `api-client.ts` ใช้ `/v1/...` ถูกต้อง หน้า `/affiliates`, `/integrations`, `/marketing` และบางส่วนของ `/email-campaigns` จึงน่าจะได้ 404
- `email-campaigns/page.tsx` มี `fetch('/api/go/v1/email/ab-tests')` และ fallback เป็น `defaultABTests` เมื่อ error: หน้าจะแสดงข้อมูลตัวอย่างแทน error
- Sidebar มีลิงก์ `/social` และ `/docs` ที่ไม่มี `page.tsx`
- `web/public/index.html`, `app.html`, `dashboard.html` (title "MarketPro") เป็นหน้า static รุ่นเก่า เสิร์ฟที่ `/index.html` ฯลฯ ไม่ได้เชื่อมกับ App Router

## Open questions

- ต้องการหน้า `/verify-email` หรือให้ลิงก์ verify ชี้ไปที่อื่น
- `/social` และ `/docs` ควรมีหน้าหรือเอาออกจาก Sidebar
- ไฟล์ static ใน `web/public/*.html` ยังต้องใช้หรือไม่
