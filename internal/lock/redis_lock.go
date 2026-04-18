package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (string, error)
	Release(ctx context.Context, key string, token string) (bool, error)
}

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) Locker {
	return &Redis{client: client}
}

func randomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (r *Redis) Acquire(ctx context.Context, key string, ttl time.Duration) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}

	res, err := r.client.SetArgs(ctx, key, token, redis.SetArgs{
		Mode: "NX",
		TTL:  ttl,
	}).Result()
	if err != nil {
		return "", err
	}

	if res != "OK" {
		return "", nil
	}

	return token, nil
}

func (r *Redis) Release(ctx context.Context, key string, token string) (bool, error) {
	if token == "" {
		return false, errors.New("empty lock token")
	}

	const releaseScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`

	deleted, err := r.client.Eval(ctx, releaseScript, []string{key}, token).Int64()
	if err != nil {
		return false, err
	}

	return deleted == 1, nil
}
