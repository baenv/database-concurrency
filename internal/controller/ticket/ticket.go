package ticket

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/controller/utils"
	"database-concurrency/internal/repository"
	"database-concurrency/internal/transducer"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ticket struct {
	repo repository.Repositoy
	log  *logrus.Logger
}

func (t ticket) Book(ctx context.Context, ticketID, userID uuid.UUID, locks utils.Locks) (*ent.Ticket, error) {
	// Get ticket

	// Check if ticket is owned by given user

	// Create new ticket event

	// Update ticket

	// Return value

	if locks.SessionAdvisoryLock {
		if err := t.repo.AdvisoryLockTable("tickets"); err != nil {
			t.log.Error("failed to acquire advisory lock", err)
			return nil, err
		}

		defer func() {
			if err := t.repo.AdvisoryUnlockTable("tickets"); err != nil {
				t.log.Error("failed to release advisory lock", err)
			}
		}()
	}

	var result ent.Ticket
	return &result, repository.WithTx(ctx, t.repo.Raw(), t.repo.Pg(), func(txRepo repository.Repositoy) error {
		var (
			ticket *ent.Ticket
			err    error
		)

		// Handle Lock For Update Flag
		if locks.ForUpdate {
			ticket, err = txRepo.Ticket().OneForUpdate(ctx, ticketID)
		} else {
			ticket, err = txRepo.Ticket().One(ctx, ticketID)
		}

		if err != nil {
			t.log.Error("failed to get ticket", err)
			return err
		}

		if ticket.UserID.String() != userID.String() {
			t.log.Error("ticket is not owned by given user")
			return errors.New("ticket is not owned by given user")
		}

		if ticket.Status != transducer.Reserved.String() {
			t.log.Error("ticket is not reserved")
			return errors.New("ticket is not reserved")
		}

		config, machine := transducer.NewBookingMachine(ticket.Status)
		outputs := machine.Transduce(config, transducer.Book)

		state := outputs.GetState()
		if state == transducer.Invalid {
			t.log.Error("machine has reached an invalid state")
			return errors.New("machine has reached an invalid state")
		}
		if state != transducer.Booked {
			t.log.Error("machine has failed to transition to booked")
			return errors.New("machine has failed to transition to booked")
		}

		var ticketEvent *ent.TicketEvent

		for _, effect := range outputs.Effects {
			switch effect {
			case transducer.CreateBookingEvent:
				ticketEvent, err = txRepo.TicketEvent().Create(ctx, &ent.TicketEvent{
					TicketID: ticketID,
					UserID:   userID,
					Type:     transducer.Book.String(),
				})
				if err != nil {
					t.log.Error("failed to create ticket event", err)
					return err
				}

			case transducer.UpdateBookingStatus:
				ticket.Status = state.String()
				ticket.LastEventID = ticketEvent.ID
				ticket, err = txRepo.Ticket().Update(ctx, ticket)
				if err != nil {
					t.log.Error("failed to update ticket", err)
					return err
				}

			case transducer.EmailUser:
				// TODO: create mock for email user
				t.log.Info("email user")

			case transducer.CallClient:
				// TODO: create mock for call client
				t.log.Info("call client")

			case transducer.SMSUser:
				// TODO: create mock for SMS user
				t.log.Info("SMS user")

			default:
				t.log.Error("not all effects have been covered")
				return errors.New("not all effects have been covered")
			}
		}

		ticket.Edges.LastEvent = ticketEvent
		result = *ticket
		return nil
	})
}

