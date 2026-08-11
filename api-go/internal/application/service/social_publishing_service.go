package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// SocialPublishService handles social media publishing
type SocialPublishService struct {
	db           *sql.DB
	facebook     *FacebookClient
	instagram    *InstagramClient
	tiktok       *TikTokClient
	twitter      *TwitterClient
	linkedin     *LinkedInClient
	scheduler    *SocialScheduler
}

// NewSocialPublishService creates a new social publishing service
func NewSocialPublishService(db *sql.DB, fbAppID, fbAppSecret string) *SocialPublishService {
	fb := NewFacebookClient(fbAppID, fbAppSecret)
	ig := NewInstagramClient(fbAppID, fbAppSecret)

	svc := &SocialPublishService{
		db:        db,
		facebook:  fb,
		instagram: ig,
		scheduler: NewSocialScheduler(),
	}

	return svc
}

// RegisterTikTokClient registers the TikTok client
func (s *SocialPublishService) RegisterTikTokClient(client *TikTokClient) {
	s.tiktok = client
}

// RegisterTwitterClient registers the Twitter/X client
func (s *SocialPublishService) RegisterTwitterClient(client *TwitterClient) {
	s.twitter = client
}

// RegisterLinkedInClient registers the LinkedIn client
func (s *SocialPublishService) RegisterLinkedInClient(client *LinkedInClient) {
	s.linkedin = client
}

