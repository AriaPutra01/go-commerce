package cache

import (
	"context"
	"time"
)

type Cache interface {
	Save(ctx context.Context, key, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
