package autotask

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.InitialInterval != 100*time.Millisecond {
		t.Errorf("Expected InitialInterval to be 100ms, got %v", config.InitialInterval)
	}

	if config.MaxInterval != 2*time.Second {
		t.Errorf("Expected MaxInterval to be 2s, got %v", config.MaxInterval)
	}

	if config.Multiplier != 2.0 {
		t.Errorf("Expected Multiplier to be 2.0, got %f", config.Multiplier)
	}

	if config.Jitter != 0.1 {
		t.Errorf("Expected Jitter to be 0.1, got %f", config.Jitter)
	}
}

func TestRetryableError(t *testing.T) {
	originalErr := errors.New("original error")
	resp := &http.Response{StatusCode: 500}
	retryableErr := &RetryableError{
		Err:      originalErr,
		Response: resp,
	}

	// Test Error() method
	if retryableErr.Error() != "retryable error: original error" {
		t.Errorf("Unexpected error message: %s", retryableErr.Error())
	}

	// Test Unwrap() method
	if errors.Unwrap(retryableErr) != originalErr {
		t.Errorf("Unwrap() did not return the original error")
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-retryable error",
			err:      errors.New("regular error"),
			expected: false,
		},
		{
			name:     "retryable error type",
			err:      &RetryableError{Err: errors.New("test")},
			expected: true,
		},
		{
			name: "429 too many requests",
			err: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusTooManyRequests},
			},
			expected: true,
		},
		{
			name: "500 internal server error",
			err: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusInternalServerError},
			},
			expected: true,
		},
		{
			name: "502 bad gateway",
			err: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusBadGateway},
			},
			expected: true,
		},
		{
			name: "503 service unavailable",
			err: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusServiceUnavailable},
			},
			expected: true,
		},
		{
			name: "504 gateway timeout",
			err: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusGatewayTimeout},
			},
			expected: true,
		},
		{
			name: "400 bad request - not retryable",
			err: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusBadRequest},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsRetryable(tc.err)
			if result != tc.expected {
				t.Errorf("IsRetryable(%v) = %v, expected %v", tc.err, result, tc.expected)
			}
		})
	}
}

func TestRetryWithBackoff(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:      3,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0.1,
	}

	t.Run("successful operation", func(t *testing.T) {
		attempts := 0
		err := RetryWithBackoff(ctx, config, func() error {
			attempts++
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if attempts != 1 {
			t.Errorf("Expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("non-retryable error", func(t *testing.T) {
		attempts := 0
		expectedErr := errors.New("non-retryable error")

		err := RetryWithBackoff(ctx, config, func() error {
			attempts++
			return expectedErr
		})

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}

		if attempts != 1 {
			t.Errorf("Expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("retryable error with eventual success", func(t *testing.T) {
		attempts := 0

		err := RetryWithBackoff(ctx, config, func() error {
			attempts++
			if attempts < 3 {
				return &RetryableError{Err: errors.New("temporary error")}
			}
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("retryable error with max retries exceeded", func(t *testing.T) {
		attempts := 0
		retryableErr := &RetryableError{Err: errors.New("persistent error")}

		err := RetryWithBackoff(ctx, config, func() error {
			attempts++
			return retryableErr
		})

		if err == nil {
			t.Error("Expected an error, got nil")
		}

		if attempts != config.MaxRetries {
			t.Errorf("Expected %d attempts, got %d", config.MaxRetries, attempts)
		}

		if !errors.Is(errors.Unwrap(err), retryableErr.Err) {
			t.Errorf("Expected error to wrap %v, got %v", retryableErr.Err, err)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		attempts := 0
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(20 * time.Millisecond)
			cancel()
		}()

		err := RetryWithBackoff(ctx, config, func() error {
			attempts++
			return &RetryableError{Err: errors.New("temporary error")}
		})

		if err == nil {
			t.Error("Expected an error, got nil")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
	})
}

func TestWithRetry(t *testing.T) {
	ctx := context.Background()

	t.Run("with nil config", func(t *testing.T) {
		attempts := 0

		err := WithRetry(ctx, nil, func() error {
			attempts++
			if attempts < 2 {
				return &RetryableError{Err: errors.New("temporary error")}
			}
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if attempts != 2 {
			t.Errorf("Expected 2 attempts, got %d", attempts)
		}
	})

	t.Run("with custom config", func(t *testing.T) {
		attempts := 0
		config := &RetryConfig{
			MaxRetries:      1,
			InitialInterval: 10 * time.Millisecond,
			MaxInterval:     100 * time.Millisecond,
			Multiplier:      2.0,
			Jitter:          0.1,
		}

		err := WithRetry(ctx, config, func() error {
			attempts++
			return &RetryableError{Err: errors.New("temporary error")}
		})

		if err == nil {
			t.Error("Expected an error, got nil")
		}

		if attempts != config.MaxRetries {
			t.Errorf("Expected %d attempts, got %d", config.MaxRetries, attempts)
		}
	})
}
