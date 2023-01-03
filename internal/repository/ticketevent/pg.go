package ticketevent

import (
	"database-concurrency/ent"
)

type pg struct {
	client *ent.TicketEventClient
}
