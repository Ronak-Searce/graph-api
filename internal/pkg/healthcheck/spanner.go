package healthcheck

import (
	"cloud.google.com/go/spanner"
	"github.com/Ronak-Searce/graph-api/internal/pkg/healthcheck/internal/pkg/healthcheck"
)

type spannerCheck struct {
	c *healthcheck.SpannerSubCheck
}

func (c *spannerCheck) apply(hc HealthChecker) {
	hc.AddSubcheck(c.c)
}

// WithSpannerCheck adds Spanner check (SELECT 1).
func WithSpannerCheck(client *spanner.Client) HealthCheckerOption {
	return &spannerCheck{c: healthcheck.WithSpannerCheck(client)}
}
