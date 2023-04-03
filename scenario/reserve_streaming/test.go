package main

import (
	"database-concurrency/scenario"
	"fmt"
	"os"

	pewpew "github.com/bengadbois/pewpew/lib"
)

type TestCase struct {
	ReqAPIURL string `json:"req_url"`
	ReqBody   string `json:"req_body"`
}

func main() {
	tcs := map[string]TestCase{
		"normal": {
			ReqAPIURL: "v1/tickets/reserve",
			ReqBody: `{
				"ticket_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
				"user_id": "91272a62-c537-42ed-948c-bb2a91af2051"
			}`,
		},
		"with_redis_stream": {
			ReqAPIURL: "v2/tickets/reserve",
			ReqBody: `{
				"ticket_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
				"user_id": "91272a62-c537-42ed-948c-bb2a91af2051"
			}`,
		},
	}

	for name, tc := range tcs {
		scenario.Seed("./scenario/reserve_streaming/testdata/seed.up.sql")
		runScenario(name, tc)
		scenario.Seed("./scenario/reserve_streaming/testdata/seed.down.sql")
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
			URL:     fmt.Sprintf("http://127.0.0.1:8083/api/%s", testCase.ReqAPIURL),
			Timeout: "2s",
			Method:  "POST",
			Body:    testCase.ReqBody,
		}},
	}

	logFile, err := os.Create(fmt.Sprintf("scenario/reserve_streaming/logs/%s.log", scrName))
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
