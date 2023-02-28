package gen

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/repository"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type gen struct {
	repo repository.Repository
	log  *logrus.Logger
}

func (g gen) CreateTicketID(ctx context.Context, uniqueID uuid.UUID) (string, error) {
	ticketID := uuid.New()

	log, err := g.repo.CreateTicketLog().Create(ctx, &ent.CreateTicketLog{
		ID:       ticketID,
		UniqueID: uniqueID,
	})
	if err != nil {
		if sqlgraph.IsConstraintError(err) {
			if err.Error() == "ent: constraint failed: pq: duplicate key value violates unique constraint \"create_ticket_logs_unique_id_key\"" {
				log, err = g.repo.CreateTicketLog().GetByUniqueID(ctx, uniqueID)
				if err != nil {
					return "", err
				}
				return log.ID.String(), nil
			}
		}
		return "", err
	}
	return log.ID.String(), nil
}
