package domain

import "time"

// ABTest represents an A/B test entity
type ABTest struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	CampaignID     string     `json:"campaign_id,omitempty"`
	Name           string     `json:"name"`
	Variant        string     `json:"variant"` // A or B
	Status         string     `json:"status"`  // running, completed, draft
	Visitors       int64      `json:"visitors"`
	Conversions    int64      `json:"conversions"`
	ConversionRate float64    `json:"conversion_rate"`
	Confidence     float64    `json:"confidence"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ABTestRepository defines the interface for A/B test repository
type ABTestRepository interface {
	Create(test *ABTest) error
	GetByID(id string) (*ABTest, error)
	GetAll(userID string, limit, offset int) ([]*ABTest, error)
	Update(test *ABTest) error
	Delete(id string) error
	GetByStatus(userID, status string) ([]*ABTest, error)
}