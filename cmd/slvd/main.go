package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Useles5/slvd/internal/config"
	"github.com/Useles5/slvd/internal/filter"
	"github.com/Useles5/slvd/internal/models"
	"github.com/Useles5/slvd/internal/platform/atcoder"
	"github.com/Useles5/slvd/internal/platform/codeforces"
)

func main() {
	opts := config.Parse()

	var allSubmissions []models.Submission

	cfSubs, err := codeforces.FetchSubmissions(opts.Handle)
	if err != nil {
		log.Fatalf("Failed to fetch Codeforces submissions: %v", err)
	}

	allSubmissions = append(allSubmissions, cfSubs...)

	acFrom := opts.GetAtCoderSecond()
	acSubs, err := atcoder.FetchSubmissions(opts.Handle, acFrom)
	if err != nil {
		log.Fatalf("Failed to fetch AtCoder submissions: %v", err)
	}

	allSubmissions = append(allSubmissions, acSubs...)

	// Safety check
	if len(allSubmissions) == 0 {
		log.Fatalf("Critical: Could not fetch data from any platform")
	}

	//sort.Slice(allSubmissions, func(i, j int) bool {
	//	return allSubmissions[i].SubmittedAt.After(allSubmissions[j].SubmittedAt)
	//})

	slices.SortFunc(allSubmissions, func(a, b models.Submission) int {
		return b.SubmittedAt.Compare(a.SubmittedAt)
	})

	solvedProblems, processed := filter.GetSolvedProblems(allSubmissions, opts)
	for _, solvedProblem := range solvedProblems {
		fmt.Println(solvedProblem)
	}

	if len(solvedProblems) == 0 {
		fmt.Fprintf(os.Stderr, "No problems found\n")
	}

	modeStr := "Today"
	if opts.Date != "" {
		modeStr = "Date: " + opts.Date
	}
	if opts.Last != -1 {
		modeStr = fmt.Sprintf("Last %d", opts.Last)
	}

	fmt.Fprintf(os.Stderr, "[Mode: %s] fetched=%d processed=%d unique_solved=%d\n", modeStr, len(allSubmissions), processed, len(solvedProblems))
}
