package ticket

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/controller/utils"
	"database-concurrency/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Book(ctx context.Context, ticketID, userID uuid.UUID, locks utils.Locks) (*ent.Ticket, error)
	// Reserve the ticket
	Reserve(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error)
	// Cancel the ticket
	Cancel(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error)
}

func New(repo repository.Repositoy, log *logrus.Logger) Controller {
	return ticket{
		repo: repo,
		log:  log,
	}
}
