package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// lock:resource:{resource_id}

type RedisLock struct {
	Client *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{Client: client}
}

func (r *RedisLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {

	ok, err := r.Client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (r *RedisLock) Release(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
