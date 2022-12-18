package main

import (
	"database-concurrency/config"
	"database-concurrency/internal/queue"
	"database-concurrency/internal/repository"

	"github.com/hibiken/asynq"
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

	tqueue := queue.NewAsynqQueue(repo, cfg, log)

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		queue.TypeBlockHeader,
		tqueue.HandleNewBlockHeader,
	)

	if err := tqueue.NewWorker().Run(mux); err != nil {
		log.WithError(err).Error("failed to run queue worker")
	}
}
