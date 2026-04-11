package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// CalendarService handles business logic for marketing calendar
type CalendarService struct {
	eventRepo domain.CalendarEventRepository
	viewRepo  domain.CalendarViewRepository
}

// NewCalendarService creates a new calendar service
func NewCalendarService(
	eventRepo domain.CalendarEventRepository,
	viewRepo domain.CalendarViewRepository,
) *CalendarService {
	return &CalendarService{
		eventRepo: eventRepo,
		viewRepo:  viewRepo,
	}
}

// CreateEvent creates a new calendar event
func (s *CalendarService) CreateEvent(ctx context.Context, event *domain.CalendarEvent) error {
	return s.eventRepo.Create(ctx, event)
}

// GetEvent retrieves a calendar event by ID
func (s *CalendarService) GetEvent(ctx context.Context, id string) (*domain.CalendarEvent, error) {
	return s.eventRepo.GetByID(ctx, id)
}

// ListEvents retrieves calendar events filtered by date range and optional filters
func (s *CalendarService) ListEvents(ctx context.Context, organizationID string, start, end *time.Time, channel, eventType, status *string, limit, offset int) ([]*domain.CalendarEvent, error) {
	return s.eventRepo.List(ctx, organizationID, start, end, channel, eventType, status, limit, offset)
}

// UpdateEvent updates an existing calendar event
func (s *CalendarService) UpdateEvent(ctx context.Context, event *domain.CalendarEvent) error {
	return s.eventRepo.Update(ctx, event)
}

// DeleteEvent deletes a calendar event
func (s *CalendarService) DeleteEvent(ctx context.Context, id string) error {
	return s.eventRepo.Delete(ctx, id)
}

// CreateView creates a new calendar view
func (s *CalendarService) CreateView(ctx context.Context, view *domain.CalendarView) error {
	return s.viewRepo.Create(ctx, view)
}

// ListViews retrieves all calendar views for an organization
func (s *CalendarService) ListViews(ctx context.Context, organizationID string) ([]*domain.CalendarView, error) {
	return s.viewRepo.List(ctx, organizationID)
}

// GetDefaultView retrieves the default calendar view for an organization
func (s *CalendarService) GetDefaultView(ctx context.Context, organizationID, userID string) (*domain.CalendarView, error) {
	return s.viewRepo.GetDefault(ctx, organizationID, userID)
}

// SetDefaultView sets a calendar view as the default for an organization
func (s *CalendarService) SetDefaultView(ctx context.Context, id, organizationID string) error {
	view, err := s.viewRepo.GetDefault(ctx, organizationID, "")
	if err != nil {
		return fmt.Errorf("failed to get default view: %w", err)
	}
	if view == nil || view.ID != id {
		return s.viewRepo.SetDefault(ctx, id, organizationID)
	}
	return nil
}
