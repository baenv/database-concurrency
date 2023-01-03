package serviceprovider

import (
	"database-concurrency/ent"
)

// Repository is used to interact with service_providers table
type Repository interface {
}

// NewPGRepo is used to generate new ServiceProvider repository
func NewPGRepo(client *ent.Client) Repository {
	return &pg{
		client: client.ServiceProdiver,
	}
}
