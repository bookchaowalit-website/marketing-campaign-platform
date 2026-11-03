# 05 — C4 Context

ระบบเดียว (BookMarketing dashboard + APIs) ยังไม่มี host จริง ภาพนี้คือ context ตามโค้ด: external systems จะถูกเรียกก็ต่อเมื่อมี env var ของ provider นั้น (ไม่มี = ปิดหรือ stub)

```mermaid
C4Context
  title BookMarketing system context
  Person(user, "Marketer", "role user, JWT")
  Person(admin, "Admin", "role admin")
  Person(recipient, "Email recipient", "opens and clicks campaign mail")
  Person_Ext(partner, "Referral partner site", "posts affiliate track events")
  System(bm, "BookMarketing", "Next.js dashboard, Go core API, Python analytics API, PostgreSQL")
  System_Ext(mail, "SendGrid or Amazon SES", "EMAIL_PROVIDER, default stub")
  System_Ext(twilio, "Twilio", "SMS, TWILIO_ACCOUNT_SID")
  System_Ext(social, "Facebook, Instagram, TikTok, X, LinkedIn APIs", "OAuth publish, FACEBOOK_APP_ID")
  System_Ext(market, "Shopee and TikTok Shop affiliate APIs", "SHOPEE_PARTNER_ID, TIKTOK_APP_ID")
  System_Ext(seo, "Google Search Console, Ahrefs, SEMrush", "GSC_CREDENTIALS_PATH, AHREFS_API_TOKEN, SEMRUSH_API_KEY")
  System_Ext(sentry, "Sentry", "web errors, NEXT_PUBLIC_SENTRY_DSN")
  System_Ext(gh, "GitHub Actions", "CI and deploy-smoke")
  Rel(user, bm, "uses", "HTTPS browser")
  Rel(admin, bm, "manages email, SMS, PPC, integrations")
  Rel(recipient, bm, "pixel and click", "GET /api/v1/email/track")
  Rel(partner, bm, "track click, signup, conversion", "POST /api/v1/affiliates/track")
  Rel(bm, mail, "sends auth and campaign mail")
  Rel(bm, twilio, "sends SMS")
  Rel(bm, social, "OAuth and publish posts")
  Rel(bm, market, "OAuth and sync orders")
  Rel(bm, seo, "reads SEO metrics")
  Rel(bm, sentry, "error events")
  Rel(gh, bm, "builds and tests; smoke GET when PRODUCTION_URL set")
```

Sources: `api-go/internal/delivery/http/router.go` (env gates), `api-go/internal/application/service/email_factory.go`, `twilio_service.go`, `facebook_client.go`, `instagram_client.go`, `tiktok_client.go`, `twitter_client.go`, `linkedin_client.go`, `shopee_integration.go`, `tiktok_integration.go`, `api-python/app/services/seo_data_sources.py`, `web/src/lib/sentry.ts`, `.github/workflows/ci.yml`, `.github/workflows/deploy-smoke.yml`

หมายเหตุ: Instagram ใช้ `graph.facebook.com` ผ่าน Facebook app เดียวกัน; TikTok มีสอง client (`tiktok_client.go` สำหรับ publish, `tiktok_integration.go` สำหรับ affiliate)

## Open questions

- Provider ใดที่ owner ตั้งใจเปิดจริงเมื่อ deploy (ตอนนี้ไม่มี env ใดถูกตั้ง เพราะไม่มี environment)
- Partner site ที่จะเรียก `/affiliates/track/*` คือใคร: endpoint เปิดสาธารณะไม่มี auth หรือ signature
