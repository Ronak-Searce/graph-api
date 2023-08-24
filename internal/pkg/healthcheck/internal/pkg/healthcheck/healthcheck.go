package healthcheck

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/Ronak-Searce/graph-api/internal/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// StatusOK represents normal status.
	StatusOK CheckStatus = iota
	// StatusWarn represents warning status.
	StatusWarn
	// StatusError represents errored status.
	StatusError
)

// CheckStatus is a type to represent check status.
type CheckStatus int

// String returns status as a string.
func (c CheckStatus) String() string {
	switch c {
	case StatusOK:
		return "ok"
	case StatusWarn:
		return "warn"
	case StatusError:
		return "error"
	default:
		return "invalid"
	}
}

// MarshalJSON returns json of the status as string.
func (c CheckStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// HealthCheck is a subchecks runner and exposer as HTTP and GRPC handlers.
type HealthCheck struct {
	name     string
	checks   []Subcheck
	metrics  bool
	shutdown bool

	healthgrpc.UnimplementedHealthServer
}

// Response is a struct for storing complete healthcheck response.
type Response struct {
	Status   CheckStatus         `json:"status"`
	Checks   []*SubcheckResponse `json:"checks"`
	Duration float64             `json:"duration_seconds"`
}

// SubcheckResponse is a struct for storing subcheck response.
type SubcheckResponse struct {
	Name     string      `json:"name"`
	Status   CheckStatus `json:"status"`
	Error    string      `json:"error,omitempty"`
	Duration float64     `json:"duration_seconds"`
}

// Subcheck is an interface to run individual check.
type Subcheck interface {
	name() string
	isWarning() bool
	run(ctx context.Context) error
}

// Run runs all subchecks and returns Response.
func (hc *HealthCheck) Run(ctx context.Context) *Response {
	resp := &Response{
		Status: StatusOK,
	}

	timer := prometheus.NewTimer(nil)
	resp.Checks = hc.runChecks(ctx)
	duration := timer.ObserveDuration().Seconds()

	resp.Duration = duration

	for _, subcheck := range resp.Checks {
		if subcheck.Status > resp.Status {
			resp.Status = subcheck.Status
		}
	}

	if hc.metrics {
		latestStatus.WithLabelValues(hc.name).Set(float64(resp.Status))
		responsesTotal.WithLabelValues(hc.name, resp.Status.String()).Inc()
		responsesDurations.WithLabelValues(hc.name, resp.Status.String()).Observe(resp.Duration)
	}

	return resp
}

func (hc *HealthCheck) runChecks(ctx context.Context) []*SubcheckResponse {
	resp := make([]*SubcheckResponse, 0)

	if hc.shutdown {
		resp = append(resp, &SubcheckResponse{
			Name:     "shutdown",
			Status:   StatusError,
			Error:    "healthcheck is shutdown",
			Duration: 0,
		})
		return resp
	}

	workerChan := make(chan Subcheck)
	go func() {
		for _, check := range hc.checks {
			workerChan <- check
		}
		close(workerChan)
	}()

	respChan := make(chan *SubcheckResponse, len(hc.checks))

	wg := &sync.WaitGroup{}
	for range hc.checks {
		wg.Add(1)
		go hc.runCheckWorker(ctx, wg, workerChan, respChan)
	}
	wg.Wait()
	close(respChan)

	for r := range respChan {
		if r.Status != StatusOK {
			logger.Errorf(ctx, "failed healthcheck %s: %s", r.Name, r.Error)
		}
		resp = append(resp, r)
	}

	return resp
}

func (hc *HealthCheck) runCheckWorker(ctx context.Context, wg *sync.WaitGroup, checkChan <-chan Subcheck, respChan chan<- *SubcheckResponse) {
	defer wg.Done()

	resp := &SubcheckResponse{}

	select {
	case check := <-checkChan:
		timer := prometheus.NewTimer(nil)
		err := check.run(ctx)
		duration := timer.ObserveDuration().Seconds()

		resp.Name = check.name()
		resp.Duration = duration

		if err == nil {
			resp.Status = StatusOK
		} else {
			resp.Error = err.Error()

			if check.isWarning() {
				resp.Status = StatusWarn
			} else {
				resp.Status = StatusWarn
			}
		}

		if hc.metrics {
			subcheckLatestStatus.WithLabelValues(hc.name, check.name()).Set(float64(resp.Status))
			subcheckResponsesTotal.WithLabelValues(hc.name, check.name(), resp.Status.String()).Inc()
			subcheckResponsesDurations.WithLabelValues(hc.name, check.name(), resp.Status.String()).Observe(duration)
		}
	case <-ctx.Done():
		resp.Name = "internal"
		resp.Status = StatusError
		resp.Error = ctx.Err().Error()
	}

	respChan <- resp
}

// AddSubcheck adds new subcheck.
func (hc *HealthCheck) AddSubcheck(check Subcheck) {
	hc.checks = append(hc.checks, check)
}

// SetName sets healthcheck name (used for metrics).
func (hc *HealthCheck) SetName(name string) {
	hc.name = name
}

// Shutdown forces healthcheck to fail always. Used for graceful shutdown.
func (hc *HealthCheck) Shutdown(_ context.Context) error {
	hc.shutdown = true

	return nil
}

// String returns healthcheck name.
func (hc *HealthCheck) String() string {
	return hc.name
}
