package codeforces

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Useles5/slvd/internal/models"
)

type cfResponse struct {
	Status string         `json:"status"`
	Result []cfSubmission `json:"result"`
}

type cfSubmission struct {
	CreationTimeSeconds int64     `json:"creationTimeSeconds"`
	Verdict             string    `json:"verdict"`
	Problem             cfProblem `json:"problem"`
}

type cfProblem struct {
	ContestID int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
}

func FetchSubmissions(handle string) ([]models.Submission, error) {
	safeUserName := url.QueryEscape(handle)
	apiURL := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=2000", safeUserName)
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
		return nil, fmt.Errorf("codeforces API returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}

	var data cfResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var submissions []models.Submission
	for _, sub := range data.Result {
		submissions = append(submissions, models.Submission{
			Platform:    "Codeforces",
			ProblemKey:  fmt.Sprintf("%d%s", sub.Problem.ContestID, sub.Problem.Index),
			ProblemName: sub.Problem.Name,
			IsAccepted:  sub.Verdict == "OK",
			SubmittedAt: time.Unix(sub.CreationTimeSeconds, 0),
		})
	}

	return submissions, nil
}
