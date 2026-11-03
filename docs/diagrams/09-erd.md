# 09 — ERD

Schema หลักคือ `api-go/migrations/001..016_*.up.sql` (ADR-0002) รวมราว 140 ตาราง ภาพด้านล่างวาดเฉพาะตารางที่ repository/service ใน `api-go` อ่านหรือเขียน แสดงเฉพาะ field สำคัญ (key, status, FK) ตารางที่ไม่มีโค้ดใช้แสดงเป็นรายการท้ายไฟล์
หมายเหตุสำคัญ: คอลัมน์ `organization_id` / `org_id` ในตาราง feature ส่วนใหญ่ **ไม่มี FK** ไป `organizations` (มี FK เฉพาะ `organization_members` และ `campaigns.organization_id` จาก 016)

## Identity และ tenancy

```mermaid
erDiagram
  users ||--o{ refresh_tokens : has
  users ||--o{ password_reset_tokens : has
  users ||--o{ email_verification_tokens : has
  users ||--o{ organization_members : joins
  organizations ||--o{ organization_members : has
  users {
    uuid id PK
    varchar email UK
    varchar password_hash
    varchar full_name "added 008"
    varchar role "user or admin"
    boolean is_active
    boolean email_verified "added 010"
  }
  organizations {
    uuid id PK
    varchar name
    varchar slug UK
    uuid plan_id FK
    jsonb settings
  }
  organization_members {
    uuid id PK
    uuid organization_id FK
    uuid user_id FK
    varchar role "owner admin member viewer"
  }
  refresh_tokens {
    uuid id PK
    uuid user_id FK
    varchar token_hash "sha256"
    timestamp expires_at
    boolean revoked
  }
  token_blacklist {
    uuid id PK
    varchar token_hash
    timestamp expires_at
  }
  password_reset_tokens {
    uuid id PK
    uuid user_id FK
    varchar token
    timestamp expires_at
    boolean used
  }
  email_verification_tokens {
    uuid id PK
    uuid user_id FK
    varchar token
    timestamp expires_at
    boolean used
  }
```

Sources: `api-go/migrations/001_init.up.sql`, `005_organizations.up.sql`, `006_billing.up.sql` (`fk_organizations_plan_id`), `008_fix_users_full_name.up.sql`, `009_refresh_tokens.up.sql`, `010_password_reset_email_verification.up.sql`, `api-go/internal/infrastructure/repository/auth_repository.go`, `user_repository.go`

## Campaigns, A/B, approvals, budget, calendar, PPC

```mermaid
erDiagram
  users ||--o{ campaigns : owns
  organizations ||--o{ campaigns : "org, added 016"
  campaigns ||--o{ ab_tests : has
  users ||--o{ ab_tests : owns
  campaign_groups ||--o{ unified_campaignes : groups
  users ||--o{ unified_campaignes : "owner, approver"
  unified_campaignes ||--o{ marketing_calendar : schedules
  unified_campaignes ||--o{ budget_transactions : spends
  budget_periods ||--o{ budget_allocations : splits
  budget_allocations ||--o{ budget_transactions : records
  approval_workflows ||--o{ approval_steps : has
  approval_workflows ||--o{ approval_requests : governs
  approval_requests ||--o{ approval_responses : collects
  approval_steps ||--o{ approval_responses : at
  ad_platforms ||--o{ ad_campaigns : runs
  unified_campaignes ||--o{ ad_campaigns : links
  campaigns {
    uuid id PK
    uuid user_id FK
    uuid organization_id FK
    varchar type "SEO SEM Social Email Display"
    varchar status "active paused completed draft"
    decimal budget
    decimal spent
  }
  ab_tests {
    uuid id PK
    uuid user_id FK
    uuid campaign_id FK
    varchar variant "A or B"
    varchar status "running completed draft"
  }
  unified_campaignes {
    uuid id PK
    uuid organization_id "no FK"
    uuid group_id FK
    varchar channel_category "online offline hybrid"
    varchar channel
    varchar status
    varchar approval_status
    uuid approved_by FK
  }
  approval_requests {
    uuid id PK
    uuid organization_id "no FK"
    uuid workflow_id FK
    varchar entity_type
    uuid entity_id
    varchar status
  }
  approval_responses {
    uuid id PK
    uuid request_id FK
    uuid step_id FK
    uuid approver_id FK
    varchar action
  }
  budget_periods {
    uuid id PK
    uuid organization_id "no FK"
    varchar period_type
    varchar status "draft approved active closed"
  }
  ad_campaigns {
    uuid id PK
    uuid org_id "no FK"
    uuid unified_campaign_id FK
    uuid platform_id FK
    varchar status "draft active paused ended archived"
  }
```

Sources: `api-go/migrations/001_init.up.sql`, `012_phase1_unified_campaigns.up.sql`, `013_ppc_ads.up.sql`, `016_campaigns_organization_id.up.sql`, `api-go/internal/infrastructure/repository/campaign_repository.go`, `abtest_repository.go`, `unified_campaign_repository.go`, `approval_repository.go`, `budget_repository.go`, `calendar_repository.go`, `ads_repository.go`

ชื่อตาราง `unified_campaignes` (สะกดแบบนี้จริงทั้งใน migration และ SQL ของ repository)

## Email, SMS, social

