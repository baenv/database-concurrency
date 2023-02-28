package createticketlog

import (
	"context"
	"database-concurrency/ent"

	"github.com/google/uuid"
)

// Repository is used to interact with ticket_events table
type Repository interface {
	Create(ctx context.Context, createTicketEvent *ent.CreateTicketLog) (*ent.CreateTicketLog, error)
	GetByUniqueID(ctx context.Context, uniqueID uuid.UUID) (*ent.CreateTicketLog, error)
}

// NewPGRepo is used to generate new TicketEvent repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.CreateTicketLog,
	}
}
