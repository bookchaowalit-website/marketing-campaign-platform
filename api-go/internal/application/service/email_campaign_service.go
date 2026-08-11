package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// EmailCampaignService handles email campaign sending, tracking, and A/B testing
type EmailCampaignService struct {
	db           *sql.DB
	emailService domain.EmailService
}

// NewEmailCampaignService creates a new email campaign service
func NewEmailCampaignService(db *sql.DB, emailService domain.EmailService) *EmailCampaignService {
	return &EmailCampaignService{
		db:           db,
		emailService: emailService,
	}
}

// CampaignAnalytics represents analytics data for a campaign
type CampaignAnalytics struct {
	CampaignID       string  `json:"campaign_id"`
	CampaignName     string  `json:"campaign_name"`
	Status           string  `json:"status"`
	TotalSent        int     `json:"total_sent"`
	Delivered        int     `json:"delivered"`
	Opened           int     `json:"opened"`
	Clicked          int     `json:"clicked"`
	Bounced          int     `json:"bounced"`
	Unsubscribed     int     `json:"unsubscribed"`
	Complained       int     `json:"complained"`
	Pending          int     `json:"pending"`
	OpenRate         float64 `json:"open_rate"`
	ClickRate        float64 `json:"click_rate"`
	BounceRate       float64 `json:"bounce_rate"`
	UnsubscribeRate  float64 `json:"unsubscribe_rate"`
}

