package healthcheck

const (
	// HTTPLivenessRoute is a standard route to use for our liveness probes.
	HTTPLivenessRoute = "/healthz/live"
	// HTTPReadinessRoute is a standard route to use for our readiness probes.
	HTTPReadinessRoute = "/healthz/ready"
)
