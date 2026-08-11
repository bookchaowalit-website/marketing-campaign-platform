package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// =====================================================
// EMAIL MARKETING SERVICE
// =====================================================

// EmailMarketingService handles email marketing operations
type EmailMarketingService struct {
	db *sql.DB
}

// NewEmailMarketingService creates a new email marketing service
func NewEmailMarketingService(db *sql.DB) *EmailMarketingService {
	return &EmailMarketingService{db: db}
}

// CreateList creates a new email list
func (s *EmailMarketingService) CreateList(ctx context.Context, orgID, name, description string) (*domain.EmailList, error) {
	list := &domain.EmailList{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		Description:    description,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO email_lists (id, organization_id, name, description, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(ctx, query, list.ID, list.OrganizationID, list.Name, list.Description, list.Status, list.CreatedAt, list.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetLists retrieves all email lists for an organization
func (s *EmailMarketingService) GetLists(ctx context.Context, orgID string) ([]*domain.EmailList, error) {
	query := `SELECT id, organization_id, name, description, status, subscriber_count, created_at, updated_at 
	          FROM email_lists WHERE organization_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := make([]*domain.EmailList, 0)
	for rows.Next() {
		list := &domain.EmailList{}
		err := rows.Scan(&list.ID, &list.OrganizationID, &list.Name, &list.Description, &list.Status, &list.SubscriberCount, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	return lists, nil
}

// AddSubscriber adds a subscriber to an organization
func (s *EmailMarketingService) AddSubscriber(ctx context.Context, orgID, email, firstName, lastName, source string, listIDs []string) (*domain.EmailSubscriber, error) {
	subscriber := &domain.EmailSubscriber{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		Source:         source,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert subscriber
	query := `INSERT INTO email_subscribers (id, organization_id, email, first_name, last_name, source, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (organization_id, email) DO UPDATE SET first_name = $4, last_name = $5, updated_at = $8 RETURNING id`

	err = tx.QueryRowContext(ctx, query, subscriber.ID, subscriber.OrganizationID, subscriber.Email, subscriber.FirstName, subscriber.LastName, subscriber.Source, subscriber.Status, subscriber.CreatedAt, subscriber.UpdatedAt).Scan(&subscriber.ID)
	if err != nil {
		return nil, err
	}

	// Add to lists
	for _, listID := range listIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO email_list_subscribers (list_id, subscriber_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, listID, subscriber.ID)
		if err != nil {
			return nil, err
		}
		_, err = tx.ExecContext(ctx, `UPDATE email_lists SET subscriber_count = subscriber_count + 1 WHERE id = $1`, listID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return subscriber, nil
}

// GetSubscribers retrieves subscribers for an organization
func (s *EmailMarketingService) GetSubscribers(ctx context.Context, orgID string, limit, offset int) ([]*domain.EmailSubscriber, int, error) {
	countQuery := `SELECT COUNT(*) FROM email_subscribers WHERE organization_id = $1`
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, organization_id, email, first_name, last_name, status, source, created_at 
	          FROM email_subscribers WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	subscribers := make([]*domain.EmailSubscriber, 0)
	for rows.Next() {
		sub := &domain.EmailSubscriber{}
		err := rows.Scan(&sub.ID, &sub.OrganizationID, &sub.Email, &sub.FirstName, &sub.LastName, &sub.Status, &sub.Source, &sub.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		subscribers = append(subscribers, sub)
	}

	return subscribers, total, nil
}

// CreateCampaign creates a new email campaign
func (s *EmailMarketingService) CreateCampaign(ctx context.Context, orgID, name, subject, fromName, fromEmail, htmlContent string, listIDs []string) (*domain.EmailCampaign, error) {
	campaign := &domain.EmailCampaign{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		Subject:        subject,
		FromName:       fromName,
		FromEmail:      fromEmail,
		HTMLContent:    htmlContent,
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO email_campaigns (id, organization_id, name, subject, from_name, from_email, html_content, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.ExecContext(ctx, query, campaign.ID, campaign.OrganizationID, campaign.Name, campaign.Subject, campaign.FromName, campaign.FromEmail, campaign.HTMLContent, campaign.Status, campaign.CreatedAt, campaign.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Link lists
	for _, listID := range listIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO email_campaign_lists (campaign_id, list_id) VALUES ($1, $2)`, campaign.ID, listID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return campaign, nil
}

// GetCampaigns retrieves campaigns for an organization
func (s *EmailMarketingService) GetCampaigns(ctx context.Context, orgID string, limit, offset int) ([]*domain.EmailCampaign, int, error) {
	countQuery := `SELECT COUNT(*) FROM email_campaigns WHERE organization_id = $1`
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, organization_id, name, subject, status, recipient_count, sent_count, opened_count, clicked_count, open_rate, click_rate, created_at 
	          FROM email_campaigns WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	campaigns := make([]*domain.EmailCampaign, 0)
	for rows.Next() {
		camp := &domain.EmailCampaign{}
		err := rows.Scan(&camp.ID, &camp.OrganizationID, &camp.Name, &camp.Subject, &camp.Status, &camp.RecipientCount, &camp.SentCount, &camp.OpenedCount, &camp.ClickedCount, &camp.OpenRate, &camp.ClickRate, &camp.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		campaigns = append(campaigns, camp)
	}

	return campaigns, total, nil
}

// =====================================================
// SOCIAL MEDIA SERVICE
// =====================================================

// SocialMediaService handles social media operations
type SocialMediaService struct {
	db *sql.DB
}

// NewSocialMediaService creates a new social media service
func NewSocialMediaService(db *sql.DB) *SocialMediaService {
	return &SocialMediaService{db: db}
}

// CreatePost creates a new social media post
func (s *SocialMediaService) CreatePost(ctx context.Context, orgID, content string, mediaURLs []string, hashtags []string) (*domain.SocialPost, error) {
	mediaJSON, _ := json.Marshal(mediaURLs)
	hashtagsJSON, _ := json.Marshal(hashtags)

	post := &domain.SocialPost{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Content:        content,
		MediaURLs:      mediaURLs,
		Hashtags:       hashtags,
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO social_posts (id, organization_id, content, media_urls, hashtags, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query, post.ID, post.OrganizationID, post.Content, mediaJSON, hashtagsJSON, post.Status, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// SchedulePost schedules a post for publishing
func (s *SocialMediaService) SchedulePost(ctx context.Context, postID, accountID string, scheduledTime time.Time) error {
	schedule := &domain.SocialPostSchedule{
		ID:            uuid.New().String(),
		PostID:        postID,
		AccountID:     accountID,
		ScheduledTime: scheduledTime,
		Status:        "scheduled",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO social_post_schedules (id, post_id, account_id, scheduled_time, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, query, schedule.ID, schedule.PostID, schedule.AccountID, schedule.ScheduledTime, schedule.Status, schedule.CreatedAt, schedule.UpdatedAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE social_posts SET status = 'scheduled', updated_at = $1 WHERE id = $2`, time.Now(), postID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetPosts retrieves posts for an organization
func (s *SocialMediaService) GetPosts(ctx context.Context, orgID string, limit, offset int) ([]map[string]interface{}, error) {
	query := `SELECT p.id, p.content, p.media_urls, p.hashtags, p.status, p.created_at,
	          COALESCE(json_agg(json_build_object('account_id', ps.account_id, 'scheduled_time', ps.scheduled_time, 'status', ps.status)) FILTER (WHERE ps.id IS NOT NULL), '[]') as schedules
	          FROM social_posts p
	          LEFT JOIN social_post_schedules ps ON p.id = ps.post_id
	          WHERE p.organization_id = $1
	          GROUP BY p.id
	          ORDER BY p.created_at DESC
	          LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, content, status string
		var mediaURLs, hashtags, schedules []byte
		var createdAt time.Time
		err := rows.Scan(&id, &content, &mediaURLs, &hashtags, &status, &createdAt, &schedules)
		if err != nil {
			return nil, err
		}

		var mediaURLsList, hashtagsList []string
		var schedulesList []map[string]interface{}
		json.Unmarshal(mediaURLs, &mediaURLsList)
		json.Unmarshal(hashtags, &hashtagsList)
		json.Unmarshal(schedules, &schedulesList)

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"content":    content,
			"media_urls": mediaURLsList,
			"hashtags":   hashtagsList,
			"status":     status,
			"schedules":  schedulesList,
			"created_at": createdAt,
		})
	}

	return posts, nil
}

// =====================================================
// LANDING PAGE SERVICE
// =====================================================

// LandingPageService handles landing page operations
type LandingPageService struct {
	db *sql.DB
}

// NewLandingPageService creates a new landing page service
func NewLandingPageService(db *sql.DB) *LandingPageService {
	return &LandingPageService{db: db}
}

// CreatePage creates a new landing page
func (s *LandingPageService) CreatePage(ctx context.Context, orgID, name, slug, title, htmlContent string) (*domain.LandingPage, error) {
	page := &domain.LandingPage{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		Slug:           slug,
		Title:          title,
		HTMLContent:    htmlContent,
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO landing_pages (id, organization_id, name, slug, title, html_content, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := s.db.ExecContext(ctx, query, page.ID, page.OrganizationID, page.Name, page.Slug, page.Title, page.HTMLContent, page.Status, page.CreatedAt, page.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// GetPages retrieves landing pages for an organization
func (s *LandingPageService) GetPages(ctx context.Context, orgID string, limit, offset int) ([]*domain.LandingPage, error) {
	query := `SELECT id, organization_id, name, slug, title, status, view_count, conversion_count, created_at 
	          FROM landing_pages WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pages := make([]*domain.LandingPage, 0)
	for rows.Next() {
		page := &domain.LandingPage{}
		err := rows.Scan(&page.ID, &page.OrganizationID, &page.Name, &page.Slug, &page.Title, &page.Status, &page.ViewCount, &page.ConversionCount, &page.CreatedAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

// CreateLead creates a new lead from a landing page
func (s *LandingPageService) CreateLead(ctx context.Context, orgID, pageID, email, firstName, lastName, phone string, customData map[string]interface{}) (*domain.LandingPageLead, error) {
	customDataJSON, _ := json.Marshal(customData)

	lead := &domain.LandingPageLead{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		PageID:         pageID,
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		Phone:          phone,
		CustomData:     string(customDataJSON),
		Status:         "new",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO landing_page_leads (id, organization_id, page_id, email, first_name, last_name, phone, custom_data, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := s.db.ExecContext(ctx, query, lead.ID, lead.OrganizationID, lead.PageID, lead.Email, lead.FirstName, lead.LastName, lead.Phone, lead.CustomData, lead.Status, lead.CreatedAt, lead.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Update page conversion count
	s.db.ExecContext(ctx, `UPDATE landing_pages SET conversion_count = conversion_count + 1 WHERE id = $1`, pageID)

	return lead, nil
}

// =====================================================
// MARKETING AUTOMATION SERVICE
// =====================================================

// AutomationService handles marketing automation operations
type AutomationService struct {
	db *sql.DB
}

// NewAutomationService creates a new automation service
func NewAutomationService(db *sql.DB) *AutomationService {
	return &AutomationService{db: db}
}

// CreateWorkflow creates a new automation workflow
func (s *AutomationService) CreateWorkflow(ctx context.Context, orgID, name, triggerType string, triggerConfig map[string]interface{}) (*domain.AutomationWorkflow, error) {
	triggerConfigJSON, _ := json.Marshal(triggerConfig)

	workflow := &domain.AutomationWorkflow{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		TriggerType:    triggerType,
		TriggerConfig:  string(triggerConfigJSON),
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO automation_workflows (id, organization_id, name, trigger_type, trigger_config, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query, workflow.ID, workflow.OrganizationID, workflow.Name, workflow.TriggerType, workflow.TriggerConfig, workflow.Status, workflow.CreatedAt, workflow.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetWorkflows retrieves workflows for an organization
func (s *AutomationService) GetWorkflows(ctx context.Context, orgID string) ([]*domain.AutomationWorkflow, error) {
	query := `SELECT id, organization_id, name, trigger_type, status, total_contacts, active_contacts, created_at 
	          FROM automation_workflows WHERE organization_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := make([]*domain.AutomationWorkflow, 0)
	for rows.Next() {
		wf := &domain.AutomationWorkflow{}
		err := rows.Scan(&wf.ID, &wf.OrganizationID, &wf.Name, &wf.TriggerType, &wf.Status, &wf.TotalContacts, &wf.ActiveContacts, &wf.CreatedAt)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

// =====================================================
// INFLUENCER MARKETING SERVICE
// =====================================================

// InfluencerService handles influencer marketing operations
type InfluencerService struct {
	db *sql.DB
}

// NewInfluencerService creates a new influencer service
func NewInfluencerService(db *sql.DB) *InfluencerService {
	return &InfluencerService{db: db}
}

// CreateInfluencer creates a new influencer
func (s *InfluencerService) CreateInfluencer(ctx context.Context, orgID, name, email, category string, niches []string) (*domain.Influencer, error) {
	nichesJSON, _ := json.Marshal(niches)

	influencer := &domain.Influencer{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		Email:          email,
		Category:       category,
		Niches:         niches,
		Status:         "prospect",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO influencers (id, organization_id, name, email, category, niches, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := s.db.ExecContext(ctx, query, influencer.ID, influencer.OrganizationID, influencer.Name, influencer.Email, influencer.Category, nichesJSON, influencer.Status, influencer.CreatedAt, influencer.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return influencer, nil
}

// GetInfluencers retrieves influencers for an organization
func (s *InfluencerService) GetInfluencers(ctx context.Context, orgID string, limit, offset int) ([]*domain.Influencer, error) {
	query := `SELECT id, organization_id, name, email, category, niches, status, rating, total_collaborations, created_at 
	          FROM influencers WHERE organization_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	influencers := make([]*domain.Influencer, 0)
	for rows.Next() {
		inf := &domain.Influencer{}
		var nichesJSON []byte
		err := rows.Scan(&inf.ID, &inf.OrganizationID, &inf.Name, &inf.Email, &inf.Category, &nichesJSON, &inf.Status, &inf.Rating, &inf.TotalCollaborations, &inf.CreatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(nichesJSON, &inf.Niches)
		influencers = append(influencers, inf)
	}

	return influencers, nil
}

// CreateCampaign creates a new influencer campaign
func (s *InfluencerService) CreateCampaign(ctx context.Context, orgID, name, campaignType string, budget float64) (*domain.InfluencerCampaign, error) {
	campaign := &domain.InfluencerCampaign{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		CampaignType:   campaignType,
		Budget:         budget,
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO influencer_campaigns (id, organization_id, name, campaign_type, budget, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query, campaign.ID, campaign.OrganizationID, campaign.Name, campaign.CampaignType, campaign.Budget, campaign.Status, campaign.CreatedAt, campaign.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

// =====================================================
// REPORTING SERVICE
// =====================================================

// ReportingService handles reporting operations
type ReportingService struct {
	db *sql.DB
}

// NewReportingService creates a new reporting service
func NewReportingService(db *sql.DB) *ReportingService {
	return &ReportingService{db: db}
}

// GetMarketingOverview retrieves marketing overview data
func (s *ReportingService) GetMarketingOverview(ctx context.Context, orgID string) (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// Email stats
	var emailStats struct {
		TotalCampaigns   int
		TotalSent        int
		TotalOpens       int
		TotalClicks      int
		AvgOpenRate      float64
		AvgClickRate     float64
		TotalSubscribers int
	}

	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(sent_count), 0), COALESCE(SUM(opened_count), 0), COALESCE(SUM(clicked_count), 0),
		       COALESCE(AVG(open_rate), 0), COALESCE(AVG(click_rate), 0)
		FROM email_campaigns WHERE organization_id = $1`, orgID).Scan(
		&emailStats.TotalCampaigns, &emailStats.TotalSent, &emailStats.TotalOpens, &emailStats.TotalClicks,
		&emailStats.AvgOpenRate, &emailStats.AvgClickRate)

	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM email_subscribers WHERE organization_id = $1`, orgID).Scan(&emailStats.TotalSubscribers)
	
	if err == nil {
		overview["email"] = emailStats
	} else {
		// Return empty stats instead of nil
		overview["email"] = map[string]interface{}{
			"total_campaigns":    0,
			"total_sent":         0,
			"total_opens":        0,
			"total_clicks":       0,
			"avg_open_rate":      0.0,
			"avg_click_rate":     0.0,
			"total_subscribers":  0,
		}
	}

	// Social stats
	var socialStats struct {
		TotalPosts     int
		ScheduledPosts int
		PublishedPosts int
		TotalImpressions int64
		TotalEngagement  int64
	}

	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM social_posts WHERE organization_id = $1`, orgID).Scan(&socialStats.TotalPosts)
	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM social_posts WHERE organization_id = $1 AND status = 'scheduled'`, orgID).Scan(&socialStats.ScheduledPosts)
	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM social_posts WHERE organization_id = $1 AND status = 'published'`, orgID).Scan(&socialStats.PublishedPosts)
	s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(impressions), 0), COALESCE(SUM(engagement_count), 0) FROM social_analytics sa
		JOIN social_accounts sac ON sa.account_id = sac.id WHERE sac.organization_id = $1`, orgID).Scan(&socialStats.TotalImpressions, &socialStats.TotalEngagement)

	overview["social"] = socialStats

	// Landing page stats
	var landingStats struct {
		TotalPages       int
		ActivePages      int
		TotalViews       int64
		TotalLeads       int
		AvgConversionRate float64
	}

	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM landing_pages WHERE organization_id = $1`, orgID).Scan(&landingStats.TotalPages)
	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM landing_pages WHERE organization_id = $1 AND status = 'published'`, orgID).Scan(&landingStats.ActivePages)
	s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(view_count), 0) FROM landing_pages WHERE organization_id = $1`, orgID).Scan(&landingStats.TotalViews)
	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM landing_page_leads WHERE organization_id = $1`, orgID).Scan(&landingStats.TotalLeads)

	overview["landing_pages"] = landingStats

	// Influencer stats
	var influencerStats struct {
		TotalInfluencers     int
		ActiveCampaigns      int
		TotalCollaborations  int
		TotalReach           int64
		TotalSpend           float64
	}

	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM influencers WHERE organization_id = $1`, orgID).Scan(&influencerStats.TotalInfluencers)
	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM influencer_campaigns WHERE organization_id = $1 AND status = 'active'`, orgID).Scan(&influencerStats.ActiveCampaigns)
	s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM influencer_collaborations ic
		JOIN influencer_campaigns c ON ic.campaign_id = c.id WHERE c.organization_id = $1`, orgID).Scan(&influencerStats.TotalCollaborations)
	s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total_reach), 0) FROM influencer_campaigns WHERE organization_id = $1`, orgID).Scan(&influencerStats.TotalReach)
	s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(spent), 0) FROM influencer_campaigns WHERE organization_id = $1`, orgID).Scan(&influencerStats.TotalSpend)

	overview["influencers"] = influencerStats

	return overview, nil
}

// TrackEvent tracks an analytics event
func (s *ReportingService) TrackEvent(ctx context.Context, orgID, eventType, eventCategory, eventAction string, eventData map[string]interface{}) error {
	customDataJSON, _ := json.Marshal(eventData)

	query := `INSERT INTO analytics_events_unified (id, organization_id, event_type, event_category, event_action, custom_data, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.ExecContext(ctx, query, uuid.New().String(), orgID, eventType, eventCategory, eventAction, customDataJSON, time.Now())
	return err
}

// Helper functions
func generateRandomSlug() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Errors
var (
	ErrListNotFound       = errors.New("email list not found")
	ErrSubscriberNotFound = errors.New("subscriber not found")
	ErrCampaignNotFound   = errors.New("campaign not found")
	ErrPostNotFound       = errors.New("post not found")
	ErrPageNotFound       = errors.New("landing page not found")
	ErrWorkflowNotFound   = errors.New("workflow not found")
	ErrInfluencerNotFound = errors.New("influencer not found")
)
