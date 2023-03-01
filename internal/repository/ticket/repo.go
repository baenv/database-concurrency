package ticket

import (
	"context"
	"database-concurrency/ent"

	"github.com/google/uuid"
)

// Repository is used to interact with tickets table
type Repository interface {
	One(ctx context.Context, id uuid.UUID) (*ent.Ticket, error)
	OneForUpdate(ctx context.Context, id uuid.UUID) (*ent.Ticket, error)
	Update(ctx context.Context, ticket *ent.Ticket) (*ent.Ticket, error)
	Create(ctx context.Context, ticket *ent.Ticket) (*ent.Ticket, error)
}

// NewPGRepo is used to generate new Ticket repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.Ticket,
	}
}
