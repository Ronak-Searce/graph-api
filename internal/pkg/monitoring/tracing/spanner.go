package tracing

import (
	"context"
	"strings"

	"cloud.google.com/go/spanner"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// StartSpannerTransaction is a helper to start span with some spanner attributes.
func StartSpannerTransaction(ctx context.Context, name string, client *spanner.Client, statement *spanner.Statement) (context.Context, trace.Span) {
	return Start(
		ctx,
		name,
		trace.WithAttributes(
			semconv.DBSystemKey.String("spanner"),
			semconv.DBNameKey.String(client.DatabaseName()),
			semconv.DBStatementKey.String(statement.SQL),
			semconv.DBOperationKey.String(strings.Fields(statement.SQL)[0]),
		),
	)
}

//func StartSpannerMutations(ctx context.Context, name string, client *spanner.Client, table string, columns []string) (context.Context, trace.Span) {
//	return Start(
//		ctx,
//		name,
//		trace.WithAttributes(
//			semconv.DBSystemKey.String("spanner"),
//			semconv.DBNameKey.String(client.DatabaseName()),
//			semconv.DBOperationKey.String("mutation"),
//			semconv.DBSQLTableKey.String(table),
//			dbSpannerColumnsKey.StringSlice(columns),
//		),
//	)
//}
