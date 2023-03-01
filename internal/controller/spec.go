package controller

import (
	"database-concurrency/config"
	"database-concurrency/internal/controller/gen"
	"database-concurrency/internal/controller/ticket"
	"database-concurrency/internal/controller/transaction"
	"database-concurrency/internal/repository"

	"github.com/sirupsen/logrus"
)

// Controller is the controller interface
type Controller interface {
	TransactionCtrl() transaction.Controller
	TicketCtrl() ticket.Controller
	GenCtrl() gen.Controller
}

// New is used to create new controller
func New(repo repository.Repository, log *logrus.Logger, cfg config.Config) Controller {
	return controller{
		transactionCtrl: transaction.New(repo, log),
		ticketCtrl:      ticket.New(repo, log, cfg),
		genCtrl:         gen.New(repo, log),
	}
}
