package dto

// CreateCalendarEventRequest represents the request to create a calendar event
type CreateCalendarEventRequest struct {
	CampaignID  *string  `json:"campaign_id,omitempty"`
	EventType   string   `json:"event_type" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	StartTime   string   `json:"start_time" binding:"required"`
	EndTime     *string  `json:"end_time,omitempty"`
	AllDay      bool     `json:"all_day"`
	Color       string   `json:"color"`
	Channel     *string  `json:"channel,omitempty"`
	Status      string   `json:"status"`
	Location    *string  `json:"location,omitempty"`
	LocationURL *string  `json:"location_url,omitempty"`
	Attendees   []string `json:"attendees,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Reminders   string   `json:"reminders,omitempty"`
	Recurrence  string   `json:"recurrence,omitempty"`
}

// UpdateCalendarEventRequest represents the request to update a calendar event
type UpdateCalendarEventRequest struct {
	CampaignID  *string  `json:"campaign_id,omitempty"`
	EventType   string   `json:"event_type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	StartTime   string   `json:"start_time"`
	EndTime     *string  `json:"end_time,omitempty"`
	AllDay      *bool    `json:"all_day"`
	Color       string   `json:"color"`
	Channel     *string  `json:"channel,omitempty"`
	Status      string   `json:"status"`
	Location    *string  `json:"location,omitempty"`
	LocationURL *string  `json:"location_url,omitempty"`
	Attendees   []string `json:"attendees,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Reminders   string   `json:"reminders,omitempty"`
	Recurrence  string   `json:"recurrence,omitempty"`
}

// CalendarEventResponse represents the calendar event response
type CalendarEventResponse struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	CampaignID     *string  `json:"campaign_id,omitempty"`
	EventType      string   `json:"event_type"`
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	StartTime      string   `json:"start_time"`
	EndTime        *string  `json:"end_time,omitempty"`
	AllDay         bool     `json:"all_day"`
	Color          string   `json:"color"`
	Channel        *string  `json:"channel,omitempty"`
	Status         string   `json:"status"`
	Location       *string  `json:"location,omitempty"`
	LocationURL    *string  `json:"location_url,omitempty"`
	Attendees      []string `json:"attendees,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// CreateCalendarViewRequest represents the request to create a calendar view
type CreateCalendarViewRequest struct {
	Name      string   `json:"name" binding:"required"`
	ViewType  string   `json:"view_type" binding:"required"`
	Filters   string   `json:"filters,omitempty"`
	Channels  []string `json:"channels,omitempty"`
	IsDefault bool     `json:"is_default"`
}

// CalendarViewResponse represents the calendar view response
type CalendarViewResponse struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	UserID         *string  `json:"user_id,omitempty"`
	Name           string   `json:"name"`
	ViewType       string   `json:"view_type"`
	Channels       []string `json:"channels,omitempty"`
	IsDefault      bool     `json:"is_default"`
	CreatedAt      string   `json:"created_at"`
}
