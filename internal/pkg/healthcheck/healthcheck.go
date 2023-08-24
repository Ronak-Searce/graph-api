package healthcheck

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ronak-Searce/graph-api/internal/pkg/healthcheck/internal/pkg/healthcheck"
	"google.golang.org/grpc"
)

var (
	_ HealthChecker = (*healthcheck.HealthCheck)(nil)
)

// HealthChecker is an HTTP healthcheck handler that runs checks.
type HealthChecker interface {
	http.Handler
	setupper
	fmt.Stringer

	Shutdown(context.Context) error
	GRPCRegisterService(grpc.ServiceRegistrar)
	Register(grpc.ServiceRegistrar)
}

type setupper interface {
	AddSubcheck(healthcheck.Subcheck)
	EnableMetrics()
	SetName(string)
}

// HealthCheckerOption is an option to apply to HealthChecker.
type HealthCheckerOption interface {
	apply(HealthChecker)
}

// NewHealthCheck returns new healthcheck.
func NewHealthCheck(opts ...HealthCheckerOption) HealthChecker {
	hc := &healthcheck.HealthCheck{}
	for _, opt := range opts {
		opt.apply(hc)
	}
	return hc
}

type nameChanger struct {
	name string
}

func (n *nameChanger) apply(hc HealthChecker) {
	hc.SetName(n.name)
}

// WithName allows to change healthcheck name in metrics.
func WithName(name string) HealthCheckerOption {
	return &nameChanger{name: name}
}

type metricsEnabler struct{}

func (*metricsEnabler) apply(hc HealthChecker) {
	hc.EnableMetrics()
}

// WithMetrics enables metrics collection.
func WithMetrics() HealthCheckerOption {
	return &metricsEnabler{}
}
