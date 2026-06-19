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
	ID                  int     `json:"id"`
	CreationTimeSeconds int     `json:"creationTimeSeconds"`
	Verdict             Verdict `json:"verdict,omitempty"`
	Problem             Problem `json:"problem"`
}

type Problem struct {
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

	foundSolved := false

	for _, sub := range cfData.Result {
		if sub.Verdict == VerdictOK {
			foundSolved = true

			fmt.Printf("%d%s - %s\n",
				sub.Problem.ContestID,
				sub.Problem.Index,
				sub.Problem.Name)
		}
	}

	if !foundSolved {
		fmt.Println("No Successful submissions found in given range")
	}
}
