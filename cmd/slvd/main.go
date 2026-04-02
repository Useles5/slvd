package main

import (
	"fmt"
	"log"

	"github.com/Useles5/slvd/internal/platform/codeforces"
)

func main() {
	handle := "tourist"
	//fmt.Println("fetching recent codeforces submissions...")
	solves, err := codeforces.FetchRecent(handle)
	if err != nil {
		log.Fatalf("Critical error: %v\n", err)
	}

	if len(solves) == 0 {
		fmt.Printf("No problems solved today for %s. Time to grind!\n", handle)
		return
	}

	fmt.Printf("\n--- Today's Solved Problems for %s ---\n", handle)
	for _, problemName := range solves {
		fmt.Printf("✅ [Codeforces] %s\n", problemName)
	}
}
