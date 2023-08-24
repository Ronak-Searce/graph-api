package tracing

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/opencensus"
	"google.golang.org/grpc"
)

// ServerUnaryInterceptor is a server unary interceptor for generating tracing spans.
func ServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithPropagators(opencensus.Binary{}),
	)
}

// ServerStreamInterceptor is a server stream interceptor for generating tracing spans.
func ServerStreamInterceptor() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor(
		otelgrpc.WithPropagators(opencensus.Binary{}),
	)
}

// ClientUnaryInterceptor is a client unary interceptor for generating tracing spans.
func ClientUnaryInterceptor() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithPropagators(opencensus.Binary{}),
	)
}

// ClientStreamInterceptor is a client stream interceptor for generating tracing spans.
func ClientStreamInterceptor() grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(
		otelgrpc.WithPropagators(opencensus.Binary{}),
	)
}
