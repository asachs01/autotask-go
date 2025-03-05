package autotask

import (
	"sync"
	"time"
)

// RateLimiter handles API rate limiting
type RateLimiter struct {
	requestsPerMinute int
	lastRequest       time.Time
	mu                sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		requestsPerMinute: requestsPerMinute,
	}
}

// Wait waits if necessary to respect rate limits
func (r *RateLimiter) Wait() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRequest)
	waitTime := time.Minute/time.Duration(r.requestsPerMinute) - elapsed

	if waitTime > 0 {
		time.Sleep(waitTime)
	}

	r.lastRequest = time.Now()
	return waitTime
}
