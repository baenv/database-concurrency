package redis

import (
	"context"
	"database-concurrency/config"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// RedisWrapper Redis wrapper to handle pub/sub calls
type RedisWrapper interface {
	Ping(ctx context.Context) (string, error)
	// CreateStream creates a stream that sends incoming messages into a buffered channel
	CreateStream(streamName string, bufferSize int) RedisStreamWrapper

	CreateConsumer(ctx context.Context, stream, group, consumer string) redisStreamConsumerWrapper
	CreateTable(tableName string) RedisTableWrapper
}

type redisWapper struct {
	C *redis.Client
}

// InitRedisClient get a redis wrapper instance
func InitRedisClient(cfg config.Config) (RedisWrapper, error) {
	c := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%v:%v", cfg.REDIS_HOST, cfg.REDIS_PORT),
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

func (w *redisWapper) CreateConsumer(ctx context.Context, stream, group, consumer string) redisStreamConsumerWrapper {
	_, err := w.C.XGroupCreate(ctx, stream, group, "$").Result() // Ignore err if group already created
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalln(err)
	}
	return redisStreamConsumerWrapper{
		c:        w.C,
		group:    group,
		consumer: consumer,
	}
}

// CreateTable creates a table in the form of key value
func (w *redisWapper) CreateTable(tableName string) RedisTableWrapper {
	return redisTableWrapper{
		c:     w.C,
		table: tableName,
	}
}
