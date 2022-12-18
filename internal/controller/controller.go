package controller

import (
	"database-concurrency/internal/controller/transaction"
)

type controller struct {
	transactionCtrl transaction.Controller
}

func (c controller) TransactionCtrl() transaction.Controller {
	return c.transactionCtrl
}
