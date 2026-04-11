package domain

import (
	"time"
)

// ===========================
// SMS LIST
// ===========================

type SMSList struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	SubscriberCount int       `json:"subscriber_count"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ===========================
// SMS SUBSCRIBER
// ===========================

type SMSSubscriber struct {
	ID            string                 `json:"id"`
	OrganizationID string                `json:"organization_id"`
	ListID        *string                `json:"list_id,omitempty"`
	Phone         string                 `json:"phone"`
	FirstName     string                 `json:"first_name,omitempty"`
	LastName      string                 `json:"last_name,omitempty"`
	CustomFields  map[string]interface{} `json:"custom_fields,omitempty"`
	Status        string                 `json:"status"`
	SubscribedAt  time.Time              `json:"subscribed_at"`
	OptedOutAt    *time.Time             `json:"opted_out_at,omitempty"`
	OptOutReason  string                 `json:"opt_out_reason,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// ===========================
// SMS CAMPAIGN
// ===========================

type SMSCampaign struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	ListID         *string    `json:"list_id,omitempty"`
	Name           string     `json:"name"`
	MessageTemplate string    `json:"message_template"`
	Status         string     `json:"status"`
	ScheduledAt    *time.Time `json:"scheduled_at,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	SentCount      int        `json:"sent_count"`
	DeliveredCount int        `json:"delivered_count"`
	FailedCount    int        `json:"failed_count"`
	OptOutCount    int        `json:"opt_out_count"`
	Cost           float64    `json:"cost"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ===========================
// SMS LOG
// ===========================

type SMSLog struct {
	ID           string     `json:"id"`
	CampaignID   string     `json:"campaign_id"`
	SubscriberID *string    `json:"subscriber_id,omitempty"`
	Phone        string     `json:"phone"`
	Status       string     `json:"status"`
	MessageSID   *string    `json:"message_sid,omitempty"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	Cost         float64    `json:"cost,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ===========================
// REPOSITORY INTERFACES
// ===========================

type SMSListRepository interface {
	Create(l *SMSList) error
	GetByID(id string) (*SMSList, error)
	ListByOrg(orgID string) ([]*SMSList, error)
	Update(l *SMSList) error
	Delete(id string) error
	UpdateSubscriberCount(id string, count int) error
}

type SMSSubscriberRepository interface {
	Create(s *SMSSubscriber) error
	GetByID(id string) (*SMSSubscriber, error)
	GetByPhone(orgID, phone string) (*SMSSubscriber, error)
	ListByOrg(orgID, status string, limit, offset int) ([]*SMSSubscriber, int, error)
	ListByList(listID string, status string, limit, offset int) ([]*SMSSubscriber, int, error)
	Update(s *SMSSubscriber) error
	Delete(id string) error
	OptOut(id string, reason string) error
}

type SMSCampaignRepository interface {
	Create(c *SMSCampaign) error
	GetByID(id string) (*SMSCampaign, error)
	ListByOrg(orgID, status string, limit, offset int) ([]*SMSCampaign, int, error)
	Update(c *SMSCampaign) error
	Delete(id string) error
	UpdateStatus(id, status string) error
	UpdateSendStats(id string, sentCount, deliveredCount, failedCount int, cost float64) error
}

type SMSLogRepository interface {
	Create(l *SMSLog) error
	GetByID(id string) (*SMSLog, error)
	ListByCampaign(campaignID string) ([]*SMSLog, error)
	ListBySubscriber(subscriberID string, limit, offset int) ([]*SMSLog, int, error)
	UpdateStatus(id, status string, messageSID string, errMsg string) error
}

// ===========================
// VALIDATION CONSTANTS
// ===========================

var ValidSubscriberStatuses = map[string]bool{
	"active":    true,
	"opted_out": true,
	"bounced":   true,
}

var ValidSMSCampaignStatuses = map[string]bool{
	"draft":     true,
	"scheduled": true,
	"sending":   true,
	"sent":      true,
	"failed":    true,
	"cancelled": true,
}
