package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/asachs01/autotask-go/pkg/autotask"
)

// captureOutput captures the output of the logger
func captureOutput(logger *Logger, fn func()) string {
	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	originalOutput := logger.output
	logger.output = w
	logger.logger.SetOutput(w)

	// Execute the function
	fn()

	// Restore original output
	logger.output = originalOutput
	logger.logger.SetOutput(originalOutput)
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestNew(t *testing.T) {
	// Create a new logger
	logger := New(autotask.LogLevelInfo, false)

	// Verify logger is not nil
	if logger == nil {
		t.Fatal("logger should not be nil")
	}

	// Verify logger properties
	if logger.level != autotask.LogLevelInfo {
		t.Errorf("expected level to be LogLevelInfo, got %v", logger.level)
	}
	if logger.debugMode != false {
		t.Errorf("expected debugMode to be false, got %v", logger.debugMode)
	}
	if logger.output != os.Stdout {
		t.Error("expected output to be os.Stdout")
	}
	if logger.logger == nil {
		t.Error("logger.logger should not be nil")
	}
}

func TestLogLevels(t *testing.T) {
	// Test Debug level
	t.Run("Debug", func(t *testing.T) {
		logger := New(autotask.LogLevelDebug, true)
		output := captureOutput(logger, func() {
			logger.Debug("debug message", map[string]interface{}{"key": "value"})
		})

		// Verify output contains JSON with the expected fields
		if !strings.Contains(output, "debug message") {
			t.Errorf("expected output to contain 'debug message', got %s", output)
		}
		if !strings.Contains(output, "DEBUG") {
			t.Errorf("expected output to contain 'DEBUG', got %s", output)
		}
		if !strings.Contains(output, "key") || !strings.Contains(output, "value") {
			t.Errorf("expected output to contain field 'key: value', got %s", output)
		}

		// Verify debug message is not logged at higher levels
		logger = New(autotask.LogLevelInfo, false)
		output = captureOutput(logger, func() {
			logger.Debug("debug message", map[string]interface{}{"key": "value"})
		})
		if output != "" {
			t.Errorf("expected no output at Info level, got %s", output)
		}
	})

	// Test Info level
	t.Run("Info", func(t *testing.T) {
		logger := New(autotask.LogLevelInfo, false)
		output := captureOutput(logger, func() {
			logger.Info("info message", map[string]interface{}{"key": "value"})
		})

		// Verify output
		if !strings.Contains(output, "info message") {
			t.Errorf("expected output to contain 'info message', got %s", output)
		}
		if !strings.Contains(output, "INFO") {
			t.Errorf("expected output to contain 'INFO', got %s", output)
		}

		// Verify info message is not logged at higher levels
		logger = New(autotask.LogLevelWarn, false)
		output = captureOutput(logger, func() {
			logger.Info("info message", map[string]interface{}{"key": "value"})
		})
		if output != "" {
			t.Errorf("expected no output at Warn level, got %s", output)
		}
	})

	// Test Warn level
	t.Run("Warn", func(t *testing.T) {
		logger := New(autotask.LogLevelWarn, false)
		output := captureOutput(logger, func() {
			logger.Warn("warn message", map[string]interface{}{"key": "value"})
		})

		// Verify output
		if !strings.Contains(output, "warn message") {
			t.Errorf("expected output to contain 'warn message', got %s", output)
		}
		if !strings.Contains(output, "WARN") {
			t.Errorf("expected output to contain 'WARN', got %s", output)
		}

		// Verify warn message is not logged at higher levels
		logger = New(autotask.LogLevelError, false)
		output = captureOutput(logger, func() {
			logger.Warn("warn message", map[string]interface{}{"key": "value"})
		})
		if output != "" {
			t.Errorf("expected no output at Error level, got %s", output)
		}
	})

	// Test Error level
	t.Run("Error", func(t *testing.T) {
		logger := New(autotask.LogLevelError, false)
		output := captureOutput(logger, func() {
			logger.Error("error message", map[string]interface{}{"key": "value"})
		})

		// Verify output
		if !strings.Contains(output, "error message") {
			t.Errorf("expected output to contain 'error message', got %s", output)
		}
		if !strings.Contains(output, "ERROR") {
			t.Errorf("expected output to contain 'ERROR', got %s", output)
		}
	})
}

func TestSetters(t *testing.T) {
	// Test SetLevel
	t.Run("SetLevel", func(t *testing.T) {
		logger := New(autotask.LogLevelInfo, false)
		logger.SetLevel(autotask.LogLevelError)

		// Verify level is set
		if logger.level != autotask.LogLevelError {
			t.Errorf("expected level to be LogLevelError, got %v", logger.level)
		}

		// Verify info message is not logged at error level
		output := captureOutput(logger, func() {
			logger.Info("info message", nil)
		})
		if output != "" {
			t.Errorf("expected no output at Error level, got %s", output)
		}

		// Verify error message is logged
		output = captureOutput(logger, func() {
			logger.Error("error message", nil)
		})
		if !strings.Contains(output, "error message") {
			t.Errorf("expected output to contain 'error message', got %s", output)
		}
	})

	// Test SetDebugMode
	t.Run("SetDebugMode", func(t *testing.T) {
		logger := New(autotask.LogLevelInfo, false)
		logger.SetDebugMode(true)

		// Verify debug mode is set
		if !logger.debugMode {
			t.Error("expected debugMode to be true")
		}
	})

	// Test SetOutput
	t.Run("SetOutput", func(t *testing.T) {
		logger := New(autotask.LogLevelInfo, false)

		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "logger-test")
		if err != nil {
			t.Fatalf("error creating temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		// Set output to the temp file
		logger.SetOutput(tmpfile)

		// Verify output is set
		if logger.output != tmpfile {
			t.Error("expected output to be the temp file")
		}

		// Log a message
		logger.Info("test message", nil)

		// Close the file
		tmpfile.Close()

		// Read the file
		content, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			t.Fatalf("error reading temp file: %v", err)
		}

		// Verify message was written to the file
		if !strings.Contains(string(content), "test message") {
			t.Errorf("expected file to contain 'test message', got %s", string(content))
		}
	})
}

func TestSpecializedLogging(t *testing.T) {
	// Test LogRequest
	t.Run("LogRequest", func(t *testing.T) {
		logger := New(autotask.LogLevelDebug, true)
		output := captureOutput(logger, func() {
			headers := map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "test-agent",
			}
			logger.LogRequest("GET", "https://example.com/api", headers)
		})

		// Verify output
		if !strings.Contains(output, "GET") {
			t.Errorf("expected output to contain 'GET', got %s", output)
		}
		if !strings.Contains(output, "https://example.com/api") {
			t.Errorf("expected output to contain URL, got %s", output)
		}
		if !strings.Contains(output, "Content-Type") {
			t.Errorf("expected output to contain header name, got %s", output)
		}
		if !strings.Contains(output, "application/json") {
			t.Errorf("expected output to contain header value, got %s", output)
		}
	})

	// Test LogResponse
	t.Run("LogResponse", func(t *testing.T) {
		logger := New(autotask.LogLevelDebug, true)
		output := captureOutput(logger, func() {
			headers := map[string]string{
				"Content-Type": "application/json",
			}
			logger.LogResponse(200, headers)
		})

		// Verify output
		if !strings.Contains(output, "200") {
			t.Errorf("expected output to contain status code, got %s", output)
		}
		if !strings.Contains(output, "Content-Type") {
			t.Errorf("expected output to contain header name, got %s", output)
		}
	})

	// Test LogError
	t.Run("LogError", func(t *testing.T) {
		logger := New(autotask.LogLevelDebug, true)
		output := captureOutput(logger, func() {
			err := &testError{message: "test error"}
			logger.LogError(err)
		})

		// Verify output
		if !strings.Contains(output, "test error") {
			t.Errorf("expected output to contain error message, got %s", output)
		}
	})

	// Test LogRateLimit
	t.Run("LogRateLimit", func(t *testing.T) {
		logger := New(autotask.LogLevelDebug, true)
		output := captureOutput(logger, func() {
			logger.LogRateLimit(time.Second * 5)
		})

		// Verify output contains rate limit information
		if !strings.Contains(output, "Rate Limit Applied") {
			t.Errorf("expected output to contain 'Rate Limit Applied', got %s", output)
		}
		if !strings.Contains(output, "wait_time_ms") {
			t.Errorf("expected output to contain wait_time_ms field, got %s", output)
		}
		if !strings.Contains(output, "5000") {
			t.Errorf("expected output to contain wait time value 5000, got %s", output)
		}
	})

	// Test LogZoneInfo
	t.Run("LogZoneInfo", func(t *testing.T) {
		logger := New(autotask.LogLevelDebug, true)
		output := captureOutput(logger, func() {
			zoneInfo := &autotask.ZoneInfo{
				ZoneName: "Test Zone",
				URL:      "https://example.com",
				WebURL:   "https://web.example.com",
				CI:       123,
			}
			logger.LogZoneInfo(zoneInfo)
		})

		// Verify output
		if !strings.Contains(output, "Test Zone") {
			t.Errorf("expected output to contain zone name, got %s", output)
		}
		if !strings.Contains(output, "https://example.com") {
			t.Errorf("expected output to contain URL, got %s", output)
		}
	})
}

// testError is a simple error implementation for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
