package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Verdict string

const (
	VerdictOK                Verdict = "OK"
	VerdictTimeLimitExceeded Verdict = "TIME_LIMIT_EXCEEDED"
	VerdictWrongAnswer       Verdict = "WRONG_ANSWER"
	// for now just these
)

type CFResponse struct {
	Status string         `json:"status"`
	Result []CFSubmission `json:"result"`
}

type CFSubmission struct {
	ID                  int       `json:"id"`
	CreationTimeSeconds int64     `json:"creationTimeSeconds"`
	Verdict             Verdict   `json:"verdict"`
	Problem             CFProblem `json:"problem"`
}

type CFProblem struct {
	ContestID int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go userName")
	}

	safeUserName := url.QueryEscape(os.Args[1])
	apiURL := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=100", safeUserName)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}
	// Passing Body as an argument takes a snapshot of the variable right now.
	// This prevents dangerous "changing variable" bugs if used inside a loop.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}

	var cfData CFResponse
	err = json.NewDecoder(resp.Body).Decode(&cfData)
	if err != nil {
		log.Fatalf("Failed to parse json: %v", err)
	}

	seen := make(map[string]struct{})

	now := time.Now()
	midnightUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	processed := 0

	for _, sub := range cfData.Result {

		// API sends newest first.
		// Once we hit yesterday's problem, we can safely break.
		if sub.CreationTimeSeconds < midnightUnix {
			break
		}

		processed++
		
		if sub.Verdict != VerdictOK {
			continue
		}

		problemKey := fmt.Sprintf("%d%s", sub.Problem.ContestID, sub.Problem.Index)

		if _, exists := seen[problemKey]; exists {
			continue
		}

		seen[problemKey] = struct{}{}

		fmt.Printf("%s - %s\n", problemKey, sub.Problem.Name)

	}

	if len(seen) == 0 {
		fmt.Fprintf(os.Stderr, "No problems found\n")
	}

	fmt.Fprintf(os.Stderr, "fetched=%d processed_today=%d unique_solved=%d\n", len(cfData.Result), processed, len(seen))
}
