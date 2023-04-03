package stream

import (
	"context"
	"database-concurrency/config"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// RedisWrapper Redis wrapper to handle pub/sub calls
type RedisWrapper interface {
	Ping(ctx context.Context) (string, error)
	// CreateStream creates a stream that sends incoming messages into a buffered channel
	CreateStream(streamName string, bufferSize int) RedisStreamWrapper
}

type redisWapper struct {
	C *redis.Client
}

// InitRedisClient get a redis wrapper instance
func InitRedisClient(cfg config.Config) (RedisWrapper, error) {
	c := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%v:%v", "127.0.0.1", "6379"),
		MaxRetries: 3,
	})

	return &redisWapper{
		C: c,
	}, nil
}

// Ping ping redis server
func (w *redisWapper) Ping(ctx context.Context) (string, error) {
	return w.C.Ping(ctx).Result()
}

// CreateStream creates a stream that sends incoming messages into a buffered channel
func (w *redisWapper) CreateStream(streamName string, bufferSize int) RedisStreamWrapper {
	return &redisStreamWrapper{
		c:           w.C,
		stream:      streamName,
		bufferSize:  bufferSize,
		messageChan: make(chan interface{}, bufferSize),
		errChan:     make(chan error),
	}
}
