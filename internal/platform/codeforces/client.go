package codeforces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CFResponse struct {
	Status string         `json:"status"`
	Result []CFSubmission `json:"result"`
}

type CFSubmission struct {
	CreationTimeSeconds int64     `json:"creation_time_seconds"`
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
	defer resp.Body.Close()

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

	fmt.Printf("Recent Submissions for %s:\n", handle)
	fmt.Println("--------------------------------")

	for _, sub := range data.Result {
		fmt.Printf("Problem: %s | Verdict: %s\n", sub.Problem.Name, sub.Verdict)
	}
}
