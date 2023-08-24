package healthcheck

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	latestStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "healthcheck_status",
		Help: "Latest healthcheck status",
	}, []string{"healthcheck"})
	responsesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "healthcheck_responses_total",
		Help: "Total number of healthcheck responses",
	}, []string{"healthcheck", "status"})
	responsesDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "healthcheck_responses_duration_seconds",
		Help:    "Latency of healthcheck",
		Buckets: prometheus.ExponentialBucketsRange(0.0001, 10, 6),
	}, []string{"healthcheck", "status"})

	subcheckLatestStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "healthcheck_subcheck_status",
		Help: "Latest subcheck status",
	}, []string{"healthcheck", "subcheck"})
	subcheckResponsesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "healthcheck_subcheck_responses_total",
		Help: "Total number of subcheck responses",
	}, []string{"healthcheck", "subcheck", "status"})
	subcheckResponsesDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "healthcheck_subcheck_responses_duration_seconds",
		Help:    "Latency of subcheck",
		Buckets: prometheus.ExponentialBucketsRange(0.0001, 10, 6),
	}, []string{"healthcheck", "subcheck", "status"})
)

// EnableMetrics enables metrics collection.
func (hc *HealthCheck) EnableMetrics() {
	hc.metrics = true
}
