package repository

import (
	"database-concurrency/ent"
	serviceproviderRepo "database-concurrency/internal/repository/serviceprovider"
	ticketRepo "database-concurrency/internal/repository/ticket"
	ticketeventRepo "database-concurrency/internal/repository/ticketevent"
	transactionRepo "database-concurrency/internal/repository/transaction"
	userRepo "database-concurrency/internal/repository/user"
)

type db struct {
	pg          *ent.Client
	transaction transactionRepo.Repository

	// Booking service
	user            userRepo.Repository
	serviceProvider serviceproviderRepo.Repository
	ticket          ticketRepo.Repository
	ticketEvent     ticketeventRepo.Repository
}

func (d *db) Pg() *ent.Client {
	return d.pg
}

func (d *db) Transaction() transactionRepo.Repository {
	return d.transaction
}

func (d *db) User() userRepo.Repository {
	return d.user
}

func (d *db) ServiceProvider() serviceproviderRepo.Repository {
	return d.serviceProvider
}

func (d *db) Ticket() ticketRepo.Repository {
	return d.ticket
}

func (d *db) TicketEvent() ticketeventRepo.Repository {
	return d.ticketEvent
}
