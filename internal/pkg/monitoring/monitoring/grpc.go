package monitoring

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcResponsesLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_msg_latency",
			Help:    "Latency of GRPC responses in seconds",
			Buckets: prometheus.ExponentialBucketsRange(0.001, 30, 10),
		},
		[]string{"grpc_service", "grpc_method", "grpc_type", "grpc_code"},
	)
	_ grpc.UnaryServerInterceptor  = MetricsServerUnaryInterceptor
	_ grpc.StreamServerInterceptor = MetricsServerStreamInterceptor
)

const unknown = "unknown"

// MetricsServerUnaryInterceptor is a server unary interceptor for generating prometheus metrics
func MetricsServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	serviceName, methodName := splitMethodName(info.FullMethod)
	timer := prometheus.NewTimer(nil)

	res, err := handler(ctx, req)

	st, _ := status.FromError(err)
	grpcResponsesLatency.
		WithLabelValues(
			serviceName,
			methodName,
			"unary",
			st.Code().String(),
		).
		Observe(timer.ObserveDuration().Seconds())

	return res, err
}

// MetricsServerStreamInterceptor is a server stream interceptor for generating prometheus metrics
func MetricsServerStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	serviceName, methodName := splitMethodName(info.FullMethod)
	timer := prometheus.NewTimer(nil)

	err := handler(srv, stream)

	st, _ := status.FromError(err)
	grpcResponsesLatency.
		WithLabelValues(
			serviceName,
			methodName,
			"stream",
			st.Code().String(),
		).
		Observe(timer.ObserveDuration().Seconds())

	return err
}

func splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return unknown, unknown
}
