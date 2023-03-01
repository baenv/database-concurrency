package ticket

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/ent/ticket"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

type pg struct {
	client *ent.TicketClient
}

// One retrieves a ticket by its ID
func (pg pg) One(ctx context.Context, id uuid.UUID) (*ent.Ticket, error) {
	return pg.client.Query().Where(ticket.ID(id)).Only(ctx)
}

// OneForUpdate retrieves a ticket by its ID with an exclusive lock
func (pg pg) OneForUpdate(ctx context.Context, id uuid.UUID) (*ent.Ticket, error) {
	return pg.client.Query().Where(ticket.ID(id)).ForUpdate(
		sql.WithLockAction(sql.NoWait),
	).Only(ctx)
}

// Update updates a ticket
func (pg pg) Update(ctx context.Context, ticket *ent.Ticket) (*ent.Ticket, error) {
	return pg.client.UpdateOne(ticket).
		SetLastEventID(ticket.LastEventID).
		SetStatus(ticket.Status).
		SetVersions(ticket.Versions).
		Save(ctx)
}

// Update updates a ticket
func (pg pg) Create(ctx context.Context, ticket *ent.Ticket) (*ent.Ticket, error) {
	return pg.client.Create().
		SetID(ticket.ID).
		SetUserID(ticket.UserID).
		SetStatus(ticket.Status).
		SetVersions(ticket.Versions).
		Save(ctx)
}
