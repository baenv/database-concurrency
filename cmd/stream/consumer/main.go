package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"database-concurrency/config"
	"database-concurrency/ent"
	"database-concurrency/internal/handler/payload"
	"database-concurrency/internal/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	data, _ := json.Marshal(payload.ConsumerRegisterRequest{
		ConsumerName: cfg.CONSUMER_NAME,
		HealthURL:    fmt.Sprintf("http://localhost:%v/healthz", cfg.CONSUMER_PORT),
	})
	bodyReader := bytes.NewReader(data)

	masterConsumerRegisterURL := fmt.Sprintf("%v/api/v1/consumers/register", cfg.MASTER_URL)
	_, err = http.Post(masterConsumerRegisterURL, "application/json", bodyReader)
	if err != nil {
		log.Fatal(err)
		return
	}

	// _, err = repository.Init(cfg)
	// if err != nil {
	// 	log.WithError(err).Error("failed to init repo")
	// 	return
	// }

	redis, err := redis.InitRedisClient(cfg)
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
		// c := controller.New(repo, nil, log, cfg)

		consumer := redis.CreateConsumer(ctx, "concurrency_stream", "concurrency_stream_group", cfg.CONSUMER_NAME)

		go func() {
			mockFailureCount := 2
			for {
				if mockFailureCount == 0 {
					log.Infoln("mock failure, shutting down consumer")
					time.Sleep(time.Minute)
					log.Fatal("consumer disconnected")
				}

				messages, err := consumer.Read(ctx, 1, "ticket_stream", "0")
				if err != nil {
					log.Fatal(err)
				}
				if len(messages) == 0 {
					time.Sleep(500 * time.Millisecond)
					continue
				}

				data := []byte(messages[0].Values["data"].(string))
				messageID := messages[0].ID

				var event ent.TicketEvent

				err = json.Unmarshal(data, &event)
				if err != nil {
					log.Fatal("can not parse ticket event")
					return
				}
				log.Infof("processing eventID: %s, ticketID: %s", messageID, event.TicketID)

				// TODO: handling of ticket event
				time.Sleep(time.Minute)

				consumer.Acknowledge(ctx, "ticket_stream", messageID)

				data, _ = json.Marshal(payload.TicketAcknowledgeRequest{
					TicketID:  event.TicketID.String(),
					MessageID: messageID,
				})
				bodyReader := bytes.NewReader(data)

				masterTicketAcknowledgeURL := fmt.Sprintf("%v/api/v1/tickets/acknowledge", cfg.MASTER_URL)
				_, err = http.Post(masterTicketAcknowledgeURL, "application/json", bodyReader)
				if err != nil {
					log.Fatal(err)
					return
				}

				mockFailureCount--
			}
		}()
	}
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/healthz", healthz)

	if err := e.Start(fmt.Sprintf(":%s", cfg.CONSUMER_PORT)); err != nil {
		log.WithError(err).Error("failed to start server")
	}
}

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}
