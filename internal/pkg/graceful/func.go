package graceful

import "context"

var (
	_ Shutdowner = (*ShutdownContextErrorFunc)(nil)
	_ Shutdowner = (*ShutdownContextFunc)(nil)
	_ Shutdowner = (*ShutdownErrorFunc)(nil)
)

const anonymousFuncName = "func"

// ShutdownContextErrorFunc is a function that implements Shutdowner interface.
type ShutdownContextErrorFunc func(context.Context) error

// Shutdown runs the function itself.
func (s ShutdownContextErrorFunc) Shutdown(ctx context.Context) error {
	if s != nil {
		return s(ctx)
	}

	return nil
}

// String returns just dummy value.
func (s ShutdownContextErrorFunc) String() string {
	return anonymousFuncName
}

// ShutdownContextFunc is a function that implements Shutdowner interface.
type ShutdownContextFunc func(ctx context.Context)

// Shutdown runs the function itself. It always returns nil as the original function doesn't return any error.
func (s ShutdownContextFunc) Shutdown(ctx context.Context) error {
	if s != nil {
		s(ctx)
	}

	return nil
}

// String returns just dummy value.
func (s ShutdownContextFunc) String() string {
	return anonymousFuncName
}

// ShutdownErrorFunc is a function that implements Shutdowner interface.
type ShutdownErrorFunc func() error

// Shutdown runs the function itself. It returns early in case of context timeout but original function is not canceled
// and will run even after canceling context.
func (s ShutdownErrorFunc) Shutdown(ctx context.Context) error {
	if s != nil {
		return makeShutdownFunc(s)(ctx)
	}

	return nil
}

// String returns just dummy value.
func (s ShutdownErrorFunc) String() string {
	return anonymousFuncName
}

func makeShutdownFunc(f func() error) ShutdownContextErrorFunc {
	return func(ctx context.Context) error {
		ch := make(chan error, 1)
		go func() {
			defer recoverer(ctx)

			ch <- f()
			close(ch)
		}()

		select {
		case err := <-ch:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func makeErrorFunc(f func()) ShutdownErrorFunc {
	return func() error {
		f()
		return nil
	}
}
