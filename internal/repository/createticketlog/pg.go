package createticketlog

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/ent/createticketlog"

	"github.com/google/uuid"
)

type pg struct {
	client *ent.CreateTicketLogClient
}

func (pg pg) Create(ctx context.Context, createTicketLog *ent.CreateTicketLog) (*ent.CreateTicketLog, error) {
	return pg.client.Create().
		SetID(createTicketLog.ID).
		SetUniqueID(createTicketLog.UniqueID).
		Save(ctx)
}

func (pg pg) GetByUniqueID(ctx context.Context, uniqueID uuid.UUID) (*ent.CreateTicketLog, error) {
	return pg.client.Query().Where(createticketlog.UniqueID(uniqueID)).Only(ctx)
}
