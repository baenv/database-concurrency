package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"database-concurrency/config"
	"database-concurrency/ent"
	transactionRepo "database-concurrency/internal/repository/transaction"

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

	return New(client)
}
