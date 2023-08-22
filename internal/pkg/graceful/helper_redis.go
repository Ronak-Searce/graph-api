package graceful

// IRedisClient is an interface for shutting down *redis.Client.
type IRedisClient interface {
	Close() error
}

// Redis returns a ShutdownErrorFunc that gracefully stops redis client.
func Redis(client IRedisClient) ShutdownErrorFunc {
	if client != nil {
		return client.Close
	}

	return nil
}
