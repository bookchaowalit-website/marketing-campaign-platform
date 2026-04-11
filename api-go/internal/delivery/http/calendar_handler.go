package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/dto"
	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CalendarHandler handles HTTP requests for marketing calendar
type CalendarHandler struct {
	service *service.CalendarService
}

// NewCalendarHandler creates a new calendar handler
func NewCalendarHandler(service *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

// CreateEvent handles POST /api/v1/calendar/events
func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	var req dto.CreateCalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format, use RFC3339"})
		return
	}

	var endTime *time.Time
	if req.EndTime != nil {
		t, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format, use RFC3339"})
			return
		}
		endTime = &t
	}

	status := req.Status
	if status == "" {
		status = "scheduled"
	}

	color := req.Color
	if color == "" {
		color = "#1976d2"
	}

	userID, exists := c.Get("user_id")
	var createdBy *string
	if exists {
		uid := userID.(string)
		createdBy = &uid
	}

	event := &domain.CalendarEvent{
		ID:             uuid.New().String(),
		OrganizationID: defaultOrganizationID,
		CampaignID:     req.CampaignID,
		EventType:      req.EventType,
		Title:          req.Title,
		Description:    req.Description,
		StartTime:      startTime,
		EndTime:        endTime,
		AllDay:         req.AllDay,
		Color:          color,
		Channel:        req.Channel,
		Status:         status,
		Location:       req.Location,
		LocationURL:    req.LocationURL,
		Attendees:      req.Attendees,
		Tags:           req.Tags,
		CreatedBy:      createdBy,
	}

	if req.Reminders != "" {
		event.Reminders = []byte(req.Reminders)
	}
	if req.Recurrence != "" {
		event.Recurrence = []byte(req.Recurrence)
	}

	if err := h.service.CreateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CalendarEventResponse{
		ID:             event.ID,
		OrganizationID: event.OrganizationID,
		CampaignID:     event.CampaignID,
		EventType:      event.EventType,
		Title:          event.Title,
		Description:    event.Description,
		StartTime:      event.StartTime.Format(time.RFC3339),
		EndTime:        formatTimePtr(event.EndTime),
		AllDay:         event.AllDay,
		Color:          event.Color,
		Channel:        event.Channel,
		Status:         event.Status,
		Location:       event.Location,
		LocationURL:    event.LocationURL,
		Attendees:      event.Attendees,
		Tags:           event.Tags,
		CreatedAt:      event.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      event.UpdatedAt.Format(time.RFC3339),
	})
}

// ListEvents handles GET /api/v1/calendar/events
func (h *CalendarHandler) ListEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var startTime, endTime *time.Time
	if s := c.Query("start"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start format, use RFC3339"})
			return
		}
		startTime = &t
	}
	if e := c.Query("end"); e != "" {
		t, err := time.Parse(time.RFC3339, e)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end format, use RFC3339"})
			return
		}
		endTime = &t
	}

	var channel, eventType, status *string
	if ch := c.Query("channel"); ch != "" {
		channel = &ch
	}
	if et := c.Query("event_type"); et != "" {
		eventType = &et
	}
	if st := c.Query("status"); st != "" {
		status = &st
	}

	events, err := h.service.ListEvents(c.Request.Context(), defaultOrganizationID, startTime, endTime, channel, eventType, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.CalendarEventResponse, len(events))
	for i, e := range events {
		response[i] = dto.CalendarEventResponse{
			ID:             e.ID,
			OrganizationID: e.OrganizationID,
			CampaignID:     e.CampaignID,
			EventType:      e.EventType,
			Title:          e.Title,
			Description:    e.Description,
			StartTime:      e.StartTime.Format(time.RFC3339),
			EndTime:        formatTimePtr(e.EndTime),
			AllDay:         e.AllDay,
			Color:          e.Color,
			Channel:        e.Channel,
			Status:         e.Status,
			Location:       e.Location,
			LocationURL:    e.LocationURL,
			Attendees:      e.Attendees,
			Tags:           e.Tags,
			CreatedAt:      e.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      e.UpdatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, gin.H{"events": response, "total": len(events)})
}

