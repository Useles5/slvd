package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sync"

	"github.com/Useles5/slvd/internal/config"
	"github.com/Useles5/slvd/internal/filter"
	"github.com/Useles5/slvd/internal/models"
	"github.com/Useles5/slvd/internal/platform/atcoder"
	"github.com/Useles5/slvd/internal/platform/codeforces"
	"github.com/Useles5/slvd/internal/platform/leetcode"
)

func main() {
	opts := config.Parse()

	fetchAll := !opts.CF && !opts.ATC && !opts.LC

	var wg sync.WaitGroup

	var cfSubmissions []models.Submission
	var atcSubmissions []models.Submission
	var lcSubmissions []models.Submission

	if fetchAll || opts.CF {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfSubs, err := codeforces.FetchSubmissions(opts.Handle)
			if err != nil {
				log.Printf("Warning: Failed to fetch Codeforces submissions: %v", err)
				return
			}

			cfSubmissions = cfSubs
		}()

	}

	if fetchAll || opts.ATC {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acFrom := opts.GetAtCoderSecond()
			acSubs, err := atcoder.FetchSubmissions(opts.Handle, acFrom)
			if err != nil {
				log.Printf("Warning: Failed to fetch AtCoder submissions: %v", err)
			}

			atcSubmissions = acSubs
		}()

	}

	if fetchAll || opts.LC {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lcSubs, err := leetcode.FetchSubmissions(opts.Handle)
			if err != nil {
				log.Printf("Warning: Failed to fetch LeetCode submissions: %v", err)
			}

			lcSubmissions = lcSubs
		}()
	}

	wg.Wait()

	var allSubmissions []models.Submission
	allSubmissions = append(allSubmissions, cfSubmissions...)
	allSubmissions = append(allSubmissions, atcSubmissions...)
	allSubmissions = append(allSubmissions, lcSubmissions...)

	// Safety check
	if len(allSubmissions) == 0 {
		log.Fatalf("Critical: Could not fetch data from any platform")
	}

	//sort.Slice(allSubmissions, func(i, j int) bool {
	//	return allSubmissions[i].SubmittedAt.After(allSubmissions[j].SubmittedAt)
	//})

	// Sorting
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
