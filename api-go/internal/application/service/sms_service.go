package service

import (
	"context"
	"errors"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type SMSService struct {
	listRepo       domain.SMSListRepository
	subscriberRepo domain.SMSSubscriberRepository
	campaignRepo   domain.SMSCampaignRepository
	logRepo        domain.SMSLogRepository
	twilioClient   *TwilioClient
}

func NewSMSService(
	listRepo domain.SMSListRepository,
	subscriberRepo domain.SMSSubscriberRepository,
	campaignRepo domain.SMSCampaignRepository,
	logRepo domain.SMSLogRepository,
	twilioClient *TwilioClient,
) *SMSService {
	return &SMSService{
		listRepo:       listRepo,
		subscriberRepo: subscriberRepo,
		campaignRepo:   campaignRepo,
		logRepo:        logRepo,
		twilioClient:   twilioClient,
	}
}

// ===========================
// LIST OPERATIONS
// ===========================

func (s *SMSService) CreateList(ctx context.Context, l *domain.SMSList) error {
	l.Status = "active"
	return s.listRepo.Create(l)
}

func (s *SMSService) GetList(ctx context.Context, id string) (*domain.SMSList, error) {
	return s.listRepo.GetByID(id)
}

func (s *SMSService) ListLists(ctx context.Context, orgID string) ([]*domain.SMSList, error) {
	return s.listRepo.ListByOrg(orgID)
}

// ===========================
// SUBSCRIBER OPERATIONS
// ===========================

func (s *SMSService) AddSubscriber(ctx context.Context, sub *domain.SMSSubscriber) error {
	if !domain.ValidSubscriberStatuses[sub.Status] {
		sub.Status = "active"
	}
	return s.subscriberRepo.Create(sub)
}

func (s *SMSService) GetSubscriber(ctx context.Context, id string) (*domain.SMSSubscriber, error) {
	return s.subscriberRepo.GetByID(id)
}

func (s *SMSService) ListSubscribers(ctx context.Context, orgID, status string, limit, offset int) ([]*domain.SMSSubscriber, int, error) {
	return s.subscriberRepo.ListByOrg(orgID, status, limit, offset)
}

func (s *SMSService) OptOutSubscriber(ctx context.Context, id, reason string) error {
	sub, err := s.subscriberRepo.GetByID(id)
	if err != nil {
		return err
	}

	err = s.subscriberRepo.OptOut(id, reason)
	if err != nil {
		return err
	}

	// Update list subscriber count
	if sub.ListID != nil {
		list, err := s.listRepo.GetByID(*sub.ListID)
		if err == nil && list.SubscriberCount > 0 {
			s.listRepo.UpdateSubscriberCount(*sub.ListID, list.SubscriberCount-1)
		}
	}

	return nil
}

// ===========================
// CAMPAIGN OPERATIONS
// ===========================

func (s *SMSService) CreateCampaign(ctx context.Context, c *domain.SMSCampaign) error {
	c.Status = "draft"
	return s.campaignRepo.Create(c)
}

func (s *SMSService) GetCampaign(ctx context.Context, id string) (*domain.SMSCampaign, error) {
	return s.campaignRepo.GetByID(id)
}

func (s *SMSService) ListCampaigns(ctx context.Context, orgID, status string, limit, offset int) ([]*domain.SMSCampaign, int, error) {
	return s.campaignRepo.ListByOrg(orgID, status, limit, offset)
}

func (s *SMSService) SendCampaign(ctx context.Context, id string) error {
	campaign, err := s.campaignRepo.GetByID(id)
	if err != nil {
		return err
	}

	if campaign.Status != "draft" {
		return errors.New("campaign must be in draft status to send")
	}

	// Update campaign status to sending
	if err := s.campaignRepo.UpdateStatus(id, "sending"); err != nil {
		return err
	}

	// Get subscribers from the list
	if campaign.ListID == nil {
		return errors.New("campaign must have a list assigned")
	}

	subscribers, _, err := s.subscriberRepo.ListByList(*campaign.ListID, "active", 10000, 0)
	if err != nil {
		s.campaignRepo.UpdateStatus(id, "failed")
		return err
	}

	if len(subscribers) == 0 {
		s.campaignRepo.UpdateStatus(id, "failed")
		return errors.New("no active subscribers in the list")
	}

	// Send SMS to each subscriber
	sentCount := 0
	deliveredCount := 0
	failedCount := 0
	totalCost := 0.0

	for _, sub := range subscribers {
		logEntry := &domain.SMSLog{
			CampaignID:   id,
			SubscriberID: &sub.ID,
			Phone:        sub.Phone,
			Status:       "pending",
		}

		if err := s.logRepo.Create(logEntry); err != nil {
			failedCount++
			continue
		}

		// Replace template variables
		message := campaign.MessageTemplate
		message = replaceTemplateVar(message, "first_name", sub.FirstName)
		message = replaceTemplateVar(message, "last_name", sub.LastName)

		// Send via Twilio
		if s.twilioClient != nil {
			sid, cost, err := s.twilioClient.SendMessage(sub.Phone, message)
			if err != nil {
				s.logRepo.UpdateStatus(logEntry.ID, "failed", "", err.Error())
				failedCount++
				continue
			}

			s.logRepo.UpdateStatus(logEntry.ID, "sent", sid, "")
			sentCount++

			if cost > 0 {
				totalCost += cost
			}

			// Simulate delivery (in production, handle webhook callbacks)
			deliveredCount++
		} else {
			// Stub mode: simulate sending
			s.logRepo.UpdateStatus(logEntry.ID, "sent", "stub_"+logEntry.ID, "")
			sentCount++
			deliveredCount++
		}
	}

	// Update campaign stats
	status := "sent"
	if failedCount > 0 && sentCount == 0 {
		status = "failed"
	}
	if err := s.campaignRepo.UpdateSendStats(id, sentCount, deliveredCount, failedCount, totalCost); err != nil {
		return err
	}

	return nil
}

// replaceTemplateVar replaces {{var_name}} with the given value
func replaceTemplateVar(template, varName, value string) string {
	if value == "" {
		return template
	}
	placeholder := "{{" + varName + "}}"
	for i := 0; i < 10; i++ { // Safety limit
		newStr := replaceFirst(template, placeholder, value)
		if newStr == template {
			break
		}
		template = newStr
	}
	return template
}

func replaceFirst(s, old, new string) string {
	idx := indexOf(s, old)
	if idx == -1 {
		return s
	}
	return s[:idx] + new + s[idx+len(old):]
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
