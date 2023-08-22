package graceful

import (
	"context"
)

// IHTTPServer is an interface for shutting down *http.Server.
type IHTTPServer interface {
	Shutdown(ctx context.Context) error
}

// HTTPServer returns a ShutdownContextErrorFunc that gracefully stops HTTP server.
func HTTPServer(server IHTTPServer) ShutdownContextErrorFunc {
	if server != nil {
		return server.Shutdown
	}

	return nil
}
