package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisStreamWrapper interface to handle streams
type RedisTableWrapper interface {
	SetValue(ctx context.Context, data interface{}) error
	GetValue(ctx context.Context) (string, error)
}

type redisTableWrapper struct {
	c     *redis.Client
	table string
}

func (r redisTableWrapper) SetValue(ctx context.Context, data interface{}) error {
	_, err := r.c.Set(ctx, r.table, data, 0).Result()
	return err
}

func (r redisTableWrapper) GetValue(ctx context.Context) (string, error) {
	return r.c.Get(ctx, r.table).Result()
}
