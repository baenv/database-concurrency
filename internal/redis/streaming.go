package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

// RedisStreamWrapper interface to handle streams
type RedisStreamWrapper interface {
	// SetChannels set the message and error channels
	SetChannels(messageChan chan interface{}, errChan chan error)
	// Publish publish data into the stream
	Publish(ctx context.Context, message interface{}) (string, error)
	// Consume consume messages from the stream with a count limit. If 0 it will consume all messages
	Consume(ctx context.Context, count int64)
	// MessageChannel get the message channel
	MessageChannel() chan interface{}
	// ErrorChannel get the error channel
	ErrorChannel() chan error
	// FinishedChannel get the finished notification channel
	FinishedChannel() chan bool
}

type redisStreamWrapper struct {
	c            *redis.Client
	stream       string
	bufferSize   int
	messageChan  chan interface{} // Channel where the consumed messages are send
	errChan      chan error
	finishedChan chan bool
}

// SetChannels set the message and error channels
func (s *redisStreamWrapper) SetChannels(messageChan chan interface{}, errChan chan error) {
	if messageChan != nil {
		s.messageChan = messageChan
	}
	if errChan != nil {
		s.errChan = errChan
	}
}

// MessageChannel get the message channel
func (s *redisStreamWrapper) MessageChannel() chan interface{} {
	return s.messageChan
}

// ErrorChannel get the error channel
func (s *redisStreamWrapper) ErrorChannel() chan error {
	return s.errChan
}

// FinishedChannel get the finished notification channel
func (s *redisStreamWrapper) FinishedChannel() chan bool {
	return s.finishedChan
}

// Publish publish data into the stream
func (s *redisStreamWrapper) Publish(ctx context.Context, message interface{}) (string, error) {
	args := redis.XAddArgs{
		Stream: s.stream,
		Values: map[string]interface{}{
			"data": message,
		},
	}
	return s.c.XAdd(ctx, &args).Result()
}

// Consume consume messages from the stream with a count limit. If 0 it will consume all messages
func (s *redisStreamWrapper) Consume(ctx context.Context, count int64) {
	go func() {
		for {
			var err error
			var data []redis.XMessage
			if count > 0 {
				data, err = s.c.XRangeN(ctx, s.stream, "-", "+", count).Result()
			} else {
				data, err = s.c.XRange(ctx, s.stream, "-", "+").Result()
			}
			if err != nil {
				s.errChan <- err
			}
			for _, element := range data {
				data := []byte(element.Values["data"].(string)) // Get pack message
				var message interface{}
				err := json.Unmarshal(data, &message)
				if err != nil {
					s.errChan <- err
					continue
				}
				s.messageChan <- message
				s.c.XDel(ctx, s.stream, element.ID) // Remove consumed message
			}
		}
	}()
}
