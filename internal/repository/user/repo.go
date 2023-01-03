package user

import (
	"database-concurrency/ent"
)

// Repository is used to interact with users table
type Repository interface {
}

// NewPGRepo is used to generate new User repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.User,
	}
}
