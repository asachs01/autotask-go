package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Telemetry handles OpenTelemetry integration
type Telemetry struct {
	tracer  trace.Tracer
	meter   metric.Meter
	enabled bool
}

// New creates a new telemetry instance
func New(enabled bool) *Telemetry {
	if !enabled {
		return &Telemetry{enabled: false}
	}

	tracer := otel.Tracer("autotask-go")
	meter := otel.Meter("autotask-go")

	return &Telemetry{
		tracer:  tracer,
		meter:   meter,
		enabled: true,
	}
}

// StartRequestSpan starts a new span for an API request
func (t *Telemetry) StartRequestSpan(ctx context.Context, method, url string) (context.Context, trace.Span) {
	if !t.enabled {
		return ctx, nil
	}

	ctx, span := t.tracer.Start(ctx, fmt.Sprintf("autotask.%s", method))
	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("http.url", url),
	)

	return ctx, span
}

// EndRequestSpan ends a span and records metrics
func (t *Telemetry) EndRequestSpan(ctx context.Context, span trace.Span, method string, statusCode int, err error) {
	if !t.enabled {
		return
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else if statusCode >= 200 && statusCode < 300 {
		span.SetStatus(codes.Ok, "success")
	} else {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
	}

	span.SetAttributes(attribute.Int("http.status_code", statusCode))
	span.End()

	// Record metrics
	t.recordMetrics(ctx, method, statusCode, err)
}

// recordMetrics records request metrics
func (t *Telemetry) recordMetrics(ctx context.Context, method string, statusCode int, err error) {
	// Record request count
	requestCounter, _ := t.meter.Int64Counter(
		"autotask.requests",
		metric.WithDescription("Number of API requests"),
	)
	requestCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("method", method),
		attribute.Int("status_code", statusCode),
		attribute.Bool("error", err != nil),
	))

	// Record request duration
	durationHistogram, _ := t.meter.Float64Histogram(
		"autotask.request_duration",
		metric.WithDescription("Request duration in milliseconds"),
	)
	if span := trace.SpanFromContext(ctx); span != nil {
		// Use current time since we can't reliably get span start time
		duration := time.Since(time.Now()).Milliseconds()
		durationHistogram.Record(ctx, float64(duration), metric.WithAttributes(
			attribute.String("method", method),
			attribute.Int("status_code", statusCode),
		))
	}
}

// StartBatchSpan starts a new span for a batch operation
func (t *Telemetry) StartBatchSpan(ctx context.Context, operation string, count int) (context.Context, trace.Span) {
	if !t.enabled {
		return ctx, nil
	}

	ctx, span := t.tracer.Start(ctx, fmt.Sprintf("autotask.batch.%s", operation))
	span.SetAttributes(
		attribute.String("operation", operation),
		attribute.Int("count", count),
	)

	return ctx, span
}

// RecordRateLimit records rate limiting metrics
func (t *Telemetry) RecordRateLimit(ctx context.Context, waitTime time.Duration) {
	if !t.enabled {
		return
	}

	rateLimitCounter, _ := t.meter.Int64Counter(
		"autotask.rate_limit",
		metric.WithDescription("Number of rate limit events"),
	)
	rateLimitCounter.Add(ctx, 1)

	rateLimitWaitHistogram, _ := t.meter.Float64Histogram(
		"autotask.rate_limit_wait",
		metric.WithDescription("Rate limit wait time in milliseconds"),
	)
	rateLimitWaitHistogram.Record(ctx, float64(waitTime.Milliseconds()))
}
