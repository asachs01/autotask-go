package autotask

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New(LogLevelInfo, true)

	if logger == nil {
		t.Fatal("Expected non-nil Logger")
	}

	if logger.level != LogLevelInfo {
		t.Errorf("Expected level to be %d, got %d", LogLevelInfo, logger.level)
	}

	if !logger.debugMode {
		t.Errorf("Expected debugMode to be true, got false")
	}
}

func TestSetLevel(t *testing.T) {
	logger := New(LogLevelInfo, false)
	logger.SetLevel(LogLevelError)

	if logger.level != LogLevelError {
		t.Errorf("Expected level to be %d, got %d", LogLevelError, logger.level)
	}
}

func TestLoggerSetDebugMode(t *testing.T) {
	logger := New(LogLevelInfo, false)
	logger.SetDebugMode(true)

	if !logger.debugMode {
		t.Errorf("Expected debugMode to be true, got false")
	}
}

func TestSetOutput(t *testing.T) {
	logger := New(LogLevelInfo, false)
	buffer := &bytes.Buffer{}
	logger.SetOutput(buffer)

	if logger.output != buffer {
		t.Errorf("Expected output to be set to the buffer")
	}
}

func TestDebug(t *testing.T) {
	t.Run("debug mode enabled", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelDebug, true)
		logger.SetOutput(buffer)

		logger.Debug("test debug message", map[string]interface{}{"key": "value"})

		output := buffer.String()
		if !strings.Contains(output, "DEBUG") || !strings.Contains(output, "test debug message") || !strings.Contains(output, "key") || !strings.Contains(output, "value") {
			t.Errorf("Expected debug message with fields, got: %s", output)
		}
	})

	t.Run("debug mode disabled", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelDebug, false)
		logger.SetOutput(buffer)

		logger.Debug("test debug message", map[string]interface{}{"key": "value"})

		output := buffer.String()
		if output != "" {
			t.Errorf("Expected no output when debug mode is disabled, got: %s", output)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("level allows info", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelInfo, false)
		logger.SetOutput(buffer)

		logger.Info("test info message", map[string]interface{}{"key": "value"})

		output := buffer.String()
		if !strings.Contains(output, "INFO") || !strings.Contains(output, "test info message") || !strings.Contains(output, "key") || !strings.Contains(output, "value") {
			t.Errorf("Expected info message with fields, got: %s", output)
		}
	})

	t.Run("level blocks info", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelWarn, false)
		logger.SetOutput(buffer)

		logger.Info("test info message", map[string]interface{}{"key": "value"})

		output := buffer.String()
		if output != "" {
			t.Errorf("Expected no output when level is higher than info, got: %s", output)
		}
	})
}

func TestWarn(t *testing.T) {
	t.Run("level allows warn", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelWarn, false)
		logger.SetOutput(buffer)

		logger.Warn("test warn message", map[string]interface{}{"key": "value"})

		output := buffer.String()
		if !strings.Contains(output, "WARN") || !strings.Contains(output, "test warn message") || !strings.Contains(output, "key") || !strings.Contains(output, "value") {
			t.Errorf("Expected warn message with fields, got: %s", output)
		}
	})

	t.Run("level blocks warn", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelError, false)
		logger.SetOutput(buffer)

		logger.Warn("test warn message", map[string]interface{}{"key": "value"})

		output := buffer.String()
		if output != "" {
			t.Errorf("Expected no output when level is higher than warn, got: %s", output)
		}
	})
}

func TestError(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := New(LogLevelError, false)
	logger.SetOutput(buffer)

	logger.Error("test error message", map[string]interface{}{"key": "value"})

	output := buffer.String()
	if !strings.Contains(output, "ERROR") || !strings.Contains(output, "test error message") || !strings.Contains(output, "key") || !strings.Contains(output, "value") {
		t.Errorf("Expected error message with fields, got: %s", output)
	}
}

func TestLogRequest(t *testing.T) {
	t.Run("debug mode enabled", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelDebug, true)
		logger.SetOutput(buffer)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		logger.LogRequest("GET", "https://api.example.com", headers)

		output := buffer.String()
		if !strings.Contains(output, "HTTP Request") || !strings.Contains(output, "GET") || !strings.Contains(output, "https://api.example.com") || !strings.Contains(output, "Content-Type") {
			t.Errorf("Expected request log with method, URL and headers, got: %s", output)
		}
	})

	t.Run("debug mode disabled", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelDebug, false)
		logger.SetOutput(buffer)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		logger.LogRequest("GET", "https://api.example.com", headers)

		output := buffer.String()
		if output != "" {
			t.Errorf("Expected no output when debug mode is disabled, got: %s", output)
		}
	})
}

func TestLogResponse(t *testing.T) {
	t.Run("debug mode enabled", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelDebug, true)
		logger.SetOutput(buffer)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		logger.LogResponse(200, headers)

		output := buffer.String()
		if !strings.Contains(output, "HTTP Response") || !strings.Contains(output, "200") || !strings.Contains(output, "Content-Type") {
			t.Errorf("Expected response log with status code and headers, got: %s", output)
		}
	})

	t.Run("debug mode disabled", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelDebug, false)
		logger.SetOutput(buffer)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		logger.LogResponse(200, headers)

		output := buffer.String()
		if output != "" {
			t.Errorf("Expected no output when debug mode is disabled, got: %s", output)
		}
	})
}

func TestLogError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelError, false)
		logger.SetOutput(buffer)

		err := errors.New("test error")
		logger.LogError(err)

		output := buffer.String()
		if !strings.Contains(output, "ERROR") || !strings.Contains(output, "Error") || !strings.Contains(output, "test error") {
			t.Errorf("Expected error log, got: %s", output)
		}
	})

	t.Run("with nil error", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		logger := New(LogLevelError, false)
		logger.SetOutput(buffer)

		logger.LogError(nil)

		output := buffer.String()
		if output != "" {
			t.Errorf("Expected no output for nil error, got: %s", output)
		}
	})
}

func TestGetLevelString(t *testing.T) {
	logger := New(LogLevelInfo, false)

	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevel(99), "UNKNOWN"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			result := logger.getLevelString(tc.level)
			if result != tc.expected {
				t.Errorf("Expected level string %s, got %s", tc.expected, result)
			}
		})
	}
}
