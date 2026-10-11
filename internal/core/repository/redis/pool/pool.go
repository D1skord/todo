package core_redis_pool

import (
	"context"
	"time"
)

type Pool interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error

	HGet(ctx context.Context, key, field string) ([]byte, error)
	HSet(ctx context.Context, key, field string, value []byte) error
	Expire(ctx context.Context, key string, ttl time.Duration) error

	Close() error
	TTL() time.Duration
}
