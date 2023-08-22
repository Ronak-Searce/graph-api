package graceful

import (
	"context"
	"sync"
	"time"

	"gitlab.com/picnic-app/backend/libs/golang/logger"
)

type parallelShutdown struct {
	name    string
	funcs   []Shutdowner
	timeout time.Duration
}

// ParallelShutdownOptions is an options structure for parallel Shutdowner.
type ParallelShutdownOptions struct {
	Name    string
	Timeout time.Duration
}

// Parallel returns a Shutdowner that executes underlying Shutdowner's in parallel.
func Parallel(opts *ParallelShutdownOptions, funcs ...Shutdowner) Shutdowner {
	return &parallelShutdown{
		name:    opts.Name,
		funcs:   funcs,
		timeout: opts.Timeout,
	}
}

func (s *parallelShutdown) Shutdown(ctx context.Context) error {
	var cancel context.CancelFunc
	if s.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(s.funcs))

	for _, f := range s.funcs {
		go func(f Shutdowner) {
			defer wg.Done()
			defer recoverer(ctx)

			err := f.Shutdown(ctx)
			if err != nil {
				logger.Errorf(ctx, "error during graceful shutdown %s: %v", s, err)
			}
		}(f)
	}

	return WaitGroup(wg).Shutdown(ctx)
}

func (s *parallelShutdown) String() string {
	return s.name
}
