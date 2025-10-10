package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Increment(key string, expiration time.Duration) (int, error)
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStore{client: rdb}
}

func (r *RedisStore) Increment(key string, expiration time.Duration) (int, error) {
	ctx := context.Background()
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.client.Expire(ctx, key, expiration)
	}
	return int(count), nil
}
