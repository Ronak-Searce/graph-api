package graceful

import (
	"context"
	"fmt"
)

// Shutdowner is an interface that gracefully shuts down something.
type Shutdowner interface {
	Shutdown(context.Context) error
	fmt.Stringer
}
