package indexer

import (
	"context"
	"database-concurrency/config"
	"database-concurrency/internal/queue"
	"database-concurrency/internal/repository"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type evm struct {
	chain int64
	repo  repository.Repository
	cfg   config.Config
	log   *logrus.Logger
	tq    queue.Queue
}

// Sub is used to subscribe to specific chain and listen for new blocks
func (c *evm) Sub(ctx context.Context) error {
	client, err := c.Client()
	if err != nil {
		return err
	}
	defer client.Close()

	// Chan to listen for new blocks
	headers := make(chan *types.Header)

	queueClient := c.tq.NewClient()
	defer queueClient.Close()

	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return errors.WithStack(err)
	}

	for {
		select {
		case <-ctx.Done():
			c.log.Info("context canceled, stopping subscription")
			return nil
		case err := <-sub.Err():
			return errors.WithStack(err)
		case header := <-headers:
			c.log.WithField("block", header.Number.String()).Info("new block")
			// TODO: Push to global task queue

			task, err := c.tq.NewBlockHeaderTask(c.chain, *header)
			if err != nil {
				c.log.Error(err, "failed to create new block header task")
			}
			queueClient.Enqueue(
				task,
				asynq.Queue("default"),
				asynq.TaskID(fmt.Sprintf("%d:%d", c.chain, header.Number.Int64())),
				asynq.ProcessIn(2*time.Minute),
			)
		}
	}
}

// Client is used to get ethclient.Client of configured chain
func (c *evm) Client() (*ethclient.Client, error) {
	switch c.chain {
	case 1:
		client, err := ethclient.Dial(c.cfg.ETHEREUM_RPC)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return client, err
	default:
		return nil, ErrChainNotSupported
	}
}
