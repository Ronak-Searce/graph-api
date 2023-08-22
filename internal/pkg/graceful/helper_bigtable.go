package graceful

// IBigtableClient is an interface for shutting down *bigtable.Client.
type IBigtableClient interface {
	Close() error
}

// Bigtable returns a ShutdownErrorFunc that stops Bigtable client.
func Bigtable(client IBigtableClient) ShutdownErrorFunc {
	if client != nil {
		return client.Close
	}

	return nil
}
