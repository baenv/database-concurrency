package repository

import (
	"database-concurrency/ent"
	createticketlogRepo "database-concurrency/internal/repository/createticketlog"
	serviceproviderRepo "database-concurrency/internal/repository/serviceprovider"
	ticketRepo "database-concurrency/internal/repository/ticket"
	ticketeventRepo "database-concurrency/internal/repository/ticketevent"
	transactionRepo "database-concurrency/internal/repository/transaction"
	userRepo "database-concurrency/internal/repository/user"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type db struct {
	pg          *ent.Client
	raw         *sql.DB
	transaction transactionRepo.Repository

	// Booking service
	user            userRepo.Repository
	serviceProvider serviceproviderRepo.Repository
	ticket          ticketRepo.Repository
	ticketEvent     ticketeventRepo.Repository
	createTicketLog createticketlogRepo.Repository
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

func (d *db) CreateTicketLog() createticketlogRepo.Repository {
	return d.createTicketLog
}

func (d *db) Raw() *sql.DB {
	return d.raw
}

func (d *db) AdvisoryLockTable(table string) error {
	var oid int64
	if err := d.raw.QueryRow(fmt.Sprintf(`
		SELECT oid
		FROM pg_class
		WHERE relname = '%s'
	`, table)).Scan(&oid); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get oid for table %s", table))
	}

	var success bool
	if err := d.raw.QueryRow(fmt.Sprintf(`
		SELECT pg_try_advisory_lock(%d)
	`, oid)).Scan(&success); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get advisory lock for table %s", table))
	}

	if !success {
		return errors.New(fmt.Sprintf("have no capacity to get advisory lock for table %s", table))
	}

	return nil
}

func (d *db) AdvisoryUnlockTable(table string) error {
	var oid int64
	if err := d.raw.QueryRow(fmt.Sprintf(`
		SELECT oid
		FROM pg_class
		WHERE relname = '%s'
	`, table)).Scan(&oid); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get oid for table %s", table))
	}

	if _, err := d.raw.Exec(fmt.Sprintf(`
		SELECT pg_advisory_unlock(%d)
	`, oid)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to unlock advisory lock for table %s", table))
	}

	return nil
}