// ABTest represents an email A/B test
type ABTest struct {
	ID               string     `json:"id"`
	CampaignIDA      string     `json:"campaign_id_a"`
	CampaignIDB      string     `json:"campaign_id_b"`
	CampaignNameA    string     `json:"campaign_name_a"`
	CampaignNameB    string     `json:"campaign_name_b"`
	Status           string     `json:"status"`
	Winner           *string    `json:"winner"`
	TestMetric       string     `json:"test_metric"`
	SampleSizePercent int       `json:"sample_size_percent"`
	StartedAt        time.Time  `json:"started_at"`
	CompletedAt      *time.Time `json:"completed_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

// CreateABTestRequest represents the request to create an A/B test
type CreateABTestRequest struct {
	CampaignIDA       string `json:"campaign_id_a" binding:"required"`
	CampaignIDB       string `json:"campaign_id_b" binding:"required"`
	TestMetric        string `json:"test_metric"`
	SampleSizePercent int    `json:"sample_size_percent"`
}

// ===========================
// SEND EMAIL CAMPAIGN
// ===========================

// SendEmailCampaign sends an email campaign to all subscribers in linked lists
func (s *EmailCampaignService) SendEmailCampaign(ctx context.Context, campaignID string) error {
	// Get campaign
	var campaign domain.EmailCampaign
	err := s.db.QueryRowContext(ctx, `
		SELECT id, organization_id, name, subject, preview_text, from_name, from_email, reply_to,
			template_id, html_content, text_content, campaign_type, status, recipient_count,
			sent_count, delivered_count, opened_count, clicked_count, bounced_count,
			unsubscribed_count, complained_count, open_rate, click_rate, bounce_rate, unsubscribe_rate,
			schedule_time, sent_at, tracking_enabled, google_analytics, utm_source, utm_medium, utm_campaign
		FROM email_campaigns WHERE id = $1
	`, campaignID).Scan(
		&campaign.ID, &campaign.OrganizationID, &campaign.Name, &campaign.Subject,
		&campaign.PreviewText, &campaign.FromName, &campaign.FromEmail, &campaign.ReplyTo,
		&campaign.TemplateID, &campaign.HTMLContent, &campaign.TextContent, &campaign.CampaignType,
		&campaign.Status, &campaign.RecipientCount, &campaign.SentCount, &campaign.DeliveredCount,
		&campaign.OpenedCount, &campaign.ClickedCount, &campaign.BouncedCount,
		&campaign.UnsubscribedCount, &campaign.ComplainedCount, &campaign.OpenRate,
		&campaign.ClickRate, &campaign.BounceRate, &campaign.UnsubscribeRate,
		&campaign.ScheduleTime, &campaign.SentAt, &campaign.TrackingEnabled,
		&campaign.GoogleAnalytics, &campaign.UTMSource, &campaign.UTMMedium, &campaign.UTMCampaign,
	)
	if err != nil {
		return fmt.Errorf("campaign not found: %w", err)
	}

	if campaign.Status != "draft" && campaign.Status != "scheduled" {
		return errors.New("campaign must be in draft or scheduled status to send")
	}

	// Get linked lists for this campaign
	rows, err := s.db.QueryContext(ctx, `
		SELECT ecl.list_id, el.name
		FROM email_campaign_lists ecl
		JOIN email_lists el ON ecl.list_id = el.id
		WHERE ecl.campaign_id = $1
	`, campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign lists: %w", err)
	}
	defer rows.Close()

	var listIDs []string
	for rows.Next() {
		var listID, listName string
		if err := rows.Scan(&listID, &listName); err != nil {
			return fmt.Errorf("failed to scan list: %w", err)
		}
		listIDs = append(listIDs, listID)
	}

	if len(listIDs) == 0 {
		return errors.New("campaign has no linked lists")
	}

	// Update campaign status to sending
	_, err = s.db.ExecContext(ctx, `
		UPDATE email_campaigns SET status = 'sending', updated_at = NOW() WHERE id = $1
	`, campaignID)
	if err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	// Get all active subscribers from the linked lists (deduplicated)
	subscriberRows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT es.id, es.email, es.first_name, es.last_name
		FROM email_subscribers es
		JOIN email_list_subscribers els ON es.id = els.subscriber_id
		WHERE els.list_id = ANY($1) AND es.status = 'active'
	`, listIDs)
	if err != nil {
		s.db.ExecContext(ctx, `UPDATE email_campaigns SET status = 'failed', updated_at = NOW() WHERE id = $1`, campaignID)
		return fmt.Errorf("failed to get subscribers: %w", err)
	}
	defer subscriberRows.Close()

	type subscriber struct {
		id        string
		email     string
		firstName string
		lastName  string
	}

	var subscribers []subscriber
	for subscriberRows.Next() {
		var sub subscriber
		if err := subscriberRows.Scan(&sub.id, &sub.email, &sub.firstName, &sub.lastName); err != nil {
			continue
		}
		subscribers = append(subscribers, sub)
	}

	if len(subscribers) == 0 {
		s.db.ExecContext(ctx, `UPDATE email_campaigns SET status = 'failed', updated_at = NOW() WHERE id = $1`, campaignID)
		return errors.New("no active subscribers in linked lists")
	}

	// Update recipient count
	_, err = s.db.ExecContext(ctx, `
		UPDATE email_campaigns SET recipient_count = $1, updated_at = NOW() WHERE id = $2
	`, len(subscribers), campaignID)
	if err != nil {
		return fmt.Errorf("failed to update recipient count: %w", err)
	}

	// Send emails and create email_sends records
	sentCount := 0
	failedCount := 0
	now := time.Now()

	for _, sub := range subscribers {
		// Create email_sends record
		sendID := uuid.New().String()
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO email_sends (id, campaign_id, subscriber_id, status, created_at)
			VALUES ($1, $2, $3, 'pending', $4)
		`, sendID, campaignID, sub.id, now)
		if err != nil {
			// Skip if already exists (unique constraint)
			continue
		}

		// Send email via email service
		if s.emailService != nil {
			err = s.emailService.SendEmail(
				sub.email,
				sub.firstName+" "+sub.lastName,
				campaign.Subject,
				campaign.HTMLContent,
				campaign.TextContent,
			)
		}

		if err != nil {
			// Update send record as failed
			s.db.ExecContext(ctx, `
				UPDATE email_sends SET status = 'bounced', send_error = $1, bounced_at = NOW()
				WHERE id = $2
			`, err.Error(), sendID)
			failedCount++
			continue
		}

		// Update send record as sent
		s.db.ExecContext(ctx, `
			UPDATE email_sends SET status = 'sent', sent_at = NOW() WHERE id = $1
		`, sendID)
		sentCount++

		// Create sent event
		s.db.ExecContext(ctx, `
			INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, created_at)
			VALUES ($1, $2, $3, 'sent', $4, NOW())
		`, campaign.OrganizationID, campaignID, sub.id, sub.email)
	}

	// Update campaign status and stats
	status := "sent"
	if failedCount > 0 && sentCount == 0 {
		status = "failed"
	}

	_, err = s.db.ExecContext(ctx, `
		UPDATE email_campaigns
		SET status = $1, sent_count = $2, sent_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`, status, sentCount, campaignID)
	if err != nil {
		return fmt.Errorf("failed to update campaign stats: %w", err)
	}

	return nil
}

// ===========================
// TRACK EMAIL EVENT
// ===========================

// TrackEmailEvent tracks an email event (open, click, bounce, unsubscribe)
func (s *EmailCampaignService) TrackEmailEvent(ctx context.Context, sendID string, eventType string, metadata map[string]string) error {
	// Get the send record
	var campaignID, subscriberID, currentStatus string
	err := s.db.QueryRowContext(ctx, `
		SELECT campaign_id, subscriber_id, status FROM email_sends WHERE id = $1
	`, sendID).Scan(&campaignID, &subscriberID, &currentStatus)
	if err != nil {
		return fmt.Errorf("send record not found: %w", err)
	}

	// Get campaign and subscriber details
	var orgID, email string
	err = s.db.QueryRowContext(ctx, `
		SELECT ec.organization_id, es.email
		FROM email_campaigns ec
		JOIN email_subscribers es ON es.id = $1
		WHERE ec.id = $2
	`, subscriberID, campaignID).Scan(&orgID, &email)
	if err != nil {
		return fmt.Errorf("failed to get campaign/subscriber details: %w", err)
	}

	var updateQuery string
	var eventQuery string

	switch eventType {
	case "delivered":
		updateQuery = `UPDATE email_sends SET status = 'delivered', delivered_at = NOW() WHERE id = $1`
		eventQuery = `INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, created_at) VALUES ($1, $2, $3, 'delivered', $4, NOW())`

	case "opened":
		updateQuery = `UPDATE email_sends SET status = 'opened', opened_at = NOW() WHERE id = $1`
		eventQuery = `INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, ip_address, user_agent, created_at) VALUES ($1, $2, $3, 'opened', $4, $5, $6, NOW())`
		ip := metadata["ip_address"]
		ua := metadata["user_agent"]
		s.db.ExecContext(ctx, eventQuery, orgID, campaignID, subscriberID, email, ip, ua)

	case "clicked":
		updateQuery = `UPDATE email_sends SET status = 'clicked', clicked_at = NOW() WHERE id = $1`
		eventQuery = `INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, url, created_at) VALUES ($1, $2, $3, 'clicked', $4, $5, NOW())`
		url := metadata["url"]
		s.db.ExecContext(ctx, eventQuery, orgID, campaignID, subscriberID, email, url)

	case "bounced":
		updateQuery = `UPDATE email_sends SET status = 'bounced', bounced_at = NOW(), bounce_type = $2, bounce_reason = $3 WHERE id = $1`
		eventQuery = `INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, bounce_type, bounce_reason, created_at) VALUES ($1, $2, $3, 'bounced', $4, $5, $6, NOW())`
		bounceType := metadata["bounce_type"]
		bounceReason := metadata["bounce_reason"]
		s.db.ExecContext(ctx, eventQuery, orgID, campaignID, subscriberID, email, bounceType, bounceReason)

	case "unsubscribed":
		updateQuery = `UPDATE email_sends SET status = 'unsubscribed', unsubscribed_at = NOW() WHERE id = $1`
		eventQuery = `INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, created_at) VALUES ($1, $2, $3, 'unsubscribed', $4, NOW())`
		// Also update subscriber status
		s.db.ExecContext(ctx, `UPDATE email_subscribers SET status = 'unsubscribed', unsubscribed_at = NOW() WHERE id = $1`, subscriberID)

	case "complained":
		updateQuery = `UPDATE email_sends SET status = 'complained' WHERE id = $1`
		eventQuery = `INSERT INTO email_events (organization_id, campaign_id, subscriber_id, event_type, email, created_at) VALUES ($1, $2, $3, 'complained', $4, NOW())`

	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	// For events that don't need extra metadata in update, execute directly
	if eventType == "delivered" || eventType == "unsubscribed" || eventType == "complained" {
		_, err = s.db.ExecContext(ctx, updateQuery, sendID)
		if err != nil {
			return fmt.Errorf("failed to update send record: %w", err)
		}
		_, err = s.db.ExecContext(ctx, eventQuery, orgID, campaignID, subscriberID, email)
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}
	}

	// Update campaign counters
	s.updateCampaignCounters(ctx, campaignID, eventType)

	return nil
}

func (s *EmailCampaignService) updateCampaignCounters(ctx context.Context, campaignID string, eventType string) {
	var query string
	switch eventType {
	case "delivered":
		query = `UPDATE email_campaigns SET delivered_count = delivered_count + 1, updated_at = NOW() WHERE id = $1`
	case "opened":
		query = `UPDATE email_campaigns SET opened_count = opened_count + 1, updated_at = NOW() WHERE id = $1`
	case "clicked":
		query = `UPDATE email_campaigns SET clicked_count = clicked_count + 1, updated_at = NOW() WHERE id = $1`
	case "bounced":
		query = `UPDATE email_campaigns SET bounced_count = bounced_count + 1, updated_at = NOW() WHERE id = $1`
	case "unsubscribed":
		query = `UPDATE email_campaigns SET unsubscribed_count = unsubscribed_count + 1, updated_at = NOW() WHERE id = $1`
	case "complained":
		query = `UPDATE email_campaigns SET complained_count = complained_count + 1, updated_at = NOW() WHERE id = $1`
	default:
		return
	}
	s.db.ExecContext(ctx, query, campaignID)
}

// ===========================
// GET CAMPAIGN ANALYTICS
// ===========================

// GetCampaignAnalytics returns analytics data for a campaign
func (s *EmailCampaignService) GetCampaignAnalytics(ctx context.Context, campaignID string) (*CampaignAnalytics, error) {
	var analytics CampaignAnalytics

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, status,
			COALESCE(sent_count, 0), COALESCE(delivered_count, 0),
			COALESCE(opened_count, 0), COALESCE(clicked_count, 0),
			COALESCE(bounced_count, 0), COALESCE(unsubscribed_count, 0),
			COALESCE(complained_count, 0),
			COALESCE(open_rate, 0), COALESCE(click_rate, 0),
			COALESCE(bounce_rate, 0), COALESCE(unsubscribe_rate, 0)
		FROM email_campaigns WHERE id = $1
	`, campaignID).Scan(
		&analytics.CampaignID, &analytics.CampaignName, &analytics.Status,
		&analytics.TotalSent, &analytics.Delivered,
		&analytics.Opened, &analytics.Clicked,
		&analytics.Bounced, &analytics.Unsubscribed,
		&analytics.Complained,
		&analytics.OpenRate, &analytics.ClickRate,
		&analytics.BounceRate, &analytics.UnsubscribeRate,
	)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Also get pending count from email_sends
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM email_sends WHERE campaign_id = $1 AND status = 'pending'
	`, campaignID).Scan(&analytics.Pending)
	if err != nil {
		analytics.Pending = 0
	}

	return &analytics, nil
}

// ===========================
// A/B TESTING
// ===========================

// CreateABTest creates an A/B test between two campaigns
func (s *EmailCampaignService) CreateABTest(ctx context.Context, req CreateABTestRequest) (*ABTest, error) {
	// Validate campaigns exist and are in draft status
	var statusA, statusB string
	err := s.db.QueryRowContext(ctx, `SELECT status FROM email_campaigns WHERE id = $1`, req.CampaignIDA).Scan(&statusA)
	if err != nil {
		return nil, fmt.Errorf("campaign A not found: %w", err)
	}
	err = s.db.QueryRowContext(ctx, `SELECT status FROM email_campaigns WHERE id = $1`, req.CampaignIDB).Scan(&statusB)
	if err != nil {
		return nil, fmt.Errorf("campaign B not found: %w", err)
	}

	if statusA != "draft" || statusB != "draft" {
		return nil, errors.New("both campaigns must be in draft status")
	}

	testMetric := req.TestMetric
	if testMetric == "" {
		testMetric = "open_rate"
	}
	sampleSize := req.SampleSizePercent
	if sampleSize <= 0 || sampleSize > 100 {
		sampleSize = 20
	}

	abTest := &ABTest{
		ID:                uuid.New().String(),
		CampaignIDA:       req.CampaignIDA,
		CampaignIDB:       req.CampaignIDB,
		Status:            "running",
		TestMetric:        testMetric,
		SampleSizePercent: sampleSize,
		StartedAt:         time.Now(),
		CreatedAt:         time.Now(),
	}

	// Get campaign names
	s.db.QueryRowContext(ctx, `SELECT name FROM email_campaigns WHERE id = $1`, req.CampaignIDA).Scan(&abTest.CampaignNameA)
	s.db.QueryRowContext(ctx, `SELECT name FROM email_campaigns WHERE id = $1`, req.CampaignIDB).Scan(&abTest.CampaignNameB)

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO email_campaign_ab_tests (id, campaign_id_a, campaign_id_b, status, test_metric, sample_size_percent, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, abTest.ID, abTest.CampaignIDA, abTest.CampaignIDB, abTest.Status, abTest.TestMetric, abTest.SampleSizePercent, abTest.StartedAt, abTest.CreatedAt, abTest.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create A/B test: %w", err)
	}

	return abTest, nil
}

// ListABTests returns all A/B tests for an organization
func (s *EmailCampaignService) ListABTests(ctx context.Context, orgID string) ([]*ABTest, error) {
	query := `
		SELECT abt.id, abt.campaign_id_a, abt.campaign_id_b, abt.status, abt.winner,
			abt.test_metric, abt.sample_size_percent, abt.started_at, abt.completed_at, abt.created_at,
			ca.name, cb.name
		FROM email_campaign_ab_tests abt
		JOIN email_campaigns ca ON abt.campaign_id_a = ca.id
		JOIN email_campaigns cb ON abt.campaign_id_b = cb.id
		WHERE ca.organization_id = $1
		ORDER BY abt.created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*ABTest
	for rows.Next() {
		var t ABTest
		err := rows.Scan(&t.ID, &t.CampaignIDA, &t.CampaignIDB, &t.Status, &t.Winner,
			&t.TestMetric, &t.SampleSizePercent, &t.StartedAt, &t.CompletedAt, &t.CreatedAt,
			&t.CampaignNameA, &t.CampaignNameB)
		if err != nil {
			return nil, err
		}
		tests = append(tests, &t)
	}
	return tests, nil
}

