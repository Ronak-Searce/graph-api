package tracing

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"runtime"

	"github.com/Ronak-Searce/graph-api/internal/pkg/config"
	"github.com/Ronak-Searce/graph-api/internal/pkg/logger"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// Start is a wrapper to start a new span. SetupExporter must be called before using this function.
func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// If for some reason exporter was not registered
	if DefaultTracer == nil {
		DefaultTracer = GetTracer(config.ServiceName())
	}

	return DefaultTracer.Start(ctx, spanName, opts...)
}

// RecordError is a wrapper to add an error to current span in context.
// It also sets span to error status if updateStatus is true.
func RecordError(ctx context.Context, err error, updateStatus bool, options ...trace.EventOption) {
	if err == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	options = append(options, trace.WithAttributes(
		semconv.ExceptionMessageKey.String(err.Error()),
	))
	span.AddEvent("error", options...)
	if updateStatus {
		span.SetStatus(codes.Error, err.Error())
	}
}

// WithCallMetadata adds caller file and line number to span or event attributes.
func WithCallMetadata() trace.SpanStartEventOption {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return trace.WithAttributes()
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return trace.WithAttributes()
	}

	return trace.WithAttributes(
		semconv.CodeFilepathKey.String(file),
		semconv.CodeFunctionKey.String(f.Name()),
		semconv.CodeLineNumberKey.Int(line),
	)
}

// MarshalJSONFromContext gets current span from context and marshals it to JSON.
// It is useful to pass tracing through Pub/Sub. It also drops TraceState because Pub/Sub
// has limitation of 1024 bytes for attribute value.
func MarshalJSONFromContext(ctx context.Context) ([]byte, error) {
	data, err := trace.SpanContextFromContext(ctx).
		WithTraceState(trace.TraceState{}).
		MarshalJSON()
	if err != nil {
		logger.Debugf(ctx, "failed to marshal span to json: %v", err)
		return nil, err
	}

	return data, nil
}

type spanContextConfig struct {
	TraceID    string
	SpanID     string
	TraceFlags string
	Remote     bool
}

// UnmarshalJSONToContext unmarshals span from JSON and adds it to context.
func UnmarshalJSONToContext(ctx context.Context, data []byte) context.Context {
	cfg := spanContextConfig{}
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		logger.Debugf(ctx, "failed to unmarshal span from json: %v", err)
		return ctx
	}
	traceID, err := trace.TraceIDFromHex(cfg.TraceID)
	if err != nil {
		logger.Debugf(ctx, "failed to decode traceid %v: %v", cfg.TraceID, err)
		return ctx
	}
	spanID, err := trace.SpanIDFromHex(cfg.SpanID)
	if err != nil {
		logger.Debugf(ctx, "failed to decode spanid %v: %v", cfg.SpanID, err)
		return ctx
	}
	traceFlags, err := hex.DecodeString(cfg.TraceFlags)
	if err != nil {
		logger.Debugf(ctx, "failed to decode hex traceflags %v: %v", cfg.TraceFlags, err)
		return ctx
	}

	spanCtxConfig := trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.TraceFlags(traceFlags[0]),
		Remote:     cfg.Remote,
	}

	spanCtx := trace.NewSpanContext(spanCtxConfig)
	return trace.ContextWithSpanContext(ctx, spanCtx)
}

// CopyToContext copies span from src context to dst context. Returns dst context with span.
func CopyToContext(src context.Context, dst context.Context) context.Context {
	span := trace.SpanFromContext(src)
	return trace.ContextWithSpan(dst, span)
}
