package autotask

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	requestsPerMinute := 60
	limiter := NewRateLimiter(requestsPerMinute)

	if limiter == nil {
		t.Fatal("Expected non-nil RateLimiter")
	}

	if limiter.requestsPerMinute != requestsPerMinute {
		t.Errorf("Expected requestsPerMinute to be %d, got %d", requestsPerMinute, limiter.requestsPerMinute)
	}

	// Initial lastRequest should be zero time
	if !limiter.lastRequest.IsZero() {
		t.Errorf("Expected lastRequest to be zero time, got %v", limiter.lastRequest)
	}
}

func TestRateLimiterWait(t *testing.T) {
	t.Run("first request should not wait", func(t *testing.T) {
		limiter := NewRateLimiter(60) // 1 request per second

		start := time.Now()
		waitTime := limiter.Wait()
		elapsed := time.Since(start)

		if waitTime > 0 {
			t.Errorf("Expected no wait time for first request, got %v", waitTime)
		}

		if elapsed > 10*time.Millisecond {
			t.Errorf("First request took too long: %v", elapsed)
		}

		// lastRequest should be updated
		if limiter.lastRequest.IsZero() {
			t.Error("Expected lastRequest to be updated")
		}
	})

	t.Run("subsequent requests should respect rate limit", func(t *testing.T) {
		requestsPerMinute := 60 // 1 request per second
		limiter := NewRateLimiter(requestsPerMinute)

		// First request (no wait)
		limiter.Wait()

		// Second request immediately after should wait
		start := time.Now()
		waitTime := limiter.Wait()
		elapsed := time.Since(start)

		expectedWaitTime := time.Minute / time.Duration(requestsPerMinute)

		// Allow for small timing variations
		if waitTime < time.Duration(float64(expectedWaitTime)*0.8) || waitTime > time.Duration(float64(expectedWaitTime)*1.2) {
			t.Errorf("Expected wait time around %v, got %v", expectedWaitTime, waitTime)
		}

		// Actual elapsed time should be close to expected wait time
		if elapsed < time.Duration(float64(expectedWaitTime)*0.8) || elapsed > time.Duration(float64(expectedWaitTime)*1.2) {
			t.Errorf("Expected elapsed time around %v, got %v", expectedWaitTime, elapsed)
		}
	})

	t.Run("high rate limit should not cause wait", func(t *testing.T) {
		limiter := NewRateLimiter(6000) // 100 requests per second

		// First request
		limiter.Wait()

		// Second request immediately after
		start := time.Now()
		waitTime := limiter.Wait()
		elapsed := time.Since(start)

		// With such a high rate limit, wait time should be very small
		// But it might still be non-zero due to implementation details
		// So we'll just check that it's reasonably small
		if waitTime > 20*time.Millisecond {
			t.Errorf("Expected small wait time, got %v", waitTime)
		}

		if elapsed > 25*time.Millisecond {
			t.Errorf("Request took too long: %v", elapsed)
		}
	})

	t.Run("wait after delay", func(t *testing.T) {
		requestsPerMinute := 60 // 1 request per second
		limiter := NewRateLimiter(requestsPerMinute)

		// First request
		limiter.Wait()

		// Wait half the rate limit interval
		halfInterval := time.Minute / time.Duration(requestsPerMinute) / 2
		time.Sleep(halfInterval)

		// Second request should wait for remaining half
		start := time.Now()
		waitTime := limiter.Wait()
		elapsed := time.Since(start)

		// Expected wait time is approximately half the interval
		if waitTime < time.Duration(float64(halfInterval)*0.5) || waitTime > time.Duration(float64(halfInterval)*1.5) {
			t.Errorf("Expected wait time around %v, got %v", halfInterval, waitTime)
		}

		// Actual elapsed time should be close to expected wait time
		if elapsed < time.Duration(float64(halfInterval)*0.5) || elapsed > time.Duration(float64(halfInterval)*1.5) {
			t.Errorf("Expected elapsed time around %v, got %v", halfInterval, elapsed)
		}
	})
}
