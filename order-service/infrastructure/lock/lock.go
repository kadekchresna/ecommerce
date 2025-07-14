package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type DistributedLock interface {
	Acquire(ctx context.Context, key string, ttlSeconds int) (bool, error)
	Release(ctx context.Context, key string) error
}

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{client: client}
}

func (r *RedisLock) Acquire(ctx context.Context, key string, ttlSeconds int) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, "locked", time.Duration(ttlSeconds)*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("redis lock acquire failed: %w", err)
	}
	return ok, nil
}

func (r *RedisLock) Release(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis lock release failed: %w", err)
	}
	return nil
}
