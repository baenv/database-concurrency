//go:build ignore

package main

import (
	"context"

	"fmt"
	"log"
	"os"

	"database-concurrency/ent/migrate"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {

	fmt.Println(os.Args)
	ctx := context.Background()

	dir, err := atlas.NewLocalDir("ent/migrate/migrations")
	if err != nil {
		logrus.Fatal(err, "failed to create local dir")
	}

	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                          // provide migration directory
		schema.WithMigrationMode(schema.ModeInspect), // provide migration mode
		schema.WithDialect(dialect.Postgres),         // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	if len(os.Args) != 3 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <db-uri> <name>'")
	}

	if err = migrate.NamedDiff(ctx, os.Args[1], os.Args[2], opts...); err != nil {
		logrus.Fatal(err, "failed to migrate")
	}
}