// CompleteABTest analyzes the results and declares a winner
func (s *EmailCampaignService) CompleteABTest(ctx context.Context, testID string) (*ABTest, error) {
	var abTest ABTest
	err := s.db.QueryRowContext(ctx, `
		SELECT id, campaign_id_a, campaign_id_b, status, winner, test_metric,
			sample_size_percent, started_at, completed_at, created_at
		FROM email_campaign_ab_tests WHERE id = $1
	`, testID).Scan(
		&abTest.ID, &abTest.CampaignIDA, &abTest.CampaignIDB, &abTest.Status,
		&abTest.Winner, &abTest.TestMetric, &abTest.SampleSizePercent,
		&abTest.StartedAt, &abTest.CompletedAt, &abTest.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("A/B test not found: %w", err)
	}

	if abTest.Status == "completed" {
		return nil, errors.New("A/B test is already completed")
	}

	// Get stats for both campaigns
	var openA, clickA, sentA, deliveredA int
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(opened_count, 0), COALESCE(clicked_count, 0), COALESCE(sent_count, 0), COALESCE(delivered_count, 0)
		FROM email_campaigns WHERE id = $1
	`, abTest.CampaignIDA).Scan(&openA, &clickA, &sentA, &deliveredA)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign A stats: %w", err)
	}

	var openB, clickB, sentB, deliveredB int
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(opened_count, 0), COALESCE(clicked_count, 0), COALESCE(sent_count, 0), COALESCE(delivered_count, 0)
		FROM email_campaigns WHERE id = $1
	`, abTest.CampaignIDB).Scan(&openB, &clickB, &sentB, &deliveredB)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign B stats: %w", err)
	}

	// Determine winner based on test metric
	var winner string
	switch abTest.TestMetric {
	case "open_rate":
		if deliveredA > 0 && deliveredB > 0 {
			rateA := float64(openA) / float64(deliveredA)
			rateB := float64(openB) / float64(deliveredB)
			if rateA > rateB {
				winner = "A"
			} else {
				winner = "B"
			}
		} else if deliveredA > 0 {
			winner = "A"
		} else if deliveredB > 0 {
			winner = "B"
		} else {
			winner = "A" // Default to A if no deliveries
		}
	case "click_rate":
		if deliveredA > 0 && deliveredB > 0 {
			rateA := float64(clickA) / float64(deliveredA)
			rateB := float64(clickB) / float64(deliveredB)
			if rateA > rateB {
				winner = "A"
			} else {
				winner = "B"
			}
		} else if deliveredA > 0 {
			winner = "A"
		} else if deliveredB > 0 {
			winner = "B"
		} else {
			winner = "A"
		}
	default:
		winner = "A"
	}

	now := time.Now()
	winnerCampaignID := abTest.CampaignIDA
	if winner == "B" {
		winnerCampaignID = abTest.CampaignIDB
	}

	// Update A/B test
	_, err = s.db.ExecContext(ctx, `
		UPDATE email_campaign_ab_tests
		SET status = 'completed', winner = $1, completed_at = $2, updated_at = NOW()
		WHERE id = $3
	`, winner, now, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to update A/B test: %w", err)
	}

	// Update winner campaign status to "sent" and loser to "cancelled"
	winnerStatus := "sent"
	s.db.ExecContext(ctx, `
		UPDATE email_campaigns SET status = $1, updated_at = NOW() WHERE id = $2
	`, winnerStatus, winnerCampaignID)

	loserStatus := "cancelled"
	loserID := abTest.CampaignIDB
	if winner == "B" {
		loserID = abTest.CampaignIDA
	}
	s.db.ExecContext(ctx, `
		UPDATE email_campaigns SET status = $1, updated_at = NOW() WHERE id = $2
	`, loserStatus, loserID)

	abTest.Status = "completed"
	abTest.Winner = &winner
	abTest.CompletedAt = &now

	return &abTest, nil
}

