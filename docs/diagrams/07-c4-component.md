# 07 — C4 Component (api-go)

`api-go` ใช้ hexagonal layout: `delivery` (HTTP Gin + gRPC) → `application/service` → `domain` (entities + repository interfaces) ← `infrastructure/repository` (SQL) ทุกอย่างถูกประกอบใน `NewRouter(db)` (`router.go`) ไม่มี DI container

```mermaid
C4Component
  title api-go components
  Container_Ext(web, "web", "Next.js rewrites")
  Container_Ext(py, "api-python", "gRPC client")
  ContainerDb_Ext(db, "PostgreSQL")
  System_Ext(ext, "Providers")
  Container_Boundary(go, "api-go") {
    Component(mw, "Middleware", "Gin", "CORS reflect Origin, RateLimiter 10 per min on /auth, AuthMiddleware JWT, AdminOnlyMiddleware")
    Component(h, "HTTP handlers", "delivery/http", "auth, campaign, unified, abtest, affiliate, external affiliate, email campaign, marketing features, social, budget, calendar, approval, ads, sms, dashboard")
    Component(g, "gRPC server", "delivery/grpc", "CampaignService, AnalyticsService, ABTestService registered, reflection on")
    Component(auth, "AuthService", "application/service", "bcrypt, HS256 JWT 24 h, refresh 7 d sha256, blacklist, reset and verify tokens")
    Component(svc, "Domain services", "application/service", "Campaign, UnifiedCampaign, ABTest, Affiliate, Budget, Calendar, Approval, PPCAds, SMS, EmailCampaign, marketing features")
    Component(clients, "Provider clients", "application/service", "EmailService factory, TwilioClient, Facebook, Instagram, TikTok, Twitter, LinkedIn, Shopee, TikTok integration, SocialScheduler in-memory timers")
    Component(repo, "Repositories", "infrastructure/repository", "SQL per aggregate via database/sql")
  }
  Rel(web, mw, "HTTP /api/v1")
  Rel(mw, h, "passes context user_id, email, role")
  Rel(py, g, "gRPC")
  Rel(h, auth, "register, login, validate")
  Rel(h, svc, "calls")
  Rel(g, svc, "campaign and dashboard services")
  Rel(svc, repo, "domain repository interfaces")
  Rel(svc, clients, "send and publish")
  Rel(auth, repo, "users, refresh_tokens, token_blacklist")
  Rel(repo, db, "SQL")
  Rel(clients, ext, "HTTPS")
```

Sources: `api-go/cmd/main.go`, `api-go/internal/delivery/http/router.go`, `middleware.go`, `rate_limiter.go`, `api-go/internal/delivery/grpc/server.go`, `api-go/internal/application/service/*.go`, `api-go/internal/domain/repository.go`, `api-go/internal/infrastructure/repository/*.go`

หมายเหตุ

- บาง service (`EmailMarketingService`, `SocialMediaService`, `EmailCampaignService`, `SocialPublishService`) รับ `*sql.DB` ตรง ๆ แทน repository interface
- `main.go` สร้าง `CampaignService` และ `DashboardService` ชุดที่สองสำหรับ gRPC แยกจากชุดใน `NewRouter`
- gRPC `ABTestService` ถูก register แต่ `GRPCServer` embed แค่ `UnimplementedABTestServiceServer`: ทุก RPC ตอบ `Unimplemented`

## Component ของ api-python และ web (ย่อ)

```mermaid
flowchart LR
  subgraph PY["api-python"]
    MAIN["app/main.py create_app, CORS allow all"] --> DEP["core/auth.py get_current_user_id, HS256 JWT_SECRET"]
    DEP --> R1["api/competitors.py"]
    DEP --> R2["api/seo.py"]
    DEP --> R3["api/analytics.py"]
    DEP --> R4["api/grpc_client.py"]
    R1 --> S1["services/competitor_service.py, SQLAlchemy models"]
    R2 --> S2["services/seo_data_sources.py GSC, Ahrefs, SEMrush"]
    R2 --> S3["services/seo_analyzer.py mock fallback"]
    R3 --> S4["services/analytics_service.py random data"]
    R3 --> S5["services/predictive_analytics.py"]
  end
  subgraph WEB["web"]
    LAYOUT["app/(dashboard)/layout.tsx client auth guard"] --> PAGES["app/(dashboard)/*/page.tsx"]
    PAGES --> RQ["@tanstack/react-query"]
    RQ --> CLIENT["lib/api-client.ts axios apiGo, apiPython"]
    AUTHH["hooks/useAuth.tsx localStorage token"] --> CLIENT
    INST["instrumentation.ts, instrumentation-client.ts"] --> SCRUB["lib/sentry.ts scrubEvent"]
    HEALTH["app/api/health/route.ts"]
  end
```

Sources: `api-python/app/main.py`, `api-python/app/core/auth.py`, `api-python/app/api/*.py`, `api-python/app/services/*.py`, `web/src/app/(dashboard)/layout.tsx`, `web/src/hooks/useAuth.tsx`, `web/src/lib/api-client.ts`, `web/src/lib/sentry.ts`, `web/src/instrumentation.ts`, `web/src/instrumentation-client.ts`, `web/src/app/api/health/route.ts`

## Open questions

- ต้องการให้ gRPC `ABTestService` และ `CompetitorService` (มีใน proto แต่ Go ไม่ register CompetitorService) ทำงานจริงหรือเอาออกจาก proto
