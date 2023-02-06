package ticket

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/ent/ticket"

	"github.com/google/uuid"
)

type pg struct {
	client *ent.TicketClient
}

func (pg pg) One(ctx context.Context, id uuid.UUID) (*ent.Ticket, error) {
	return pg.client.Query().Where(ticket.ID(id)).Only(ctx)
}

func (pg pg) Update(ctx context.Context, ticket *ent.Ticket) (*ent.Ticket, error) {
	return pg.client.UpdateOne(ticket).
		SetLastEventID(ticket.LastEventID).
		SetStatus(ticket.Status).
		SetVersions(ticket.Versions).
		Save(ctx)
}
