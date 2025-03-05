package autotask

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger handles logging for the Autotask client
type Logger struct {
	level      LogLevel
	debugMode  bool
	output     io.Writer
	timeFormat string
}

// New creates a new logger
func New(level LogLevel, debugMode bool) *Logger {
	return &Logger{
		level:      level,
		debugMode:  debugMode,
		output:     os.Stdout,
		timeFormat: time.RFC3339,
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetDebugMode enables or disables debug logging
func (l *Logger) SetDebugMode(debug bool) {
	l.debugMode = debug
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if !l.debugMode {
		return
	}
	l.log(LogLevelDebug, msg, fields)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	if l.level > LogLevelInfo {
		return
	}
	l.log(LogLevelInfo, msg, fields)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	if l.level > LogLevelWarn {
		return
	}
	l.log(LogLevelWarn, msg, fields)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	if l.level > LogLevelError {
		return
	}
	l.log(LogLevelError, msg, fields)
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(method, url string, headers map[string]string) {
	if !l.debugMode {
		return
	}
	fields := map[string]interface{}{
		"method":  method,
		"url":     url,
		"headers": headers,
	}
	l.log(LogLevelDebug, "HTTP Request", fields)
}

// LogResponse logs an HTTP response
func (l *Logger) LogResponse(statusCode int, headers map[string]string) {
	if !l.debugMode {
		return
	}
	fields := map[string]interface{}{
		"status_code": statusCode,
		"headers":     headers,
	}
	l.log(LogLevelDebug, "HTTP Response", fields)
}

// LogError logs an error
func (l *Logger) LogError(err error) {
	if err == nil {
		return
	}
	fields := map[string]interface{}{
		"error": err.Error(),
	}
	l.log(LogLevelError, "Error", fields)
}

func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	timestamp := time.Now().Format(l.timeFormat)
	levelStr := l.getLevelString(level)

	// Build the log message
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, msg)
	if len(fields) > 0 {
		logMsg += fmt.Sprintf(" %v", fields)
	}

	// Write to output
	fmt.Fprintln(l.output, logMsg)
}

func (l *Logger) getLevelString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
