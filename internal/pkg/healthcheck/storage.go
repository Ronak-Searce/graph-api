package healthcheck

import (
	"cloud.google.com/go/storage"
	"github.com/Ronak-Searce/graph-api/internal/pkg/healthcheck/internal/pkg/healthcheck"
)

type bucketCheck struct {
	c *healthcheck.BucketSubCheck
}

func (c *bucketCheck) apply(hc HealthChecker) {
	hc.AddSubcheck(c.c)
}

// WithBucketCheck adds GCS bucket check (tries to read nonexistent test object).
func WithBucketCheck(bucket *storage.BucketHandle) HealthCheckerOption {
	return &bucketCheck{c: healthcheck.WithBucketCheck(bucket)}
}
