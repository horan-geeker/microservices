package store

import (
	"context"
	"time"
)

// DataFactory .
type DataFactory interface {
	Users() UserStore
}

// CacheFactory .
type CacheFactory interface {
	GetKey(ctx context.Context, key string) (string, error)
	SetKey(ctx context.Context, key, value string, expire time.Duration) error
}
