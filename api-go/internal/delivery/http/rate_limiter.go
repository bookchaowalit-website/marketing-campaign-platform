package http

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	limit   int
	window  time.Duration
}

type client struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, exists := rl.clients[key]

	if !exists {
		rl.clients[key] = &client{tokens: rl.limit - 1, lastRefill: now}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(c.lastRefill)
	tokensToAdd := int(elapsed / rl.window) * rl.limit
	if tokensToAdd > 0 {
		c.tokens = min(c.tokens+tokensToAdd, rl.limit)
		c.lastRefill = now
	}

	if c.tokens <= 0 {
		return false
	}

	c.tokens--
	return true
}

// Cleanup removes expired entries
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, c := range rl.clients {
		if now.Sub(c.lastRefill) > rl.window*2 {
			delete(rl.clients, key)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if !limiter.Allow(key) {
			c.JSON(429, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}
		c.Next()
	}
}
