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

	solvedProblems, processed := filter.GetSolvedProblems(subs, opts.Last)
	for _, solvedProblem := range solvedProblems {
		fmt.Println(solvedProblem)
	}

	if len(solvedProblems) == 0 {
		fmt.Fprintf(os.Stderr, "No problems found\n")
	}

	if opts.Last > 0 {
		fmt.Fprintf(os.Stderr, "[Mode: Last %d] fetched=%d processed=%d unique_solved=%d\n", opts.Last, len(subs), processed, len(solvedProblems))
	} else {
		fmt.Fprintf(os.Stderr, "[Mode: Today]  fetched=%d processed=%d unique_solved=%d\n", len(subs), processed, len(solvedProblems))
	}

}
