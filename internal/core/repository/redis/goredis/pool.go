// Package redispool provides a concrete implentation of the redis.Pool interface
// using the github.com/redis/go-redis/v9 client.
// It acts as an adapter, translating abstract database operations into goredis-specific calls,
// manages configuration for Redis database.
package redispool

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/redis"
)

// Pool wraps the standard redis.CLient to implement the [redis.Pool] interface.
// It also stores TTL parameter that accompanies all SET operations.
type Pool struct {
	client *goredis.Client
	ttl    time.Duration
}

// NewPool establishes a client for the Redis database
// using the provided configuration. It verifies the connection with a Ping.
func NewPool(ctx context.Context, cfg config) (*Pool, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redispool ping: %w", err)
	}

	return &Pool{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

func (p *Pool) Get(ctx context.Context, key string) redis.StringCmd {
	cmd := p.client.Get(ctx, key)
	return goredisStringCmd{cmd}
}

func (p *Pool) Set(ctx context.Context, key string, value any, ttl time.Duration) redis.StatusCmd {
	cmd := p.client.Set(ctx, key, value, ttl)
	return goredisStatusCmd{cmd}
}

func (p *Pool) Del(ctx context.Context, keys ...string) redis.IntCmd {
	cmd := p.client.Del(ctx, keys...)
	return goredisIntCmd{cmd}
}

func (p *Pool) HGet(ctx context.Context, key string, field string) redis.StringCmd {
	cmd := p.client.HGet(ctx, key, field)
	return goredisStringCmd{cmd}
}

func (p *Pool) HSet(ctx context.Context, key string, values ...any) redis.IntCmd {
	cmd := p.client.HSet(ctx, key, values...)
	return goredisIntCmd{cmd}
}

func (p *Pool) Close() error {
	return p.client.Close()
}
func (p *Pool) TTL() time.Duration {
	return p.ttl
}
