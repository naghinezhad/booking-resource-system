package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(addr string, password string) *Redis {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &Redis{
		Client: rdb,
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
