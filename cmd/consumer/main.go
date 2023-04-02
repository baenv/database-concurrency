package main

import (
	"context"
	"encoding/json"
	"net/http"

	"database-concurrency/config"
	"database-concurrency/ent"
	"database-concurrency/internal/controller"
	"database-concurrency/internal/repository"
	"database-concurrency/internal/stream"

	"github.com/labstack/echo/v4"
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

	redis, err := stream.InitRedisClient(cfg)
	if err != nil {
		log.WithError(err).Error("failed to init redis")
		return
	}

	ctx := context.Background()
	_, err = redis.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: should move to a consumer handler pkg
	{
		c := controller.New(repo, nil, log, cfg)

		reserveStream := redis.CreateStream("reserve_ticket", 10)

		for {
			reserveStream.Consume(ctx, 0)
			e := <-reserveStream.MessageChannel()

			j, _ := json.Marshal(e)

			var event ent.TicketEvent

			err := json.Unmarshal(j, &event)
			if err != nil {
				log.Fatal("can not parse ticket event")
				return
			}

			c.TicketCtrl().ConsumeReserve(ctx, event)
		}
	}
}

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}
