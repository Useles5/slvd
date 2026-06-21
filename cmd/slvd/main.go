package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Useles5/slvd/internal/config"
	"github.com/Useles5/slvd/internal/filter"
	"github.com/Useles5/slvd/internal/platform/codeforces"
)

func main() {
	opts := config.Parse()

	subs, err := codeforces.FetchSubmissions(opts.Handle)
	if err != nil {
		log.Fatalf("Failed to fetch Codeforces submissions: %v", err)
	}

	solvedProblems, processed := filter.GetSolvedProblems(subs, opts)
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

	fmt.Fprintf(os.Stderr, "[Mode: %s] fetched=%d processed=%d unique_solved=%d\n", modeStr, len(subs), processed, len(solvedProblems))
}
