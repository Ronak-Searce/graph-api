package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextKey int

const (
	loggerContextKey contextKey = iota

	// https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
	traceIDFieldName      = "logging.googleapis.com/trace"
	traceSampledFieldName = "logging.googleapis.com/trace_sampled"
	spanIDFieldName       = "logging.googleapis.com/spanId"
)

// ToContext ...
func ToContext(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

// FromContext ...
func FromContext(ctx context.Context) *zap.SugaredLogger {
	l := global

	if logger, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger); ok {
		l = logger
	}

	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		l = loggerWithSpanContext(l, span.SpanContext())
	}

	return l
}

func loggerWithSpanContext(l *zap.SugaredLogger, sc trace.SpanContext) *zap.SugaredLogger {
	return l.Desugar().With(
		zap.Stringer(traceIDFieldName, sc.TraceID()),
		zap.Bool(traceSampledFieldName, sc.IsSampled()),
		zap.Stringer(spanIDFieldName, sc.SpanID()),
	).Sugar()
}
