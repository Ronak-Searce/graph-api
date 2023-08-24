package graceful

// IPubsubTopic is an interface for shutting down *pubsub.Topic.
type IPubsubTopic interface {
	Stop()
}

// PubSubTopic returns a ShutdownErrorFunc that gracefully stops publishing into PubSub Topic.
func PubSubTopic(topic IPubsubTopic) ShutdownErrorFunc {
	if topic != nil {
		return makeErrorFunc(topic.Stop)
	}

	return nil
}
