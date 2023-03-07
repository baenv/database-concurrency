package main

import (
	"database-concurrency/scenario"
	"fmt"
	"os"

	pewpew "github.com/bengadbois/pewpew/lib"
)

type TestCase struct {
	ReqBody string `json:"req_body"`
	API_URL string `json:"api_url"`
}

func main() {
	tcs := map[string]TestCase{
		"without_unique_id": {
			ReqBody: `{}`,
			API_URL: "api/v1/tickets/create",
		},
		"with_unique_id": {
			ReqBody: `{
				"unique_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"
			}`,
			API_URL: "api/v2/tickets/create",
		},
	}

	for name, tc := range tcs {
		runScenario(name, tc)
		scenario.Seed("./scenario/spam_creating/testdata/seed.down.sql")
	}
}

func runScenario(scrName string, testCase TestCase) {
	fmt.Printf("test scenario SPAM CREATING for case <%s> started \n", scrName)
	defer fmt.Printf("test scenario SPAM CREATING for case <%s> completed \n", scrName)
	stressCfg := pewpew.StressConfig{
		Count:       100,
		Concurrency: 100,
		Verbose:     false,
		Targets: []pewpew.Target{{
			URL:     fmt.Sprintf("http://127.0.0.1:8083/%s", testCase.API_URL),
			Timeout: "5s",
			Method:  "POST",
			Body:    testCase.ReqBody,
		}},
	}

	logFile, err := os.Create(fmt.Sprintf("scenario/spam_creating/logs/%s.log", scrName))
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
