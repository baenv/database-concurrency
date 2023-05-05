package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

type ticket struct {
	ConsumerName string   `json:"consumer_name"`
	Events       []string `json:"event_ids"`
}

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	cfg, err := config.NewConfig()
	if err != nil {
		log.WithError(err).Error("failed to load config")
		return
	}

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

	// TODO: should move to a master handler pkg
	{
		masterConsumer := redis.CreateConsumer(ctx, "ticket_stream", "concurrency_stream_group", "master")

		// tickets, consumers used as traffic control info to redirect a message to a consumer that currently processing the same ticket_id
		var tickets *map[string]ticket
		var consumers *map[string][]string

		// For failure recovery, tickets, consumers are created as backup, saved in Redis
		ticketTable := redis.CreateTable("tickets")
		consumerTable := redis.CreateTable("consumers")

		{
			ticketData, err := ticketTable.GetValue(ctx)
			if err != nil {
				log.Errorln("get ticket table err:", err)
				ticketData = "{}"
			}

			err = json.Unmarshal([]byte(ticketData), &tickets)
			if err != nil {
				log.Fatal("fail to parse tickets table data err:", err)
				return
			}

			consumerData, err := consumerTable.GetValue(ctx)
			if err != nil {
				log.Errorln("get ticket table err:", err)
				consumerData = "{}"
			}

			err = json.Unmarshal([]byte(consumerData), &consumers)
			if err != nil {
				log.Fatal("fail to parse consumers table data err:", err)
				return
			}
		}

		// process ticket stream
		go func() {
			for {
				if len(*consumers) == 0 {
					log.Warningln("no consumer online, sleep for 5s")
					time.Sleep(5 * time.Second)
					continue
				}

				// with ">", the master consumer always the first to read the incoming message and able to "XCLAIM" for other consumers
				messages, err := masterConsumer.Read(ctx, 1, "ticket_stream", ">")
				if err != nil {
					log.Fatal("fail to read ticket stream")
					return
				}
				if len(messages) == 0 {
					log.Warningln("no message, sleep for 1s")
					time.Sleep(time.Second)
					continue
				}

				data := []byte(messages[0].Values["data"].(string))
				eventID := messages[0].ID

				var event ent.TicketEvent

				err = json.Unmarshal(data, &event)
				if err != nil {
					log.Fatal("can not parse ticket event")
					return
				}

				ticketID := event.TicketID.String()

				var consumerName string
				var consumerProcessingTicketIDs []string
				var consumerProcessingEventIDs []string

				// check if a ticket_ID is being process by some consumers. If not then choose a random consumer
				{
					t, isExist := (*tickets)[ticketID]
					if isExist {
						consumerName = t.ConsumerName
						consumerProcessingTicketIDs = []string{ticketID}
						consumerProcessingEventIDs = append(t.Events, eventID)
					} else {
						magicNumber := rand.Intn(len(*consumers))
						count := 0
						for key, value := range *consumers {
							if count == magicNumber {
								consumerName = key
								consumerProcessingTicketIDs = append(value, ticketID)
								consumerProcessingEventIDs = []string{eventID}
								break
							}
							count++
						}
					}
				}

				consumer := redis.CreateConsumer(ctx, "ticket_stream", "concurrency_stream_group", consumerName)
				err = consumer.Claim(ctx, "ticket_stream", eventID)
				if err != nil {
					log.Fatal("fail to claim event from ticket_stream, ", err)
					return
				}

				{
					(*tickets)[ticketID] = ticket{
						ConsumerName: consumerName,
						Events:       consumerProcessingEventIDs,
					}
					(*consumers)[consumerName] = consumerProcessingTicketIDs

					c, _ := json.Marshal(consumers)
					err = consumerTable.SetValue(ctx, c)
					if err != nil {
						log.Fatal("can not save consumers table")
					}
					t, _ := json.Marshal(tickets)
					err = ticketTable.SetValue(ctx, t)
					if err != nil {
						log.Fatal("can not save tickets table")
					}
				}
			}
		}()

		// Echo instance
		e := echo.New()

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		// Routes
		e.GET("/healthz", healthz)

		apiV1 := e.Group("/api/v1")

		handler := masterHandler{
			consumers: consumers,
			tickets:   tickets,
		}

		consumerRouter := apiV1.Group("/consumers")
		consumerRouter.Add(http.MethodPost, "/register", handler.register)

		ticketRouter := apiV1.Group("/tickets")
		ticketRouter.Add(http.MethodPost, "/acknowledge", handler.acknowledge)

		if err := e.Start(fmt.Sprintf(":%s", cfg.MASTER_PORT)); err != nil {
			log.WithError(err).Error("failed to start server")
		}
	}
}

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}

// TODO: should move this to a handler pkg
type masterHandler struct {
	tickets   *map[string]ticket
	consumers *map[string][]string

	ticketTable   redis.RedisTableWrapper
	consumerTable redis.RedisTableWrapper
}

// register, newly started consumer will call this to register itself to the consumers map
func (h masterHandler) register(ctx echo.Context) error {
	var req payload.ConsumerRegisterRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	_, ok := (*h.consumers)[req.ConsumerName]
	if !ok {
		(*h.consumers)[req.ConsumerName] = []string{}
	}

	redisCtx := context.Background()

	c, _ := json.Marshal(h.consumers)
	err := h.consumerTable.SetValue(redisCtx, c)
	if err != nil {
		log.Fatal("can not save consumers table")
	}

	return ctx.JSON(200, "OK")
}

// acknowledge, a consumer call this to acknowledge a message is successfully processed, so remove the message from consumers, tickets maps
func (h masterHandler) acknowledge(ctx echo.Context) error {
	var req payload.TicketAcknowledgeRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	t := (*h.tickets)[req.TicketID]
	eventIDs := []string{}
	for _, eventID := range t.Events {
		if eventID != req.EventID {
			eventIDs = append(eventIDs, eventID)
		}
	}

	if len(eventIDs) == 0 {
		delete((*h.tickets), req.TicketID)
		(*h.consumers)[t.ConsumerName] = []string{}
	} else {
		(*h.tickets)[req.TicketID] = ticket{
			ConsumerName: t.ConsumerName,
			Events:       eventIDs,
		}
	}

	redisCtx := context.Background()

	c, _ := json.Marshal(h.consumers)
	err := h.consumerTable.SetValue(redisCtx, c)
	if err != nil {
		log.Fatal("can not save consumers table")
	}
	td, _ := json.Marshal(h.tickets)
	err = h.ticketTable.SetValue(redisCtx, td)
	if err != nil {
		log.Fatal("can not save tickets table")
	}

	return ctx.JSON(200, "OK")
}
