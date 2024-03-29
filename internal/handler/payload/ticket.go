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

type ReserveResponseV2 struct {
	Data string `json:"data"`
}
type CancelRequest struct {
	TicketID uuid.UUID `json:"ticket_id"`
	UserID   uuid.UUID `json:"user_id"` // Mock, actually, it is fetched from JWT
}
type CancelResponse struct {
	Ticket *ent.Ticket `json:"ticket"`
}

type GenTicketIDRequest struct {
	UniqueID uuid.UUID `json:"unique_id"`
}

type GenTicketIDResponse struct {
	TicketID string `json:"ticket_id"`
}

type CreateTicketRequest struct {
	UniqueID uuid.UUID `json:"unique_id"`
}

type CreateTicketResponse struct {
	Ticket *ent.Ticket `json:"ticket"`
}
