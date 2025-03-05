package logger

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/asachs01/autotask-go/pkg/autotask"
)

// LogLevel represents the level of logging
type LogLevel int

const (
	// LogLevelDebug represents debug level logging
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents info level logging
	LogLevelInfo
	// LogLevelWarn represents warning level logging
	LogLevelWarn
	// LogLevelError represents error level logging
	LogLevelError
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger handles logging for the Autotask client
type Logger struct {
	level     autotask.LogLevel
	logger    *log.Logger
	debugMode bool
	output    *os.File
}

// New creates a new logger instance
func New(level autotask.LogLevel, debugMode bool) *Logger {
	return &Logger{
		level:     level,
		logger:    log.New(os.Stdout, "", 0), // Remove default flags as we handle timestamp ourselves
		debugMode: debugMode,
		output:    os.Stdout,
	}
}

// log writes a structured log entry
func (l *Logger) log(level string, msg string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		// Fallback to basic logging if JSON marshaling fails
		l.logger.Printf("[%s] %s %v", level, msg, fields)
		return
	}

	l.logger.Println(string(data))
}

// Debug logs a debug message with optional fields
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if l.level <= autotask.LogLevelDebug && l.debugMode {
		l.log("DEBUG", msg, fields)
	}
}

// Info logs an info message with optional fields
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	if l.level <= autotask.LogLevelInfo {
		l.log("INFO", msg, fields)
	}
}

// Warn logs a warning message with optional fields
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	if l.level <= autotask.LogLevelWarn {
		l.log("WARN", msg, fields)
	}
}

// Error logs an error message with optional fields
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	if l.level <= autotask.LogLevelError {
		l.log("ERROR", msg, fields)
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level autotask.LogLevel) {
	l.level = level
}

// SetDebugMode sets the debug mode
func (l *Logger) SetDebugMode(debug bool) {
	l.debugMode = debug
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(output *os.File) {
	l.output = output
	l.logger.SetOutput(output)
}

// LogRequest logs an API request with structured fields
func (l *Logger) LogRequest(method, url string, headers map[string]string) {
	if l.debugMode {
		fields := map[string]interface{}{
			"method": method,
			"url":    url,
		}

		// Mask sensitive information in headers
		safeHeaders := make(map[string]string)
		for key, value := range headers {
			if strings.ToLower(key) == "authorization" || strings.ToLower(key) == "secret" {
				safeHeaders[key] = value[:10] + "..."
			} else {
				safeHeaders[key] = value
			}
		}
		fields["headers"] = safeHeaders

		l.Debug("API Request", fields)
	}
}

// LogResponse logs an API response with structured fields
func (l *Logger) LogResponse(statusCode int, headers map[string]string) {
	if l.debugMode {
		fields := map[string]interface{}{
			"status_code": statusCode,
			"headers":     headers,
		}
		l.Debug("API Response", fields)
	}
}

// LogError logs an API error with structured fields
func (l *Logger) LogError(err error) {
	fields := map[string]interface{}{
		"error": err.Error(),
	}
	if apiErr, ok := err.(*autotask.ErrorResponse); ok {
		fields["status_code"] = apiErr.Response.StatusCode
		fields["errors"] = apiErr.Errors
	}
	l.Error("API Error", fields)
}

// LogRateLimit logs rate limiting information
func (l *Logger) LogRateLimit(waitTime time.Duration) {
	if l.debugMode {
		fields := map[string]interface{}{
			"wait_time_ms": waitTime.Milliseconds(),
		}
		l.Debug("Rate Limit Applied", fields)
	}
}

// LogZoneInfo logs zone information
func (l *Logger) LogZoneInfo(zoneInfo *autotask.ZoneInfo) {
	if l.debugMode {
		fields := map[string]interface{}{
			"zone_name": zoneInfo.ZoneName,
			"url":       zoneInfo.URL,
			"web_url":   zoneInfo.WebURL,
			"ci":        zoneInfo.CI,
		}
		l.Debug("Zone Information", fields)
	}
}
