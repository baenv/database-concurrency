package ticket

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ticket struct {
	repo repository.Repositoy
	log  *logrus.Logger
}

func (t ticket) Book(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error) {
	// Get ticket

	// Check if ticket is owned by given user

	// Create new ticket event

	// Update ticket

	// Return value

	var result ent.Ticket
	return &result, repository.WithTx(ctx, t.repo.Pg(), func(txRepo repository.Repositoy) error {
		ticket, err := txRepo.Ticket().One(ctx, ticketID)
		if err != nil {
			t.log.Error("failed to get ticket", err)
			return err
		}

		if ticket.UserID.String() != userID.String() {
			t.log.Error("ticket is not owned by given user")
			return errors.New("ticket is not owned by given user")
		}

		if ticket.Status != "reserved" {
			t.log.Error("ticket is not reserved")
			return errors.New("ticket is not reserved")
		}

		ticketEvent, err := txRepo.TicketEvent().Create(ctx, &ent.TicketEvent{
			TicketID: ticketID,
			UserID:   userID,
			Type:     "book",
		})
		if err != nil {
			t.log.Error("failed to create ticket event", err)
			return err
		}

		ticket.Status = "booked"
		ticket.LastEventID = ticketEvent.ID
		ticket, err = txRepo.Ticket().Update(ctx, ticket)
		if err != nil {
			t.log.Error("failed to update ticket", err)
			return err
		}

		ticket.Edges.LastEvent = ticketEvent
		result = *ticket
		return nil
	})
}
