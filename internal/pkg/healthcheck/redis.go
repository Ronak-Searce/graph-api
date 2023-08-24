package healthcheck

import "github.com/Ronak-Searce/graph-api/internal/pkg/healthcheck/internal/pkg/healthcheck"

type redisCheck struct {
	c *healthcheck.RedisSubCheck
}

func (c *redisCheck) apply(hc HealthChecker) {
	hc.AddSubcheck(c.c)
}

// WithRedisCheck adds Redis check (PING command).
func WithRedisCheck(client healthcheck.IRedisClient, warn bool) HealthCheckerOption {
	return &redisCheck{c: healthcheck.WithRedisCheck(client, warn)}
}
