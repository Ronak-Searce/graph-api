package tracing

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
)

// HTTPMiddleware is an HTTP middleware that:
// 1. Checks if HTTP request has tracing (by checking B3 headers);
// 2. Starts new span (and sets parent span if any was in HTTP request) and stores it in context.
func HTTPMiddleware(h http.Handler) http.Handler {
	return otelhttp.NewHandler(
		&HeaderMiddleware{next: h},
		"",
		otelhttp.WithPropagators(b3.New()),
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
	)
}

// HTTPClient adds middleware transport to given client and returns it.
// The middleware transport propagates tracing info to request HTTP headers.
// If client is nil default client is used.
func HTTPClient(client *http.Client) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseTransport := client.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	client.Transport = otelhttp.NewTransport(
		baseTransport,
		otelhttp.WithPropagators(b3.New()),
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
	)

	return client
}

func spanNameFormatter(_ string, r *http.Request) string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.String())
}