// ConnectAccount handles OAuth callback and stores account
func (s *SocialPublishService) ConnectAccount(ctx context.Context, orgID, platform, code string) (*domain.SocialPublishAccount, error) {
	var client domain.SocialPublishClient
	switch platform {
	case "facebook":
		client = s.facebook
	case "instagram":
		client = s.instagram
	case "tiktok":
		client = s.tiktok
	case "twitter":
		client = s.twitter
	case "linkedin":
		client = s.linkedin
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	// Exchange code for tokens
	accessToken, refreshToken, expiresIn, err := client.ExchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get account info
	accountID, accountName, err := client.GetAccountInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	// Calculate expiry
	tokenExpiry := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Store account
	account := &domain.SocialPublishAccount{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Platform:       platform,
		PlatformID:     accountID,
		AccountName:    accountName,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		TokenExpiry:    tokenExpiry,
		Status:         "connected",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `
		INSERT INTO social_publish_accounts (id, organization_id, platform, platform_id, account_name, 
			access_token, refresh_token, token_expiry, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = s.db.Exec(query, account.ID, account.OrganizationID, account.Platform,
		account.PlatformID, account.AccountName, account.AccessToken, account.RefreshToken,
		account.TokenExpiry, account.Status, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	return account, nil
}

// GetAccounts returns all connected accounts for an org
func (s *SocialPublishService) GetAccounts(ctx context.Context, orgID string) ([]*domain.SocialPublishAccount, error) {
	query := `
		SELECT id, organization_id, platform, platform_id, account_name, 
			token_expiry, status, created_at, updated_at
		FROM social_publish_accounts
		WHERE organization_id = $1 AND status = 'connected'
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.SocialPublishAccount
	for rows.Next() {
		var a domain.SocialPublishAccount
		err := rows.Scan(&a.ID, &a.OrganizationID, &a.Platform, &a.PlatformID,
			&a.AccountName, &a.TokenExpiry, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

// DisconnectAccount disconnects a social account
func (s *SocialPublishService) DisconnectAccount(ctx context.Context, accountID string) error {
	query := `UPDATE social_publish_accounts SET status = 'disconnected', updated_at = NOW() WHERE id = $1`
	_, err := s.db.Exec(query, accountID)
	return err
}

// CreatePost creates a new social post
func (s *SocialPublishService) CreatePost(ctx context.Context, orgID, accountID, content string, mediaURLs []string, scheduledAt *time.Time) (*domain.SocialPublishPost, error) {
	// Get account
	var account domain.SocialPublishAccount
	err := s.db.QueryRow(`
		SELECT id, organization_id, platform, platform_id, account_name, access_token, 
			refresh_token, token_expiry, status
		FROM social_publish_accounts WHERE id = $1 AND organization_id = $2
	`, accountID, orgID).Scan(&account.ID, &account.OrganizationID, &account.Platform,
		&account.PlatformID, &account.AccountName, &account.AccessToken, &account.RefreshToken,
		&account.TokenExpiry, &account.Status)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	post := &domain.SocialPublishPost{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		AccountID:      accountID,
		Content:        content,
		MediaURLs:      mediaURLs,
		Platform:       account.Platform,
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if scheduledAt != nil {
		post.Status = "scheduled"
		post.ScheduledAt = scheduledAt
	}

	query := `
		INSERT INTO social_publish_posts (id, organization_id, account_id, content, media_urls, 
			platform, status, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = s.db.Exec(query, post.ID, post.OrganizationID, post.AccountID, post.Content,
		post.MediaURLs, post.Platform, post.Status, post.ScheduledAt, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// If scheduled, add to scheduler
	if scheduledAt != nil && !scheduledAt.IsZero() {
		s.scheduler.SchedulePost(post.ID, *scheduledAt, func() {
			s.PublishPost(ctx, post.ID)
		})
	}

	return post, nil
}

// PublishPost publishes a post immediately
func (s *SocialPublishService) PublishPost(ctx context.Context, postID string) error {
	var post domain.SocialPublishPost
	var account domain.SocialPublishAccount
	
	err := s.db.QueryRow(`
		SELECT sp.id, sp.organization_id, sp.account_id, sp.content, sp.media_urls, 
			sp.platform, sp.status
		FROM social_publish_posts sp WHERE sp.id = $1
	`, postID).Scan(&post.ID, &post.OrganizationID, &post.AccountID, &post.Content,
		&post.MediaURLs, &post.Platform, &post.Status)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	err = s.db.QueryRow(`
		SELECT access_token, refresh_token, token_expiry, platform_id
		FROM social_publish_accounts WHERE id = $1
	`, post.AccountID).Scan(&account.AccessToken, &account.RefreshToken, &account.TokenExpiry, &account.PlatformID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Check token expiry
	if time.Now().After(account.TokenExpiry) {
		// Try refresh token
		newToken, newExpiry, err := s.refreshToken(post.Platform, account.RefreshToken)
		if err != nil {
			return fmt.Errorf("token expired and refresh failed: %w", err)
		}
		account.AccessToken = newToken
		account.TokenExpiry = newExpiry
	}

	// Publish
	var platformClient domain.SocialPublishClient
	switch post.Platform {
	case "facebook":
		platformClient = s.facebook
	case "instagram":
		platformClient = s.instagram
	case "tiktok":
		platformClient = s.tiktok
	case "twitter":
		platformClient = s.twitter
	case "linkedin":
		platformClient = s.linkedin
	default:
		return fmt.Errorf("unsupported platform: %s", post.Platform)
	}

	var platformPostID string
	if post.ScheduledAt != nil && !post.ScheduledAt.IsZero() {
		platformPostID, err = platformClient.SchedulePost(account.AccessToken, post.Content, post.MediaURLs, post.ScheduledAt.Format(time.RFC3339))
	} else {
		platformPostID, err = platformClient.PublishPost(account.AccessToken, post.Content, post.MediaURLs)
	}
	
	if err != nil {
		// Update post status to failed
		s.db.Exec(`UPDATE social_publish_posts SET status = 'failed', error_msg = $1, updated_at = NOW() WHERE id = $2`, err.Error(), postID)
		return fmt.Errorf("publish failed: %w", err)
	}

	// Update post status
	now := time.Now()
	_, err = s.db.Exec(`
		UPDATE social_publish_posts 
		SET status = 'published', platform_post_id = $1, published_at = $2, updated_at = NOW()
		WHERE id = $3
	`, platformPostID, now, postID)
	if err != nil {
		log.Printf("Warning: failed to update post status: %v", err)
	}

	return nil
}

func (s *SocialPublishService) refreshToken(platform, refreshToken string) (string, time.Time, error) {
	var client domain.SocialPublishClient
	switch platform {
	case "facebook":
		client = s.facebook
	case "instagram":
		client = s.instagram
	case "tiktok":
		client = s.tiktok
	case "twitter":
		client = s.twitter
	case "linkedin":
		client = s.linkedin
	default:
		return "", time.Time{}, fmt.Errorf("unsupported platform")
	}

	newToken, expiresIn, err := client.RefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	return newToken, time.Now().Add(time.Duration(expiresIn) * time.Second), nil
}

// GetPostAnalytics retrieves analytics for a post
func (s *SocialPublishService) GetPostAnalytics(ctx context.Context, postID string) (*domain.SocialPostAnalytics, error) {
	var post domain.SocialPublishPost
	var account domain.SocialPublishAccount
	
	err := s.db.QueryRow(`
		SELECT platform, platform_post_id FROM social_publish_posts WHERE id = $1 AND status = 'published'
	`, postID).Scan(&post.Platform, &post.PlatformPostID)
	if err != nil {
		return nil, fmt.Errorf("published post not found")
	}

	err = s.db.QueryRow(`
		SELECT access_token FROM social_publish_accounts WHERE id = (
			SELECT account_id FROM social_publish_posts WHERE id = $1
		)
	`, postID).Scan(&account.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}

	var platformClient domain.SocialPublishClient
	switch post.Platform {
	case "facebook":
		platformClient = s.facebook
	case "instagram":
		platformClient = s.instagram
	case "tiktok":
		platformClient = s.tiktok
	case "twitter":
		platformClient = s.twitter
	case "linkedin":
		platformClient = s.linkedin
	default:
		return nil, fmt.Errorf("unsupported platform")
	}

	analytics, err := platformClient.GetAnalytics(account.AccessToken, post.PlatformPostID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analytics: %w", err)
	}

	// Save analytics
	postAnalytics := &domain.SocialPostAnalytics{
		PostID:         postID,
		Platform:       post.Platform,
		Impressions:    analytics.Impressions,
		Reach:          analytics.Reach,
		Likes:          analytics.Likes,
		Comments:       analytics.Comments,
		Shares:         analytics.Shares,
		Clicks:         analytics.Clicks,
		EngagementRate: analytics.EngagementRate,
		FetchedAt:      time.Now(),
	}

	query := `
		INSERT INTO social_publish_post_analytics (post_id, platform, impressions, reach, likes, 
			comments, shares, clicks, engagement_rate, fetched_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (post_id) DO UPDATE SET
			impressions = EXCLUDED.impressions,
			reach = EXCLUDED.reach,
			likes = EXCLUDED.likes,
			comments = EXCLUDED.comments,
			shares = EXCLUDED.shares,
			clicks = EXCLUDED.clicks,
			engagement_rate = EXCLUDED.engagement_rate,
			fetched_at = EXCLUDED.fetched_at
	`
	s.db.Exec(query, postAnalytics.PostID, postAnalytics.Platform, postAnalytics.Impressions,
		postAnalytics.Reach, postAnalytics.Likes, postAnalytics.Comments, postAnalytics.Shares,
		postAnalytics.Clicks, postAnalytics.EngagementRate, postAnalytics.FetchedAt)

	return postAnalytics, nil
}

// RefreshExpiredTokens refreshes tokens for all accounts with expiring tokens
func (s *SocialPublishService) RefreshExpiredTokens(ctx context.Context) error {
	query := `
		SELECT id, platform, refresh_token 
		FROM social_publish_accounts 
		WHERE status = 'connected' 
			AND token_expiry < NOW() + INTERVAL '24 hours'
			AND refresh_token IS NOT NULL
			AND refresh_token != ''
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var refreshed int
	for rows.Next() {
		var accountID, platform, refreshToken string
		if err := rows.Scan(&accountID, &platform, &refreshToken); err != nil {
			continue
		}

		newToken, newExpiry, err := s.refreshToken(platform, refreshToken)
		if err != nil {
			log.Printf("Failed to refresh token for account %s: %v", accountID, err)
			continue
		}

		s.db.Exec(`
			UPDATE social_publish_accounts 
			SET access_token = $1, token_expiry = $2, updated_at = NOW()
			WHERE id = $3
		`, newToken, newExpiry, accountID)
		refreshed++
	}

	log.Printf("Refreshed tokens for %d accounts", refreshed)
	return nil
}

// GetPosts returns posts for an org
func (s *SocialPublishService) GetPosts(ctx context.Context, orgID string, limit, offset int) ([]*domain.SocialPublishPost, int, error) {
	var total int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM social_publish_posts WHERE organization_id = $1`, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, organization_id, account_id, content, media_urls, hashtags,
			platform, status, scheduled_at, published_at, platform_post_id, error_msg,
			created_at, updated_at
		FROM social_publish_posts
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(query, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*domain.SocialPublishPost
	for rows.Next() {
		var p domain.SocialPublishPost
		err := rows.Scan(&p.ID, &p.OrganizationID, &p.AccountID, &p.Content, &p.MediaURLs,
			&p.Hashtags, &p.Platform, &p.Status, &p.ScheduledAt, &p.PublishedAt,
			&p.PlatformPostID, &p.ErrorMsg, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, &p)
	}

	return posts, total, nil
}

// GetScheduledPosts returns scheduled posts
func (s *SocialPublishService) GetScheduledPosts(ctx context.Context, orgID string) ([]*domain.SocialPublishPost, error) {
	query := `
		SELECT id, organization_id, account_id, content, media_urls, hashtags,
			platform, status, scheduled_at, published_at, platform_post_id, error_msg,
			created_at, updated_at
		FROM social_publish_posts
		WHERE organization_id = $1 AND status = 'scheduled'
		ORDER BY scheduled_at ASC
	`
	rows, err := s.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.SocialPublishPost
	for rows.Next() {
		var p domain.SocialPublishPost
		err := rows.Scan(&p.ID, &p.OrganizationID, &p.AccountID, &p.Content, &p.MediaURLs,
			&p.Hashtags, &p.Platform, &p.Status, &p.ScheduledAt, &p.PublishedAt,
			&p.PlatformPostID, &p.ErrorMsg, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

// DeletePost deletes a post
func (s *SocialPublishService) DeletePost(ctx context.Context, postID string) error {
	// Get post to check if it's published
	var platform, platformPostID, accountID string
	err := s.db.QueryRow(`
		SELECT platform, platform_post_id, account_id FROM social_publish_posts WHERE id = $1
	`, postID).Scan(&platform, &platformPostID, &accountID)
	if err != nil {
		return fmt.Errorf("post not found")
	}

	// If published, try to delete from platform
	if platformPostID != "" {
		var accessToken string
		s.db.QueryRow(`SELECT access_token FROM social_publish_accounts WHERE id = $1`, accountID).Scan(&accessToken)
		
		if accessToken != "" {
			var client domain.SocialPublishClient
			switch platform {
			case "facebook":
				client = s.facebook
			case "instagram":
				client = s.instagram
			case "tiktok":
				client = s.tiktok
			case "twitter":
				client = s.twitter
			case "linkedin":
				client = s.linkedin
			}
			if client != nil {
				client.DeletePost(accessToken, platformPostID)
			}
		}
	}

	// Delete from DB
	_, err = s.db.Exec(`DELETE FROM social_publish_posts WHERE id = $1`, postID)
	return err
}

// GetFacebookOAuthURL returns Facebook OAuth URL
func (s *SocialPublishService) GetFacebookOAuthURL(state string) string {
	return s.facebook.GetOAuthURL(state)
}

// GetInstagramOAuthURL returns Instagram OAuth URL (same as Facebook)
func (s *SocialPublishService) GetInstagramOAuthURL(state string) string {
	return s.instagram.GetOAuthURL(state)
}
