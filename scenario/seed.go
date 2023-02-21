package scenario

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

func Seed(filePath string) {
	// Connect to the database
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load the seed file
	err = loadSeedFiles(db, filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Seed file %s loaded successfully", filePath)
}

func loadSeedFiles(db *sql.DB, filePath string) error {
	// Read the seed file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Execute the seed file in a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(content))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
