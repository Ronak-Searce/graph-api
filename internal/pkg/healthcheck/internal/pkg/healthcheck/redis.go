package healthcheck

import (
	"context"
	"fmt"
)

var _ Subcheck = (*RedisSubCheck)(nil)

// RedisSubCheck is a subcheck for Redis.
type RedisSubCheck struct {
	client IRedisClient
	warn   bool
}

// IRedisClient is an interface for redis client.
type IRedisClient interface {
	fmt.Stringer

	Ping(ctx context.Context) IRedisCmd
}

// IRedisCmd is an interface for redis command.
type IRedisCmd interface {
	Err() error
}

// WithRedisCheck returns a RedisSubCheck.
func WithRedisCheck(client IRedisClient, warn bool) *RedisSubCheck {
	return &RedisSubCheck{
		client: client,
		warn:   warn,
	}
}

func (check *RedisSubCheck) name() string {
	return fmt.Sprintf("redis (%s)", check.client.String())
}

func (check *RedisSubCheck) run(ctx context.Context) error {
	cmd := check.client.Ping(ctx)
	return cmd.Err()
}

func (check *RedisSubCheck) isWarning() bool {
	return check.warn
}
