package main

import (
	"database-concurrency/scenario"
	"fmt"
	"os"

	pewpew "github.com/bengadbois/pewpew/lib"
)

func main() {
	scenario.Seed("./scenario/spam_booking/testdata/seed.up.sql")
	defer scenario.Seed("./scenario/spam_booking/testdata/seed.down.sql")

	runScenario()
}

func runScenario() {
	fmt.Println("test scenario SPAM BOOKING started")
	defer fmt.Println("test scenario SPAM BOOKING completed")
	stressCfg := pewpew.StressConfig{
		Count:       100,
		Concurrency: 100,
		Verbose:     false,
		Targets: []pewpew.Target{{
			URL:     "http://127.0.0.1:8083/api/v1/tickets/book",
			Timeout: "2s",
			Method:  "POST",
			Body:    `{"ticket_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d", "user_id": "91272a62-c537-42ed-948c-bb2a91af2051"}`,
		}},
	}

	logFile, err := os.Create("scenario/spam_booking/test.log")
	if err != nil {
		fmt.Printf("failed to create log file: %s", err.Error())
		return
	}
	defer logFile.Close()

	_, err = pewpew.RunStress(stressCfg, logFile)
	if err != nil {
		fmt.Printf("pewpew stress failed: %s", err.Error())
		return
	}
}
