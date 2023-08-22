package graceful

import (
	"context"
)

// ICentrifugeNode is an interface for shutting down *centrifuge.Node.
type ICentrifugeNode interface {
	Shutdown(context.Context) error
}

// Centrifuge returns a ShutdownContextErrorFunc that gracefully stops centrifuge node.
func Centrifuge(node ICentrifugeNode) ShutdownContextErrorFunc {
	if node != nil {
		return node.Shutdown
	}

	return nil
}
