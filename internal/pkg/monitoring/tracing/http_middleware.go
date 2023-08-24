package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

var (
	_ http.Handler = (*HeaderMiddleware)(nil)
)

const httpHeader = "x-picnic-traceid"

// HeaderMiddleware is an http.Handler middleware that exposes trace ID as X-Picnic-TraceID response header.
type HeaderMiddleware struct {
	next http.Handler
}

func (p *HeaderMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	span := trace.SpanFromContext(req.Context())
	w.Header().Set(httpHeader, span.SpanContext().TraceID().String())

	p.next.ServeHTTP(w, req)
}
