package healthcheck

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Check runs the healthcheck.
func (hc *HealthCheck) Check(ctx context.Context, _ *healthgrpc.HealthCheckRequest) (*healthgrpc.HealthCheckResponse, error) {
	checkResponse := hc.Run(ctx)
	resp := &healthgrpc.HealthCheckResponse{}

	var err error

	if checkResponse.Status != StatusError {
		resp.Status = healthgrpc.HealthCheckResponse_SERVING
	} else {
		resp.Status = healthgrpc.HealthCheckResponse_NOT_SERVING
		err = status.Error(codes.Internal, "healthcheck failed")
	}

	return resp, err
}

// GRPCRegisterService registers GRPC method implementations in GRPC server.
func (hc *HealthCheck) GRPCRegisterService(s grpc.ServiceRegistrar) {
	healthgrpc.RegisterHealthServer(s, hc)
}

// Register registers GRPC method implementations in GRPC server.
func (hc *HealthCheck) Register(s grpc.ServiceRegistrar) {
	hc.GRPCRegisterService(s)
}
