package filter

import (
	"fmt"
	"time"

	"github.com/Useles5/slvd/internal/config"
	"github.com/Useles5/slvd/internal/models"
)

// GetSolvedProblems is the pipeline Orchestrator
func GetSolvedProblems(submissions []models.Submission, opts *config.Options) ([]models.Submission, int) {

	// Clean
	data := keepOnlyAccepted(submissions)

	// Time bound
	if opts.Date != "" || opts.Last == -1 {
		data = filterByDate(data, opts.Date)
	}

	// Limit and Format
	return applyLimitAndDeDuplicate(data, opts.Last)
}

// keepOnlyAccepted removes failed submissions
func keepOnlyAccepted(submission []models.Submission) []models.Submission {
	var validSubmissions []models.Submission
	for _, sub := range submission {
		if sub.IsAccepted {
			validSubmissions = append(validSubmissions, sub)
		}
	}
	return validSubmissions
}

// filterByDate restricts the data to a specific Date
func filterByDate(submissions []models.Submission, dateStr string) []models.Submission {
	now := time.Now()
	var startBound, endBound time.Time

	if dateStr != "" {
		parsedDate, err := time.ParseInLocation("02-01-2006", dateStr, now.Location())
		if err != nil {
			fmt.Println("Warning: Invalid date format. Defaulting to Today.")
			startBound = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		} else {
			startBound = parsedDate
		}
	} else {
		startBound = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	endBound = startBound.AddDate(0, 0, 1)

	var filteredSubmissions []models.Submission
	for _, sub := range submissions {

		// new -------------> old data
		// A B C || D E F || G H
		//        ^        ^
		//        |        |
		//       end      start

		// thank u cf

		if sub.SubmittedAt.Before(startBound) {
			break
		}

		if sub.SubmittedAt.After(endBound) || sub.SubmittedAt.Equal(endBound) {
			continue
		}
		filteredSubmissions = append(filteredSubmissions, sub)
	}
	return filteredSubmissions
}

// applyLimitAndDeDuplicate handles unique submissions and the N cutoff (--last flag)
func applyLimitAndDeDuplicate(submissions []models.Submission, limit int) ([]models.Submission, int) {
	seen := make(map[string]struct{})
	var solved []models.Submission
	processed := 0

	for _, sub := range submissions {

		// Stop if limit is provided and we reach it
		if limit != -1 && len(seen) == limit {
			break
		}

		processed++

		// Incase diff sources have same problemKey
		uniqueKey := fmt.Sprintf("%s:%s", sub.Platform, sub.ProblemKey)
		// Skip duplicates
		if _, exists := seen[uniqueKey]; exists {
			continue
		}

		seen[uniqueKey] = struct{}{}

		solved = append(solved, sub)
	}

	return solved, processed

}
