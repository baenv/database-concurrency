package transaction

import (
	"context"
	"database-concurrency/ent"
)

// Repository is used to interact with transaction table
type Repository interface {
	OneByHash(ctx context.Context, hash string) (*ent.Transaction, error)
}

// NewPGRepo is used to generate new Transaction repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.Transaction,
	}
}
