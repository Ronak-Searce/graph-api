package graceful

// IBigQueryClient is an interface for shutting down *bigquery.Client.
type IBigQueryClient interface {
	Close() error
}

// BigQuery returns a ShutdownErrorFunc that stops BigQuery client.
func BigQuery(client IBigQueryClient) ShutdownErrorFunc {
	if client != nil {
		return client.Close
	}

	return nil
}
