package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisStreamWrapper interface to handle streams
type RedisStreamGroupConsumerWrapper interface {
	Claim(ctx context.Context, stream, eventID string) error
	Read(ctx context.Context, count int64, streams ...string) ([]redis.XMessage, error)
	Acknowledge(ctx context.Context, stream, eventID string) error
}

type redisStreamConsumerWrapper struct {
	c        *redis.Client
	group    string
	consumer string
}

func (s redisStreamConsumerWrapper) Read(ctx context.Context, count int64, streams ...string) ([]redis.XMessage, error) {
	stream, err := s.c.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    s.group,
		Streams:  streams,
		Consumer: s.consumer,
		Count:    count,
	}).Result()
	if err != nil {
		return nil, err
	}

	messages := []redis.XMessage{}
	for _, message := range stream {
		messages = append(messages, message.Messages...)
	}
	return messages, nil
}

func (s redisStreamConsumerWrapper) Claim(ctx context.Context, stream, eventID string) error {
	_, err := s.c.XClaimJustID(ctx, &redis.XClaimArgs{
		Stream:   stream,
		Group:    s.group,
		Consumer: s.consumer,
		MinIdle:  0,
		Messages: []string{eventID},
	}).Result()
	return err
}

func (s redisStreamConsumerWrapper) Acknowledge(ctx context.Context, stream, eventID string) error {
	_, err := s.c.XAck(ctx, stream, s.group, eventID).Result()
	return err
}
