package healthcheck

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
)

var (
	statement          = spanner.NewStatement("SELECT 1")
	_         Subcheck = (*SpannerSubCheck)(nil)
)

// SpannerSubCheck is a subcheck for Cloud Spanner.
type SpannerSubCheck struct {
	client *spanner.Client
}

// WithSpannerCheck returns a SpannerSubCheck.
func WithSpannerCheck(client *spanner.Client) *SpannerSubCheck {
	return &SpannerSubCheck{client: client}
}

func (check *SpannerSubCheck) name() string {
	return fmt.Sprintf("spanner (%s)", check.client.DatabaseName())
}

func (check *SpannerSubCheck) run(ctx context.Context) error {
	iter := check.client.Single().Query(ctx, statement)
	defer iter.Stop()

	_, err := iter.Next()
	return err
}

func (check *SpannerSubCheck) isWarning() bool {
	return false
}
