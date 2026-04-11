package domain

import "time"

// =====================================================
// EMAIL MARKETING
// =====================================================

// EmailList represents an email subscriber list
type EmailList struct {
	ID                 string    `json:"id"`
	OrganizationID     string    `json:"organization_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	SubscriberCount    int       `json:"subscriber_count"`
	BounceCount        int       `json:"bounce_count"`
	UnsubscribeCount   int       `json:"unsubscribe_count"`
	DoubleOptin        bool      `json:"double_optin"`
	WelcomeEmailEnabled bool     `json:"welcome_email_enabled"`
	WelcomeEmailID     string    `json:"welcome_email_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// EmailSubscriber represents an email subscriber
type EmailSubscriber struct {
	ID            string     `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Phone         string     `json:"phone"`
	Company       string     `json:"company"`
	Status        string     `json:"status"`
	Source        string     `json:"source"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	CustomFields  string     `json:"custom_fields"` // JSON
	Tags          []string   `json:"tags"`
	SubscribedAt  *time.Time `json:"subscribed_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at"`
	BouncedAt     *time.Time `json:"bounced_at"`
	BounceReason  string     `json:"bounce_reason"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID           string    `json:"id"`
	OrganizationID string  `json:"organization_id"`
	Name         string    `json:"name"`
	Subject      string    `json:"subject"`
	PreviewText  string    `json:"preview_text"`
	FromName     string    `json:"from_name"`
	FromEmail    string    `json:"from_email"`
	ReplyTo      string    `json:"reply_to"`
	HTMLContent  string    `json:"html_content"`
	TextContent  string    `json:"text_content"`
	DesignJSON   string    `json:"design_json"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Category     string    `json:"category"`
	TemplateType string    `json:"template_type"`
	IsPublic     bool      `json:"is_public"`
	UsageCount   int       `json:"usage_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmailCampaign represents an email campaign
type EmailCampaign struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	Name             string     `json:"name"`
	Subject          string     `json:"subject"`
	PreviewText      string     `json:"preview_text"`
	FromName         string     `json:"from_name"`
	FromEmail        string     `json:"from_email"`
	ReplyTo          string     `json:"reply_to"`
	TemplateID       string     `json:"template_id"`
	HTMLContent      string     `json:"html_content"`
	TextContent      string     `json:"text_content"`
	CampaignType     string     `json:"campaign_type"`
	Status           string     `json:"status"`
	RecipientCount   int        `json:"recipient_count"`
	SentCount        int        `json:"sent_count"`
	DeliveredCount   int        `json:"delivered_count"`
	OpenedCount      int        `json:"opened_count"`
	ClickedCount     int        `json:"clicked_count"`
	BouncedCount     int        `json:"bounced_count"`
	UnsubscribedCount int       `json:"unsubscribed_count"`
	ComplainedCount  int        `json:"complained_count"`
	OpenRate         float64    `json:"open_rate"`
	ClickRate        float64    `json:"click_rate"`
	BounceRate       float64    `json:"bounce_rate"`
	UnsubscribeRate  float64    `json:"unsubscribe_rate"`
	ScheduleTime     *time.Time `json:"schedule_time"`
	SentAt           *time.Time `json:"sent_at"`
	TrackingEnabled  bool       `json:"tracking_enabled"`
	GoogleAnalytics  bool       `json:"google_analytics"`
	UTMSource        string     `json:"utm_source"`
	UTMMedium        string     `json:"utm_medium"`
	UTMCampaign      string     `json:"utm_campaign"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// EmailEvent represents an email event (open, click, etc.)
type EmailEvent struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	CampaignID     string     `json:"campaign_id"`
	SubscriberID   string     `json:"subscriber_id"`
	EventType      string     `json:"event_type"`
	Email          string     `json:"email"`
	IPAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
	URL            string     `json:"url"`
	LinkText       string     `json:"link_text"`
	BounceType     string     `json:"bounce_type"`
	BounceReason   string     `json:"bounce_reason"`
	UserData       string     `json:"user_data"` // JSON
	CreatedAt      time.Time  `json:"created_at"`
}

// =====================================================
// SOCIAL MEDIA
// =====================================================

// SocialPlatform types
type SocialPlatform string

const (
	SocialPlatformFacebook  SocialPlatform = "facebook"
	SocialPlatformInstagram SocialPlatform = "instagram"
	SocialPlatformTwitter   SocialPlatform = "twitter"
	SocialPlatformLinkedIn  SocialPlatform = "linkedin"
	SocialPlatformTikTok    SocialPlatform = "tiktok"
	SocialPlatformYouTube   SocialPlatform = "youtube"
	SocialPlatformPinterest SocialPlatform = "pinterest"
	SocialPlatformThreads   SocialPlatform = "threads"
)

// SocialAccount represents a connected social media account
type SocialAccount struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organization_id"`
	Platform       SocialPlatform `json:"platform"`
	AccountID      string        `json:"account_id"`
	AccountName    string        `json:"account_name"`
	AccountUsername string       `json:"account_username"`
	AccountAvatar  string        `json:"account_avatar"`
	AccountType    string        `json:"account_type"`
	AccessToken    string        `json:"-"` // Never expose
	RefreshToken   string        `json:"-"`
	TokenExpiresAt *time.Time    `json:"token_expires_at"`
	Status         string        `json:"status"`
	AutoPublish    bool          `json:"auto_publish"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// SocialPost represents a social media post
type SocialPost struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	MediaURLs        []string   `json:"media_urls"`
	MediaType        string     `json:"media_type"`
	LinkURL          string     `json:"link_url"`
	Hashtags         []string   `json:"hashtags"`
	Mentions         []string   `json:"mentions"`
	PlatformSpecific string     `json:"platform_specific"` // JSON
	Status           string     `json:"status"`
	IsRecurring      bool       `json:"is_recurring"`
	RecurringPattern string     `json:"recurring_pattern"`
	RecurringEndDate *time.Time `json:"recurring_end_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	PublishedAt      *time.Time `json:"published_at"`
}

