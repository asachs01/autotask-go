package telemetry

import (
	"context"
	"testing"
	"time"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Test with telemetry disabled
	telemetry := New(false)
	if telemetry == nil {
		t.Fatal("telemetry should not be nil")
	}
	if telemetry.enabled {
		t.Error("telemetry should be disabled")
	}
	if telemetry.tracer != nil {
		t.Error("tracer should be nil when disabled")
	}
	if telemetry.meter != nil {
		t.Error("meter should be nil when disabled")
	}

	// Test with telemetry enabled
	telemetry = New(true)
	if !telemetry.enabled {
		t.Error("telemetry should be enabled")
	}
	if telemetry.tracer == nil {
		t.Error("tracer should not be nil when enabled")
	}
	if telemetry.meter == nil {
		t.Error("meter should not be nil when enabled")
	}
}

// TestDisabledTelemetry tests that disabled telemetry doesn't cause errors
func TestDisabledTelemetry(t *testing.T) {
	// Create a telemetry instance with disabled flag
	telemetry := &Telemetry{
		enabled: false,
	}

	// Test StartRequestSpan
	ctx := context.Background()
	ctx, span := telemetry.StartRequestSpan(ctx, "GET", "https://example.com/api")
	if span != nil {
		t.Error("expected nil span when telemetry is disabled")
	}

	// Test StartBatchSpan
	ctx, span = telemetry.StartBatchSpan(ctx, "create", 10)
	if span != nil {
		t.Error("expected nil span when telemetry is disabled")
	}

	// Test RecordRateLimit (should not panic)
	telemetry.RecordRateLimit(ctx, time.Second*5)
}
