package handler

import (
	"database-concurrency/internal/controller"
	"database-concurrency/internal/repository"

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

	// GenTicketID gen new ticket id for creating ticket, id gen API ONLY
	GenTicketID(ctx echo.Context) error
}

// New is used to create new controller
func New(repo repository.Repository, log *logrus.Logger) Handler {
	return handler{
		ctrl: controller.New(repo, log),
	}
}
