package domain

import (
	"context"
	"time"
)

// CalendarEvent represents a marketing calendar event
type CalendarEvent struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	CampaignID     *string    `json:"campaign_id,omitempty"`
	EventType      string     `json:"event_type"`
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	AllDay         bool       `json:"all_day"`
	Color          string     `json:"color"`
	Channel        *string    `json:"channel,omitempty"`
	Status         string     `json:"status"`
	Location       *string    `json:"location,omitempty"`
	LocationURL    *string    `json:"location_url,omitempty"`
	Attendees      []string   `json:"attendees,omitempty"`
	Tags           []string   `json:"tags,omitempty"`
	Reminders      []byte     `json:"reminders,omitempty"`
	Recurrence     []byte     `json:"recurrence,omitempty"`
	CreatedBy      *string    `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CalendarView represents a saved calendar view for a user
type CalendarView struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	UserID         *string   `json:"user_id,omitempty"`
	Name           string    `json:"name"`
	ViewType       string    `json:"view_type"`
	Filters        []byte    `json:"filters,omitempty"`
	Channels       []string  `json:"channels,omitempty"`
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
}

// CalendarEventRepository defines the interface for calendar event repository
type CalendarEventRepository interface {
	Create(ctx context.Context, event *CalendarEvent) error
	GetByID(ctx context.Context, id string) (*CalendarEvent, error)
	List(ctx context.Context, organizationID string, start, end *time.Time, channel, eventType, status *string, limit, offset int) ([]*CalendarEvent, error)
	Update(ctx context.Context, event *CalendarEvent) error
	Delete(ctx context.Context, id string) error
}

// CalendarViewRepository defines the interface for calendar view repository
type CalendarViewRepository interface {
	Create(ctx context.Context, view *CalendarView) error
	List(ctx context.Context, organizationID string) ([]*CalendarView, error)
	GetDefault(ctx context.Context, organizationID, userID string) (*CalendarView, error)
	SetDefault(ctx context.Context, id, organizationID string) error
}
