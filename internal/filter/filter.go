package filter

import (
	"fmt"
	"time"

	"github.com/Useles5/slvd/internal/models"
)

func GetSolvedProblems(submission []models.Submission, lastN int) ([]string, int) {
	seen := make(map[string]struct{})
	var solved []string

	now := time.Now()
	midnightUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	processed := 0

	for _, sub := range submission {

		if lastN > 0 {
			if len(seen) == lastN {
				break
			}
		} else {
			// API sends newest first.
			// Once we hit yesterday's problem, we can safely break.
			if sub.SubmittedAt.Before(midnightUnix) {
				break
			}
		}

		processed++

		if !sub.IsAccepted {
			continue
		}

		if _, exists := seen[sub.ProblemKey]; exists {
			continue
		}

		seen[sub.ProblemKey] = struct{}{}

		formattedString := fmt.Sprintf("%s - %s", sub.ProblemKey, sub.ProblemName)
		solved = append(solved, formattedString)

	}

	return solved, processed
}
