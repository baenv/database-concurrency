package controller

import (
	"database-concurrency/internal/controller/transaction"
	"database-concurrency/internal/repository"

	"github.com/sirupsen/logrus"
)

// Controller is the controller interface
type Controller interface {
	TransactionCtrl() transaction.Controller
}

// New is used to create new controller
func New(repo repository.Repositoy, log *logrus.Logger) Controller {
	return controller{
		transactionCtrl: transaction.New(repo, log),
	}
}
