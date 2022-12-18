package repository

import (
	"context"
	"database/sql"
	"fmt"

	"database-concurrency/config"
	"database-concurrency/ent"
	transactionRepo "database-concurrency/internal/repository/transaction"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

// The general Repository that contains all other model-repositories
type Repositoy interface {
	Pg() *ent.Client
	Transaction() transactionRepo.Repository
}

// New is used to create new repository
func New(pg *ent.Client) (Repositoy, error) {
	return &db{
		pg:          pg,
		transaction: transactionRepo.NewPGRepo(pg),
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

	// Migrate UP to latest
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return New(client)
}
