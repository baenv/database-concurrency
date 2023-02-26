package payload

import (
	"database-concurrency/ent"
	"database-concurrency/internal/controller/utils"

	"github.com/google/uuid"
)

type BookRequest struct {
	TicketID uuid.UUID   `json:"ticket_id"`
	UserID   uuid.UUID   `json:"user_id"` // Mock, actually, it is fetched from JWT
	Locks    utils.Locks `json:"locks"`
}

type BookResponse struct {
	Ticket *ent.Ticket `json:"ticket"`
}

type ReserveRequest struct {
	TicketID uuid.UUID `json:"ticket_id"`
	UserID   uuid.UUID `json:"user_id"` // Mock, actually, it is fetched from JWT
}

type ReserveResponse struct {
	Ticket *ent.Ticket `json:"ticket"`
}
type CancelRequest struct {
	TicketID uuid.UUID `json:"ticket_id"`
	UserID   uuid.UUID `json:"user_id"` // Mock, actually, it is fetched from JWT
}
type CancelResponse struct {
	Ticket *ent.Ticket `json:"ticket"`
}
