package healthcheck

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
)

const nonexistingObjectName = "healthz-ready"

var (
	_ Subcheck = (*BucketSubCheck)(nil)
)

// BucketSubCheck is a subcheck for GCS bucket.
type BucketSubCheck struct {
	bucket *storage.BucketHandle
}

// WithBucketCheck returns a BucketSubCheck.
func WithBucketCheck(bucket *storage.BucketHandle) *BucketSubCheck {
	return &BucketSubCheck{bucket: bucket}
}

func (check *BucketSubCheck) name() string {
	return fmt.Sprintf("gcs (%s)", check.bucket.Object("").BucketName())
}

func (check *BucketSubCheck) run(ctx context.Context) error {
	obj := check.bucket.Object(nonexistingObjectName)
	reader, err := obj.NewRangeReader(ctx, 0, 1)
	_ = reader.Close()

	if errors.Is(err, storage.ErrObjectNotExist) {
		return nil
	}

	return err
}

func (check *BucketSubCheck) isWarning() bool {
	return false
}
