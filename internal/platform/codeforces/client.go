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

func FetchRecent(handle string) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=10", handle)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch codeforces data: %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("warning: failed to close body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v\n", err)
		return
	}

	var data CFResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("failed to unmarshal response body: %v\n", err)
	}

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfDayUnix := startOfDay.Unix()

	solvedCount := 0

	for _, sub := range data.Result {
		if sub.CreationTimeSeconds < startOfDayUnix {
			break
		}

		if sub.Verdict == "OK" {
			if solvedCount == 0 {
				fmt.Printf("Today's Solved Problem for %s:\n", handle)
			}
			fmt.Printf("Problem: %s\n", sub.Problem.Name)
			solvedCount++
		}
	}

	if solvedCount == 0 {
		fmt.Println("No problems solved today")
	}
}
