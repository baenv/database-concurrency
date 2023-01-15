package controller

import (
	"database-concurrency/internal/controller/ticket"
	"database-concurrency/internal/controller/transaction"
)

type controller struct {
	transactionCtrl transaction.Controller
	ticketCtrl      ticket.Controller
}

func (c controller) TransactionCtrl() transaction.Controller {
	return c.transactionCtrl
}

func (c controller) TicketCtrl() ticket.Controller {
	return c.ticketCtrl
}