// SocialPostSchedule represents a scheduled post
type SocialPostSchedule struct {
	ID                string     `json:"id"`
	PostID            string     `json:"post_id"`
	AccountID         string     `json:"account_id"`
	ScheduledTime     time.Time  `json:"scheduled_time"`
	PlatformPostID    string     `json:"platform_post_id"`
	PlatformPermalink string     `json:"platform_permalink"`
	Status            string     `json:"status"`
	ErrorMessage      string     `json:"error_message"`
	PublishedAt       *time.Time `json:"published_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// SocialAnalytics represents social media analytics
type SocialAnalytics struct {
	ID              string    `json:"id"`
	PostScheduleID  string    `json:"post_schedule_id"`
	AccountID       string    `json:"account_id"`
	Platform        string    `json:"platform"`
	Impressions     int64     `json:"impressions"`
	Reach           int64     `json:"reach"`
	Likes           int       `json:"likes"`
	Comments        int       `json:"comments"`
	Shares          int       `json:"shares"`
	Saves           int       `json:"saves"`
	Clicks          int       `json:"clicks"`
	VideoViews      int       `json:"video_views"`
	VideoViewRate   float64   `json:"video_view_rate"`
	EngagementRate  float64   `json:"engagement_rate"`
	EngagementCount int       `json:"engagement_count"`
	AnalyticsData   string    `json:"analytics_data"` // JSON
	FetchedAt       time.Time `json:"fetched_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// =====================================================
// LANDING PAGES
// =====================================================

// LandingPage represents a landing page
type LandingPage struct {
	ID              string     `json:"id"`
	OrganizationID  string     `json:"organization_id"`
	Name            string     `json:"name"`
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	MetaDescription string     `json:"meta_description"`
	FaviconURL      string     `json:"favicon_url"`
	CustomDomain    string     `json:"custom_domain"`
	HTMLContent     string     `json:"html_content"`
	DesignJSON      string     `json:"design_json"`
	CSSStyles       string     `json:"css_styles"`
	HeadScripts     string     `json:"head_scripts"`
	BodyScripts     string     `json:"body_scripts"`
	Status          string     `json:"status"`
	IsTemplate      bool       `json:"is_template"`
	TemplateCategory string    `json:"template_category"`
	ThumbnailURL    string     `json:"thumbnail_url"`
	ViewCount       int64      `json:"view_count"`
	UniqueVisitors  int64      `json:"unique_visitors"`
	ConversionCount int        `json:"conversion_count"`
	ConversionRate  float64    `json:"conversion_rate"`
	PublishedAt     *time.Time `json:"published_at"`
	SEOIndexable    bool       `json:"seo_indexable"`
	PasswordProtected bool      `json:"password_protected"`
	Password        string     `json:"password"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// LandingPageLead represents a lead from a landing page
type LandingPageLead struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	PageID         string     `json:"page_id"`
	FormID         string     `json:"form_id"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Phone          string     `json:"phone"`
	Company        string     `json:"company"`
	JobTitle       string     `json:"job_title"`
	CustomData     string     `json:"custom_data"` // JSON
	UTMSource      string     `json:"utm_source"`
	UTMMedium      string     `json:"utm_medium"`
	UTMCampaign    string     `json:"utm_campaign"`
	UTMTerm        string     `json:"utm_term"`
	UTMContent     string     `json:"utm_content"`
	IPAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
	ReferrerURL    string     `json:"referrer_url"`
	Status         string     `json:"status"`
	Value          float64    `json:"value"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// =====================================================
// MARKETING AUTOMATION
// =====================================================

// AutomationWorkflow represents an automation workflow
type AutomationWorkflow struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	TriggerType    string    `json:"trigger_type"`
	TriggerConfig  string    `json:"trigger_config"` // JSON
	Status         string    `json:"status"`
	TotalContacts  int       `json:"total_contacts"`
	ActiveContacts int       `json:"active_contacts"`
	CompletedContacts int    `json:"completed_contacts"`
	GoalsMet       int       `json:"goals_met"`
	Revenue        float64   `json:"revenue"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AutomationStep represents a step in an automation workflow
type AutomationStep struct {
	ID           string    `json:"id"`
	WorkflowID   string    `json:"workflow_id"`
	ParentStepID string    `json:"parent_step_id"`
	StepType     string    `json:"step_type"`
	Name         string    `json:"name"`
	OrderIndex   int       `json:"order_index"`
	Config       string    `json:"config"` // JSON
	DelayMinutes int       `json:"delay_minutes"`
	DelayType    string    `json:"delay_type"`
	Conditions   string    `json:"conditions"` // JSON
	Actions      string    `json:"actions"` // JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AutomationContact represents a contact in a workflow
type AutomationContact struct {
	ID            string     `json:"id"`
	WorkflowID    string     `json:"workflow_id"`
	SubscriberID  string     `json:"subscriber_id"`
	CurrentStepID string     `json:"current_step_id"`
	Status        string     `json:"status"`
	EnteredAt     time.Time  `json:"entered_at"`
	ExitedAt      *time.Time `json:"exited_at"`
	GoalMet       bool       `json:"goal_met"`
	GoalMetAt     *time.Time `json:"goal_met_at"`
	Revenue       float64    `json:"revenue"`
	StepHistory   string     `json:"step_history"` // JSON
	CustomData    string     `json:"custom_data"` // JSON
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ContactTag represents a tag for contacts
type ContactTag struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Color          string    `json:"color"`
	Description    string    `json:"description"`
	ContactCount   int       `json:"contact_count"`
	CreatedAt      time.Time `json:"created_at"`
}

// =====================================================
// INFLUENCER MARKETING
// =====================================================

// Influencer represents an influencer
type Influencer struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	Phone               string    `json:"phone"`
	Bio                 string    `json:"bio"`
	ProfileImage        string    `json:"profile_image"`
	Category            string    `json:"category"`
	Niches              []string  `json:"niches"`
	Location            string    `json:"location"`
	Country             string    `json:"country"`
	Language            string    `json:"language"`
	Gender              string    `json:"gender"`
	AgeRange            string    `json:"age_range"`
	Status              string    `json:"status"`
	Rating              float64   `json:"rating"`
	Notes               string    `json:"notes"`
	CustomFields        string    `json:"custom_fields"` // JSON
	TotalCollaborations int       `json:"total_collaborations"`
	TotalReach          int64     `json:"total_reach"`
	TotalEngagement     int64     `json:"total_engagement"`
	TotalEarnings       float64   `json:"total_earnings"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// InfluencerProfile represents an influencer's social profile
type InfluencerProfile struct {
	ID                      string    `json:"id"`
	InfluencerID            string    `json:"influencer_id"`
	Platform                string    `json:"platform"`
	Username                string    `json:"username"`
	ProfileURL              string    `json:"profile_url"`
	ProfileImage            string    `json:"profile_image"`
	Bio                     string    `json:"bio"`
	Followers               int64     `json:"followers"`
	Following               int64     `json:"following"`
	Posts                   int       `json:"posts"`
	AvgLikes                int       `json:"avg_likes"`
	AvgComments             int       `json:"avg_comments"`
	AvgShares               int       `json:"avg_shares"`
	EngagementRate          float64   `json:"engagement_rate"`
	Verified                bool      `json:"verified"`
	AudienceGenderMale      float64   `json:"audience_gender_male"`
	AudienceGenderFemale    float64   `json:"audience_gender_female"`
	AudienceAgeDistribution string    `json:"audience_age_distribution"` // JSON
	AudienceLocations       string    `json:"audience_locations"` // JSON
	LastSyncedAt            *time.Time `json:"last_synced_at"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// InfluencerCampaign represents an influencer campaign
type InfluencerCampaign struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	CampaignType     string     `json:"campaign_type"`
	Budget           float64    `json:"budget"`
	Spent            float64    `json:"spent"`
	Status           string     `json:"status"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	Goals            string     `json:"goals"` // JSON
	Requirements     string     `json:"requirements"`
	Deliverables     string     `json:"deliverables"`
	Hashtags         []string   `json:"hashtags"`
	Mentions         []string   `json:"mentions"`
	TrackingLinks    []string   `json:"tracking_links"`
	TotalReach       int64      `json:"total_reach"`
	TotalImpressions int64      `json:"total_impressions"`
	TotalEngagement  int64      `json:"total_engagement"`
	TotalClicks      int64      `json:"total_clicks"`
	TotalConversions int        `json:"total_conversions"`
	ROI              float64    `json:"roi"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// InfluencerCollaboration represents a collaboration with an influencer
type InfluencerCollaboration struct {
	ID                string     `json:"id"`
	CampaignID        string     `json:"campaign_id"`
	InfluencerID      string     `json:"influencer_id"`
	Status            string     `json:"status"`
	NegotiationAmount float64    `json:"negotiation_amount"`
	AgreedAmount      float64    `json:"agreed_amount"`
	Currency          string     `json:"currency"`
	PaymentStatus     string     `json:"payment_status"`
	PaidAmount        float64    `json:"paid_amount"`
	Deliverables      string     `json:"deliverables"` // JSON
	ContentDeadline   *time.Time `json:"content_deadline"`
	PublishDate       *time.Time `json:"publish_date"`
	Notes             string     `json:"notes"`
	OutreachDate      *time.Time `json:"outreach_date"`
	ResponseDate      *time.Time `json:"response_date"`
	AcceptedDate      *time.Time `json:"accepted_date"`
	CompletedDate     *time.Time `json:"completed_date"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// =====================================================
// REPORTING & ANALYTICS
// =====================================================

// CustomReport represents a custom report
type CustomReport struct {
	ID                string     `json:"id"`
	OrganizationID    string     `json:"organization_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	ReportType        string     `json:"report_type"`
	DataSources       string     `json:"data_sources"` // JSON
	Metrics           string     `json:"metrics"` // JSON
	Dimensions        string     `json:"dimensions"` // JSON
	Filters           string     `json:"filters"` // JSON
	DateRangeType     string     `json:"date_range_type"`
	CustomDateStart   *time.Time `json:"custom_date_start"`
	CustomDateEnd     *time.Time `json:"custom_date_end"`
	ComparisonEnabled bool       `json:"comparison_enabled"`
	ComparisonType    string     `json:"comparison_type"`
	ChartConfig       string     `json:"chart_config"` // JSON
	TableConfig       string     `json:"table_config"` // JSON
	IsScheduled       bool       `json:"is_scheduled"`
	ScheduleFrequency string     `json:"schedule_frequency"`
	ScheduleRecipients []string  `json:"schedule_recipients"`
	LastRunAt         *time.Time `json:"last_run_at"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CustomDashboard represents a custom dashboard
type CustomDashboard struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Layout         string    `json:"layout"` // JSON
	IsDefault      bool      `json:"is_default"`
	IsPublic       bool      `json:"is_public"`
	SharedWith     string    `json:"shared_with"` // JSON
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AnalyticsEvent represents a unified analytics event
type AnalyticsEvent struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	EventType      string    `json:"event_type"`
	EventCategory  string    `json:"event_category"`
	EventAction    string    `json:"event_action"`
	EventLabel     string    `json:"event_label"`
	EventValue     float64   `json:"event_value"`
	UserID         string    `json:"user_id"`
	SessionID      string    `json:"session_id"`
	DeviceType     string    `json:"device_type"`
	Browser        string    `json:"browser"`
	OS             string    `json:"os"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	IPAddress      string    `json:"ip_address"`
	PageURL        string    `json:"page_url"`
	ReferrerURL    string    `json:"referrer_url"`
	UTMSource      string    `json:"utm_source"`
	UTMMedium      string    `json:"utm_medium"`
	UTMCampaign    string    `json:"utm_campaign"`
	UTMTerm        string    `json:"utm_term"`
	UTMContent     string    `json:"utm_content"`
	CustomData     string    `json:"custom_data"` // JSON
	CreatedAt      time.Time `json:"created_at"`
}