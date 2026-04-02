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
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=10", handle)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codeforces data: %v\n", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v\n", err)
	}

	var data CFResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v\n", err)
	}

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfDayUnix := startOfDay.Unix()

	var solvedToday []string

	for _, sub := range data.Result {
		if sub.CreationTimeSeconds < startOfDayUnix {
			break
		}

		if sub.Verdict == "OK" {
			solvedToday = append(solvedToday, sub.Problem.Name)
		}
	}

	return solvedToday, nil

}
