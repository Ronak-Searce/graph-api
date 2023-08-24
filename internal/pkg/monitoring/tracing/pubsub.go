package tracing

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// PubSubAttributeKey is a Pub/Sub attribute key for span.
	PubSubAttributeKey = "trace-context"
)

// PubSubAttributeValue returns Pub/Sub attribute value with current span from context.
func PubSubAttributeValue(ctx context.Context) string {
	data, err := MarshalJSONFromContext(ctx)
	if err != nil {
		return ""
	}

	return string(data)
}

// StartPubSubConsumer extracts dumped span from Pub/Sub message attributes and injects it into context.
func StartPubSubConsumer(ctx context.Context, sub *pubsub.Subscription, msg *pubsub.Message) (context.Context, trace.Span) {
	span, ok := msg.Attributes[PubSubAttributeKey]
	if ok {
		ctx = UnmarshalJSONToContext(ctx, []byte(span))
	}

	deliveryAttempt := 0
	if msg.DeliveryAttempt != nil {
		deliveryAttempt = *msg.DeliveryAttempt
	}

	return Start(
		ctx,
		fmt.Sprintf("%s receive", sub.ID()),
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("pubsub"),
			semconv.MessagingURLKey.String(sub.String()),
			semconv.MessagingMessageIDKey.String(msg.ID),
			semconv.MessagingMessagePayloadSizeBytesKey.Int(len(msg.Data)),
			semconv.MessagingOperationReceive,
			pubsubDeliveryAttemptKey.Int(deliveryAttempt),
			pubsubOrderingKeyKey.String(msg.OrderingKey),
		),
	)
}

// StartPubSubPublisher is a helper to start Pub/Sub publisher span and automatically add tracing info to message attributes.
func StartPubSubPublisher(ctx context.Context, topic *pubsub.Topic, msg *pubsub.Message) (context.Context, trace.Span) {
	ctx, span := Start(
		ctx,
		fmt.Sprintf("%s send", topic.ID()),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystemKey.String("pubsub"),
			semconv.MessagingDestinationKey.String(topic.ID()),
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingURLKey.String(topic.String()),
			semconv.MessagingMessagePayloadSizeBytesKey.Int(len(msg.Data)),
			pubsubOrderingKey.Bool(topic.EnableMessageOrdering),
		),
	)

	msg.Attributes[PubSubAttributeKey] = PubSubAttributeValue(ctx)

	return ctx, span
}
