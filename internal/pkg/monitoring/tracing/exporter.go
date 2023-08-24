package tracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Ronak-Searce/graph-api/internal/pkg/config"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SetupExporter prepares OpenTelemetry batch exporter that periodically sends batch of spans to OpenTelemetry collector.
func SetupExporter(ctx context.Context) error {
	if !config.Bool("tracing.enabled") {
		return nil
	}

	res, err := getResource(ctx)
	if err != nil {
		return fmt.Errorf("failed to set up tracer resource: %w", err)
	}

	exporter, err := createExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to set up opentelemetry exporter: %w", err)
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	tracerProvider := createTracerProvider(res, bsp)
	otel.SetTracerProvider(tracerProvider)

	DefaultTracer = GetTracer(config.ServiceName())

	return nil
}

func getResource(ctx context.Context) (*resource.Resource, error) {
	hostname, _ := os.Hostname()

	return resource.New(
		ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName()),
			semconv.ServiceNamespaceKey.String("backend"),
			semconv.ServiceInstanceIDKey.String(hostname),
			semconv.DeploymentEnvironmentKey.String(config.EnvID()),
		),
	)
}

func createExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		config.String("tracing.collector_addr"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func createTracerProvider(res *resource.Resource, processor trace.SpanProcessor) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(processor),
	)
}
