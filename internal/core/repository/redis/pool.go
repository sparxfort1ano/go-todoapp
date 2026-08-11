// Package redis defines the core contract for Redis database operations
// used by the repository layer, allowing the underlying implentation (e.g. goredis) to be easily swapped or mocked.
package redis

import (
	"context"
	"time"
)

// Pool defines the client-agnostic contract for executing queries against Redis DB.
// It abstracts the underlying connection management and query execution.
type Pool interface {
	Get(ctx context.Context, key string) StringCmd
	Set(ctx context.Context, key string, value any, ttl time.Duration) StatusCmd
	Del(ctx context.Context, keys ...string) IntCmd
	HGet(ctx context.Context, key string, field string) StringCmd
	HSet(ctx context.Context, key string, values ...any) IntCmd
	Close() error
	TTL() time.Duration
}

// StringCmd is a command (i.e. GET) result returning a string.
type StringCmd interface {
	Bytes() ([]byte, error)
}

// StatusCmd is a command (i.e. SET) result returning an operation status.
type StatusCmd interface {
	Err() error
}

// IntCmd is a command (i.e. DEL) result returning an integer.
type IntCmd interface {
	Err() error
}
