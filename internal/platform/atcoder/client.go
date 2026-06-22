package atcoder

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

type acSubmission struct {
	EpochSecond int64  `json:"epoch_second"`
	ProblemID   string `json:"problem_id"`
	Result      string `json:"result"`
}

type acProblem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FetchSubmissions(handle string, fromSecond int64) ([]models.Submission, error) {
	problemMap, err := fetchProblemMap()
	if err != nil {
		fmt.Printf("Warning: Failed to load AtCoder problem names: %v\n", err)
	}

	safeUserName := url.QueryEscape(handle)
	
	apiURL := fmt.Sprintf("https://kenkoooo.com/atcoder/atcoder-api/v3/user/submissions?user=%s&from_second=%d", safeUserName, fromSecond)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("warning: failed to cleanly close: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}

	var data []acSubmission
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var submissions []models.Submission

	for i := len(data) - 1; i >= 0; i-- {
		sub := data[i]

		problemName := sub.ProblemID
		if name, exists := problemMap[problemName]; exists {
			problemName = name
		}

		submissions = append(submissions, models.Submission{
			Platform:    "AtCoder",
			ProblemKey:  sub.ProblemID,
			ProblemName: problemName,
			IsAccepted:  sub.Result == "AC",
			SubmittedAt: time.Unix(sub.EpochSecond, 0),
		})
	}

	return submissions, nil
}

func fetchProblemMap() (map[string]string, error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get("https://kenkoooo.com/atcoder/resources/problems.json")
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
		return nil, fmt.Errorf("API returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}

	var problems []acProblem
	err = json.NewDecoder(resp.Body).Decode(&problems)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	problemMap := make(map[string]string)
	for _, p := range problems {
		problemMap[p.ID] = p.Name
	}

	return problemMap, nil
}
