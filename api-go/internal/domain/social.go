package domain

import "time"

// SocialPublishAccount represents a connected social media account
type SocialPublishAccount struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Platform       string    `json:"platform"`
	PlatformID     string    `json:"platform_id"`
	AccountName    string    `json:"account_name"`
	PageID         string    `json:"page_id,omitempty"`
	AccessToken    string    `json:"-"`
	RefreshToken   string    `json:"-"`
	TokenExpiry    time.Time `json:"-"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SocialPublishPost represents a post to be published
type SocialPublishPost struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	AccountID      string     `json:"account_id"`
	Content        string     `json:"content"`
	MediaURLs      []string   `json:"media_urls,omitempty"`
	Hashtags       []string   `json:"hashtags,omitempty"`
	Platform       string     `json:"platform"`
	Status         string     `json:"status"`
	ScheduledAt    *time.Time `json:"scheduled_at,omitempty"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`
	PlatformPostID string     `json:"platform_post_id,omitempty"`
	ErrorMsg       string     `json:"error_msg,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// SocialPostAnalytics represents analytics for a published post
type SocialPostAnalytics struct {
	PostID         string    `json:"post_id"`
	Platform       string    `json:"platform"`
	Impressions    int       `json:"impressions"`
	Reach          int       `json:"reach"`
	Likes          int       `json:"likes"`
	Comments       int       `json:"comments"`
	Shares         int       `json:"shares"`
	Saves          int       `json:"saves,omitempty"`
	Clicks         int       `json:"clicks"`
	EngagementRate float64   `json:"engagement_rate"`
	FetchedAt      time.Time `json:"fetched_at"`
}

// SocialPublishClient defines the interface for platform-specific publishing
type SocialPublishClient interface {
	GetOAuthURL(state string) string
	ExchangeCode(code string) (accessToken, refreshToken string, expiresIn int, err error)
	GetAccountInfo(accessToken string) (accountID, accountName string, err error)
	PublishPost(accessToken, content string, mediaURLs []string) (postID string, err error)
	SchedulePost(accessToken, content string, mediaURLs []string, scheduledAt string) (postID string, err error)
	DeletePost(accessToken, postID string) error
	GetAnalytics(accessToken, postID string) (*SocialPostAnalytics, error)
	RefreshToken(refreshToken string) (newAccessToken string, newExpiresIn int, err error)
}
