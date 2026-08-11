package http

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/bookmarketing/marketing-saas/api-go/internal/infrastructure/repository"
	"github.com/gin-gonic/gin"
)

// Router holds all handlers and sets up routes
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new router with all handlers
func NewRouter(db *sql.DB) *Router {
	// Initialize repositories
	campaignRepo := repository.NewCampaignRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	userRepo := repository.NewUserRepository(db)
	abTestRepo := repository.NewABTestRepository(db)
	affiliateRepo := repository.NewAffiliateRepository(db)
	affiliateLinkRepo := repository.NewAffiliateLinkRepository(db)
	referralRepo := repository.NewReferralRepository(db)
	commissionRepo := repository.NewCommissionRepository(db)
	affiliateSettingsRepo := repository.NewAffiliateSettingsRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	passwordResetRepo := repository.NewPasswordResetTokenRepository(db)
	emailVerificationRepo := repository.NewEmailVerificationTokenRepository(db)

	// Budget repositories
	budgetPeriodRepo := repository.NewBudgetPeriodRepository(db)
	budgetAllocRepo := repository.NewBudgetAllocationRepository(db)
	budgetTxnRepo := repository.NewBudgetTransactionRepository(db)
	budgetAlertRepo := repository.NewBudgetAlertRepository(db)

	// External affiliate repositories
	extAccountRepo := repository.NewExternalAffiliateAccountRepository(db)
	extProductRepo := repository.NewExternalAffiliateProductRepository(db)
	extOrderRepo := repository.NewExternalAffiliateOrderRepository(db)
	extEarningsRepo := repository.NewExternalAffiliateEarningsRepository(db)
	extSyncLogRepo := repository.NewExternalAffiliateSyncLogRepository(db)

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// Get frontend URL for email links
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Get platform API credentials from environment
	shopeePartnerID := os.Getenv("SHOPEE_PARTNER_ID")
	shopeePartnerKey := os.Getenv("SHOPEE_PARTNER_KEY")
	shopeeRedirectURL := os.Getenv("SHOPEE_REDIRECT_URL")
	
	tiktokAppID := os.Getenv("TIKTOK_APP_ID")
	tiktokAppSecret := os.Getenv("TIKTOK_APP_SECRET")
	tiktokRedirectURL := os.Getenv("TIKTOK_REDIRECT_URL")

	// Initialize services
	campaignService := service.NewCampaignService(campaignRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)
	
	// Initialize email service (supports SendGrid, SES, or stub)
	emailService, err := service.NewEmailService(frontendURL)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}
	
	// Initialize social publishing service
	fbAppID := os.Getenv("FACEBOOK_APP_ID")
	fbAppSecret := os.Getenv("FACEBOOK_APP_SECRET")
	var socialPublishService *service.SocialPublishService
	if fbAppID != "" && fbAppSecret != "" {
		socialPublishService = service.NewSocialPublishService(db, fbAppID, fbAppSecret)

		// Register TikTok client
		tiktokClientID := os.Getenv("TIKTOK_CLIENT_ID")
		tiktokClientSecret := os.Getenv("TIKTOK_CLIENT_SECRET")
		if tiktokClientID != "" && tiktokClientSecret != "" {
			socialPublishService.RegisterTikTokClient(
				service.NewTikTokClient(tiktokClientID, tiktokClientSecret, tiktokRedirectURL),
			)
			log.Println("TikTok client registered")
		}

		// Register Twitter/X client
		twitterAPIKey := os.Getenv("TWITTER_API_KEY")
		twitterAPISecret := os.Getenv("TWITTER_API_SECRET")
		twitterBearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
		if twitterAPIKey != "" && twitterAPISecret != "" {
			socialPublishService.RegisterTwitterClient(
				service.NewTwitterClient(twitterAPIKey, twitterAPISecret, twitterBearerToken, "http://localhost:3000/settings/social/callback?platform=twitter"),
			)
			log.Println("Twitter/X client registered")
		}

		// Register LinkedIn client
		linkedinClientID := os.Getenv("LINKEDIN_CLIENT_ID")
		linkedinClientSecret := os.Getenv("LINKEDIN_CLIENT_SECRET")
		if linkedinClientID != "" && linkedinClientSecret != "" {
			socialPublishService.RegisterLinkedInClient(
				service.NewLinkedInClient(linkedinClientID, linkedinClientSecret, "http://localhost:3000/settings/social/callback?platform=linkedin"),
			)
			log.Println("LinkedIn client registered")
		}

		log.Println("Social publishing service initialized")
	} else {
		log.Println("Social publishing service disabled (set FACEBOOK_APP_ID and FACEBOOK_APP_SECRET to enable)")
	}
	
	authService := service.NewAuthService(userRepo, refreshTokenRepo, passwordResetRepo, emailVerificationRepo, emailService, jwtSecret, frontendURL)
	abTestService := service.NewABTestService(abTestRepo)
	affiliateService := service.NewAffiliateService(affiliateRepo, affiliateLinkRepo, referralRepo, commissionRepo, affiliateSettingsRepo)
	
	// External affiliate service
	extAffiliateService := service.NewExternalAffiliateService(
		extAccountRepo, extProductRepo, extOrderRepo, extEarningsRepo, extSyncLogRepo,
	)
	
	// Register platform integrations
	if shopeePartnerID != "" && shopeePartnerKey != "" {
		shopeeIntegration := service.NewShopeeIntegration(
			shopeePartnerID, shopeePartnerKey, shopeeRedirectURL,
			extAccountRepo, extProductRepo, extOrderRepo, extEarningsRepo, extSyncLogRepo,
		)
		extAffiliateService.RegisterIntegration(domain.PlatformShopee, shopeeIntegration)
	}
	
	if tiktokAppID != "" && tiktokAppSecret != "" {
		tiktokIntegration := service.NewTikTokIntegration(
			tiktokAppID, tiktokAppSecret, tiktokRedirectURL,
			extAccountRepo, extProductRepo, extOrderRepo, extEarningsRepo, extSyncLogRepo,
		)
		extAffiliateService.RegisterIntegration(domain.PlatformTikTok, tiktokIntegration)
	}

	// Initialize handlers
	campaignHandler := NewCampaignHandler(campaignService)
	dashboardHandler := NewDashboardHandler(dashboardService)
	authHandler := NewAuthHandler(authService)
	abTestHandler := NewABTestHandler(abTestService)
	affiliateHandler := NewAffiliateHandler(affiliateService)
	extAffiliateHandler := NewExternalAffiliateHandler(extAffiliateService)
	
	// Initialize unified campaign handler
	unifiedCampaignRepo := repository.NewUnifiedCampaignRepository(db)
	unifiedCampaignSvc := service.NewUnifiedCampaignService(unifiedCampaignRepo)
	unifiedCampaignHandler := NewUnifiedCampaignHandler(unifiedCampaignSvc)
	
	// Initialize social handler (only if social service is available)
	var socialHandler *SocialHandler
	if socialPublishService != nil {
		socialHandler = NewSocialHandler(socialPublishService)
	}

	// Initialize marketing features services
	// Note: emailMarketingService manages lists/subscribers/campaigns (DB only);
	// it does not implement domain.EmailService. Actual sending uses the real
	// mailer (SendGrid/SES/stub) constructed above as `emailService`.
	emailMarketingService := service.NewEmailMarketingService(db)
	socialService := service.NewSocialMediaService(db)
	landingService := service.NewLandingPageService(db)
	automationService := service.NewAutomationService(db)
	influencerService := service.NewInfluencerService(db)
	reportingService := service.NewReportingService(db)

	// Initialize email campaign service (send, track, A/B testing)
	emailCampaignSvc := service.NewEmailCampaignService(db, emailService)
	emailCampaignHandler := NewEmailCampaignHandler(emailCampaignSvc)

	// Initialize marketing features handler
	marketingHandler := NewMarketingFeaturesHandler(emailMarketingService, socialService, landingService, automationService, influencerService, reportingService)

	// Initialize budget service and handler
	budgetService := service.NewBudgetService(budgetPeriodRepo, budgetAllocRepo, budgetTxnRepo, budgetAlertRepo)
	budgetHandler := NewBudgetHandler(budgetService)

	// Initialize calendar service and handler
	calendarEventRepo := repository.NewCalendarEventRepository(db)
	calendarViewRepo := repository.NewCalendarViewRepository(db)
	calendarService := service.NewCalendarService(calendarEventRepo, calendarViewRepo)
	calendarHandler := NewCalendarHandler(calendarService)

	// Initialize approval service and handler
	approvalWorkflowRepo := repository.NewApprovalWorkflowRepository(db)
	approvalStepRepo := repository.NewApprovalStepRepository(db)
	approvalRequestRepo := repository.NewApprovalRequestRepository(db)
	approvalResponseRepo := repository.NewApprovalResponseRepository(db)
	approvalService := service.NewApprovalService(approvalWorkflowRepo, approvalStepRepo, approvalRequestRepo, approvalResponseRepo)
	approvalHandler := NewApprovalHandler(approvalService)

	// Initialize PPC Ads service and handler
	adPlatformRepo := repository.NewAdPlatformRepository(db)
	adCampaignRepo := repository.NewAdCampaignRepository(db)
	adGroupRepo := repository.NewAdGroupRepository(db)
	adCreativeRepo := repository.NewAdCreativeRepository(db)
	adKeywordRepo := repository.NewAdKeywordRepository(db)
	adPerformanceRepo := repository.NewAdPerformanceRepository(db)
	adsService := service.NewPPCAdsService(adPlatformRepo, adCampaignRepo, adGroupRepo, adCreativeRepo, adKeywordRepo, adPerformanceRepo)
	adsHandler := NewPPCAdsHandler(adsService)

	// Initialize SMS service and handler
	smsListRepo := repository.NewSMSListRepository(db)
	smsSubscriberRepo := repository.NewSMSSubscriberRepository(db)
	smsCampaignRepo := repository.NewSMSCampaignRepository(db)
	smsLogsRepo := repository.NewSMSLogRepository(db)
	twilioClient := service.NewTwilioClient()
	smsService := service.NewSMSService(smsListRepo, smsSubscriberRepo, smsCampaignRepo, smsLogsRepo, twilioClient)
	smsHandler := NewSMSHandler(smsService)

	// Set up Gin
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	// Set up CORS
	engine.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Rate limiter for auth endpoints (10 requests per minute per IP)
		authLimiter := NewRateLimiter(10, time.Minute)
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			for range ticker.C {
				authLimiter.Cleanup()
			}
		}()

		// Auth routes (public, rate limited)
		auth := v1.Group("/auth")
		auth.Use(RateLimitMiddleware(authLimiter, func(c *gin.Context) string {
			return c.ClientIP()
		}))
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-email", authHandler.VerifyEmail)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(AuthMiddleware(authService))
		{
			// Auth protected routes
			protected.GET("/auth/me", authHandler.GetMe)
			protected.POST("/auth/logout", authHandler.Logout)
			protected.POST("/auth/resend-verification", authHandler.ResendVerificationEmail)

			// Dashboard
			protected.GET("/dashboard", dashboardHandler.GetOverview)

			// Campaigns (legacy)
			campaigns := protected.Group("/campaigns")
			{
				campaigns.GET("", campaignHandler.GetAllCampaigns)
				campaigns.GET("/:id", campaignHandler.GetCampaign)
				campaigns.POST("", campaignHandler.CreateCampaign)
				campaigns.PUT("/:id", campaignHandler.UpdateCampaign)
				campaigns.DELETE("/:id", campaignHandler.DeleteCampaign)
				campaigns.POST("/:id/pause", campaignHandler.PauseCampaign)
				campaigns.POST("/:id/activate", campaignHandler.ActivateCampaign)
			}

			// Unified Campaigns (Phase 1 - supports ALL online + offline channels)
			unified := protected.Group("/unified-campaigns")
			{
				unified.GET("", unifiedCampaignHandler.ListCampaigns)
				unified.GET("/:id", unifiedCampaignHandler.GetCampaign)
				unified.POST("", unifiedCampaignHandler.CreateCampaign)
				unified.PUT("/:id", unifiedCampaignHandler.UpdateCampaign)
				unified.DELETE("/:id", unifiedCampaignHandler.DeleteCampaign)
				unified.POST("/:id/submit", unifiedCampaignHandler.SubmitForApproval)
				unified.POST("/:id/approve", unifiedCampaignHandler.ApproveCampaign)
				unified.POST("/:id/reject", unifiedCampaignHandler.RejectCampaign)
			}

			// A/B Tests
			abTests := protected.Group("/ab-tests")
			{
				abTests.GET("", abTestHandler.GetAllABTests)
				abTests.GET("/:id", abTestHandler.GetABTest)
				abTests.POST("", abTestHandler.CreateABTest)
				abTests.PUT("/:id", abTestHandler.UpdateABTest)
				abTests.DELETE("/:id", abTestHandler.DeleteABTest)
				abTests.POST("/:id/start", abTestHandler.StartABTest)
				abTests.POST("/:id/stop", abTestHandler.StopABTest)
			}

			// Affiliate Program
			affiliates := protected.Group("/affiliates")
			{
				// Affiliate registration and management (user)
				affiliates.POST("/register", affiliateHandler.RegisterAffiliate)
				affiliates.GET("/me", affiliateHandler.GetMyAffiliate)
				affiliates.GET("/dashboard", affiliateHandler.GetAffiliateDashboard)
				affiliates.GET("/stats", affiliateHandler.GetAffiliateStats)
				affiliates.GET("/links", affiliateHandler.GetAffiliateLinks)
				affiliates.POST("/links", affiliateHandler.CreateAffiliateLink)
				affiliates.DELETE("/links/:id", affiliateHandler.DeleteAffiliateLink)
				affiliates.GET("/commissions", affiliateHandler.GetCommissions)
			}

			// Affiliate admin routes
			affiliatesAdmin := protected.Group("/affiliates")
			affiliatesAdmin.Use(AdminOnlyMiddleware())
			{
				affiliatesAdmin.GET("", affiliateHandler.GetAllAffiliates)
				affiliatesAdmin.GET("/:id", affiliateHandler.GetAffiliate)
				affiliatesAdmin.PUT("/:id", affiliateHandler.UpdateAffiliate)
				affiliatesAdmin.POST("/:id/approve", affiliateHandler.ApproveAffiliate)
				affiliatesAdmin.POST("/:id/suspend", affiliateHandler.SuspendAffiliate)
				affiliatesAdmin.GET("/settings", affiliateHandler.GetAffiliateSettings)
				affiliatesAdmin.PUT("/settings", affiliateHandler.UpdateAffiliateSettings)
			}

			// External Affiliate Integrations (admin only)
			integrations := protected.Group("/integrations")
			integrations.Use(AdminOnlyMiddleware())
			{
				// Platform connections
				integrations.POST("/connect", extAffiliateHandler.ConnectPlatform)
				integrations.GET("/accounts", extAffiliateHandler.GetConnectedAccounts)
				integrations.GET("/accounts/:id", extAffiliateHandler.GetAccount)
				integrations.POST("/accounts/:id/disconnect", extAffiliateHandler.DisconnectPlatform)
				integrations.POST("/accounts/:id/sync", extAffiliateHandler.SyncAccount)
				integrations.GET("/accounts/:id/products", extAffiliateHandler.GetProducts)
				integrations.GET("/accounts/:id/orders", extAffiliateHandler.GetOrders)
				integrations.GET("/accounts/:id/logs", extAffiliateHandler.GetSyncLogs)
				integrations.POST("/accounts/:id/link", extAffiliateHandler.GenerateAffiliateLink)
				
				// Unified dashboard
				integrations.GET("/dashboard", extAffiliateHandler.GetUnifiedDashboard)
				
				// OAuth endpoints
				integrations.GET("/oauth/:platform", extAffiliateHandler.GetOAuthURL)
				integrations.GET("/oauth/:platform/callback", extAffiliateHandler.HandleOAuthCallback)
			}

			// Email Marketing (admin only)
			email := protected.Group("/email")
			email.Use(AdminOnlyMiddleware())
			{
				email.GET("/lists", marketingHandler.GetEmailLists)
				email.POST("/lists", marketingHandler.CreateEmailList)
				email.GET("/subscribers", marketingHandler.GetSubscribers)
				email.POST("/subscribers", marketingHandler.AddSubscriber)
				email.GET("/campaigns", marketingHandler.GetEmailCampaigns)
				email.POST("/campaigns", marketingHandler.CreateEmailCampaign)

				// Email Campaign Complete (send, analytics, A/B testing)
				email.POST("/campaigns/:id/send", emailCampaignHandler.SendCampaignEmail)
				email.GET("/campaigns/:id/analytics", emailCampaignHandler.GetCampaignAnalytics)
				email.POST("/ab-tests", emailCampaignHandler.CreateABTest)
				email.GET("/ab-tests", emailCampaignHandler.ListABTests)
				email.POST("/ab-tests/:id/complete", emailCampaignHandler.CompleteABTest)
				email.POST("/track", emailCampaignHandler.TrackEvent)
			}

			// Email tracking endpoints (public - for pixel and link tracking)
			emailTrack := v1.Group("/email/track")
			{
				emailTrack.GET("/open", emailCampaignHandler.TrackOpen)
				emailTrack.GET("/click", emailCampaignHandler.TrackClick)
			}

			// Social Media (Legacy - marketing_features_handler)
			social := protected.Group("/social")
			{
				social.GET("/posts", marketingHandler.GetSocialPosts)
				social.POST("/posts", marketingHandler.CreateSocialPost)
				social.POST("/posts/:id/schedule", marketingHandler.ScheduleSocialPost)
			}

			// Social Publishing (Facebook/Instagram - new)
			if socialHandler != nil {
				socialPub := protected.Group("/social-pub")
				{
					socialPub.GET("/oauth/:platform", socialHandler.GetOAuthURL)
					socialPub.POST("/connect", socialHandler.ConnectAccount)
					socialPub.GET("/accounts", socialHandler.GetAccounts)
					socialPub.POST("/accounts/:id/disconnect", socialHandler.DisconnectAccount)
					
					socialPub.GET("/posts", socialHandler.GetPosts)
					socialPub.GET("/posts/scheduled", socialHandler.GetScheduledPosts)
					socialPub.POST("/posts", socialHandler.CreatePost)
					socialPub.POST("/posts/:id/publish", socialHandler.PublishPost)
					socialPub.DELETE("/posts/:id", socialHandler.DeletePost)
					socialPub.GET("/posts/:id/analytics", socialHandler.GetAnalytics)
				}
			}

			// Landing Pages
			landingPages := protected.Group("/landing-pages")
			{
				landingPages.GET("", marketingHandler.GetLandingPages)
				landingPages.POST("", marketingHandler.CreateLandingPage)
				landingPages.POST("/:id/leads", marketingHandler.CreateLead)
			}

			// Marketing Automation
			automation := protected.Group("/automation")
			{
				automation.GET("/workflows", marketingHandler.GetWorkflows)
				automation.POST("/workflows", marketingHandler.CreateWorkflow)
			}

			// Influencer Marketing
			influencers := protected.Group("/influencers")
			{
				influencers.GET("", marketingHandler.GetInfluencers)
				influencers.POST("", marketingHandler.CreateInfluencer)
				influencers.POST("/campaigns", marketingHandler.CreateInfluencerCampaign)
			}

			// Reporting
			reporting := protected.Group("/reporting")
			{
				reporting.GET("/overview", marketingHandler.GetMarketingOverview)
				reporting.POST("/events", marketingHandler.TrackEvent)
			}

			// Budget Planning
			budget := protected.Group("/budget")
			{
				budget.POST("/periods", budgetHandler.CreateBudgetPeriod)
				budget.GET("/periods", budgetHandler.ListBudgetPeriods)
				budget.GET("/periods/:id", budgetHandler.GetBudgetPeriod)
				budget.POST("/periods/:id/allocations", budgetHandler.CreateBudgetAllocation)
				budget.PUT("/allocations/:id", budgetHandler.UpdateBudgetAllocation)
				budget.POST("/transactions", budgetHandler.RecordTransaction)
				budget.GET("/transactions", budgetHandler.ListTransactions)
				budget.GET("/alerts", budgetHandler.ListAlerts)
			}

			// Marketing Calendar
			calendar := protected.Group("/calendar")
			{
				calendar.POST("/events", calendarHandler.CreateEvent)
				calendar.GET("/events", calendarHandler.ListEvents)
				calendar.GET("/events/:id", calendarHandler.GetEvent)
				calendar.PUT("/events/:id", calendarHandler.UpdateEvent)
				calendar.DELETE("/events/:id", calendarHandler.DeleteEvent)
				calendar.POST("/views", calendarHandler.CreateView)
				calendar.GET("/views", calendarHandler.ListViews)
			}

			// Approval Workflows
			approval := protected.Group("/approval")
			{
				approval.POST("/workflows", approvalHandler.CreateWorkflow)
				approval.GET("/workflows", approvalHandler.ListWorkflows)
				approval.POST("/requests", approvalHandler.CreateRequest)
				approval.GET("/requests", approvalHandler.ListRequests)
				approval.GET("/requests/:id", approvalHandler.GetRequest)
				approval.POST("/requests/:id/respond", approvalHandler.Respond)
			}

			// PPC/Paid Ads Management (admin only)
			ads := protected.Group("/ads")
			ads.Use(AdminOnlyMiddleware())
			{
				ads.POST("/platforms/connect", adsHandler.ConnectPlatform)
				ads.GET("/platforms", adsHandler.ListPlatforms)
				ads.POST("/campaigns", adsHandler.CreateCampaign)
				ads.GET("/campaigns", adsHandler.ListCampaigns)
				ads.GET("/campaigns/:id", adsHandler.GetCampaign)
				ads.POST("/campaigns/:id/sync", adsHandler.SyncCampaign)
				ads.POST("/campaigns/:id/pause", adsHandler.PauseCampaign)
				ads.POST("/campaigns/:id/activate", adsHandler.ActivateCampaign)
			}

			// SMS Marketing (admin only)
			sms := protected.Group("/sms")
			sms.Use(AdminOnlyMiddleware())
			{
				sms.GET("/lists", smsHandler.ListSMSLists)
				sms.POST("/lists", smsHandler.CreateSMSList)
				sms.GET("/subscribers", smsHandler.ListSubscribers)
				sms.POST("/subscribers", smsHandler.AddSubscriber)
				sms.GET("/campaigns", smsHandler.ListCampaigns)
				sms.POST("/campaigns", smsHandler.CreateCampaign)
				sms.POST("/campaigns/:id/send", smsHandler.SendCampaign)
			}
		}

		// Public affiliate tracking endpoints
		v1.POST("/affiliates/track/click", affiliateHandler.TrackClick)
		v1.POST("/affiliates/track/signup", affiliateHandler.TrackSignup)
		v1.POST("/affiliates/track/conversion", affiliateHandler.TrackConversion)
	}

	return &Router{engine: engine}
}

// Run starts the HTTP server
func (r *Router) Run(port string) error {
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}
	return r.engine.Run(":" + port)
}

// GetEngine returns the gin engine for testing
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
