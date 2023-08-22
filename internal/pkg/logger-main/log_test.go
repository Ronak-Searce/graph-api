package logger

import (
	"context"
	"testing"

	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogLevel(t *testing.T) {
	t.Run("init log level from string", func(t *testing.T) {
		const (
			in = "DEBUG"
			ex = "debug"
		)
		InitWithString(in)
		l := Logger().Level().String()
		if l != ex {
			t.Errorf("InitWithString level = %s, but want %s", l, ex)
		}
	})
}

func TestSpanContext(t *testing.T) {
	t.Run("check existing span in context", func(t *testing.T) {
		observerCore, observedLogs := observer.New(zap.ErrorLevel)
		ctx := context.WithValue(
			context.Background(),
			loggerContextKey,
			zap.New(observerCore).Sugar(),
		)

		ctx, _ = sdk_trace.NewTracerProvider().Tracer("").Start(ctx, "testspan")
		Error(ctx, "test error")

		if observedLogs.Len() < 1 {
			t.Errorf("no errors were logged")
			t.FailNow()
		}

		if observedLogs.Len() > 1 {
			t.Errorf("more than one entry were logged")
			t.FailNow()
		}

		l := observedLogs.All()[0]

		if _, ok := l.ContextMap()[traceIDFieldName]; !ok {
			t.Errorf("%s is not logged when it should be", traceIDFieldName)
		}

		if _, ok := l.ContextMap()[spanIDFieldName]; !ok {
			t.Errorf("%s is not logged when it should be", spanIDFieldName)
		}
	})

	t.Run("check missing span in context", func(t *testing.T) {
		observerCore, observedLogs := observer.New(zap.ErrorLevel)
		ctx := context.WithValue(
			context.Background(),
			loggerContextKey,
			zap.New(observerCore).Sugar(),
		)

		Error(ctx, "test error")

		if observedLogs.Len() < 1 {
			t.Errorf("no errors were logged")
			t.FailNow()
		}

		if observedLogs.Len() > 1 {
			t.Errorf("more than one entry were logged")
			t.FailNow()
		}

		l := observedLogs.All()[0]

		if _, ok := l.ContextMap()[traceIDFieldName]; ok {
			t.Errorf("%s is logged when it shouldn't", traceIDFieldName)
		}

		if _, ok := l.ContextMap()[spanIDFieldName]; ok {
			t.Errorf("%s is logged when it shouldn't", spanIDFieldName)
		}
	})
}
