package main

import (
	"github.com/Useles5/slvd/internal/platform/codeforces"
)

func main() {
	//fmt.Println("fetching recent codeforces submissions...")
	codeforces.FetchRecent("tourist")
}
