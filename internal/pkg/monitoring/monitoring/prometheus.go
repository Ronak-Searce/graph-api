package monitoring

import (
	"net/http"
	"strings"

	"github.com/Ronak-Searce/graph-api/internal/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterPrometheusSuffix sets service name as suffix for default prometheus registerer.
func RegisterPrometheusSuffix() {
	prefix := strings.ReplaceAll(config.ServiceName(), "-", "_") + "_"
	prometheus.DefaultRegisterer = prometheus.WrapRegistererWithPrefix(
		prefix,
		prometheus.DefaultRegisterer,
	)
}

// HTTPHandler is an HTTP handler for serving prometheus metrics.
func HTTPHandler() http.Handler {
	return promhttp.Handler()
}
