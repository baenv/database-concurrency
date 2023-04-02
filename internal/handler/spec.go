package handler

import (
	"database-concurrency/config"
	"database-concurrency/internal/controller"
	"database-concurrency/internal/repository"
	"database-concurrency/internal/stream"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	Transaction(ctx echo.Context) error

	Book(ctx echo.Context) error
	// Reserve reserve a ticket
	Reserve(ctx echo.Context) error
	// Cancel cancel a ticket
	Cancel(ctx echo.Context) error
	// Create create a ticket
	Create(ctx echo.Context) error
	// Checkin a booked ticket
	CheckIn(ctx echo.Context) error
	// Checkout from a checked-in ticket
	CheckOut(ctx echo.Context) error

	// CreateV2 create a ticket from a unique id
	CreateV2(ctx echo.Context) error
	// ReserveV2 publish an event to reserve a ticket
	ReserveV2(ctx echo.Context) error

	// GenTicketID gen new ticket id for creating ticket, id gen API ONLY
	GenTicketID(ctx echo.Context) error
}

// New is used to create new controller
func New(repo repository.Repository, redis stream.RedisWrapper, log *logrus.Logger, cfg config.Config) Handler {
	return handler{
		ctrl: controller.New(repo, redis, log, cfg),
	}
}
