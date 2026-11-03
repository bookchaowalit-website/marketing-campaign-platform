# 10 — Sequences

สี่ request path ที่สำคัญที่สุด: login + เรียก API ที่ต้อง auth, SEO metrics ข้าม service, ส่ง email campaign + tracking และ social post แบบตั้งเวลา

## 1. Login แล้วเรียก API ที่ต้อง auth (และ refresh)

```mermaid
sequenceDiagram
  autonumber
  actor U as Browser
  participant W as web Next.js
  participant G as api-go Gin
  participant DB as PostgreSQL
  U->>W: POST /api/go/v1/auth/login
  W->>G: rewrite to /api/v1/auth/login
  G->>G: RateLimitMiddleware 10 per min per IP
  G->>DB: SELECT user by email
  G->>G: bcrypt.CompareHashAndPassword, check is_active
  G->>DB: INSERT refresh_tokens sha256 hash, expires 7 d
  G-->>U: 200 token HS256 24 h, refresh_token, user
  U->>U: localStorage token, refresh_token, user
  U->>W: GET /api/go/v1/campaigns with Bearer
  W->>G: /api/v1/campaigns
  G->>DB: SELECT token_blacklist by sha256
  G->>G: ValidateToken, set user_id, email, role
  G->>DB: SELECT campaigns WHERE organization_id = defaultOrganizationID
  G-->>U: 200 campaigns, total
  Note over U,G: Any 401 makes the axios interceptor clear tokens and go to /login
  U->>W: on mount POST /api/go/v1/auth/refresh
  W->>G: /api/v1/auth/refresh
  G->>DB: lookup hash, check revoked, revoke old, insert new
  G-->>U: new token and refresh_token
```

Sources: `web/src/lib/api-client.ts`, `web/src/hooks/useAuth.tsx`, `web/next.config.ts`, `api-go/internal/delivery/http/router.go`, `middleware.go`, `rate_limiter.go`, `auth_handler.go`, `campaign_handler.go`, `api-go/internal/application/service/auth_service.go`, `api-go/internal/infrastructure/repository/auth_repository.go`

## 2. SEO metrics (web → api-python → provider หรือ mock)

```mermaid
sequenceDiagram
  autonumber
  actor U as Browser
  participant W as web Next.js
  participant P as api-python FastAPI
  participant S as SEODataService
  participant X as GSC, Ahrefs or SEMrush
  U->>W: GET /api/python/v1/seo/metrics?domain=example.com, Bearer
  W->>P: /api/v1/seo/metrics
  P->>P: HTTPBearer, jwt.decode HS256 with JWT_SECRET
  alt invalid or expired token
    P-->>U: 401 Invalid or expired token
  else valid
    P->>S: get_seo_service
    alt no source configured
      P->>P: SEOAnalyzer.get_domain_metrics mock
    else source configured
      S->>X: httpx or Google client, per source
      X-->>S: metrics
    end
    P-->>U: 200 metrics object
    U->>U: seo page unwraps metrics, maps search_volume
  end
```

Sources: `web/src/app/(dashboard)/seo/page.tsx`, `web/src/lib/api-client.ts` (`seoApi`), `api-python/app/main.py`, `api-python/app/core/auth.py`, `api-python/app/api/seo.py`, `api-python/app/services/seo_data_sources.py`, `api-python/app/services/seo_analyzer.py`

## 3. Admin ส่ง email campaign แล้วผู้รับเปิดและคลิก

```mermaid
sequenceDiagram
  autonumber
  actor A as Admin browser
  participant G as api-go
  participant DB as PostgreSQL
  participant M as EmailService SendGrid, SES or stub
  actor R as Recipient mail client
  A->>G: POST /api/v1/email/campaigns/:id/send, Bearer, role admin
  G->>DB: load campaign and email_campaign_lists
  G->>DB: UPDATE email_campaigns status sending
  G->>DB: SELECT DISTINCT active subscribers of linked lists
  loop each subscriber, inside the HTTP request
    G->>DB: INSERT email_sends pending
    G->>M: SendEmail to, name, subject, html, text
    alt send error
      G->>DB: UPDATE email_sends bounced, send_error
    else ok
      G->>DB: UPDATE email_sends sent, INSERT email_events sent
    end
  end
  G->>DB: UPDATE email_campaigns status sent or failed, sent_count
  G-->>A: response
  R->>G: GET /api/v1/email/track/open?send_id=..., no auth
  G->>DB: TrackEmailEvent opened
  G-->>R: 1x1 GIF
  R->>G: GET /api/v1/email/track/click?send_id=...&url=...
  G->>DB: TrackEmailEvent clicked
  G-->>R: 302 to url as given
```

Sources: `api-go/internal/delivery/http/router.go` (`email` admin group, `emailTrack` public group), `api-go/internal/delivery/http/email_campaign_handler.go`, `api-go/internal/application/service/email_campaign_service.go`, `api-go/internal/application/service/email_factory.go`, `sendgrid_email_service.go`, `ses_email_service.go`

ข้อสังเกต: `email_campaigns.status` CHECK ใน `004_marketing_features.up.sql` ไม่มีค่า `failed` แต่ service เขียน `failed`; การ UPDATE นั้นจะถูก constraint ปฏิเสธ (อ่านจาก SQL ยังไม่ได้รันยืนยัน) ส่วนการ insert `email_sends` ที่ซ้ำจะถูก `continue` ข้ามเงียบ ๆ

## 4. Social post แบบตั้งเวลา (in-memory timer)

```mermaid
sequenceDiagram
  autonumber
  actor U as User browser
  participant G as api-go SocialPublishService
  participant T as SocialScheduler in-memory
  participant DB as PostgreSQL
  participant F as Facebook Graph or other platform
  U->>G: POST /api/v1/social-pub/posts with scheduled_at
  G->>DB: INSERT social_publish_posts scheduled
  G->>T: SchedulePost time.AfterFunc delay
  G-->>U: 201 post
  Note over T: Timer lives only in this process, a restart loses it
  T->>G: PublishPost when timer fires
  G->>DB: SELECT post and account access_token, token_expiry
  opt token expired
    G->>F: refresh token
  end
  G->>F: SchedulePost or PublishPost with access_token
  alt error
    G->>DB: UPDATE status failed, error_msg
  else ok
    G->>DB: UPDATE status published, platform_post_id
  end
```

Sources: `api-go/internal/delivery/http/router.go` (`/social-pub` only when `FACEBOOK_APP_ID` set), `api-go/internal/delivery/http/social_handler.go`, `api-go/internal/application/service/social_publishing_service.go`, `api-go/internal/application/service/social_scheduler.go`, `facebook_client.go`

## Open questions

- การส่ง email/SMS ทำ synchronous ใน request เดียว (SMS ดึงสูงสุด 10,000 subscribers): ต้องการ queue/worker หรือไม่
- Scheduled social posts ต้องอยู่รอดข้าม restart หรือไม่ (ตอนนี้ไม่มี job ที่อ่าน `scheduled_at` จาก DB ตอน start)
