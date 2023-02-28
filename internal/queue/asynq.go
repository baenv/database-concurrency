package queue

import (
	"context"
	"database-concurrency/config"
	"database-concurrency/internal/repository"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type asynqQueue struct {
	repo repository.Repository
	log  *logrus.Logger
	cfg  config.Config
}

const (
	TypeBlockHeader = "block:header"
)

type Payload struct {
	Chain  int64 `json:"chain"`
	Number int64 `json:"number"`
	Time   int64 `json:"time"`
}

func (q asynqQueue) NewBlockHeaderTask(chain int64, h types.Header) (*asynq.Task, error) {
	payload := Payload{
		Chain:  chain,
		Number: h.Number.Int64(),
		Time:   int64(h.Time),
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TypeBlockHeader,
		encodedPayload,
	), nil
}

func (q asynqQueue) HandleNewBlockHeader(ctx context.Context, t *asynq.Task) error {
	var payload Payload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	q.log.WithFields(
		logrus.Fields{
			"chain":  payload.Chain,
			"number": payload.Number,
			"time":   payload.Time,
		},
	).Info("new block")
	return nil
}

func (q asynqQueue) NewWorker() *asynq.Server {
	redisConn := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", q.cfg.REDIS_HOST, q.cfg.REDIS_PORT),
	}

	return asynq.NewServer(redisConn, asynq.Config{
		// Specify how many concurrent workers to use.
		Concurrency: 5,
		// Specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})
}

func (q asynqQueue) NewClient() *asynq.Client {
	redisConn := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", q.cfg.REDIS_HOST, q.cfg.REDIS_PORT),
	}

	return asynq.NewClient(redisConn)
}
