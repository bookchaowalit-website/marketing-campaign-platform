package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type calendarEventRepository struct {
	db *sql.DB
}

// NewCalendarEventRepository creates a new calendar event repository
func NewCalendarEventRepository(db *sql.DB) domain.CalendarEventRepository {
	return &calendarEventRepository{db: db}
}

func (r *calendarEventRepository) Create(ctx context.Context, event *domain.CalendarEvent) error {
	query := `
		INSERT INTO marketing_calendar (id, organization_id, campaign_id, event_type, title, description,
			start_time, end_time, all_day, color, channel, status, location, location_url, attendees, tags,
			reminders, recurrence, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	var campaignID, endTime, channel, location, locationURL, createdBy interface{}
	if event.CampaignID != nil {
		campaignID = *event.CampaignID
	}
	if event.EndTime != nil {
		endTime = *event.EndTime
	}
	if event.Channel != nil {
		channel = *event.Channel
	}
	if event.Location != nil {
		location = *event.Location
	}
	if event.LocationURL != nil {
		locationURL = *event.LocationURL
	}
	if event.CreatedBy != nil {
		createdBy = *event.CreatedBy
	}

	var reminders, recurrence interface{}
	if event.Reminders != nil {
		reminders = string(event.Reminders)
	}
	if event.Recurrence != nil {
		recurrence = string(event.Recurrence)
	}

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.OrganizationID,
		campaignID,
		event.EventType,
		event.Title,
		event.Description,
		event.StartTime,
		endTime,
		event.AllDay,
		event.Color,
		channel,
		event.Status,
		location,
		locationURL,
		event.Attendees,
		event.Tags,
		reminders,
		recurrence,
		createdBy,
		event.CreatedAt,
		event.UpdatedAt,
	)
	return err
}

func (r *calendarEventRepository) GetByID(ctx context.Context, id string) (*domain.CalendarEvent, error) {
	query := `
		SELECT id, organization_id, campaign_id, event_type, title, description,
			start_time, end_time, all_day, color, channel, status, location, location_url, attendees, tags,
			reminders, recurrence, created_by, created_at, updated_at
		FROM marketing_calendar
		WHERE id = $1
	`
	event := &domain.CalendarEvent{}
	var campaignID, endTime, channel, location, locationURL, createdBy sql.NullString
	var reminders, recurrence sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.OrganizationID,
		&campaignID,
		&event.EventType,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&endTime,
		&event.AllDay,
		&event.Color,
		&channel,
		&event.Status,
		&location,
		&locationURL,
		&event.Attendees,
		&event.Tags,
		&reminders,
		&recurrence,
		&createdBy,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if campaignID.Valid {
		event.CampaignID = &campaignID.String
	}
	if endTime.Valid {
		t, _ := time.Parse(time.RFC3339, endTime.String)
		event.EndTime = &t
	}
	if channel.Valid {
		event.Channel = &channel.String
	}
	if location.Valid {
		event.Location = &location.String
	}
	if locationURL.Valid {
		event.LocationURL = &locationURL.String
	}
	if createdBy.Valid {
		event.CreatedBy = &createdBy.String
	}
	if reminders.Valid {
		event.Reminders = []byte(reminders.String)
	}
	if recurrence.Valid {
		event.Recurrence = []byte(recurrence.String)
	}
	return event, nil
}

func (r *calendarEventRepository) List(ctx context.Context, organizationID string, start, end *time.Time, channel, eventType, status *string, limit, offset int) ([]*domain.CalendarEvent, error) {
	query := `
		SELECT id, organization_id, campaign_id, event_type, title, description,
			start_time, end_time, all_day, color, channel, status, location, location_url, attendees, tags,
			reminders, recurrence, created_by, created_at, updated_at
		FROM marketing_calendar
		WHERE organization_id = $1
	`
	args := []interface{}{organizationID}
	argIndex := 2

	if start != nil {
		query += " AND start_time >= $" + strconv.Itoa(argIndex)
		args = append(args, *start)
		argIndex++
	}
	if end != nil {
		query += " AND start_time <= $" + strconv.Itoa(argIndex)
		args = append(args, *end)
		argIndex++
	}
	if channel != nil {
		query += " AND channel = $" + strconv.Itoa(argIndex)
		args = append(args, *channel)
		argIndex++
	}
	if eventType != nil {
		query += " AND event_type = $" + strconv.Itoa(argIndex)
		args = append(args, *eventType)
		argIndex++
	}
	if status != nil {
		query += " AND status = $" + strconv.Itoa(argIndex)
		args = append(args, *status)
		argIndex++
	}

	query += " ORDER BY start_time ASC"
	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*domain.CalendarEvent, 0)
	for rows.Next() {
		e := &domain.CalendarEvent{}
		var campaignID, endTime, ch, location, locationURL, createdBy sql.NullString
		var reminders, recurrence sql.NullString

		err := rows.Scan(
			&e.ID,
			&e.OrganizationID,
			&campaignID,
			&e.EventType,
			&e.Title,
			&e.Description,
			&e.StartTime,
			&endTime,
			&e.AllDay,
			&e.Color,
			&ch,
			&e.Status,
			&location,
			&locationURL,
			&e.Attendees,
			&e.Tags,
			&reminders,
			&recurrence,
			&createdBy,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if campaignID.Valid {
			e.CampaignID = &campaignID.String
		}
		if endTime.Valid {
			t, _ := time.Parse(time.RFC3339, endTime.String)
			e.EndTime = &t
		}
		if ch.Valid {
			e.Channel = &ch.String
		}
		if location.Valid {
			e.Location = &location.String
		}
		if locationURL.Valid {
			e.LocationURL = &locationURL.String
		}
		if createdBy.Valid {
			e.CreatedBy = &createdBy.String
		}
		if reminders.Valid {
			e.Reminders = []byte(reminders.String)
		}
		if recurrence.Valid {
			e.Recurrence = []byte(recurrence.String)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *calendarEventRepository) Update(ctx context.Context, event *domain.CalendarEvent) error {
	query := `
		UPDATE marketing_calendar
		SET campaign_id = $2, event_type = $3, title = $4, description = $5,
			start_time = $6, end_time = $7, all_day = $8, color = $9, channel = $10,
			status = $11, location = $12, location_url = $13, attendees = $14, tags = $15,
			reminders = $16, recurrence = $17, updated_at = $18
		WHERE id = $1
	`
	event.UpdatedAt = time.Now()

	var campaignID, endTime, channel, location, locationURL interface{}
	if event.CampaignID != nil {
		campaignID = *event.CampaignID
	}
	if event.EndTime != nil {
		endTime = *event.EndTime
	}
	if event.Channel != nil {
		channel = *event.Channel
	}
	if event.Location != nil {
		location = *event.Location
	}
	if event.LocationURL != nil {
		locationURL = *event.LocationURL
	}

	var reminders, recurrence interface{}
	if event.Reminders != nil {
		reminders = string(event.Reminders)
	}
	if event.Recurrence != nil {
		recurrence = string(event.Recurrence)
	}

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		campaignID,
		event.EventType,
		event.Title,
		event.Description,
		event.StartTime,
		endTime,
		event.AllDay,
		event.Color,
		channel,
		event.Status,
		location,
		locationURL,
		event.Attendees,
		event.Tags,
		reminders,
		recurrence,
		event.UpdatedAt,
	)
	return err
}

func (r *calendarEventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM marketing_calendar WHERE id = $1`, id)
	return err
}

