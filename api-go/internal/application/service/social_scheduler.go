package service

import (
	"sync"
	"time"
)

// SocialScheduler handles scheduled social media posts
type SocialScheduler struct {
	timers map[string]*time.Timer
	mu     sync.RWMutex
}

// NewSocialScheduler creates a new post scheduler
func NewSocialScheduler() *SocialScheduler {
	s := &SocialScheduler{
		timers: make(map[string]*time.Timer),
	}
	return s
}

// SchedulePost schedules a post for future publishing
func (s *SocialScheduler) SchedulePost(postID string, scheduledAt time.Time, publishFunc func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cancel existing timer if any
	if existing, ok := s.timers[postID]; ok {
		existing.Stop()
	}

	delay := time.Until(scheduledAt)
	if delay <= 0 {
		// Schedule is in the past, publish immediately
		go publishFunc()
		return
	}

	timer := time.AfterFunc(delay, func() {
		s.mu.Lock()
		delete(s.timers, postID)
		s.mu.Unlock()
		publishFunc()
	})

	s.timers[postID] = timer
}

// CancelPost cancels a scheduled post
func (s *SocialScheduler) CancelPost(postID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if timer, ok := s.timers[postID]; ok {
		timer.Stop()
		delete(s.timers, postID)
	}
}
