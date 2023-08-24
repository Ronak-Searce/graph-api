package graceful

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/Ronak-Searce/graph-api/internal/pkg/logger"
)

// ShutdownManager is a manager that is a Shutdowner and can handle signals.
type ShutdownManager interface {
	Shutdowner
	RegisterSignals(...os.Signal)
	AddStep(Shutdowner)
}

type shutdownManager struct {
	timeout    time.Duration
	signalChan chan os.Signal
	steps      []Shutdowner
}

// ShutdownManagerOptions is an options structure for ShutdownManager.
type ShutdownManagerOptions struct {
	Timeout time.Duration
}

// New returns a new ShutdownManager.
func New(opts *ShutdownManagerOptions, steps ...Shutdowner) ShutdownManager {
	return &shutdownManager{
		timeout:    opts.Timeout,
		signalChan: make(chan os.Signal, 1),
		steps:      steps,
	}
}

func (s *shutdownManager) Shutdown(ctx context.Context) error {
	if s.steps == nil {
		return nil
	}

	var cancel context.CancelFunc
	if s.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}

	total := len(s.steps)

	var err error
Steps:
	for i, step := range s.steps {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			logger.Errorf(ctx, "graceful shutdown canceled: %v", err)
			break Steps
		default:
			err = step.Shutdown(ctx)
		}

		logger.Infof(ctx, "graceful shutdown step %s (%d / %d) completed", step, i+1, total)
	}

	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)

	return nil
}

func (s *shutdownManager) AddStep(step Shutdowner) {
	s.steps = append(s.steps, step)
}

func (s *shutdownManager) String() string {
	return "shutdown manager"
}

func (s *shutdownManager) RegisterSignals(sig ...os.Signal) {
	signal.Notify(s.signalChan, sig...)
	go s.handleSignals()
}

func (s *shutdownManager) handleSignals() {
	sig := <-s.signalChan
	signal.Stop(s.signalChan)

	ctx := context.Background()
	logger.Infof(ctx, "received signal %s, shutting down gracefully", sig.String())
	_ = s.Shutdown(ctx)
}
