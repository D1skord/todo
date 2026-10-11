package core_goredis_pool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	core_redis_pool "github.com/D1skord/todo/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

type Pool struct {
	client *redis.Client
	ttl    time.Duration
}

var _ core_redis_pool.Pool = &Pool{}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.Timeout,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		MaxRetries:   -1,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Pool{
		client: client,
		ttl:    config.TTL,
	}, nil
}

func (p *Pool) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := p.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, mapError(err)
	}

	return value, nil
}

func (p *Pool) Set(
	ctx context.Context,
	key string,
	value []byte,
	ttl time.Duration,
) error {
	return mapError(p.client.Set(ctx, key, value, ttl).Err())
}

func (p *Pool) Del(ctx context.Context, keys ...string) error {
	return mapError(p.client.Del(ctx, keys...).Err())
}

func (p *Pool) HGet(ctx context.Context, key, field string) ([]byte, error) {
	value, err := p.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		return nil, mapError(err)
	}

	return value, nil
}

func (p *Pool) HSet(
	ctx context.Context,
	key, field string,
	value []byte,
) error {
	return mapError(p.client.HSet(ctx, key, field, value).Err())
}

func (p *Pool) Expire(
	ctx context.Context,
	key string,
	ttl time.Duration,
) error {
	return mapError(p.client.Expire(ctx, key, ttl).Err())
}

func (p *Pool) Close() error {
	return p.client.Close()
}

func (p *Pool) TTL() time.Duration {
	return p.ttl
}

func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_redis_pool.ErrKeyNotFound
	}

	return err
}