```mermaid
erDiagram
  email_lists ||--o{ email_list_subscribers : contains
  email_subscribers ||--o{ email_list_subscribers : in
  email_campaigns ||--o{ email_campaign_lists : targets
  email_lists ||--o{ email_campaign_lists : targeted
  email_campaigns ||--o{ email_sends : produces
  email_subscribers ||--o{ email_sends : receives
  email_campaigns ||--o{ email_events : logs
  email_campaigns ||--o{ email_campaign_ab_tests : "variant a or b"
  sms_lists ||--o{ sms_subscribers : contains
  sms_lists ||--o{ sms_campaigns : targets
  sms_campaigns ||--o{ sms_logs : produces
  sms_subscribers ||--o{ sms_logs : receives
  social_publish_accounts ||--o{ social_publish_posts : publishes
  email_campaigns {
    uuid id PK
    uuid organization_id "no FK"
    varchar subject
    varchar status "draft scheduled sending sent paused cancelled"
    int sent_count
  }
  email_sends {
    uuid id PK
    uuid campaign_id FK
    uuid subscriber_id FK
    varchar status
    text send_error
  }
  email_subscribers {
    uuid id PK
    varchar email
    varchar status "active unsubscribed bounced complained pending"
  }
  sms_campaigns {
    uuid id PK
    uuid org_id "no FK"
    uuid list_id FK
    varchar status "draft scheduled sending sent failed cancelled"
  }
  sms_logs {
    uuid id PK
    uuid campaign_id FK
    uuid subscriber_id FK
    varchar status
    varchar message_sid
  }
  social_publish_accounts {
    uuid id PK
    uuid organization_id "no FK"
    varchar platform "facebook instagram twitter linkedin"
    text access_token "plaintext"
    text refresh_token "plaintext"
  }
  social_publish_posts {
    uuid id PK
    uuid account_id FK
    varchar status "draft scheduled published failed"
    timestamp scheduled_at
  }
```

Sources: `api-go/migrations/004_marketing_features.up.sql`, `011_social_media_publishing.up.sql`, `014_sms_marketing.up.sql`, `015_email_complete.up.sql`, `api-go/internal/application/service/email_campaign_service.go`, `marketing_features_service.go`, `social_publishing_service.go`, `api-go/internal/infrastructure/repository/sms_repository.go`

## Affiliates และ external affiliate integrations

```mermaid
erDiagram
  users ||--o{ affiliates : becomes
  affiliates ||--o{ affiliate_links : creates
  campaigns ||--o{ affiliate_links : promotes
  affiliates ||--o{ referrals : earns
  affiliate_links ||--o{ referrals : via
  referrals ||--o{ commissions : pays
  affiliate_payouts ||--o{ commissions : settles
  users ||--o{ external_affiliate_accounts : connects
  external_affiliate_accounts ||--o{ external_affiliate_products : syncs
  external_affiliate_accounts ||--o{ external_affiliate_orders : syncs
  external_affiliate_accounts ||--o{ external_affiliate_sync_logs : logs
  affiliates {
    uuid id PK
    uuid user_id FK
    varchar affiliate_code UK
    varchar status "active pending suspended approved"
    varchar bank_account
    varchar tax_id
  }
  referrals {
    uuid id PK
    uuid affiliate_id FK
    varchar status "pending signup qualified converted cancelled"
    varchar commission_status
  }
  commissions {
    uuid id PK
    uuid affiliate_id FK
    uuid referral_id FK
    varchar status "pending approved paid cancelled refunded"
  }
  external_affiliate_accounts {
    uuid id PK
    uuid user_id FK
    varchar platform "shopee tiktok lazada amazon custom"
    text api_secret "plaintext"
    text access_token "plaintext"
    varchar status
  }
```

Sources: `api-go/migrations/002_affiliate.up.sql`, `003_external_affiliate.up.sql`, `api-go/internal/infrastructure/repository/affiliate_repository.go`, `affiliate_link_repository.go`, `referral_repository.go`, `commission_repository.go`, `external_affiliate_*_repository.go`

## ตารางที่มีใน migration แต่ไม่พบโค้ด Go/Python อ่านหรือเขียน

| Migration | ตาราง |
|---|---|
| `001_init` | `seo_metrics`, `keyword_rankings`, `analytics_events` (Python model มีแค่ `users`, `campaigns`, `ab_tests`, `competitors`) |
| `004_marketing_features` | templates, landing page components/forms/A-B, automation steps/contacts/events/goals, contact tags, influencer profiles/collaborations/content/payments, custom reports/widgets/dashboards, attribution |
| `005_organizations` | `team_invitations`, `roles`, `api_keys`, `organization_settings`, `audit_logs` |
| `006_billing` | 13 ตาราง (`subscription_plans` … `billing_webhook_events`) ไม่มี billing code |
| `007_crm` | 15 ตาราง (`contacts`, `deals`, `pipelines`, `tasks`, `notes` …) ไม่มี CRM handler |

ข้อมูลส่วนบุคคลที่ถูกเก็บ: email/ชื่อ/เบอร์โทรของ subscriber, IP และ user agent (`email_events`, `referrals`), ข้อมูลธนาคารและ `tax_id` ของ affiliate, OAuth/API secrets ของ provider (plaintext `TEXT`, ไม่พบโค้ด encrypt)

## Open questions

- Schema ใดเป็น source of truth ระหว่าง `api-go/migrations`, `docker/migrations`, `api-python/migrations` (ADR-0002)
- ตาราง billing/CRM ที่ไม่มีโค้ดควรเก็บไว้หรือย้ายไป spec อนาคต
- ต้องการเข้ารหัส `access_token`, `api_secret`, `bank_account` at rest หรือไม่
