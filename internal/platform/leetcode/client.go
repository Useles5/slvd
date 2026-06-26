package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Useles5/slvd/internal/models"
)

type graphqlQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type leetcodeResponse struct {
	Data struct {
		RecentAcSubmissionList []struct {
			Title     string `json:"title"`
			TitleSlug string `json:"titleSlug"`
			Timestamp string `json:"timestamp"`
		} `json:"recentAcSubmissionList"`
	} `json:"data"`
}

func FetchSubmissions(handle string) ([]models.Submission, error) {
	query := `
		query recentAcSubmissions($username: String!, $limit: Int!) {
			recentAcSubmissionList(username: $username, limit: $limit) {
				title 
				titleSlug
				timestamp
			}				
		}
	`
	payload := graphqlQuery{
		Query: query,
		Variables: map[string]interface{}{
			"username": handle,
			"limit":    20,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Using NewRequest instead of http.Post as this endpoint requires custom headers (Content-Type, User-Agent) to bypass bot protection.
	// Wrap the JSON payload in an io.Reader stream so the HTTP client can stream it to the socket
	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:152.0) Gecko/20100101 Firefox/152.0")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("warning: failed to cleanly close: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("leetcode API returned non-200 status: %s", resp.Status)
	}

	var lcResp leetcodeResponse
	err = json.NewDecoder(resp.Body).Decode(&lcResp)
	if err != nil {
		return nil, err
	}

	var submissions []models.Submission
	for _, sub := range lcResp.Data.RecentAcSubmissionList {

		// Leetcode timestamp is a string Unix seconds, so convert it.
		tsInt, _ := strconv.ParseInt(sub.Timestamp, 10, 64)

		submissions = append(submissions, models.Submission{
			Platform:    "LeetCode",
			ProblemKey:  sub.TitleSlug,
			ProblemName: sub.Title,
			IsAccepted:  true,
			SubmittedAt: time.Unix(tsInt, 0),
		})
	}

	return submissions, nil
}
