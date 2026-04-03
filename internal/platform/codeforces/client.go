package codeforces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CFResponse struct {
	Status string         `json:"status"`
	Result []CFSubmission `json:"result"`
}

type CFSubmission struct {
	CreationTimeSeconds int64     `json:"creationTimeSeconds"`
	Verdict             string    `json:"verdict"`
	Problem             CFProblem `json:"problem"`
}

type CFProblem struct {
	Name string `json:"name"`
}

func FetchRecent(handle string) ([]string, error) {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfDayUnix := startOfDay.Unix()

	var solvedToday []string

	fromIndex := 1 // Cf uses 1 based indexing
	pageSize := 50 // 50 submissions per page

	for {
		url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=%d&count=%d", handle, fromIndex, pageSize)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch codeforces data: %v", err) // no need of '\n' as log.Fatalf automatically adds it
		}

		body, err := io.ReadAll(resp.Body)

		_ = resp.Body.Close() // close the connection

		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v\n", err)
		}

		var data CFResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %v\n", err)
		}

		// user has no more submissions
		if len(data.Result) == 0 {
			break
		}

		keepFetching := true

		// the list of submission is sorted, so newer appears first
		for _, sub := range data.Result {
			// If we hit a submission from yesterday, stop processing this page,
			// and tell the network loop to stop fetching more pages.
			if sub.CreationTimeSeconds < startOfDayUnix {
				keepFetching = false
				break
			}

			if sub.Verdict == "OK" {
				solvedToday = append(solvedToday, sub.Problem.Name)
			}
		}

		if !keepFetching {
			break
		}

		fromIndex += pageSize
	}

	return solvedToday, nil

}
