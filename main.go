package main

import (
	"encoding/json"
	"flag"
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

type Options struct {
	Handle string
	Last   int
}

func parseOptions() *Options {
	lastFlag := flag.Int("last", 0, "Fetches N recent successful submissions")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Usage: slvd [--last N] <codeforces-handle>")
	}

	return &Options{Handle: args[0], Last: *lastFlag}
}

func fetchCodeforcesSubmissions(opts *Options) ([]CFSubmission, error) {
	safeUserName := url.QueryEscape(opts.Handle)
	apiURL := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=100", safeUserName)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err) // %w to wrap a error
	}
	// Passing Body as an argument takes a snapshot of the variable right now.
	// This prevents dangerous "changing variable" bugs if used inside a loop.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("warning: failed to cleanly close: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}

	var cfData CFResponse
	err = json.NewDecoder(resp.Body).Decode(&cfData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return cfData.Result, nil
}

func getSolvedProblems(submission []CFSubmission, lastN int) ([]string, int) {
	seen := make(map[string]struct{})
	var solved []string

	now := time.Now()
	midnightUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	processed := 0

	for _, sub := range submission {

		if lastN > 0 {
			if len(seen) == lastN {
				break
			}
		} else {
			// API sends newest first.
			// Once we hit yesterday's problem, we can safely break.
			if sub.CreationTimeSeconds < midnightUnix {
				break
			}
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

		formattedString := fmt.Sprintf("%s - %s", problemKey, sub.Problem.Name)
		solved = append(solved, formattedString)

	}

	return solved, processed
}
func main() {
	opts := parseOptions()

	subs, err := fetchCodeforcesSubmissions(opts)
	if err != nil {
		log.Fatalf("Failed to fetch Codeforces submissions: %v", err)
	}

	solvedProblems, processed := getSolvedProblems(subs, opts.Last)
	for _, solvedProblem := range solvedProblems {
		fmt.Println(solvedProblem)
	}

	if len(solvedProblems) == 0 {
		fmt.Fprintf(os.Stderr, "No problems found\n")
	}

	if opts.Last > 0 {
		fmt.Fprintf(os.Stderr, "[Mode: Last %d] fetched=%d processed=%d unique_solved=%d\n", opts.Last, len(subs), processed, len(solvedProblems))
	} else {
		fmt.Fprintf(os.Stderr, "[Mode: Today]  fetched=%d processed=%d unique_solved=%d\n", len(subs), processed, len(solvedProblems))
	}

}
