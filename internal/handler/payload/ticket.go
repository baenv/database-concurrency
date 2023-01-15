package payload

import (
	"database-concurrency/ent"

	"github.com/google/uuid"
)

type BookRequest struct {
	TicketID uuid.UUID `json:"ticket_id"`
	UserID   uuid.UUID `json:"user_id"` // Mock, actually, it is fetched from JWT
}

type BookResponse struct {
	Ticket *ent.Ticket `json:"ticket"`
}
