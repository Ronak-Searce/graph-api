package graceful

// ISpannerClient is an interface for shutting down *spanner.Client.
type ISpannerClient interface {
	Close()
}

// Spanner returns a ShutdownErrorFunc that stops Spanner client.
func Spanner(client ISpannerClient) ShutdownErrorFunc {
	if client != nil {
		return makeErrorFunc(client.Close)
	}

	return nil
}
