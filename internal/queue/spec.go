package queue

import (
	"context"
	"database-concurrency/config"
	"database-concurrency/internal/repository"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type Queue interface {
	NewBlockHeaderTask(chain int64, h types.Header) (*asynq.Task, error)
	HandleNewBlockHeader(ctx context.Context, t *asynq.Task) error
	NewWorker() *asynq.Server
	NewClient() *asynq.Client
}

func NewAsynqQueue(repo repository.Repository, cfg config.Config, log *logrus.Logger) Queue {
	queueLog := *log

	return asynqQueue{
		repo: repo,
		cfg:  cfg,
		log:  &queueLog,
	}
}
