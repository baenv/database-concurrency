package transaction

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/ent/transaction"
)

type pg struct {
	client *ent.TransactionClient
}

func (pg pg) OneByHash(ctx context.Context, hash string) (*ent.Transaction, error) {
	return pg.client.Query().Where(transaction.HashEqualFold(hash)).Only(ctx)
}
