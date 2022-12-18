package repository

import (
	"database-concurrency/ent"
	transactionRepo "database-concurrency/internal/repository/transaction"
)

type db struct {
	pg          *ent.Client
	transaction transactionRepo.Repository
}

func (d *db) Pg() *ent.Client {
	return d.pg
}

func (d *db) Transaction() transactionRepo.Repository {
	return d.transaction
}
