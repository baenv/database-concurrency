package indexer

import (
	"context"
	"database-concurrency/config"
	"database-concurrency/internal/queue"
	"database-concurrency/internal/repository"

	"github.com/sirupsen/logrus"
)

// Indexer is used to subscribe to specific chain, listen and index data of new blocks
type Indexer interface {
	Sub(context.Context) error
}

// New is used to create new indexer
func NewEVM(chain int64, repo repository.Repositoy, cfg config.Config, tq queue.Queue, log *logrus.Logger) Indexer {
	// Create new logger with chain field
	indexerLog := *log.WithField("chain", chain).Logger

	return &evm{
		chain: chain,
		repo:  repo,
		cfg:   cfg,
		tq:    tq,
		log:   &indexerLog,
	}
}
