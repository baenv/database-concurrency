package ticketevent

import (
	"context"
	"database-concurrency/ent"
)

type pg struct {
	client *ent.TicketEventClient
}

func (pg pg) Create(ctx context.Context, ticketEvent *ent.TicketEvent) (*ent.TicketEvent, error) {
	return pg.client.Create().
		SetTicketID(ticketEvent.TicketID).
		SetUserID(ticketEvent.UserID).
		SetType(ticketEvent.Type).
		Save(ctx)
}
