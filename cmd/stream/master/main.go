package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"database-concurrency/config"
	"database-concurrency/ent"
	"database-concurrency/internal/handler/payload"
	"database-concurrency/internal/redis"

	redislib "github.com/redis/go-redis/v9"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type ticketData struct {
	ConsumerName string   `json:"consumer_name"`
	Events       []string `json:"message_ids"`
}

type consumerData struct {
	HealthURL string   `json:"health_url"`
	TicketIDs []string `json:"ticket_ids"`
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
		masterConsumer := redis.CreateConsumer(ctx, cfg.STREAM_NAME, cfg.STREAM_GROUP_NAME, cfg.STREAM_GROUP_NAME)

		// tickets, consumers used as traffic control info to redirect a message to a consumer that currently processing the same ticket_id
		var tickets *map[string]ticketData
		var consumers *map[string]consumerData

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

		service := masterService{
			redis: redis,

			tickets:   tickets,
			consumers: consumers,

			ticketTable:   ticketTable,
			consumerTable: consumerTable,
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
				messages, err := masterConsumer.Read(ctx, 1, cfg.STREAM_NAME, ">")
				if err != nil {
					log.Fatal("fail to read ticket stream")
					return
				}
				if len(messages) == 0 {
					log.Warningln("no message, sleep for 1s")
					time.Sleep(time.Second)
					continue
				}

				for _, message := range messages {
					// check if a ticket_ID is being process by some consumers. If not then choose a random consumer
					err := service.delegateMessage(ctx, message, "")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}()

		go func() {
			for {
				hasOnlineConsumer := service.recoverDownConsumers(ctx)
				if hasOnlineConsumer {
					service.recoverIdleEvent(ctx)
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

			consumerTable: consumerTable,
			ticketTable:   ticketTable,
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
	tickets   *map[string]ticketData
	consumers *map[string]consumerData

	ticketTable   redis.RedisTableWrapper
	consumerTable redis.RedisTableWrapper
}

// register, newly started consumer will call this to register itself to the consumers map
func (h masterHandler) register(ctx echo.Context) error {
	var req payload.ConsumerRegisterRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	existConsumer, ok := (*h.consumers)[req.ConsumerName]
	if !ok {
		(*h.consumers)[req.ConsumerName] = consumerData{
			HealthURL: req.HealthURL,
		}
	} else {
		(*h.consumers)[req.ConsumerName] = consumerData{
			HealthURL: req.HealthURL,
			TicketIDs: existConsumer.TicketIDs,
		}
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
	messageIDs := []string{}
	for _, messageID := range t.Events {
		if messageID != req.MessageID {
			messageIDs = append(messageIDs, messageID)
		}
	}

	if len(messageIDs) == 0 {
		delete((*h.tickets), req.TicketID)
		(*h.consumers)[t.ConsumerName] = consumerData{
			HealthURL: (*h.consumers)[t.ConsumerName].HealthURL,
			TicketIDs: []string{},
		}
	} else {
		(*h.tickets)[req.TicketID] = ticketData{
			ConsumerName: t.ConsumerName,
			Events:       messageIDs,
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

type masterService struct {
	cfg config.Config

	redis redis.RedisWrapper

	tickets   *map[string]ticketData
	consumers *map[string]consumerData

	ticketTable   redis.RedisTableWrapper
	consumerTable redis.RedisTableWrapper
}

func (s masterService) delegateMessage(ctx context.Context, message redislib.XMessage, consumerName string) error {
	data := []byte(message.Values["data"].(string))
	messageID := message.ID

	var event ent.TicketEvent

	err := json.Unmarshal(data, &event)
	if err != nil {
		return errors.New(fmt.Sprintf("fail to parse message data, %v", err))
	}

	ticketID := event.TicketID.String()

	var consumerProcessingTicketIDs []string
	var consumerProcessingMessageIDs []string

	t, isExist := (*s.tickets)[ticketID]
	if isExist {
		consumerName = t.ConsumerName
		consumerProcessingTicketIDs = []string{ticketID}
		consumerProcessingMessageIDs = append(t.Events, messageID)
	} else if consumerName == "" {
		magicNumber := rand.Intn(len(*s.consumers))
		count := 0
		for key, value := range *s.consumers {
			if count == magicNumber {
				consumerName = key
				consumerProcessingTicketIDs = append(value.TicketIDs, ticketID)
				consumerProcessingMessageIDs = []string{messageID}
				break
			}
			count++
		}
	}

	consumer := s.redis.CreateConsumer(ctx, s.cfg.STREAM_NAME, s.cfg.STREAM_GROUP_NAME, consumerName)
	err = consumer.Claim(ctx, s.cfg.STREAM_NAME, messageID)
	if err != nil {
		return errors.New(fmt.Sprintf("fail to claim event from ticket_stream, %v", err))
	}

	{
		(*s.tickets)[ticketID] = ticketData{
			ConsumerName: consumerName,
			Events:       consumerProcessingMessageIDs,
		}
		(*s.consumers)[consumerName] = consumerData{
			HealthURL: (*s.consumers)[consumerName].HealthURL,
			TicketIDs: consumerProcessingTicketIDs,
		}

		c, _ := json.Marshal(s.consumers)
		err = s.consumerTable.SetValue(ctx, c)
		if err != nil {
			return errors.New(fmt.Sprintf("can not save consumers table, %v", err))

		}
		t, _ := json.Marshal(s.tickets)
		err = s.ticketTable.SetValue(ctx, t)
		if err != nil {
			return errors.New(fmt.Sprintf("can not save tickets table, %v", err))
		}
	}
	return nil
}

func (s masterService) recoverDownConsumers(ctx context.Context) bool {
	defer time.Sleep(time.Minute)

	// check consumer health
	offlineConsumers := map[string]consumerData{}

	for consumer, data := range *s.consumers {
		_, err := http.Get(data.HealthURL)
		if err != nil {
			fmt.Printf("consumer offline: %s, healthURL: %v\n", consumer, data.HealthURL)
			offlineConsumers[consumer] = data
		}
	}

	if len(offlineConsumers) <= 0 {
		return len(*s.consumers) > 0
	}

	if len(*s.consumers)-len(offlineConsumers) <= 0 {
		fmt.Printf("no online consumer to delegate, sleep for 2m\n")
		time.Sleep(2 * time.Minute)
		return false
	}

	for consumer := range offlineConsumers {
		delete((*s.consumers), consumer)
	}

	for _, data := range offlineConsumers {
		delegateConsumerRoundRobinCount := 0
		for _, ticketID := range data.TicketIDs {
			var delegateConsumer *string
			{
				magicNumber := rand.Intn(len(*s.consumers))
				for key := range *s.consumers {
					if delegateConsumerRoundRobinCount == magicNumber {
						delegateConsumer = &key
						break
					}
					delegateConsumerRoundRobinCount++
				}
			}
			if delegateConsumer == nil {
				continue
			}

			consumer := s.redis.CreateConsumer(ctx, s.cfg.STREAM_NAME, s.cfg.STREAM_GROUP_NAME, *delegateConsumer)
			for _, messageID := range (*s.tickets)[ticketID].Events {
				err := consumer.Claim(ctx, s.cfg.STREAM_NAME, messageID)
				if err != nil {
					fmt.Printf("fail to claim event from ticket_stream, %v, %v", messageID, err)
				}
			}

			(*s.tickets)[ticketID] = ticketData{
				ConsumerName: *delegateConsumer,
				Events:       (*s.tickets)[ticketID].Events,
			}
			(*s.consumers)[*delegateConsumer] = consumerData{
				HealthURL: (*s.consumers)[*delegateConsumer].HealthURL,
				TicketIDs: append((*s.consumers)[*delegateConsumer].TicketIDs, ticketID),
			}
		}
	}

	// save to redis table
	{
		c, _ := json.Marshal(s.consumers)
		err := s.consumerTable.SetValue(ctx, c)
		if err != nil {
			log.Fatal("can not save consumers table")
		}
		t, _ := json.Marshal(s.tickets)
		err = s.ticketTable.SetValue(ctx, t)
		if err != nil {
			log.Fatal("can not save tickets table")
		}
	}

	return true
}

func (s masterService) recoverIdleEvent(ctx context.Context) {
	defer time.Sleep(time.Minute)

	masterConsumer := s.redis.CreateConsumer(ctx, s.cfg.STREAM_NAME, s.cfg.STREAM_GROUP_NAME, s.cfg.MASTER_CONSUMER_NAME)
	messages, err := masterConsumer.AutoClaim(ctx, s.cfg.STREAM_NAME, 5*time.Minute)
	if err != nil {
		fmt.Printf("fail to auto claim event from ticket_stream, %v", err)
		return
	}
	if len(messages) == 0 {
		return
	}

	for _, message := range messages {
		// check if message is processing by any consumer
		var consumerName string

		{
			var event ent.TicketEvent

			data := []byte(message.Values["data"].(string))
			err := json.Unmarshal(data, &event)
			if err != nil {
				log.Fatalln("fail to parse message data")
			}

			ticketID := event.TicketID.String()

			ticket, ok := (*s.tickets)[ticketID]
			if ok {
				consumer := (*s.consumers)[ticket.ConsumerName]
				_, err := http.Get(consumer.HealthURL)
				if err != nil {
					fmt.Printf("consumer offline: %s, healthURL: %v\n", ticket.ConsumerName, consumer.HealthURL)
				} else {
					consumerName = ticket.ConsumerName
				}
			}
		}

		//
		err = s.delegateMessage(ctx, message, consumerName)
		if err != nil {
			log.Fatal("can not parse ticket event")
			return
		}
	}
}
