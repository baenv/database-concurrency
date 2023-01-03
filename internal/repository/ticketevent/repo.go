package ticketevent

import (
	"database-concurrency/ent"
)

// Repository is used to interact with ticket_events table
type Repository interface {
}

// NewPGRepo is used to generate new TicketEvent repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.TicketEvent,
	}
}
