package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"database-concurrency/config"
	"database-concurrency/ent"
	serviceproviderRepo "database-concurrency/internal/repository/serviceprovider"
	ticketRepo "database-concurrency/internal/repository/ticket"
	ticketeventRepo "database-concurrency/internal/repository/ticketevent"
	transactionRepo "database-concurrency/internal/repository/transaction"
	userRepo "database-concurrency/internal/repository/user"

	"ariga.io/atlas/sql/migrate"
	atlas "ariga.io/atlas/sql/schema"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

// The general Repository that contains all other model-repositories
type Repositoy interface {
	Pg() *ent.Client
	Raw() *sql.DB
	Transaction() transactionRepo.Repository
	User() userRepo.Repository
	ServiceProvider() serviceproviderRepo.Repository
	Ticket() ticketRepo.Repository
	TicketEvent() ticketeventRepo.Repository
	AdvisoryLockTable(table string) error
	AdvisoryUnlockTable(table string) error
}

// New is used to create new repository
func New(rawDB *sql.DB, pg *ent.Client) (Repositoy, error) {
	return &db{
		raw:         rawDB,
		pg:          pg,
		transaction: transactionRepo.NewPGRepo(pg),

		// Booking service
		user:            userRepo.NewPGRepo(pg),
		serviceProvider: serviceproviderRepo.NewPGRepo(pg),
		ticket:          ticketRepo.NewPGRepo(pg),
		ticketEvent:     ticketeventRepo.NewPGRepo(pg),
	}, nil
}

func Init(cfg config.Config) (Repositoy, error) {
	dns := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUSER,
		cfg.DBPASS,
		cfg.DBHOST,
		cfg.DBPORT,
		cfg.DBNAME,
	)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}

	driver := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))

	err = client.Schema.Create(context.Background(), // Hook into Atlas Diff process.
		schema.WithDiffHook(func(next schema.Differ) schema.Differ {
			return schema.DiffFunc(func(current, desired *atlas.Schema) ([]atlas.Change, error) {
				// Before calculating changes.
				changes, err := next.Diff(current, desired)
				if err != nil {
					return nil, err
				}
				// After diff, you can filter
				// changes or return new ones.
				return changes, nil
			})
		}),
		// Hook into Atlas Apply process.
		schema.WithApplyHook(func(next schema.Applier) schema.Applier {
			return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *migrate.Plan) error {
				// Example to hook into the apply process, or implement
				// a custom applier. For example, write to a file.
				//
				//  for _, c := range plan.Changes {
				//      fmt.Printf("%s: %s", c.Comment, c.Cmd)
				//      if err := conn.Exec(ctx, c.Cmd, c.Args, nil); err != nil {
				//          return err
				//      }
				//  }
				//
				return next.Apply(ctx, conn, plan)
			})
		}))
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return New(db, client)
}

func WithTx(ctx context.Context, rawDB *sql.DB, client *ent.Client, fn func(repo Repositoy) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	repo, err := New(rawDB, tx.Client())
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(repo); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
