package main

import (
	"database-concurrency/scenario"
	"fmt"
	"os"

	pewpew "github.com/bengadbois/pewpew/lib"
)

type TestCase struct {
	ReqBody string `json:"req_body"`
}

func main() {
	tcs := map[string]TestCase{
		"without_lock": {
			ReqBody: `{
				"ticket_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
				"user_id": "91272a62-c537-42ed-948c-bb2a91af2051"
			}`,
		},
		"with_lock": {
			ReqBody: `{
				"ticket_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
				"user_id": "91272a62-c537-42ed-948c-bb2a91af2051",
				"locks": {
					"for_update": true
				}
			}`,
		},
		"with_session_advisory_lock": {
			ReqBody: `{
				"ticket_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
				"user_id": "91272a62-c537-42ed-948c-bb2a91af2051",
				"locks": {
					"session_advisory_lock": true
				}
			}`,
		},
	}

	for name, tc := range tcs {
		scenario.Seed("./scenario/spam_booking/testdata/seed.up.sql")
		runScenario(name, tc)
		scenario.Seed("./scenario/spam_booking/testdata/seed.down.sql")
	}
}

func runScenario(scrName string, testCase TestCase) {
	fmt.Printf("test scenario SPAM BOOKING for case <%s> started \n", scrName)
	defer fmt.Printf("test scenario SPAM BOOKING for case <%s> completed \n", scrName)
	stressCfg := pewpew.StressConfig{
		Count:       100,
		Concurrency: 100,
		Verbose:     false,
		Targets: []pewpew.Target{{
			URL:     "http://127.0.0.1:8083/api/v1/tickets/book",
			Timeout: "2s",
			Method:  "POST",
			Body:    testCase.ReqBody,
		}},
	}

	logFile, err := os.Create(fmt.Sprintf("scenario/spam_booking/logs/%s.log", scrName))
	if err != nil {
		fmt.Printf("failed to create log file: %s", err.Error())
		return
	}
	defer logFile.Close()

	logFile.WriteString("---- Requests ----\n\n")

	stats, err := pewpew.RunStress(stressCfg, logFile)
	if err != nil {
		fmt.Printf("pewpew stress failed: %s", err.Error())
		return
	}

	logFile.WriteString("\n---- Summary ----\n")

	for _, s := range stats {
		summary := pewpew.CreateRequestsStats(s)
		textSummary := pewpew.CreateTextSummary(summary)
		logFile.WriteString(textSummary)
	}
}
