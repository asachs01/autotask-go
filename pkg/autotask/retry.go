package autotask

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     2 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}
}

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err      error
	Response *http.Response
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error: %v", e.Err)
}

// Unwrap returns the underlying error
func (e *RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable checks if an error should be retried
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a retryable error type
	if _, ok := err.(*RetryableError); ok {
		return true
	}

	// Check HTTP response status codes
	if resp, ok := err.(*ErrorResponse); ok {
		switch resp.Response.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
	}

	return false
}

// RetryWithBackoff implements exponential backoff with jitter
func RetryWithBackoff(ctx context.Context, config *RetryConfig, operation func() error) error {
	var lastErr error
	interval := config.InitialInterval

	for attempt := 1; attempt <= config.MaxRetries; attempt++ {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		// Execute operation
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryable(err) {
			return err
		}

		// Calculate next interval with exponential backoff
		interval = time.Duration(float64(interval) * config.Multiplier)
		if interval > config.MaxInterval {
			interval = config.MaxInterval
		}

		// Add jitter
		if config.Jitter > 0 {
			jitter := time.Duration(float64(interval) * config.Jitter)
			interval += time.Duration(time.Now().UnixNano() % int64(jitter))
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(interval):
			continue
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// WithRetry wraps an operation with retry logic
func WithRetry(ctx context.Context, config *RetryConfig, operation func() error) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	return RetryWithBackoff(ctx, config, operation)
}
