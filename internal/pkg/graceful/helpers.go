package graceful

import (
	"context"
	"sync"
	"time"

	"github.com/Ronak-Searce/graph-api/internal/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// Context returns a ShutdownErrorFunc that cancels the context.
func Context(cancel context.CancelFunc) ShutdownErrorFunc {
	if cancel != nil {
		return makeErrorFunc(cancel)
	}

	return nil
}

// Logger returns a ShutdownErrorFunc that flushes logger. nil uses default logger.
func Logger(sugaredLogger *zap.SugaredLogger) ShutdownErrorFunc {
	if sugaredLogger == nil {
		sugaredLogger = logger.Logger()
	}
	return sugaredLogger.Sync
}

// WaitGroup returns a ShutdownErrorFunc that waits WaitGroup.
func WaitGroup(wg *sync.WaitGroup) ShutdownErrorFunc {
	if wg != nil {
		return makeErrorFunc(wg.Wait)
	}

	return nil
}

// Tracer returns a ShutdownContextErrorFunc that gracefully shutdowns tracer exporter.
func Tracer() ShutdownContextErrorFunc {
	return func(ctx context.Context) error {
		tracerProvider, ok := otel.GetTracerProvider().(*trace.TracerProvider)
		if !ok {
			return nil
		}

		err := tracerProvider.ForceFlush(ctx)
		if err != nil {
			logger.Errorf(ctx, "failed to flush tracer provider")
		}
		return tracerProvider.Shutdown(ctx)
	}
}

// Sleep returns a ShutdownErrorFunc that just sleeps. Useful for delays and debug purposes.
func Sleep(duration time.Duration) ShutdownErrorFunc {
	return func() error {
		time.Sleep(duration)
		return nil
	}
}

func recoverer(ctx context.Context) {
	r := recover()
	if r != nil {
		logger.Errorf(ctx, "got panic during shutdown: %v", r)
	}
}
