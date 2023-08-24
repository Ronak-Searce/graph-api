package tracing

import (
	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/otel/bridge/opencensus"
)

/*
Google libraries use OpenCensus library for tracing. While it is pretty good choice, we use OpenTelemetry stack which is successor to OpenCensus and OpenTracing.
And OpenTelemetry provides bridge for OpenCensus-based code with some small limitations.
*/

// RegisterGoogleTracing creates a bridge between google tracing (opencensus) and our tracing (opentelemetry).
// This function must be called after exporter is set up (SetupExporter).
func RegisterGoogleTracing() {
	tracer := opencensus.NewTracer(DefaultTracer)
	octrace.DefaultTracer = tracer
}