// Reserve the ticket
func (t ticket) Reserve(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error) {
	var result ent.Ticket
	ticket, err := t.repo.Ticket().One(ctx, ticketID)

	if err != nil {
		t.log.Error("failed to get ticket", err)
		return nil, err
	}

	if ticket.Status != transducer.Idle.String() {
		t.log.Error("ticket is not idle")
		return nil, errors.New("ticket is not idle")
	}

	config, ticketTransducer := transducer.NewBookingMachine(ticket.Status)
	output := ticketTransducer.Transduce(config, transducer.Reserve)

	resultState := output.GetState().String()
	if resultState == transducer.Invalid.String() {
		err := errors.New("invalid ticket state")
		t.log.Error("failed to get ticket", err)
		return nil, err
	}

	for _, effect := range output.Effects {
		switch effect.Int() {
		case transducer.UpdateBookingStatus.Int():
			if err := repository.WithTx(ctx, t.repo.Raw(), t.repo.Pg(), func(txRepo repository.Repositoy) error {
				ticketEvent, err := txRepo.TicketEvent().Create(ctx, &ent.TicketEvent{
					TicketID: ticketID,
					UserID:   userID,
					Type:     transducer.Reserve.String(),
				})
				if err != nil {
					t.log.Error("failed to create ticket event", err)
					return err
				}

				// ticket update
				ticket.Status = transducer.Reserved.String()
				ticket.LastEventID = ticketEvent.ID
				version, err := strconv.ParseInt(ticket.Versions, 16, 0)
				if err != nil {
					return err
				}

				ticket.Versions = strconv.FormatInt(version+1, 16)
				ticket, err = txRepo.Ticket().Update(ctx, ticket)
				if err != nil {
					t.log.Error("failed to update ticket", err)
					return err
				}

				ticket.Edges.LastEvent = ticketEvent
				result = *ticket
				return nil
			}); err != nil {
				return nil, err
			}
		case transducer.EmailUser.Int():
			// TODO: email user
		}
	}
	return &result, nil
}

// Cancel the ticket
func (t ticket) Cancel(ctx context.Context, ticketID, userID uuid.UUID) (*ent.Ticket, error) {
	var result ent.Ticket
	ticket, err := t.repo.Ticket().One(ctx, ticketID)

	if err != nil {
		t.log.Error("failed to get ticket", err)
		return nil, err
	}

	if ticket.UserID.String() != userID.String() {
		t.log.Error("ticket is not owned by given user")
		return nil, errors.New("ticket is not owned by given user")
	}

	if ticket.Status != transducer.Reserved.String() {
		t.log.Error("ticket is not reserved")
		return nil, errors.New("ticket is not reserved")
	}

	config, ticketTransducer := transducer.NewBookingMachine(ticket.Status)
	output := ticketTransducer.Transduce(config, transducer.Cancel)

	resultState := output.GetState().String()
	if resultState == transducer.Invalid.String() {
		err := errors.New("invalid ticket state")
		t.log.Error("failed to get ticket", err)
		return nil, err
	}

	for _, effect := range output.Effects {
		switch effect.Int() {
		case transducer.UpdateBookingStatus.Int():
			if err := repository.WithTx(ctx, t.repo.Raw(), t.repo.Pg(), func(txRepo repository.Repositoy) error {
				ticketEvent, err := txRepo.TicketEvent().Create(ctx, &ent.TicketEvent{
					TicketID: ticketID,
					UserID:   userID,
					Type:     transducer.Cancel.String(),
				})
				if err != nil {
					t.log.Error("failed to create ticket event", err)
					return err
				}

				// ticket update
				ticket.Status = resultState
				ticket.LastEventID = ticketEvent.ID
				version, err := strconv.ParseInt(ticket.Versions, 16, 0)
				if err != nil {
					return err
				}

				ticket.Versions = strconv.FormatInt(version+1, 16)
				ticket, err = txRepo.Ticket().Update(ctx, ticket)
				if err != nil {
					t.log.Error("failed to update ticket", err)
					return err
				}

				ticket.Edges.LastEvent = ticketEvent
				result = *ticket
				return nil
			}); err != nil {
				return nil, err
			}
		case transducer.EmailUser.Int():
			// TODO: email user
		}
	}
	return &result, nil
}
