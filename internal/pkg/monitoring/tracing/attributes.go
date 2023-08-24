package tracing

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	userIDKey = attribute.Key("picnic.user.id")

	graphqlOperationNameKey = attribute.Key("graphql.operation.name")
	graphqlOperationTypeKey = attribute.Key("graphql.operation.type")
	graphqlDocumentKey      = attribute.Key("graphql.document")

	//dbSpannerColumnsKey = attribute.Key("db.spanner.columns")

	pubsubDeliveryAttemptKey = attribute.Key("messaging.pubsub.delivery_attempt")
	pubsubOrderingKey        = attribute.Key("messaging.pubsub.ordering")
	pubsubOrderingKeyKey     = attribute.Key("messaging.pubsub.ordering_key")
)

// AddUserIDToSpan adds userID to current span's (got from context) attributes.
func AddUserIDToSpan(ctx context.Context, userID string) context.Context {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(userIDKey.String(userID))

	return trace.ContextWithSpan(ctx, span)
}

// AddGraphqlToSpan adds graphql context to span attributes and updates span name.
// https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/instrumentation/graphql/
func AddGraphqlToSpan(ctx context.Context, op *graphql.OperationContext) context.Context {
	if op == nil {
		return ctx
	}

	span := trace.SpanFromContext(ctx)

	if op.Operation == nil {
		return ctx
	}

	opName := op.Operation.Name
	opType := string(op.Operation.Operation)

	span.SetAttributes(
		graphqlOperationNameKey.String(opName),
		graphqlOperationTypeKey.String(opType),
		graphqlDocumentKey.String(op.RawQuery),
	)

	span.SetName(fmt.Sprintf("%s %s", opType, opName))

	return trace.ContextWithSpan(ctx, span)
}
