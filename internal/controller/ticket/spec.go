package ticket

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Book(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error)
	// Reserve the ticket
	Reserve(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error)
}

func New(repo repository.Repositoy, log *logrus.Logger) Controller {
	return ticket{
		repo: repo,
		log:  log,
	}
}
