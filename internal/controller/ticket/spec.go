package ticket

import (
	"context"
	"database-concurrency/config"
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
	// Create the ticket
	Create(ctx context.Context) (*ent.Ticket, error)
	// CreateV2 Create the ticket from unique id
	CreateV2(ctx context.Context, unique_id uuid.UUID) (*ent.Ticket, error)
	// CheckIn the ticket
	CheckIn(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error)
	// CheckOut the ticket
	CheckOut(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error)
}

func New(repo repository.Repository, log *logrus.Logger, cfg config.Config) Controller {
	return ticket{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}
