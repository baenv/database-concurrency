package ticketevent

import (
	"context"
	"database-concurrency/ent"
)

// Repository is used to interact with ticket_events table
type Repository interface {
	Create(ctx context.Context, ticketEvent *ent.TicketEvent) (*ent.TicketEvent, error)
}

// NewPGRepo is used to generate new TicketEvent repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.TicketEvent,
	}
}
