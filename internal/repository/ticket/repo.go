package ticket

import (
	"database-concurrency/ent"
)

// Repository is used to interact with tickets table
type Repository interface {
}

// NewPGRepo is used to generate new Ticket repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.Ticket,
	}
}
