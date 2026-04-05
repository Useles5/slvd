package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct{}

type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type graphqlResponse struct {
	Data struct {
		RecentAcSubmissionList []struct {
			Title     string `json:"title"`
			Timestamp string `json:"timestamp"`
		} `json:"recentAcSubmissionList"`
	} `json:"data"`
}

func (c *Client) FetchRecent(handle string) ([]string, error) {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfDayUnix := startOfDay.Unix()

	query := `query recentAcSubmissions($username: String!, $limit: Int!) {
		recentAcSubmissionList(username: $username, limit: $limit) {
			title
			timestamp
		}
	}`

	reqBody := graphqlRequest{
		Query: query,
		Variables: map[string]interface{}{
			"username": handle,
			"limit":    50, // Fetch up to 50 recent submissions
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result graphqlResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var solvedToday []string

	for _, sub := range result.Data.RecentAcSubmissionList {
		ts, err := strconv.ParseInt(sub.Timestamp, 10, 64)
		if err != nil {
			continue // skip q if parse failed
		}

		if ts < startOfDayUnix {
			break
		}

		solvedToday = append(solvedToday, sub.Title)
	}

	return solvedToday, nil

}
