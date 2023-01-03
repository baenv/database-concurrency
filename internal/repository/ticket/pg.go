package ticket

import (
	"database-concurrency/ent"
)

type pg struct {
	client *ent.TicketClient
}