// GetEvent handles GET /api/v1/calendar/events/:id
func (h *CalendarHandler) GetEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := h.service.GetEvent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar event not found"})
		return
	}

	c.JSON(http.StatusOK, dto.CalendarEventResponse{
		ID:             event.ID,
		OrganizationID: event.OrganizationID,
		CampaignID:     event.CampaignID,
		EventType:      event.EventType,
		Title:          event.Title,
		Description:    event.Description,
		StartTime:      event.StartTime.Format(time.RFC3339),
		EndTime:        formatTimePtr(event.EndTime),
		AllDay:         event.AllDay,
		Color:          event.Color,
		Channel:        event.Channel,
		Status:         event.Status,
		Location:       event.Location,
		LocationURL:    event.LocationURL,
		Attendees:      event.Attendees,
		Tags:           event.Tags,
		CreatedAt:      event.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      event.UpdatedAt.Format(time.RFC3339),
	})
}

// UpdateEvent handles PUT /api/v1/calendar/events/:id
func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := h.service.GetEvent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar event not found"})
		return
	}

	var req dto.UpdateCalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EventType != "" {
		event.EventType = req.EventType
	}
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format, use RFC3339"})
			return
		}
		event.StartTime = t
	}
	if req.EndTime != nil {
		t, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format, use RFC3339"})
			return
		}
		event.EndTime = &t
	}
	if req.AllDay != nil {
		event.AllDay = *req.AllDay
	}
	if req.Color != "" {
		event.Color = req.Color
	}
	if req.Channel != nil {
		event.Channel = req.Channel
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if req.Location != nil {
		event.Location = req.Location
	}
	if req.LocationURL != nil {
		event.LocationURL = req.LocationURL
	}
	if req.Attendees != nil {
		event.Attendees = req.Attendees
	}
	if req.Tags != nil {
		event.Tags = req.Tags
	}
	if req.Reminders != "" {
		event.Reminders = []byte(req.Reminders)
	}
	if req.Recurrence != "" {
		event.Recurrence = []byte(req.Recurrence)
	}
	if req.CampaignID != nil {
		event.CampaignID = req.CampaignID
	}

	if err := h.service.UpdateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CalendarEventResponse{
		ID:             event.ID,
		OrganizationID: event.OrganizationID,
		CampaignID:     event.CampaignID,
		EventType:      event.EventType,
		Title:          event.Title,
		Description:    event.Description,
		StartTime:      event.StartTime.Format(time.RFC3339),
		EndTime:        formatTimePtr(event.EndTime),
		AllDay:         event.AllDay,
		Color:          event.Color,
		Channel:        event.Channel,
		Status:         event.Status,
		Location:       event.Location,
		LocationURL:    event.LocationURL,
		Attendees:      event.Attendees,
		Tags:           event.Tags,
		CreatedAt:      event.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      event.UpdatedAt.Format(time.RFC3339),
	})
}

// DeleteEvent handles DELETE /api/v1/calendar/events/:id
func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteEvent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calendar event deleted successfully"})
}

// CreateView handles POST /api/v1/calendar/views
func (h *CalendarHandler) CreateView(c *gin.Context) {
	var req dto.CreateCalendarViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	var userIDPtr *string
	if exists {
		uid := userID.(string)
		userIDPtr = &uid
	}

	view := &domain.CalendarView{
		ID:             uuid.New().String(),
		OrganizationID: defaultOrganizationID,
		UserID:         userIDPtr,
		Name:           req.Name,
		ViewType:       req.ViewType,
		Channels:       req.Channels,
		IsDefault:      req.IsDefault,
	}

	if req.Filters != "" {
		view.Filters = []byte(req.Filters)
	}

	if err := h.service.CreateView(c.Request.Context(), view); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CalendarViewResponse{
		ID:             view.ID,
		OrganizationID: view.OrganizationID,
		UserID:         view.UserID,
		Name:           view.Name,
		ViewType:       view.ViewType,
		Channels:       view.Channels,
		IsDefault:      view.IsDefault,
		CreatedAt:      view.CreatedAt.Format(time.RFC3339),
	})
}

// ListViews handles GET /api/v1/calendar/views
func (h *CalendarHandler) ListViews(c *gin.Context) {
	views, err := h.service.ListViews(c.Request.Context(), defaultOrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.CalendarViewResponse, len(views))
	for i, v := range views {
		response[i] = dto.CalendarViewResponse{
			ID:             v.ID,
			OrganizationID: v.OrganizationID,
			UserID:         v.UserID,
			Name:           v.Name,
			ViewType:       v.ViewType,
			Channels:       v.Channels,
			IsDefault:      v.IsDefault,
			CreatedAt:      v.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, gin.H{"views": response, "total": len(views)})
}

// formatTimePtr formats a time pointer to RFC3339 string, returns nil if input is nil
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
