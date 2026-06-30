package main

import (
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/Useles5/slvd/internal/config"
	"github.com/Useles5/slvd/internal/filter"
	"github.com/Useles5/slvd/internal/models"
	"github.com/Useles5/slvd/internal/platform/atcoder"
	"github.com/Useles5/slvd/internal/platform/codeforces"
	"github.com/Useles5/slvd/internal/platform/leetcode"
	"github.com/Useles5/slvd/internal/printer"
)

func main() {
	opts := config.Parse()

	var wg sync.WaitGroup

	var mu sync.Mutex
	var allSubmissions []models.Submission

	fetch := func(platform string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handle := opts.Handles[platform]

			var subs []models.Submission
			var err error

			switch platform {
			case "codeforces":
				subs, err = codeforces.FetchSubmissions(handle)
			case "atcoder":
				acFrom := opts.GetAtCoderSecond()
				subs, err = atcoder.FetchSubmissions(handle, acFrom)
			case "leetcode":
				subs, err = leetcode.FetchSubmissions(handle)
			default:
				log.Printf("Warning: Unknown platform: %s", platform)
				return
			}

			if err != nil {
				log.Printf("Warning: Failed to fetch %s submissions: %v", platform, err)
				return
			}

			mu.Lock()
			allSubmissions = append(allSubmissions, subs...)
			mu.Unlock()

		}()
	}

	for _, platform := range opts.Platforms {
		fetch(platform)
	}
	wg.Wait()

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

	modeStr := "Today"
	if opts.Date != "" {
		modeStr = "Date: " + opts.Date
	}
	if opts.Last != -1 {
		modeStr = fmt.Sprintf("Last %d", opts.Last)
	}

	printer.PrintTable(solvedProblems, modeStr, len(allSubmissions), processed, opts.Markdown)
}
