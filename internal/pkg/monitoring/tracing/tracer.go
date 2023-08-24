package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	// DefaultTracer is configured by SetupExporter tracer instance.
	DefaultTracer Tracer
)

// Tracer is a wrapper for OpenTelemetry Tracer.
type Tracer interface {
	trace.Tracer
}

// GetTracer returns tracer instance.
func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