// SendABTestSample sends the A/B test to a random sample of subscribers
func (s *EmailCampaignService) SendABTestSample(ctx context.Context, testID string) error {
	var abTest ABTest
	err := s.db.QueryRowContext(ctx, `
		SELECT campaign_id_a, campaign_id_b, sample_size_percent, status
		FROM email_campaign_ab_tests WHERE id = $1
	`, testID).Scan(&abTest.CampaignIDA, &abTest.CampaignIDB, &abTest.SampleSizePercent, &abTest.Status)
	if err != nil {
		return fmt.Errorf("A/B test not found: %w", err)
	}

	if abTest.Status != "running" {
		return errors.New("A/B test is not running")
	}

	// Get all subscriber IDs from both campaigns' lists
	var listIDs []string
	rows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT ecl.list_id
		FROM email_campaign_lists ecl
		WHERE ecl.campaign_id = $1 OR ecl.campaign_id = $2
	`, abTest.CampaignIDA, abTest.CampaignIDB)
	if err != nil {
		return fmt.Errorf("failed to get campaign lists: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var listID string
		if err := rows.Scan(&listID); err != nil {
			continue
		}
		listIDs = append(listIDs, listID)
	}

	if len(listIDs) == 0 {
		return errors.New("no linked lists found for A/B test campaigns")
	}

	// Get all active subscribers
	type subscriber struct {
		id        string
		email     string
		firstName string
		lastName  string
	}

	subscriberRows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT es.id, es.email, es.first_name, es.last_name
		FROM email_subscribers es
		JOIN email_list_subscribers els ON es.id = els.subscriber_id
		WHERE els.list_id = ANY($1) AND es.status = 'active'
	`, listIDs)
	if err != nil {
		return fmt.Errorf("failed to get subscribers: %w", err)
	}
	defer subscriberRows.Close()

	var subscribers []subscriber
	for subscriberRows.Next() {
		var sub subscriber
		if err := subscriberRows.Scan(&sub.id, &sub.email, &sub.firstName, &sub.lastName); err != nil {
			continue
		}
		subscribers = append(subscribers, sub)
	}

	if len(subscribers) == 0 {
		return errors.New("no active subscribers found")
	}

	// Shuffle and split
	rand.Shuffle(len(subscribers), func(i, j int) {
		subscribers[i], subscribers[j] = subscribers[j], subscribers[i]
	})

	sampleSize := len(subscribers) * abTest.SampleSizePercent / 100
	if sampleSize < 2 {
		sampleSize = 2
	}
	halfSample := sampleSize / 2

	groupA := subscribers[:halfSample]
	groupB := subscribers[halfSample:sampleSize]

	// Send to group A
	for _, sub := range groupA {
		sendID := uuid.New().String()
		s.db.ExecContext(ctx, `
			INSERT INTO email_sends (id, campaign_id, subscriber_id, status, created_at)
			VALUES ($1, $2, $3, 'pending', NOW())
		`, sendID, abTest.CampaignIDA, sub.id)

		if s.emailService != nil {
			var campaign domain.EmailCampaign
			s.db.QueryRowContext(ctx, `
				SELECT id, name, subject, html_content, text_content FROM email_campaigns WHERE id = $1
			`, abTest.CampaignIDA).Scan(&campaign.ID, &campaign.Name, &campaign.Subject, &campaign.HTMLContent, &campaign.TextContent)

			s.emailService.SendEmail(sub.email, sub.firstName+" "+sub.lastName, campaign.Subject, campaign.HTMLContent, campaign.TextContent)
		}

		s.db.ExecContext(ctx, `
			UPDATE email_sends SET status = 'sent', sent_at = NOW() WHERE id = $1
		`, sendID)
	}

	// Send to group B
	for _, sub := range groupB {
		sendID := uuid.New().String()
		s.db.ExecContext(ctx, `
			INSERT INTO email_sends (id, campaign_id, subscriber_id, status, created_at)
			VALUES ($1, $2, $3, 'pending', NOW())
		`, sendID, abTest.CampaignIDB, sub.id)

		if s.emailService != nil {
			var campaign domain.EmailCampaign
			s.db.QueryRowContext(ctx, `
				SELECT id, name, subject, html_content, text_content FROM email_campaigns WHERE id = $1
			`, abTest.CampaignIDB).Scan(&campaign.ID, &campaign.Name, &campaign.Subject, &campaign.HTMLContent, &campaign.TextContent)

			s.emailService.SendEmail(sub.email, sub.firstName+" "+sub.lastName, campaign.Subject, campaign.HTMLContent, campaign.TextContent)
		}

		s.db.ExecContext(ctx, `
			UPDATE email_sends SET status = 'sent', sent_at = NOW() WHERE id = $1
		`, sendID)
	}

	return nil
}

// TrackPixelURL generates a tracking pixel URL for email opens
func (s *EmailCampaignService) TrackPixelURL(sendID string) string {
	return fmt.Sprintf("/api/v1/email/track/open?send_id=%s", sendID)
}

// TrackLinkURL generates a tracked link URL for email clicks
func (s *EmailCampaignService) TrackLinkURL(sendID string, originalURL string) string {
	return fmt.Sprintf("/api/v1/email/track/click?send_id=%s&url=%s", sendID, originalURL)
}
