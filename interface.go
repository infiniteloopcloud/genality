package genality

import (
	"context"
	"time"
)

type Descriptor interface {
	Add(ctx context.Context, apiKey string) error
	GetCountFrom(ctx context.Context, record string, start time.Time) (int, error)
	GetCountBuckets(ctx context.Context, record string, start time.Time, bucketSize time.Duration) ([]BucketResponse, error)
}
