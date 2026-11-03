# ADR-0005 — External providers เปิดด้วย env var, ไม่มี env = stub หรือปิด

- **Status:** accepted (retro-documented)
- **วันที่ตัดสินใจ:** มีตั้งแต่ `2c50c4e`; การต่อสาย mailer จริงเข้ากับ email campaigns แก้ใน `b934a69` (README §api-go: variable shadowing)

## Context

ระบบต้อง build และ demo ได้โดยไม่มี credential ของ provider ใด

## Decision

- Email: `NewEmailService` เลือกตาม `EMAIL_PROVIDER` = `sendgrid` (ต้องมี `SENDGRID_API_KEY`, `EMAIL_FROM_ADDRESS`) / `ses` (`AWS_REGION` default `us-east-1`, `EMAIL_FROM_ADDRESS`) / อื่น ๆ = stub; config ไม่ครบ = `log.Fatalf` ตอน start
- SMS: `NewTwilioClient` คืน `nil` ถ้าไม่มี `TWILIO_ACCOUNT_SID`/`TWILIO_AUTH_TOKEN`; `SendCampaign` จึงบันทึก `stub_<id>` เป็น sent
- Social publish: route `/social-pub` ถูก register เฉพาะเมื่อมี `FACEBOOK_APP_ID` + `FACEBOOK_APP_SECRET`; TikTok/X/LinkedIn ต่อเพิ่มเมื่อมี key ของตัวเอง
- Marketplace affiliate: Shopee (`SHOPEE_PARTNER_ID`, `SHOPEE_PARTNER_KEY`), TikTok (`TIKTOK_APP_ID`, `TIKTOK_APP_SECRET`)
- SEO: `SEODataService` เปิด GSC / Ahrefs / SEMrush ตาม `GSC_CREDENTIALS_PATH` / `AHREFS_API_TOKEN` / `SEMRUSH_API_KEY`; ไม่มีเลย = `SEOAnalyzer` mock
- Sentry web: `enabled: Boolean(NEXT_PUBLIC_SENTRY_DSN)`

## Consequences

- ระบบ "สำเร็จ" ได้โดยไม่ส่งอะไรจริง (stub email/SMS, mock SEO, random analytics): ต้องดู log หรือ env เพื่อรู้ว่าอยู่โหมดไหน
- การเปลี่ยน env มีผลเฉพาะตอน process start (`NewRouter` อ่าน env ครั้งเดียว) ต้อง restart/redeploy
- Callback URL ของ X และ LinkedIn hard-code เป็น `http://localhost:3000/settings/social/callback?...`

## Evidence

`api-go/internal/application/service/email_factory.go`, `twilio_service.go`, `sms_service.go`, `api-go/internal/delivery/http/router.go`, `api-python/app/services/seo_data_sources.py`, `api-python/app/api/seo.py`, `api-python/app/services/analytics_service.py`, `web/src/lib/sentry.ts`, `README.md` §api-go