// ===========================
// Calendar View Repository
// ===========================

type calendarViewRepository struct {
	db *sql.DB
}

// NewCalendarViewRepository creates a new calendar view repository
func NewCalendarViewRepository(db *sql.DB) domain.CalendarViewRepository {
	return &calendarViewRepository{db: db}
}

func (r *calendarViewRepository) Create(ctx context.Context, view *domain.CalendarView) error {
	query := `
		INSERT INTO calendar_views (id, organization_id, user_id, name, view_type, filters, channels, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now()
	view.CreatedAt = now

	var userID interface{}
	if view.UserID != nil {
		userID = *view.UserID
	}

	var filters interface{}
	if view.Filters != nil {
		filters = string(view.Filters)
	}

	_, err := r.db.ExecContext(ctx, query,
		view.ID,
		view.OrganizationID,
		userID,
		view.Name,
		view.ViewType,
		filters,
		view.Channels,
		view.IsDefault,
		view.CreatedAt,
	)
	return err
}

func (r *calendarViewRepository) List(ctx context.Context, organizationID string) ([]*domain.CalendarView, error) {
	query := `
		SELECT id, organization_id, user_id, name, view_type, filters, channels, is_default, created_at
		FROM calendar_views
		WHERE organization_id = $1
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	views := make([]*domain.CalendarView, 0)
	for rows.Next() {
		v := &domain.CalendarView{}
		var userID sql.NullString
		var filters sql.NullString

		err := rows.Scan(
			&v.ID,
			&v.OrganizationID,
			&userID,
			&v.Name,
			&v.ViewType,
			&filters,
			&v.Channels,
			&v.IsDefault,
			&v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if userID.Valid {
			v.UserID = &userID.String
		}
		if filters.Valid {
			v.Filters = []byte(filters.String)
		}
		views = append(views, v)
	}
	return views, rows.Err()
}

func (r *calendarViewRepository) GetDefault(ctx context.Context, organizationID, userID string) (*domain.CalendarView, error) {
	query := `
		SELECT id, organization_id, user_id, name, view_type, filters, channels, is_default, created_at
		FROM calendar_views
		WHERE organization_id = $1 AND is_default = true
		ORDER BY created_at DESC
		LIMIT 1
	`
	v := &domain.CalendarView{}
	var dbUserID sql.NullString
	var filters sql.NullString

	err := r.db.QueryRowContext(ctx, query, organizationID).Scan(
		&v.ID,
		&v.OrganizationID,
		&dbUserID,
		&v.Name,
		&v.ViewType,
		&filters,
		&v.Channels,
		&v.IsDefault,
		&v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if dbUserID.Valid {
		v.UserID = &dbUserID.String
	}
	if filters.Valid {
		v.Filters = []byte(filters.String)
	}
	return v, nil
}

func (r *calendarViewRepository) SetDefault(ctx context.Context, id, organizationID string) error {
	query := `
		UPDATE calendar_views
		SET is_default = (id = $1)
		WHERE organization_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, id, organizationID)
	return err
}
