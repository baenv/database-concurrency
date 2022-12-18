package main

import (
	"context"

	"database-concurrency/config"
	"database-concurrency/internal/indexer"
	"database-concurrency/internal/queue"
	"database-concurrency/internal/repository"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	cfg, err := config.NewConfig()
	if err != nil {
		log.WithError(err).Error("failed to load config")
		return
	}

	repo, err := repository.Init(cfg)
	if err != nil {
		log.WithError(err).Error("failed to init repo")
		return
	}

	// Create global context
	ctx := context.Background()

	tqueue := queue.NewAsynqQueue(repo, cfg, log)

	// Sub to ETHEREUM
	indexer := indexer.NewEVM(1, repo, cfg, tqueue, log)
	if err := indexer.Sub(ctx); err != nil {
		log.WithError(err).Error("failed to sub to ethereum")
	}
}
