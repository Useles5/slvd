package main

import (
	"fmt"
	"log"

	"github.com/Useles5/slvd/internal/platform"
	"github.com/Useles5/slvd/internal/platform/codeforces"
	"github.com/Useles5/slvd/internal/platform/leetcode"
)

func main() {
	type job struct {
		provider platform.Provider
		handle   string
	}

	jobs := []job{
		{provider: &codeforces.Client{}, handle: "tourist"},
		{provider: &leetcode.Client{}, handle: "Useles5"},
	}

	//fmt.Println("fetching recent codeforces submissions...")

	var solvedToday []string

	for _, job := range jobs {
		solves, err := job.provider.FetchRecent(job.handle)
		if err != nil {
			log.Printf("Warning: Failed to fetch for %s: %v", job.handle, err)
			continue
		}

		solvedToday = append(solvedToday, solves...)
	}

	if len(solvedToday) == 0 {
		fmt.Printf("No problems solved today. Time to grind!\n")
		return
	}

	fmt.Printf("\n--- Today's Solved Problems---\n")
	for _, problemName := range solvedToday {
		fmt.Printf("%s\n", problemName)
	}
}
